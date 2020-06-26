package quorum

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// @example
func TestAccResourceBootstrapKeyStore_whenTypical(t *testing.T) {
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
									account {
									}
									account {
										passphrase = "acc1"
									}
			                    }
			                `, tempdir),
				Check: resource.ComposeAggregateTestCheckFunc(
					//func(ts *terraform.State) error {
					//	t.Log("inspect:", ts.RootModule().Resources["quorum_bootstrap_keystore.test"].Primary.Attributes)
					//	return nil
					//},
					resource.TestCheckResourceAttrSet("quorum_bootstrap_keystore.test", "keystore_dir_abs"),
					resource.TestCheckResourceAttr("quorum_bootstrap_keystore.test", "account.#", "2"),
					resource.TestCheckResourceAttrSet("quorum_bootstrap_keystore.test", "account.0.address"),
					resource.TestCheckResourceAttrSet("quorum_bootstrap_keystore.test", "account.0.account_url"),
					resource.TestCheckResourceAttrSet("quorum_bootstrap_keystore.test", "account.1.address"),
					resource.TestCheckResourceAttrSet("quorum_bootstrap_keystore.test", "account.1.account_url"),
				),
			},
		},
	})
	// check if the tempdir is actually deleted after tf destroy
	_, err = os.Stat(tempdir)
	assert.True(t, os.IsNotExist(err))
}

func TestAccResourceBootstrapKeyStore_whenUpdate(t *testing.T) {
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
									account {
										passphrase = "1"
									}
									account {
										passphrase = "2"
									}	
			                    }
			                `, tempdir),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("quorum_bootstrap_keystore.test", "account.#", "2"),
				),
			},
			{
				Config: fmt.Sprintf(`
			                    resource "quorum_bootstrap_keystore" "test" {
			                        keystore_dir = "%s"
									use_light_weight_kdf = false
									account {
										passphrase = "3"
									}
			                    }
			                `, tempdir),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("quorum_bootstrap_keystore.test", "account.#", "1"),
				),
			},
		},
	})
	// check if the tempdir is actually deleted after tf destroy
	_, err = os.Stat(tempdir)
	assert.True(t, os.IsNotExist(err))
}
