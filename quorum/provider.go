// Provides data sources and resources to work with Quorum network
package quorum

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	elog "github.com/ethereum/go-ethereum/log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Quorum Provider is used to work with Quorum Network such as creating basic metadata required to bootstrap a new Quorum Network.
//
// This provider is best used along with other providers to successfully create a running Quorum network.
//
// Use the navigation to the left to read about the available resources.
func Provider() *schema.Provider {
	elog.Root().SetHandler(elog.StreamHandler(os.Stderr, elog.TerminalFormat(false)))
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"quorum_bootstrap_account":            resourceBootstrapAccount(),
			"quorum_bootstrap_data_dir":           resourceBootstrapDataDir(),
			"quorum_bootstrap_istanbul_extradata": resourceBootstrapIstanbulExtradata(),
			"quorum_bootstrap_keystore":           resourceBootstrapKeyStore(),
			"quorum_bootstrap_network":            resourceBootstrapNetwork(),
			"quorum_bootstrap_node_key":           resourceBootstrapNodeKey(),
			"quorum_transaction_manager_keypair":  resourceTransactionManagerKeyPair(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"quorum_bootstrap_genesis_mixhash": dataSourceBootstrapGenesisMixHash(),
			"quorum_bootstrap_node_key":        dataSourceBootstrapNodeKey(),
		},
		ConfigureFunc: func(_ *schema.ResourceData) (interface{}, error) {
			return &configurer{
				registry: newInternalRegistry(),
			}, nil
		},
	}
}

type configurer struct {
	registry            *internalRegistry
	bootstrapDataDirMux sync.Mutex // make sure we do it one by one otherwise we hit resource temporarily unavailable error
}

// this is mainly used to reference between resources/data sources
type internalRegistry struct {
	mux      sync.RWMutex
	registry map[string]interface{}
}

func newInternalRegistry() *internalRegistry {
	return &internalRegistry{
		mux:      sync.RWMutex{},
		registry: make(map[string]interface{}),
	}
}

func (ir *internalRegistry) set(key string, value interface{}) {
	ir.mux.Lock()
	defer ir.mux.Unlock()

	ir.registry[key] = value
}

func (ir *internalRegistry) get(key string) (interface{}, bool) {
	ir.mux.RLock()
	defer ir.mux.RUnlock()

	v, ok := ir.registry[key]
	return v, ok
}

func (ir *internalRegistry) delete(id string) {
	ir.mux.Lock()
	defer ir.mux.Unlock()

	delete(ir.registry, id)
}

// create new directory including parents
func createDirectory(dir string) (string, error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Println("[DEBUG] Creating new directory at ", dir)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return "", fmt.Errorf("can't create new directory due to %s", err)
		}
	}
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return "", fmt.Errorf("can't obtain absolute path due to %s", err)
	}
	return absDir, nil
}
