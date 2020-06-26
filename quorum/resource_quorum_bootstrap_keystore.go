package quorum

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Use this resource to create a keystore which maintains multiple Ethereum accounts.
func resourceBootstrapKeyStore() *schema.Resource {
	return &schema.Resource{
		Create: resourceBootstrapKeyStoreCreate,
		Read:   resourceBootstrapKeyStoreRead,
		Delete: resourceBootstrapKeyStoreDelete,
		Update: resourceBootstrapKeyStoreUpdate,

		Schema: map[string]*schema.Schema{
			"keystore_dir": {
				Type:        schema.TypeString,
				Description: "Directory contains private keys",
				Required:    true,
				ForceNew:    true,
			},
			"use_light_weight_kdf": {
				Type:        schema.TypeBool,
				Description: "True to lower the memory and CPU requirements of the key store scrypt KDF at the expense of security",
				Optional:    true,
				ForceNew:    true,
				Default:     false,
			},
			"keystore_dir_abs": {
				Type:        schema.TypeString,
				Description: "Absolute path of the keystore directory",
				Computed:    true,
			},
			"account": {
				Type:        schema.TypeList,
				Description: "Account being created under this keystore",
				Optional:    true,
				ConfigMode:  schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"passphrase": {
							Type:        schema.TypeString,
							Description: "Passphrase to lock/unlock the account. Default is empty",
							Default:     "",
							Optional:    true,
							Sensitive:   true,
						},
						"address": {
							Type:        schema.TypeString,
							Description: "Address of the newly generated account",
							Computed:    true,
						},
						"account_url": {
							Type:        schema.TypeString,
							Description: "Local path to the JSON representation of newly generated account private key",
							Computed:    true,
							Sensitive:   true,
						},
						"balance": {
							Type:        schema.TypeString,
							Description: "A place holder to keep account initial balance for referencing",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func resourceBootstrapKeyStoreCreate(d *schema.ResourceData, rawConfigurer interface{}) error {
	d.SetId(fmt.Sprintf("ks-%d", time.Now().UnixNano()))
	keystoreDir := d.Get("keystore_dir").(string)
	log.Println("[DEBUG] Keystore Directory", keystoreDir)
	absDir, err := createDirectory(keystoreDir)
	if err != nil {
		return err
	}
	if files, err := ioutil.ReadDir(absDir); err != nil {
		return err
	} else {
		if len(files) > 0 {
			return fmt.Errorf("directory [%s] is not empty", absDir)
		}
	}
	_ = d.Set("keystore_dir_abs", absDir)
	if err := resourceBootstrapKeyStoreRead(d, rawConfigurer); err != nil {
		return err
	}
	ks, err := getKeystoreInstance(d.Id(), rawConfigurer)
	if err != nil {
		return err
	}
	for _, raw := range d.Get("account").([]interface{}) {
		if err := createNewAccount(ks, raw); err != nil {
			return err
		}
	}
	return nil
}

func createNewAccount(ks *keystore.KeyStore, raw interface{}) error {
	newAccountSchema := raw.(map[string]interface{})
	newAcc, err := ks.NewAccount(newAccountSchema["passphrase"].(string))
	if err != nil {
		return err
	}
	accountAddress := strings.ToLower(newAcc.Address.String())
	newAccountSchema["address"] = accountAddress
	newAccountSchema["account_url"] = newAcc.URL.Path
	log.Println("[DEBUG] New account is created. Address", accountAddress)
	return nil
}

func getKeystoreInstance(id string, rawConfigurer interface{}) (*keystore.KeyStore, error) {
	config := rawConfigurer.(*configurer)
	ksRaw, ok := config.registry.get(id)
	if !ok {
		return nil, fmt.Errorf("can't find keystore [%s] instance in the registry", id)
	}
	return ksRaw.(*keystore.KeyStore), nil
}

func resourceBootstrapKeyStoreUpdate(d *schema.ResourceData, rawConfigurer interface{}) error {
	if err := resourceBootstrapKeyStoreRead(d, rawConfigurer); err != nil {
		return err
	}
	ks, err := getKeystoreInstance(d.Id(), rawConfigurer)
	if err != nil {
		return err
	}
	if d.HasChange("account") {
		o, n := d.GetChange("account")
		if o == nil {
			o = make([]interface{}, 0)
		}
		if n == nil {
			n = make([]interface{}, 0)
		}
		oldAccountSet := o.([]interface{})
		newAccountSet := n.([]interface{})
		existingAccountSet := make(map[string]bool)
		// create new accounts
		for _, raw := range newAccountSet {
			log.Println("[DEBUG] Diff New:", raw)
			acc := raw.(map[string]interface{})
			accountAddress := acc["address"].(string)
			if accountAddress == "" {
				if err := createNewAccount(ks, raw); err != nil {
					return err
				}
			} else {
				existingAccountSet[accountAddress] = true
			}
		}
		// delete old accounts
		for _, raw := range oldAccountSet {
			log.Println("[DEBUG] Diff Old:", raw)
			acc := raw.(map[string]interface{})
			accountAddress := acc["address"].(string)
			if _, ok := existingAccountSet[accountAddress]; !ok {
				log.Println("[DEBUG] Deleting account", accountAddress)
				if err := os.Remove(acc["account_url"].(string)); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func resourceBootstrapKeyStoreRead(d *schema.ResourceData, rawConfigurer interface{}) error {
	absDir := d.Get("keystore_dir_abs").(string)
	sn, sp := keystore.StandardScryptN, keystore.StandardScryptP
	useLightWeightKDF := d.Get("use_light_weight_kdf").(bool)
	if useLightWeightKDF {
		sn, sp = keystore.LightScryptN, keystore.LightScryptP
	}
	ks := keystore.NewKeyStore(absDir, sn, sp)
	// save into registry so if it can be retrieved later if needed
	config := rawConfigurer.(*configurer)
	config.registry.set(d.Id(), ks)
	return nil
}

func resourceBootstrapKeyStoreDelete(d *schema.ResourceData, raw interface{}) error {
	keyDir := d.Get("keystore_dir").(string)
	log.Println("[DEBUG] Deleting keystore", keyDir)
	raw.(*configurer).registry.delete(d.Id())
	d.SetId("")
	return os.RemoveAll(keyDir)
}
