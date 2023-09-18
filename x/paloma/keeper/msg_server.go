package keeper

import (
	"context"
	"os"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/palomachain/paloma/x/paloma/types"
)

const cPigeonStatusUpdateFF = "PALOMA_FF_PIGEON_STATUS_UPDATE"

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) AddStatusUpdate(goCtx context.Context, msg *types.MsgAddStatusUpdate) (*types.EmptyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	k.Logger(ctx).Info("pigeon-status-update: received call")

	// Avoid log spamming
	_, ok := os.LookupEnv(cPigeonStatusUpdateFF)
	if !ok {
		k.Logger(ctx).Info("pigeon-status-update: no env")
		return &types.EmptyResponse{}, nil
	}

	creator, _ := sdk.AccAddressFromBech32(msg.Creator)
	valAddr := sdk.ValAddress(creator.Bytes())
	status := msg.GetStatus()

	k.Logger(ctx).Info("pigeon-status-update: pre log")
	var logFn func(string, ...interface{})
	switch msg.Level {
	case types.MsgAddStatusUpdate_LEVEL_ERROR:
		logFn = k.Logger(ctx).Error
	case types.MsgAddStatusUpdate_LEVEL_INFO:
		logFn = k.Logger(ctx).Info
	case types.MsgAddStatusUpdate_LEVEL_DEBUG:
	default:
		logFn = k.Logger(ctx).Debug
	}

	k.Logger(ctx).Info("pigeon-status-update: logging")
	logFn(status,
		"component", "pigeon-status-update",
		"status", status,
		"sender", valAddr)

	return &types.EmptyResponse{}, nil
}
