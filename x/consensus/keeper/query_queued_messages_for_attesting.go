package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/palomachain/paloma/v2/x/consensus/keeper/consensus"
	"github.com/palomachain/paloma/v2/x/consensus/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) QueuedMessagesForAttesting(goCtx context.Context, req *types.QueryQueuedMessagesForAttestingRequest) (*types.QueryQueuedMessagesForAttestingResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	msgs, err := k.GetMessagesForAttesting(ctx, req.QueueTypeName, req.ValAddress)
	if err != nil {
		return nil, err
	}

	res := make([]types.MessageWithSignatures, len(msgs))
	for i, msg := range msgs {
		msgWithSignatures, err := consensus.ToMessageWithSignatures(msg, k.cdc)
		if err != nil {
			return nil, err
		}
		res[i] = msgWithSignatures
	}

	return &types.QueryQueuedMessagesForAttestingResponse{
		Messages: res,
	}, nil
}
