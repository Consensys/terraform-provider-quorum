package quorum

import (
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Use this data source to obtain `MixHash` value being used in the genesis file.
// Especially when using Istanbul consensus algorithm which defines a constant digest value.
func dataSourceBootstrapGenesisMixHash() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBootstrapGenesisMixHashRead,
		Schema: map[string]*schema.Schema{
			"istanbul": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Mix hash that is used to identify whether the block is from Istanbul consensus engine",
			},
		},
	}
}

func dataSourceBootstrapGenesisMixHashRead(d *schema.ResourceData, _ interface{}) error {
	d.SetId(fmt.Sprintf("%d", time.Now().Unix()))
	_ = d.Set("istanbul", strings.ToLower(types.IstanbulDigest.String()))
	return nil
}
