package quorum

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/p2p/enode"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Use this resource to create a node key for a new Quorum node.
//
// Node key encodes a private key that defines an identity of a Quorum node in the network. It is primarily used in P2P networking.
func resourceBootstrapNodeKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceBootstrapNodeKeyCreate,
		Read:   resourceBootstrapNodeKeyRead,
		Delete: resourceBootstrapNodeKeyDelete,

		Schema: map[string]*schema.Schema{
			"node_key_hex": {
				Type:        schema.TypeString,
				Description: "Node key as hex. This can be used to populate `--nodekeyhex` CLI parameter",
				Computed:    true,
				Sensitive:   true,
			},
			"hex_node_id": {
				Type:        schema.TypeString,
				Description: "64-byte hex value represents node ID which is seen being encoded in the username portion of enode URL. E.g.: `enode://[hex node id]@localhost:22000`",
				Computed:    true,
			},
			"node_id": {
				Type:        schema.TypeString,
				Description: "32-byte hex value represents the unique identifier for a node",
				Computed:    true,
			},
			"istanbul_address": {
				Type:        schema.TypeString,
				Description: "Address representing public key for the newly created node key. This is mainly used to construct the initial validators set encoded in `extradata` field of the genesis file",
				Computed:    true,
			},
		},
	}
}

func resourceBootstrapNodeKeyCreate(d *schema.ResourceData, meta interface{}) error {
	nodeKey, err := crypto.GenerateKey()
	if err != nil {
		return err
	}
	d.SetId(enode.PubkeyToIDV4(&nodeKey.PublicKey).String())
	_ = d.Set("node_key_hex", hex.EncodeToString(crypto.FromECDSA(nodeKey)))
	return resourceBootstrapNodeKeyRead(d, meta)
}

func resourceBootstrapNodeKeyRead(d *schema.ResourceData, _ interface{}) error {
	// recover the private key
	nodeKeyHex := d.Get("node_key_hex").(string)
	nodeKey, err := crypto.HexToECDSA(nodeKeyHex)
	if err != nil {
		return err
	}
	_ = d.Set("hex_node_id", fmt.Sprintf("%x", crypto.FromECDSAPub(&nodeKey.PublicKey)[1:]))
	_ = d.Set("node_id", d.Id())
	istanbulAddress := crypto.PubkeyToAddress(nodeKey.PublicKey)
	_ = d.Set("istanbul_address", strings.ToLower(istanbulAddress.String()))
	return nil
}

func resourceBootstrapNodeKeyDelete(d *schema.ResourceData, _ interface{}) error {
	d.SetId("")
	return nil
}
