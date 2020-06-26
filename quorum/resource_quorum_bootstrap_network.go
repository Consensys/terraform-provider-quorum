package quorum

import (
	"os"
	"path"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Use this resource to create a new directory that represents a new Quorum network.
//
// Bootstraping data will be kept in this folder.
func resourceBootstrapNetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourceBootstrapNetworkCreate,
		Read:   resourceBootstrapNetworkRead,
		Delete: resourceBootstrapNetworkDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Name of a new network. Directory name restriction applied",
				Required:    true,
				ForceNew:    true,
			},
			"target_dir": {
				Type:        schema.TypeString,
				Description: "File system path to the directory on which new directory will be created. Default is current working directory",
				Optional:    true,
				ForceNew:    true,
				Default:     ".",
			},
			"network_dir_abs": {
				Type:        schema.TypeString,
				Description: "Absolute path to a directory representing this new network",
				Computed:    true,
			},
		},
	}
}

func resourceBootstrapNetworkCreate(d *schema.ResourceData, _ interface{}) error {
	name := d.Get("name").(string)
	targetDir := d.Get("target_dir").(string)
	d.SetId(name)
	absDir, err := createDirectory(path.Join(targetDir, name))
	if err != nil {
		return err
	}
	_ = d.Set("network_dir_abs", absDir)
	return nil
}

func resourceBootstrapNetworkRead(_ *schema.ResourceData, _ interface{}) error {
	return nil
}

func resourceBootstrapNetworkDelete(d *schema.ResourceData, _ interface{}) error {
	d.SetId("")
	dir := d.Get("network_dir_abs").(string)
	_ = os.RemoveAll(dir)
	return nil
}
