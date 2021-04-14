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

//TestUnitOfWork_StateUnderTest_ExpectedBehavior
func TestAccResourceBootstrapIstanbulExtradata_GoQuorum(t *testing.T) {
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
						istanbul_addresses = [
							"0x8f1e6d8303716516cc9e562e66d09721752a1f83",
							"0x95167bde9c4c3b12180945bbee9900f69d9ea558",
							"0xa7c1d1b572f11b02cd6fadc21f1e51f399b4d4cb"
						]
					}
                `,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("quorum_bootstrap_istanbul_extradata.test", "extradata", "0x0000000000000000000000000000000000000000000000000000000000000000f885f83f948f1e6d8303716516cc9e562e66d09721752a1f839495167bde9c4c3b12180945bbee9900f69d9ea55894a7c1d1b572f11b02cd6fadc21f1e51f399b4d4cbb8410000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c0"),
				),
			},
		},
	})
}

func TestAccResourceBootstrapIstanbulExtradata_Besu(t *testing.T) {
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
						istanbul_addresses = [
							"0x8f1e6d8303716516cc9e562e66d09721752a1f83",
							"0x95167bde9c4c3b12180945bbee9900f69d9ea558",
							"0xa7c1d1b572f11b02cd6fadc21f1e51f399b4d4cb"
						]
						mode = "ibft2"
					}
               `,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("quorum_bootstrap_istanbul_extradata.test", "extradata", "0xf869a00000000000000000000000000000000000000000000000000000000000000000f83f948f1e6d8303716516cc9e562e66d09721752a1f839495167bde9c4c3b12180945bbee9900f69d9ea55894a7c1d1b572f11b02cd6fadc21f1e51f399b4d4cb808400000000c0"),
				),
			},
		},
	})
}

// TestAccResourceBootstrapQbftExtradata verifies if qbft mode generates extraData compatible with QBFT Consensus
func TestAccResourceBootstrapQbftExtradata(t *testing.T) {
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
						istanbul_addresses = [
							"0x8f1e6d8303716516cc9e562e66d09721752a1f83",
							"0x95167bde9c4c3b12180945bbee9900f69d9ea558",
							"0xa7c1d1b572f11b02cd6fadc21f1e51f399b4d4cb"
						]
						mode = "qbft"
					}
               `,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("quorum_bootstrap_istanbul_extradata.test", "extradata", "0xf865a00000000000000000000000000000000000000000000000000000000000000000f83f948f1e6d8303716516cc9e562e66d09721752a1f839495167bde9c4c3b12180945bbee9900f69d9ea55894a7c1d1b572f11b02cd6fadc21f1e51f399b4d4cbc080c0"),
				),
			},
		},
	})
}
