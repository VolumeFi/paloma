package keeper

import (
	"context"
	"fmt"
	"slices"
	"time"

	sdkmath "cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/palomachain/paloma/v2/util/blocks"
	keeperutil "github.com/palomachain/paloma/v2/util/keeper"
	"github.com/palomachain/paloma/v2/util/liblog"
	"github.com/palomachain/paloma/v2/x/consensus/keeper/consensus"
	consensustypes "github.com/palomachain/paloma/v2/x/consensus/types"
	"github.com/palomachain/paloma/v2/x/evm/types"
)

// A store, with the smart contract prefix, for all user smart contracts
func (k Keeper) userSmartContractStore(
	ctx context.Context,
) storetypes.KVStore {
	kvstore := runtime.KVStoreAdapter(k.storeKey.OpenKVStore(ctx))
	return prefix.NewStore(kvstore, []byte(types.UserSmartContractStoreKeyPrefix))
}

// A store for validator smart contracts, prefixed with the validator address
func (k Keeper) userSmartContractStoreByAddress(
	ctx context.Context,
	addr string,
) storetypes.KVStore {
	return prefix.NewStore(k.userSmartContractStore(ctx), []byte(addr))
}

func (k Keeper) UserSmartContracts(
	ctx context.Context,
	addr string,
) ([]*types.UserSmartContract, error) {
	logger := liblog.FromSDKLogger(sdk.UnwrapSDKContext(ctx).Logger())
	logger.WithFields("val_address", addr).Debug("list user smart contracts")

	st := k.userSmartContractStoreByAddress(ctx, addr)

	_, all, err := keeperutil.IterAll[*types.UserSmartContract](st, k.cdc)
	return all, err
}

func (k Keeper) SaveUserSmartContract(
	ctx context.Context,
	addr string,
	c *types.UserSmartContract,
) (uint64, error) {
	logger := liblog.FromSDKLogger(sdk.UnwrapSDKContext(ctx).Logger())
	logger.WithFields("val_address", addr).Debug("save user smart contract")

	if err := c.Validate(); err != nil {
		logger.WithError(err).Warn("save with invalid smart contract")
		return 0, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Create a new contract to make sure fields are properly initialized
	contract := &types.UserSmartContract{
		Author:               addr,
		Id:                   k.ider.IncrementNextID(ctx, types.UserSmartContractStoreKeyPrefix),
		Title:                c.Title,
		AbiJson:              c.AbiJson,
		Bytecode:             c.Bytecode,
		ConstructorInput:     c.ConstructorInput,
		CreatedAtBlockHeight: sdkCtx.BlockHeight(),
		UpdatedAtBlockHeight: sdkCtx.BlockHeight(),
	}

	key := keeperutil.Uint64ToByte(contract.Id)

	st := k.userSmartContractStoreByAddress(ctx, addr)
	return contract.Id, keeperutil.Save(st, k.cdc, key, contract)
}

func (k Keeper) DeleteUserSmartContract(
	ctx context.Context,
	addr string,
	id uint64,
) error {
	logger := liblog.FromSDKLogger(sdk.UnwrapSDKContext(ctx).Logger())
	logger.WithFields("val_address", addr, "id", id).
		Debug("delete user smart contract")

	key := keeperutil.Uint64ToByte(id)

	st := k.userSmartContractStoreByAddress(ctx, addr)

	if !st.Has(key) {
		return fmt.Errorf("contract not found %v", id)
	}

	st.Delete(key)
	return nil
}

func (k Keeper) CreateUserSmartContractDeployment(
	ctx context.Context,
	addr string,
	id uint64,
	targetChain string,
) (uint64, error) {
	logger := liblog.FromSDKLogger(sdk.UnwrapSDKContext(ctx).Logger())
	logger.WithFields("val_address", addr, "id", id, "xchain", targetChain).
		Debug("create user smart contract deployment")

	// Check if target chain is supported
	chainInfo, err := k.GetChainInfo(ctx, targetChain)
	if err != nil {
		logger.WithError(err).Warn("user smart contract deployment on invalid chain")
		return 0, err
	}

	if chainInfo.SmartContractDeployerAddr == "" {
		err = fmt.Errorf("missing smart contract deployer on chain %s", targetChain)
		logger.WithError(err).Warn("failed to deploy contract to chain")
		return 0, err
	}

	key := keeperutil.Uint64ToByte(id)
	st := k.userSmartContractStoreByAddress(ctx, addr)

	contract, err := keeperutil.Load[*types.UserSmartContract](st, k.cdc, key)
	if err != nil {
		return 0, err
	}

	blockHeight := sdk.UnwrapSDKContext(ctx).BlockHeight()

	deployment := &types.UserSmartContract_Deployment{
		ChainReferenceId:     targetChain,
		Status:               types.UserSmartContract_Deployment_IN_FLIGHT,
		CreatedAtBlockHeight: blockHeight,
		UpdatedAtBlockHeight: blockHeight,
	}

	contract.Deployments = append(contract.Deployments, deployment)
	contract.UpdatedAtBlockHeight = blockHeight

	if err := keeperutil.Save(st, k.cdc, key, contract); err != nil {
		return 0, err
	}

	// To deploy the contract, we need to concatenate the contract bytecode with
	// the ABI-encoded contructor input
	bytecode := slices.Concat(common.FromHex(contract.Bytecode),
		common.FromHex(contract.ConstructorInput))

	senderAddr, err := sdk.ValAddressFromBech32(addr)
	if err != nil {
		return 0, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	userSmartContract := &types.UploadUserSmartContract{
		Bytecode:        bytecode,
		DeployerAddress: chainInfo.SmartContractDeployerAddr,
		// 10 minute deadline, same as SLCs
		Deadline:      sdkCtx.BlockTime().Add(10 * time.Minute).Unix(),
		SenderAddress: senderAddr,
		BlockHeight:   blockHeight,
		Id:            id,
	}

	return k.AddUploadUserSmartContractToConsensus(ctx, targetChain,
		string(chainInfo.GetSmartContractUniqueID()), userSmartContract)
}

func (k Keeper) AddUploadUserSmartContractToConsensus(
	ctx context.Context,
	chainReferenceID string,
	turnstoneID string,
	userSmartContract *types.UploadUserSmartContract,
) (uint64, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	assignee, remoteAddr, err := k.PickValidatorForMessage(ctx, chainReferenceID, nil)
	if err != nil {
		return 0, err
	}

	return k.ConsensusKeeper.PutMessageInQueue(
		ctx,
		consensustypes.Queue(
			types.ConsensusTurnstoneMessage,
			xchainType,
			chainReferenceID,
		),
		&types.Message{
			ChainReferenceID: chainReferenceID,
			TurnstoneID:      turnstoneID,
			Action: &types.Message_UploadUserSmartContract{
				UploadUserSmartContract: userSmartContract,
			},
			Assignee:              assignee,
			AssigneeRemoteAddress: remoteAddr,
			AssignedAtBlockHeight: sdkmath.NewInt(sdkCtx.BlockHeight()),
		}, &consensus.PutOptions{
			RequireGasEstimation: true,
			RequireSignatures:    true,
		})
}

func (k Keeper) SetUserSmartContractDeploymentActive(
	ctx context.Context,
	addr string,
	id uint64,
	blockHeight int64,
	targetChain string,
	contractAddr string,
) error {
	logger := liblog.FromSDKLogger(sdk.UnwrapSDKContext(ctx).Logger())
	logger.WithFields(
		"val_address", addr,
		"id", id,
		"xchain", targetChain,
		"contract_addr", contractAddr,
	).Debug("user smart contract deployment success")

	return k.finishUserSmartContractDeployment(ctx, addr, id, blockHeight,
		targetChain, contractAddr, types.UserSmartContract_Deployment_ACTIVE)
}

func (k Keeper) SetUserSmartContractDeploymentError(
	ctx context.Context,
	addr string,
	id uint64,
	blockHeight int64,
	targetChain string,
) error {
	logger := liblog.FromSDKLogger(sdk.UnwrapSDKContext(ctx).Logger())
	logger.WithFields(
		"val_address", addr,
		"id", id,
		"xchain", targetChain,
	).Debug("user smart contract deployment failed")

	return k.finishUserSmartContractDeployment(ctx, addr, id, blockHeight,
		targetChain, "", types.UserSmartContract_Deployment_ERROR)
}

func (k Keeper) finishUserSmartContractDeployment(
	ctx context.Context,
	addr string,
	id uint64,
	blockHeight int64,
	targetChain string,
	contractAddr string,
	status types.UserSmartContract_Deployment_Status,
) error {
	logger := liblog.FromSDKLogger(sdk.UnwrapSDKContext(ctx).Logger())
	logger.WithFields("val_address", addr, "id", id, "xchain", targetChain).
		Debug("finish user smart contract deployment")

	key := keeperutil.Uint64ToByte(id)
	st := k.userSmartContractStoreByAddress(ctx, addr)

	contract, err := keeperutil.Load[*types.UserSmartContract](st, k.cdc, key)
	if err != nil {
		return err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	for i := range contract.Deployments {
		if contract.Deployments[i].ChainReferenceId == targetChain &&
			contract.Deployments[i].CreatedAtBlockHeight == blockHeight {

			contract.Deployments[i].Status = status
			contract.Deployments[i].Address = contractAddr
			contract.Deployments[i].UpdatedAtBlockHeight = sdkCtx.BlockHeight()

			contract.UpdatedAtBlockHeight = sdkCtx.BlockHeight()

			return keeperutil.Save(st, k.cdc, key, contract)
		}
	}

	return fmt.Errorf("contract %v not found for %v", id, targetChain)
}

// Remove contracts that are not updated nor deployed for over 30 days
func (k Keeper) PurgeStaleUserSmartContracts(ctx context.Context) error {
	logger := liblog.FromSDKLogger(sdk.UnwrapSDKContext(ctx).Logger())
	logger.Debug("purging stale user smart contracts")

	st := k.userSmartContractStore(ctx)
	cutoff := sdk.UnwrapSDKContext(ctx).BlockHeight() - blocks.MonthlyHeight

	fn := func(key []byte, contract *types.UserSmartContract) bool {
		if contract.UpdatedAtBlockHeight < cutoff {
			// If this contract was last updated before the cutoff height,
			// remove it
			logger.WithFields(
				"author", contract.Author,
				"id", contract.Id,
			).Debug("removing stale user smart contract")
			st.Delete(key)
		}

		return true
	}

	return keeperutil.IterAllFnc(st, k.cdc, fn)
}
