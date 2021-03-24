package quorum

import "github.com/ethereum/go-ethereum/common"

type BesuExtraData struct {
	Vanity      []byte
	Validators  []common.Address
	Vote        []byte
	RoundNumber []byte
	Seals       [][]byte
}
