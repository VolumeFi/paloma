package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	keeperutil "github.com/palomachain/paloma/v2/util/keeper"
	"github.com/palomachain/paloma/v2/util/slice"
	"github.com/palomachain/paloma/v2/x/valset/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) GetAlivePigeons(goCtx context.Context, req *types.QueryGetAlivePigeonsRequest) (*types.QueryGetAlivePigeonsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	vals := k.GetUnjailedValidators(ctx)

	res := slice.Map(vals, func(val stakingtypes.ValidatorI) *types.QueryGetAlivePigeonsResponse_ValidatorAlive {
		bz, err := keeperutil.ValAddressFromBech32(k.AddressCodec, val.GetOperator())
		if err != nil {
			k.Logger(ctx).Error("error while getting validator address")
			return &types.QueryGetAlivePigeonsResponse_ValidatorAlive{}
		}
		s := &types.QueryGetAlivePigeonsResponse_ValidatorAlive{
			ValAddress: bz,
		}
		data, err := k.ValidatorKeepAliveData(ctx, bz)
		if err != nil {
			s.Error = err.Error()
		} else {
			s.AliveUntilBlockHeight = data.AliveUntilBlockHeight
			s.PigeonVersion = data.PigeonVersion
		}
		return s
	})

	return &types.QueryGetAlivePigeonsResponse{
		AliveValidators: res,
	}, nil
}
