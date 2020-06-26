package quorum

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// @example
func TestAccResourceBootstrapDataDir_whenTypical(t *testing.T) {
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
					resource "quorum_bootstrap_data_dir" "test" {
						data_dir = "%s"
						genesis = <<EOF
{
  "alloc": { },
  "coinbase": "0x0000000000000000000000000000000000000000",
  "config": {
    "byzantiumBlock": 1,
    "chainId": 10,
    "eip150Block": 1,
    "eip155Block": 0,
    "eip150Hash": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "eip158Block": 1,
    "isQuorum": true
  },
  "difficulty": "0x00",
  "extraData": "0x00",
  "gasLimit": "0xE0000000",
  "mixhash": "0x0000000000000000000000000000000000000000000000000000000000000000",
  "nonce": "0x00",
  "parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
  "timestamp": "0x00",
  "number": "0x00",
  "gasUsed": "0x00"
}
EOF
					}
                `, tempdir),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("quorum_bootstrap_data_dir.test", "data_dir_abs"),
				),
			},
		},
	})
}
