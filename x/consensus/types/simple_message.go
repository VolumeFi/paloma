package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	_ ConsensusMsg = &SimpleMessage{}
	_ ConsensusMsg = &EvenSimplerMessage{}
)

func (msg *SimpleMessage) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg *SimpleMessage) Attest() {}

func (msg *EvenSimplerMessage) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg *SimpleMessage) ConsensusSignBytes() BytesToSignFunc {
	return TypedBytesToSign(func(m *SimpleMessage, salt Salt) []byte {
		return append(m.GetSignBytes(), Uint64ToByte(salt.Nonce)...)
	})
}

func (msg *SimpleMessage) Keccak256WithSignedMessage(_ *QueuedSignedMessage) ([]byte, error) {
	return crypto.Keccak256([]byte(msg.Sender), []byte(msg.Hello), []byte(msg.World)), nil
}
