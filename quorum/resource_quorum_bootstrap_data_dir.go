package quorum

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	"github.com/ethereum/go-ethereum/node"

	"github.com/ethereum/go-ethereum/core"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Use this resource to create a data dir locally. This equivalent to execute `geth init`.
func resourceBootstrapDataDir() *schema.Resource {
	return &schema.Resource{
		Create: resourceBootstrapDataDirCreate,
		Read:   resourceBootstrapDataDirRead,
		Delete: resourceBootstrapDataDirDelete,

		Schema: map[string]*schema.Schema{
			"data_dir": {
				Type:        schema.TypeString,
				Description: "Directory to intialize a genesis block",
				Required:    true,
				ForceNew:    true,
			},
			"instance_name": {
				Type:        schema.TypeString,
				Description: "The instance name of the node. This must be the same as the value in geth node config. Default is `geth`",
				Optional:    true,
				ForceNew:    true,
				Default:     "geth",
			},
			"genesis": {
				Type:        schema.TypeString,
				Description: "Genesis file content in JSON format",
				Required:    true,
				ForceNew:    true,
				ValidateFunc: func(i interface{}, s string) (ws []string, es []error) {
					jsonStr := i.(string)
					var g *core.Genesis
					if err := json.Unmarshal([]byte(jsonStr), &g); err != nil {
						es = append(es, err)
						return
					}
					return
				},
			},
			"data_dir_abs": {
				Type:        schema.TypeString,
				Description: "Absolute path to the data dir",
				Computed:    true,
			},
		},
	}
}

func resourceBootstrapDataDirCreate(d *schema.ResourceData, rawConfigurer interface{}) error {
	config := rawConfigurer.(*configurer)
	config.bootstrapDataDirMux.Lock()
	defer config.bootstrapDataDirMux.Unlock()
	targetDir := d.Get("data_dir").(string)
	absDir, err := createDirectory(targetDir)
	if err != nil {
		return err
	}
	nodeConfig := &node.DefaultConfig
	nodeConfig.DataDir = absDir
	nodeConfig.Name = d.Get("instance_name").(string)
	// check if the target dir is empty
	if files, err := ioutil.ReadDir(path.Join(absDir, nodeConfig.Name)); err != nil && !os.IsNotExist(err) {
		return err
	} else {
		if len(files) > 0 {
			return fmt.Errorf("directory [%s] is not empty", absDir)
		}
	}
	genesisJson := d.Get("genesis").(string)
	var genesis *core.Genesis
	if err := json.Unmarshal([]byte(genesisJson), &genesis); err != nil {
		return err
	}
	// init datadir
	stack, err := node.New(nodeConfig)
	if err != nil {
		return err
	}
	for _, name := range []string{"chaindata", "lightchaindata"} {
		chaindb, err := stack.OpenDatabase(name, 0, 0)
		if err != nil {
			return fmt.Errorf("can't open database for %s due to %s", name, err)
		}
		_, _, err = core.SetupGenesisBlock(chaindb, genesis)
		if err != nil {
			return fmt.Errorf("can't setup genesis for %s due to %s", name, err)
		}
		log.Printf("[DEBUG] Successfully wrote genesis state: database=%s, dir=%s", name, absDir)
	}
	_ = d.Set("data_dir_abs", absDir)
	d.SetId(fmt.Sprintf("%d", time.Now().UnixNano()))
	return nil
}

func resourceBootstrapDataDirRead(_ *schema.ResourceData, _ interface{}) error {
	return nil
}

func resourceBootstrapDataDirDelete(d *schema.ResourceData, _ interface{}) error {
	d.SetId("")
	dir := d.Get("data_dir_abs").(string)
	return os.RemoveAll(dir)
}
