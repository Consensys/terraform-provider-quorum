package quorum

import "github.com/ethereum/go-ethereum/common"

type Mode string

const (
	Ibft1 Mode = "ibft1"
	Ibft2      = "ibft2"
)

type BesuExtraData struct {
	Vanity      []byte
	Validators  []common.Address
	Vote        []byte
	RoundNumber []byte
	Seals       [][]byte
}
