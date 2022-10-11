package quorum

import "github.com/ethereum/go-ethereum/common"

type Mode string

const (
	Ibft1 Mode = "ibft1"
	Ibft2      = "ibft2"
	Qbft       = "qbft"
)

type BesuExtraData struct {
	Vanity      []byte
	Validators  []common.Address
	Vote        []byte
	RoundNumber []byte
	Seals       [][]byte
}

// QbftExtraData represents Qbft consensus compatible extraData
type QbftExtraData struct {
	Vanity      []byte
	Validators  []common.Address
	Vote        *ValidatorVote
	RoundNumber []byte
	Seals       [][]byte
}

// ValidatorVote represent a vote to add or remove a validator in Qbft Consensus
type ValidatorVote struct {
	RecipientAddress common.Address
	VoteType         byte
}
