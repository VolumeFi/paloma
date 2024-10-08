package keeper

import (
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/palomachain/paloma/v2/x/evm/types"
	metrixtypes "github.com/palomachain/paloma/v2/x/metrix/types"
	valsettypes "github.com/palomachain/paloma/v2/x/valset/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	cDummyFeeMgrAddress         = "0xb794f5ea0ba39494ce839613fffba74279579268"
	cDummySmartContractDeployer = "0xb794f5ea0ba39494ce839613fffba74279579268"
)

func addDeploymentToKeeper(t *testing.T, ctx sdk.Context, k *Keeper, mockServices mockedServices) {
	unpublishedSnapshot := &valsettypes.Snapshot{
		Id:          1,
		TotalShares: sdkmath.NewInt(75000),
		Validators: getValidators(
			3,
			[]validatorChainInfo{
				{
					chainType:        "evm",
					chainReferenceID: "test-chain",
				},
			},
		),
	}
	// test-chain mocks
	mockServices.ValsetKeeper.On("GetCurrentSnapshot", mock.Anything).Return(unpublishedSnapshot, nil)
	mockServices.ConsensusKeeper.On("PutMessageInQueue", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(uint64(10), nil)
	mockServices.MetrixKeeper.On("Validators", mock.Anything, mock.Anything).Return(&metrixtypes.QueryValidatorsResponse{
		ValMetrics: getMetrics(3),
	}, nil)
	mockServices.TreasuryKeeper.On("GetRelayerFeesByChainReferenceID", mock.Anything, mock.Anything).Return(getFees(3), nil)

	// Add a new chains for our test to use
	err := k.AddSupportForNewChain(
		ctx,
		"test-chain",
		1,
		uint64(123),
		"0x1234",
		big.NewInt(55),
	)
	require.NoError(t, err)
	err = k.SetFeeManagerAddress(ctx, "test-chain", cDummyFeeMgrAddress)
	require.NoError(t, err)

	sc, err := k.SaveNewSmartContract(ctx, contractAbi, common.FromHex(contractBytecodeStr))
	require.NoError(t, err)
	err = k.SetAsCompassContract(ctx, sc)
	require.NoError(t, err)
}

func TestKeeper_RemoveSmartContractDeployment(t *testing.T) {
	t.Run("removes a smart contract deployment", func(t *testing.T) {
		keep, mockServices, ctx := NewEvmKeeper(t)
		addDeploymentToKeeper(t, ctx, keep, mockServices)

		k := msgServer{Keeper: *keep}

		deployments, err := k.AllSmartContractsDeployments(ctx)
		require.Len(t, deployments, 1)
		require.NoError(t, err)

		_, err = k.RemoveSmartContractDeployment(ctx, &types.MsgRemoveSmartContractDeploymentRequest{
			SmartContractID:  1,
			ChainReferenceID: "test-chain",
		})
		require.NoError(t, err)

		deployments, err = k.AllSmartContractsDeployments(ctx)
		require.Len(t, deployments, 0)
		require.NoError(t, err)
	})
}
