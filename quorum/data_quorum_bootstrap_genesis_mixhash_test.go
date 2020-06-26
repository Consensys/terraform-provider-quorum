package quorum

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// @example
func TestAccDataSourceBootstrapGenesisMixHash_whenTypical(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: `
data "quorum_bootstrap_genesis_mixhash" "test" {
}

output "istanbul_mix_hash" {
  value = data.quorum_bootstrap_genesis_mixhash.test.istanbul
}
                `,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.quorum_bootstrap_genesis_mixhash.test", "istanbul"),
					resource.TestCheckOutput("istanbul_mix_hash", "0x63746963616c2062797a616e74696e65206661756c7420746f6c6572616e6365"),
				),
			},
		},
	})
}
