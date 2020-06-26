// This is to generate documentation website

//go:generate go run gen.go

package main

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path"
	"sort"
	"strings"
	"text/template"

	"github.com/hashicorp/hcl2/hclwrite"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-quorum/quorum"
)

const (
	basedir       = "../quorum"
	targetDocsDir = "docs"
	pageTemplate  = "templates/page.html.markdown.tmpl"
	navTemplate   = "templates/quorum.erb.tmpl"
	indexTemplate = "templates/index.html.markdown.tmpl"
)

func main() {
	// read all data sources and resources Go source files
	meta := quorum.Provider()
	log.Println("Building data source pages...")
	dsIndex, err := generate("data", "d", meta.DataSourcesMap)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Building resource pages...")
	resIndex, err := generate("resource", "r", meta.ResourcesMap)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Building navigation page...")
	if err := generateNavDocs(dsIndex, resIndex); err != nil {
		log.Fatal(err)
	}
	log.Println("Building index page...")
	if err := generateIndexDocs(path.Join(basedir, "provider.go")); err != nil {
		log.Fatal(err)
	}
}

func generateIndexDocs(providerSourceFile string) error {
	f, err := parser.ParseFile(token.NewFileSet(), providerSourceFile, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	t, err := template.ParseFiles(indexTemplate)
	if err != nil {
		return err
	}
	out, err := os.OpenFile(path.Join(targetDocsDir, "index.html.markdown"), os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	ctx := index{
		ShortDescription: padding(f.Doc.Text(), "   "),
	}
	for _, d := range f.Decls {
		if fd, ok := d.(*ast.FuncDecl); ok {
			if fd.Name.Name == "Provider" {
				ctx.LongDescription = strings.TrimSpace(fd.Doc.Text())
				break
			}
		}
	}
	if err := t.Execute(out, ctx); err != nil {
		return err
	}
	return nil
}

func generate(prefix string, targetDir string, m map[string]*schema.Resource) ([]*nav, error) {
	navs := make([]*nav, 0)
	for n, s := range m {
		sourceFile := fmt.Sprintf("%s_%s.go", prefix, n)
		testSourceFile := fmt.Sprintf("%s_%s_test.go", prefix, n)
		index, err := generatePageDoc(n, path.Join(basedir, sourceFile), path.Join(basedir, testSourceFile), path.Join(targetDocsDir, targetDir), s.Schema)
		if err != nil {
			return nil, err
		}
		navs = append(navs, index)
	}
	sort.Slice(navs, func(i, j int) bool {
		return navs[i].Name < navs[j].Name
	})
	return navs, nil
}

func generateNavDocs(dsIndex []*nav, resIndex []*nav) error {
	t, err := template.ParseFiles(navTemplate)
	if err != nil {
		return err
	}
	out, err := os.OpenFile("quorum.erb", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	if err := t.Execute(out, struct {
		DataSources []*nav
		Resources   []*nav
	}{
		DataSources: dsIndex,
		Resources:   resIndex,
	}); err != nil {
		return err
	}
	return nil
}

func generatePageDoc(pageName string, sourceFile string, testSourceFile string, targetDir string, dsSchema map[string]*schema.Schema) (*nav, error) {
	log.Println("Parsing", sourceFile)
	_ = os.MkdirAll(targetDir, 0755)
	shortName := strings.TrimPrefix(pageName, "quorum_")
	f, err := parser.ParseFile(token.NewFileSet(), sourceFile, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	testFile, err := parser.ParseFile(token.NewFileSet(), testSourceFile, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	t, err := template.ParseFiles(pageTemplate)
	if err != nil {
		return nil, err
	}
	out, err := os.OpenFile(path.Join(targetDir, fmt.Sprintf("%s.html.markdown", shortName)), os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	ctx := page{
		Name:    pageName,
		Title:   pageName,
		SideBar: fmt.Sprintf("docs-%s", strings.ReplaceAll(pageName, "_", "-")),
		Inputs:  make([]*pageInputSection, 0),
		Outputs: make([]*pageOutputSection, 0),
	}
	for _, d := range f.Decls {
		if f, ok := d.(*ast.FuncDecl); ok {
			desc := strings.TrimSpace(f.Doc.Text())
			if len(desc) > 0 {
				ctx.ShortDescription = padding(desc, "   ")
				ctx.LongDescription = desc
			}
		}
	}
	for _, d := range testFile.Decls {
		if f, ok := d.(*ast.FuncDecl); ok {
			desc := strings.TrimSpace(f.Doc.Text())
			if desc == "@example" {
				log.Println("Extracting example from", testSourceFile, f.Name.Name)
				v := exampleCaptor("")
				ast.Walk(&v, d)
				ctx.Example = fmt.Sprintf("```hcl"+`
%s
`+"```", v.Code())
			}
		}
	}
	for field, fschema := range dsSchema {
		if fschema.Computed { // output
			ctx.Outputs = append(ctx.Outputs, &pageOutputSection{
				Name:        field,
				Description: fschema.Description,
			})
		} else { // input
			flag := "Optional"
			if fschema.Required {
				flag = "Required"
			}
			section := &pageInputSection{
				Name:        field,
				Description: fschema.Description,
				Flag:        flag,
			}
			if fschema.Type == schema.TypeSet || fschema.Type == schema.TypeList {
				if elm, ok := fschema.Elem.(*schema.Resource); ok {
					section.Object = &pageInputObject{
						Description: "Each `" + field + "` supports the following\n",
						Fields:      make([]*pageInputSection, 0),
					}
					keys := make([]string, 0)
					for k := range elm.Schema {
						keys = append(keys, k)
					}
					sort.Strings(keys)
					for _, field := range keys {
						fschema := elm.Schema[field]
						flag := ""
						if fschema.Required {
							flag = "(Required)"
						}
						if fschema.Optional {
							flag = "(Optional)"
						}
						section.Object.Fields = append(section.Object.Fields, &pageInputSection{
							Name:        field,
							Flag:        flag,
							Description: fschema.Description,
							Object:      nil,
						})
					}
				}
			}
			ctx.Inputs = append(ctx.Inputs, section)
		}
	}
	sort.Slice(ctx.Inputs, func(i, j int) bool {
		return ctx.Inputs[i].Name < ctx.Inputs[j].Name
	})
	sort.Slice(ctx.Outputs, func(i, j int) bool {
		return ctx.Outputs[i].Name < ctx.Outputs[j].Name
	})
	if err := t.Execute(out, ctx); err != nil {
		return nil, err
	}
	return &nav{
		SideBarCurrent: ctx.SideBar,
		PageName:       shortName,
		Name:           ctx.Name,
	}, nil
}

type index struct {
	ShortDescription, LongDescription string
}

type nav struct {
	SideBarCurrent string
	PageName       string
	Name           string
}

type page struct {
	Title            string
	Name             string
	SideBar          string
	ShortDescription string
	LongDescription  string
	Example          string
	Inputs           []*pageInputSection
	Outputs          []*pageOutputSection
}

type pageInputObject struct {
	Description string
	Fields      []*pageInputSection
}

type pageInputSection struct {
	Name, Flag, Description string
	Object                  *pageInputObject
}

type pageOutputSection struct {
	Name, Description string
}

type exampleCaptor string

func (c *exampleCaptor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}
	if s, ok := n.(*ast.KeyValueExpr); ok {
		if fmt.Sprintf("%s", s.Key) == "Config" {
			switch s.Value.(type) {
			case *ast.BasicLit: // support vanilla string
				*c = exampleCaptor(s.Value.(*ast.BasicLit).Value)
				return nil
			case *ast.CallExpr: // support fmt.Sprintf("<example>", ...)
				ce := s.Value.(*ast.CallExpr)
				*c = exampleCaptor(ce.Args[0].(*ast.BasicLit).Value)
				return nil
			}
		}
	}
	return c
}

func (c *exampleCaptor) String() string {
	return string(*c)
}

func (c *exampleCaptor) Code() string {
	raw := strings.TrimSpace(strings.Trim(c.String(), "`"))
	return string(hclwrite.Format([]byte(raw)))
}

func padding(s string, p string) string {
	out := ""
	scanner := bufio.NewScanner(strings.NewReader(strings.TrimSpace(s)))
	for scanner.Scan() {
		out = fmt.Sprintf("%s\n%s%s", out, p, scanner.Text())
	}
	return strings.TrimSpace(out)
}
