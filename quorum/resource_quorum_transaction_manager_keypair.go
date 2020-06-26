package quorum

import (
	"crypto/rand"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"golang.org/x/crypto/nacl/box"
)

// Use this resource to create a key pair used in a transaction manager.
//
// This key pair provides attributes which are useful when building the configuration for a transaction manager.
func resourceTransactionManagerKeyPair() *schema.Resource {
	return &schema.Resource{
		Create: resourceTransactionManagerKeyPairCreate,
		Read:   resourceTransactionManagerKeyPairRead,
		Delete: resourceTransactionManagerKeyPairDelete,

		Schema: map[string]*schema.Schema{
			"password": {
				Type:        schema.TypeString,
				Description: "A password to protect the keypair",
				Optional:    true,
				ForceNew:    true,
				Default:     "",
				Sensitive:   true,
			},
			"config": {
				Type:        schema.TypeList,
				Description: "Key generation config",
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"variant": {
							Type:        schema.TypeString,
							Description: "Algorithm to use when hashing. Allowed values are `id` or `i`",
							Default:     defaultArgonOpts.Algorithm,
							Optional:    true,
							ValidateFunc: func(v interface{}, s string) (strings []string, errors []error) {
								value := v.(string)
								if value == "id" || value == "i" {
									return
								}
								errors = append(errors, fmt.Errorf("%s is not a valid value: [%s]. Allowed values are : id or i", s, value))
								return
							},
						},
						"iterations": {
							Type:        schema.TypeInt,
							Description: "Number of iterations to cycle through",
							Default:     defaultArgonOpts.Iterations,
							Optional:    true,
						},
						"memory": {
							Type:        schema.TypeInt,
							Description: "Memory limit",
							Default:     defaultArgonOpts.Memory,
							Optional:    true,
						},
						"parallelism": {
							Type:        schema.TypeInt,
							Description: "Number of threads to use",
							Default:     defaultArgonOpts.Parallelism,
							Optional:    true,
						},
					},
				},
			},
			"key_data": {
				Type:        schema.TypeString,
				Description: "Key Data in JSON format to be used by Private Transaction Manager",
				Computed:    true,
				Sensitive:   true,
			},
			"public_key_b64": {
				Type:        schema.TypeString,
				Description: "Public key in standard base64 encoding",
				Computed:    true,
			},
			"private_key_json": {
				Type:        schema.TypeString,
				Description: "Private key in JSON representation",
				Computed:    true,
				Sensitive:   true,
			},
		},
	}
}

func resourceTransactionManagerKeyPairCreate(d *schema.ResourceData, meta interface{}) error {
	pub, priv, err := box.GenerateKey(rand.Reader)
	if err != nil {
		return fmt.Errorf("unable to generate keypair due to %s", err)
	}
	pubB64 := toStandardBase64EncodedString(pub[:])
	d.SetId(string(pubB64))
	_ = d.Set("public_key_b64", pubB64)
	keyDataJSON, privateKeyJSON, err := toKeyDataJSON(d.Get("password").(string), toArgonOptions(d), priv[:], pubB64)
	if err != nil {
		return err
	}
	_ = d.Set("key_data", keyDataJSON)
	_ = d.Set("private_key_json", privateKeyJSON)
	return resourceTransactionManagerKeyPairRead(d, meta)
}

func toArgonOptions(d *schema.ResourceData) *argonOptions {
	if cfg, ok := d.GetOk("config"); ok {
		rawOpts := cfg.([]interface{})[0].(map[string]interface{})
		return &argonOptions{
			Algorithm:   rawOpts["variant"].(string),
			Iterations:  rawOpts["iterations"].(int),
			Memory:      rawOpts["memory"].(int),
			Parallelism: rawOpts["parallelism"].(int),
		}
	}
	return nil
}

func resourceTransactionManagerKeyPairRead(d *schema.ResourceData, _ interface{}) error {
	return nil
}

func resourceTransactionManagerKeyPairDelete(d *schema.ResourceData, _ interface{}) error {
	d.SetId("")
	return nil
}
