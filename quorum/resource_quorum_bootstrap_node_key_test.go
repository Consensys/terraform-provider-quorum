package quorum

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// @example
func TestAccResourceBootstrapNodeKey_whenTypical(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: `
                    resource "quorum_bootstrap_node_key" "test" {
                    }
                `,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("quorum_bootstrap_node_key.test", "hex_node_id"),
					resource.TestCheckResourceAttrSet("quorum_bootstrap_node_key.test", "node_key_hex"),
					resource.TestCheckResourceAttrSet("quorum_bootstrap_node_key.test", "node_id"),
					resource.TestCheckResourceAttrSet("quorum_bootstrap_node_key.test", "istanbul_address"),
				),
			},
		},
	})
}
