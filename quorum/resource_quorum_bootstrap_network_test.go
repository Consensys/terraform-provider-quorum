package quorum

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// @example
func TestAccResourceBootstrapNetwork_whenTypical(t *testing.T) {
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
					resource "quorum_bootstrap_network" "test" {
						name = "test-network"
						target_dir = "%s"
					}
                `, tempdir),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("quorum_bootstrap_network.test", "network_dir_abs"),
				),
			},
		},
	})
}
