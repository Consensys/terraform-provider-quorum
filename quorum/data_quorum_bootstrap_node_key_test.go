package quorum

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// @example
func TestAccDataSourceBootstrapNodeKey_whenTypical(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "quorum_bootstrap_node_key" "test" {
                    }

                    data "quorum_bootstrap_node_key" "test" {
						node_key_hex = quorum_bootstrap_node_key.test.node_key_hex
                    }
                `,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.quorum_bootstrap_node_key.test", "hex_node_id"),
					resource.TestCheckResourceAttrSet("data.quorum_bootstrap_node_key.test", "node_key_hex"),
					resource.TestCheckResourceAttrSet("data.quorum_bootstrap_node_key.test", "node_id"),
					resource.TestCheckResourceAttrSet("data.quorum_bootstrap_node_key.test", "istanbul_address"),
				),
			},
		},
	})
}
