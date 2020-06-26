package quorum

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// @example
func TestAccResourceTransationManagerKeyPair_whenTypical(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: `
                    resource "quorum_transaction_manager_keypair" "test" {
						password = "foo"
                    }
                `,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("quorum_transaction_manager_keypair.test", "key_data"),
					resource.TestCheckResourceAttrSet("quorum_transaction_manager_keypair.test", "public_key_b64"),
					resource.TestCheckResourceAttrSet("quorum_transaction_manager_keypair.test", "private_key_json"),
				),
			},
		},
	})
}
