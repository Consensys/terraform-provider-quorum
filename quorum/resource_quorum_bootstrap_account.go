package quorum

import (
	"fmt"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Use this resource to create a new Ethereum account
func resourceBootstrapAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceBootstrapAccountCreate,
		Read:   resourceBootstrapAccountRead,
		Update: resourceBootstrapAccountUpdate,
		Delete: resourceBootstrapAccountDelete,

		Schema: map[string]*schema.Schema{
			"wallet_id": {
				Type:        schema.TypeString,
				Description: "ID of a wallet storing the newly created account. For keystore, it's the keystore resource id",
				Required:    true,
				ForceNew:    true,
			},
			"passphrase": {
				Type:        schema.TypeString,
				Description: "Passphrase to lock/unlock the account. Default is empty",
				Optional:    true,
				ForceNew:    true,
				Default:     "",
				Sensitive:   true,
			},
			"address": {
				Type:        schema.TypeString,
				Description: "Address of the newly generated account",
				Computed:    true,
			},
			"account_url": {
				Type:        schema.TypeString,
				Description: "URL of the newly generated account",
				Computed:    true,
			},
			"balance": {
				Type:        schema.TypeString,
				Description: "A place holder to keep account initial balance for referencing",
				Optional:    true,
			},
		},
	}
}

func resourceBootstrapAccountCreate(d *schema.ResourceData, raw interface{}) error {
	registry := raw.(*configurer).registry
	walletId := d.Get("wallet_id").(string)
	wallet, ok := registry.get(walletId)
	if !ok {
		return fmt.Errorf("wallet %s does not exist nor not registered", walletId)
	}
	passphrase := d.Get("passphrase").(string)
	var newAccount accounts.Account
	var err error
	switch wallet.(type) {
	case *keystore.KeyStore:
		ks := wallet.(*keystore.KeyStore)
		newAccount, err = ks.NewAccount(passphrase)
	default:
		err = fmt.Errorf("unsupported wallet type: %s", wallet)
	}
	if err != nil {
		return err
	}
	accountAddress := strings.ToLower(newAccount.Address.String())
	d.SetId(accountAddress)
	_ = d.Set("address", accountAddress)
	_ = d.Set("account_url", newAccount.URL.Path)
	return nil
}

func resourceBootstrapAccountUpdate(d *schema.ResourceData, _ interface{}) error {
	// not required to implement due to the `balance` field is merely a place holder
	// other fields are having `forceNew` is true
	return nil
}

func resourceBootstrapAccountRead(_ *schema.ResourceData, _ interface{}) error {
	return nil
}

func resourceBootstrapAccountDelete(d *schema.ResourceData, _ interface{}) error {
	d.SetId("")
	accountUrl := d.Get("account_url").(string)
	_ = os.Remove(accountUrl) // it may or may not be a file system path so we ignore the error
	return nil
}
