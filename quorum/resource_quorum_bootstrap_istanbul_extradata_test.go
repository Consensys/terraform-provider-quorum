package quorum

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// @example
func TestAccResourceBootstrapIstanbulExtradata_whenTypical(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: `
                    resource "quorum_bootstrap_node_key" "test" {
						count = 3
                    }

					resource "quorum_bootstrap_istanbul_extradata" "test" {
						istanbul_addresses = quorum_bootstrap_node_key.test.*.istanbul_address
					}
                `,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("quorum_bootstrap_istanbul_extradata.test", "extradata"),
				),
			},
		},
	})
}
