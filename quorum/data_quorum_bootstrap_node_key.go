package quorum

import (
	"encoding/hex"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Use this data source to obtain various information parsed from an existing node key in hex format.
//
// Node key encodes a private key that defines an identity of a Quorum node in the network. It is primarily used in P2P networking.
func dataSourceBootstrapNodeKey() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBootstrapNodeKeyRead,
		Schema: map[string]*schema.Schema{
			"node_key_hex": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "This hex value encodes the private key of a node",
			},
			"hex_node_id": {
				Type:        schema.TypeString,
				Description: "64-byte hex value represents node ID which is seen being encoded in the username portion of enode URL. E.g.: `enode:/[<hex node id]@localhost:22000`",
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

func dataSourceBootstrapNodeKeyRead(d *schema.ResourceData, meta interface{}) error {
	nodeKeyHex := d.Get("node_key_hex").(string)
	nodeKey, err := crypto.HexToECDSA(nodeKeyHex)
	if err != nil {
		return err
	}
	d.SetId(enode.PubkeyToIDV4(&nodeKey.PublicKey).String())
	_ = d.Set("node_key_hex", hex.EncodeToString(crypto.FromECDSA(nodeKey)))
	return resourceBootstrapNodeKeyRead(d, meta)
}
