package quorum

import (
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/log"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testProviders = map[string]terraform.ResourceProvider{
	"quorum": Provider(),
}

func TestMain(m *testing.M) {
	if os.Getenv("TF_LOG") == "" {
		log.Root().SetHandler(log.DiscardHandler())
	}
	os.Exit(m.Run())
}
