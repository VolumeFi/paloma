package libmsg

// TODO: Remove this package and add functions to
// respective impelemtatnions as part of
// https://github.com/VolumeFi/paloma/issues/1041
import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	consensustypes "github.com/palomachain/paloma/v2/x/consensus/types"
	evmtypes "github.com/palomachain/paloma/v2/x/evm/types"
	valsettypes "github.com/palomachain/paloma/v2/x/valset/types"
)

type Envelope interface {
	GetMsg() *types.Any
}

type Metadataer interface {
	GetMetadata() valsettypes.MsgMetadata
}

type ConsensusMsgProvider interface {
	ConsensusMsg(types.AnyUnpacker) (consensustypes.ConsensusMsg, error)
}

func ToEvmMessage(c ConsensusMsgProvider, cdc types.AnyUnpacker) (*evmtypes.Message, error) {
	e, err := c.ConsensusMsg(cdc)
	if err != nil {
		return nil, err
	}

	m, ok := e.(*evmtypes.Message)
	if !ok {
		return nil, fmt.Errorf("e is not of type Message")
	}

	return m, nil
}

func GetAssignee(e Envelope, cdc types.AnyUnpacker) (string, error) {
	var unpackedMsg evmtypes.TurnstoneMsg
	if err := cdc.UnpackAny(e.GetMsg(), &unpackedMsg); err != nil {
		return "", fmt.Errorf("failed to unpack message: %w", err)
	}

	return unpackedMsg.GetAssignee(), nil
}

func GetSender(msg sdk.Msg) (string, bool) {
	m, ok := msg.(Metadataer)
	if !ok {
		return "", false
	}
	return m.GetMetadata().Creator, true
}
