package quorum

import (
	"bytes"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Use this resource to construct `extradata` field used in the genesis file.
//
// `istanbul_address` can be referenced from `quorum_bootstrap_node_key` data source or newly created from `quorum_bootstrap_node_key` resources.
func resourceBootstrapIstanbulExtradata() *schema.Resource {
	return &schema.Resource{
		Create: resourceBootstrapIstanbulExtradataCreate,
		Read:   resourceBootstrapIstanbulExtradataRead,
		Delete: resourceBootstrapIstanbulExtradataDelete,
		Schema: map[string]*schema.Schema{
			"istanbul_addresses": {
				Type:        schema.TypeList,
				Description: "list of Istanbul address to construct extradata",
				Elem:        &schema.Schema{Type: schema.TypeString},
				MinItems:    1,
				Required:    true,
				ForceNew:    true,
			},
			"mode": {
				Type:        schema.TypeString,
				Description: "generate extradata using RLP encoding mode. Supported: ibft1 and ibft2. Default is ibft1",
				Optional:    true,
				ForceNew:    true,
				Default:     Ibft1,
			},
			"vanity": {
				Type:        schema.TypeString,
				Description: "Vanity Hex Value to be included in the extradata",
				Optional:    true,
				ForceNew:    true,
				Default:     "0x00",
			},
			"extradata": {
				Type:        schema.TypeString,
				Description: "Computed value which can be used in genesis file",
				Computed:    true,
			},
		},
	}
}

func resourceBootstrapIstanbulExtradataCreate(d *schema.ResourceData, _ interface{}) error {
	addresses := d.Get("istanbul_addresses").([]interface{})
	validators := make([]common.Address, len(addresses))
	for idx, rawAddress := range addresses {
		if addr, ok := rawAddress.(string); !ok {
			return fmt.Errorf("expect string element in istanbul_addresses")
		} else {
			validators[idx] = common.HexToAddress(addr)
		}
	}
	vanity := d.Get("vanity").(string)
	mode := d.Get("mode").(string)

	// by default, Ibft1
	createFunc := createIbft1ExtraData
	switch mode {
	case Ibft2:
		createFunc = createIbft2ExtraData
	case Qbft:
		createFunc = createQbftExtraData
	}

	payload, err := createFunc(validators, vanity)
	if err != nil {
		return err
	}
	extradata := "0x" + common.Bytes2Hex(payload)
	d.Set("extradata", extradata)
	d.SetId(fmt.Sprintf("%d", time.Now().Unix()))
	return nil
}

func createIbft1ExtraData(validators []common.Address, vanity string) ([]byte, error) {
	ist := &types.IstanbulExtra{
		Validators:    validators,
		Seal:          make([]byte, types.IstanbulExtraSeal),
		CommittedSeal: [][]byte{},
	}
	payload, err := rlp.EncodeToBytes(&ist)
	if err != nil {
		return nil, err
	}

	newVanity, err := hexutil.Decode(vanity)
	if err != nil {
		return nil, err
	}
	if len(newVanity) < types.IstanbulExtraVanity {
		newVanity = append(newVanity, bytes.Repeat([]byte{0x00}, types.IstanbulExtraVanity-len(newVanity))...)
	}
	newVanity = newVanity[:types.IstanbulExtraVanity]
	return append(newVanity, payload...), nil
}

func createIbft2ExtraData(validators []common.Address, _ string) ([]byte, error) {
	data := &BesuExtraData{
		Vanity:      make([]byte, 32),
		Validators:  validators,
		RoundNumber: make([]byte, 4),
	}

	return rlp.EncodeToBytes(data)
}

// createQbftExtraData generates qbft consensus compatible extraData
func createQbftExtraData(validators []common.Address, _ string) ([]byte, error) {
	data := &QbftExtraData{
		Vanity:     make([]byte, 32),
		Validators: validators,
	}
	return rlp.EncodeToBytes(data)
}

func resourceBootstrapIstanbulExtradataRead(d *schema.ResourceData, _ interface{}) error {
	return nil
}

func resourceBootstrapIstanbulExtradataDelete(d *schema.ResourceData, _ interface{}) error {
	return nil
}
