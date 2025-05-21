package skyway

import (
	"fmt"

	sdkerrors "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsmod "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/palomachain/paloma/v2/util/liblog"
	"github.com/palomachain/paloma/v2/util/libmsg"
	"github.com/palomachain/paloma/v2/x/skyway/keeper"
	"github.com/palomachain/paloma/v2/x/skyway/types"
)

// NewHandler returns a handler for "Skyway" type messages.
func NewHandler(k keeper.Keeper) baseapp.MsgServiceHandler {
	msgServer := keeper.NewMsgServerImpl(k)

	return func(ctx sdk.Context, msg sdk.Msg) (_ *sdk.Result, err error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		lg := liblog.FromKeeper(ctx, k).WithComponent("skyway-msg-handler")
		defer func() {
			if err != nil {
				sender := "unknown"
				if s, ok := libmsg.GetSender(msg); ok {
					sender = s
				}
				lg.WithError(err).
					WithFields("msg", sdk.MsgTypeURL(msg)).
					WithFields("msg-creator", sender).
					Error("failed to handle message")
			}
		}()
		switch msg := msg.(type) {
		case *types.MsgSendToRemote:
			msg.GetMetadata()
			res, err := msgServer.SendToRemote(sdk.UnwrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgConfirmBatch:
			res, err := msgServer.ConfirmBatch(sdk.UnwrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgSendToPalomaClaim:
			res, err := msgServer.SendToPalomaClaim(sdk.UnwrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgBatchSendToRemoteClaim:
			res, err := msgServer.BatchSendToRemoteClaim(sdk.UnwrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgCancelSendToRemote:
			res, err := msgServer.CancelSendToRemote(sdk.UnwrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgSubmitBadSignatureEvidence:
			res, err := msgServer.SubmitBadSignatureEvidence(sdk.UnwrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		default:
			return nil, sdkerrors.Wrap(errorsmod.ErrUnknownRequest, fmt.Sprintf("Unrecognized Skyway Msg type: %v", sdk.MsgTypeURL(msg)))
		}
	}
}
