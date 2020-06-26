package quorum

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// @example
func TestAccResourceBootstrapAccount_whenTypical(t *testing.T) {
	tempdir, err := ioutil.TempDir("", "testacc-")
	if err != nil {
		t.Fatalf("can't create temp dir: %s", err)
	}
	defer os.RemoveAll(tempdir)
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "quorum_bootstrap_keystore" "test" {
						keystore_dir = "%s"
						use_light_weight_kdf = false
					}
                    resource "quorum_bootstrap_account" "test" {
                        wallet_id = quorum_bootstrap_keystore.test.id
						passphrase = "test"
                    }
                `, tempdir),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("quorum_bootstrap_account.test", "address"),
					resource.TestCheckResourceAttrSet("quorum_bootstrap_account.test", "account_url"),
				),
			},
		},
	})
}
