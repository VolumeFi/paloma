package keeper

import (
	"context"
	"errors"
	"fmt"

	"cosmossdk.io/math"
	sdkmath "cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	keeperutil "github.com/palomachain/paloma/v2/util/keeper"
	"github.com/palomachain/paloma/v2/util/libcons"
	"github.com/palomachain/paloma/v2/util/liblog"
	"github.com/palomachain/paloma/v2/util/slice"
	"github.com/palomachain/paloma/v2/x/consensus/keeper/consensus"
	consensustypes "github.com/palomachain/paloma/v2/x/consensus/types"
	"github.com/palomachain/paloma/v2/x/evm/types"
	metrixtypes "github.com/palomachain/paloma/v2/x/metrix/types"
)

const (
	cMaxSubmitLogicCallRetries     uint32 = 2
	cMaxUploadSmartContractRetries uint32 = 2
)

type msgAttester func(sdk.Context, consensus.Queuer, consensustypes.QueuedSignedMessageI, any) error

func (k Keeper) attestMessageWrapper(ctx context.Context, q consensus.Queuer, msg consensustypes.QueuedSignedMessageI, fn msgAttester) (retErr error) {
	if len(msg.GetEvidence()) == 0 {
		return nil
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	logger := k.Logger(ctx).WithFields(
		"component", "attest-message",
		"msg-id", msg.GetId(),
		"msg-nonce", msg.Nonce())
	logger.Debug("attest-message")

	cacheCtx, writeCache := sdkCtx.CacheContext()
	defer func() {
		// If there is no error, or we can't verify this transaction, we flush
		// the context, so we never see it again
		if retErr == nil || errors.Is(retErr, types.ErrEthTxNotVerified) ||
			errors.Is(retErr, types.ErrEthTxFailed) {
			writeCache()
		}
	}()

	result, err := k.consensusChecker.VerifyEvidence(ctx,
		slice.Map(msg.GetEvidence(), func(evidence *consensustypes.Evidence) libcons.Evidence {
			return evidence
		}),
	)
	if err != nil {
		if errors.Is(err, ErrConsensusNotAchieved) {
			logger.WithFields(
				"total-shares", result.TotalShares,
				"total-votes", result.TotalVotes,
				"distribution", result.Distribution,
			).WithError(err).Error("Consensus not achieved.")
			return nil
		}
		return err
	}

	defer func() {
		// given that there was enough evidence for a proof, regardless of the outcome,
		// we should remove this from the queue as there isn't much that we can do about it.
		if err := q.Remove(cacheCtx, msg.GetId()); err != nil {
			k.Logger(sdkCtx).Error("error removing message, attestMessage", "msg-id", msg.GetId(), "msg-nonce", msg.Nonce())
		}
	}()

	return fn(cacheCtx, q, msg, result.Winner)
}

func (k Keeper) attestRouter(ctx context.Context, q consensus.Queuer, msg consensustypes.QueuedSignedMessageI) (err error) {
	return k.attestMessageWrapper(ctx, q, msg, k.routerAttester)
}

func (k Keeper) routerAttester(sdkCtx sdk.Context, q consensus.Queuer, msg consensustypes.QueuedSignedMessageI, winner any) error {
	consensusMsg, err := msg.ConsensusMsg(k.cdc)
	if err != nil {
		k.Logger(sdkCtx).WithError(err).Error("failed to cast to consensus message")
		return err
	}

	message := consensusMsg.(*types.Message)

	defer func() {
		success := false

		// if the input type is a TX, then regardles, we want to set it as already processed
		switch winner := winner.(type) {
		case *types.TxExecutedProof:
			success = true
			tx, err := winner.GetTX()
			if err == nil {
				k.setTxAsAlreadyProcessed(sdkCtx, tx)
			}
		}

		handledAt := msg.GetHandledAtBlockHeight()
		if handledAt == nil {
			handledAt = func(i math.Int) *math.Int { return &i }(sdkmath.NewInt(sdkCtx.BlockHeight()))
		}
		publishMessageAttestedEvent(sdkCtx, &k, msg.GetId(), message.Assignee, message.AssignedAtBlockHeight, *handledAt, success)
	}()

	rawAction := message.GetAction()
	_, chainReferenceID := q.ChainInfo()
	logger := k.Logger(sdkCtx).WithFields("chain-reference-id", chainReferenceID)

	params := attestionParameters{
		originalMessage:  msg,
		msgID:            msg.GetId(),
		chainReferenceID: chainReferenceID,
		rawEvidence:      winner,
		msg:              message,
	}

	switch winner := winner.(type) {
	case *types.TxExecutedProof:
		// If we have proof of a transaction, we need to check the status on
		// the receipt
		receipt, err := winner.GetReceipt()
		if err != nil {
			logger.WithFields("message-id", msg.GetId()).
				WithError(err).Error("Failed to get transaction receipt")
			return err
		}

		if receipt.Status != ethtypes.ReceiptStatusSuccessful {
			logger.WithFields("message-id", msg.GetId(), "status", receipt.Status).
				Warn("Transaction execution failed")
			return types.ErrEthTxFailed
		}
	}

	switch rawAction.(type) {
	case *types.Message_UploadSmartContract:
		return newUploadSmartContractAttester(&k, logger, params).Execute(sdkCtx)
	case *types.Message_UploadUserSmartContract:
		return newUploadUserSmartContractAttester(&k, logger, params).Execute(sdkCtx)
	case *types.Message_UpdateValset:
		return newUpdateValsetAttester(&k, logger, q, params).Execute(sdkCtx)
	case *types.Message_SubmitLogicCall:
		return newSubmitLogicCallAttester(&k, logger, params).Execute(sdkCtx)
	case *types.Message_CompassHandover:
		return newCompassHandoverAttester(&k, logger, params).Execute(sdkCtx)
	}

	return nil
}

func publishMessageAttestedEvent(ctx context.Context, k *Keeper, msgID uint64, assignee string, assignedAt math.Int, handledAt math.Int, successful bool) {
	valAddr, err := sdk.ValAddressFromBech32(assignee)
	if err != nil {
		liblog.FromSDKLogger(k.Logger(ctx)).WithError(err).WithFields("assignee", assignee, "msg-id", msgID).Error("failed to get validator address from bech32.")
	}

	for _, v := range k.onMessageAttestedListeners {
		v.OnConsensusMessageAttested(ctx, metrixtypes.MessageAttestedEvent{
			AssignedAtBlockHeight:  assignedAt,
			HandledAtBlockHeight:   handledAt,
			Assignee:               valAddr,
			MessageID:              msgID,
			WasRelayedSuccessfully: successful,
		})
	}
}

func attestTransactionIntegrity(
	ctx context.Context,
	msg consensustypes.QueuedSignedMessageI,
	k *Keeper,
	proof *types.TxExecutedProof,
	chainReferenceID, relayer string,
	verifyTx func(context.Context, *ethtypes.Transaction, consensustypes.QueuedSignedMessageI, *types.Valset, *types.SmartContract, string) error,
) (*ethtypes.Transaction, error) {
	// check if correct thing was called
	tx, err := proof.GetTX()
	if err != nil {
		return nil, fmt.Errorf("failed to get TX: %w", err)
	}
	if k.isTxProcessed(ctx, tx) {
		// somebody submitted the old transaction that was already processed?
		// punish those validators!!
		return nil, ErrUnexpectedError.JoinErrorf("transaction %s is already processed", tx.Hash())
	}

	compass, err := k.GetLastCompassContract(ctx)
	if err != nil {
		return nil, err
	}

	var valset types.Valset
	var valsetID uint64

	if publicAccessData := msg.GetPublicAccessData(); publicAccessData != nil {
		valsetID = publicAccessData.GetValsetID()
	}

	if valsetID != 0 {
		snapshot, err := k.Valset.FindSnapshotByID(ctx, valsetID)
		// A snapshot may not yet exist if the chain is just being added, but we
		// need to continue in case we're attesting to the initial compass
		// deployment
		if err != nil {
			if !errors.Is(err, keeperutil.ErrNotFound) {
				// There was some other error accessing the store, bail
				return nil, err
			}
		} else {
			// There was no error, so convert the snapshot to valset
			logger := liblog.FromSDKLogger(k.Logger(ctx))
			valset = transformSnapshotToCompass(snapshot, chainReferenceID, logger)
		}
	}

	err = verifyTx(ctx, tx, msg, &valset, compass, relayer)
	if err != nil {
		// passed in transaction doesn't seem to be created from this smart contract
		return nil, fmt.Errorf("tx failed to verify: %w", err)
	}

	return tx, nil
}

func (k Keeper) SetSmartContractAsActive(ctx context.Context, smartContractID uint64, chainReferenceID string) (err error) {
	logger := liblog.FromSDKLogger(k.Logger(ctx))
	defer func() {
		if err == nil {
			logger.With("smart-contract-id", smartContractID).Debug("removing deployment.")
			k.DeleteSmartContractDeploymentByContractID(ctx, smartContractID, chainReferenceID)
		}
	}()

	deployment, _ := k.getSmartContractDeploymentByContractID(ctx, smartContractID, chainReferenceID)
	if deployment.GetStatus() != types.SmartContractDeployment_WAITING_FOR_ERC20_OWNERSHIP_TRANSFER {
		logger.WithError(err).Error("Deployment not awaiting transfer")
		return ErrCannotActiveSmartContractThatIsNotDeploying
	}

	smartContract, err := k.getSmartContract(ctx, deployment.GetSmartContractID())
	if err != nil {
		logger.WithError(err).Error("Failed to get contract")
		return err
	}

	err = k.ActivateChainReferenceID(
		ctx,
		chainReferenceID,
		smartContract,
		deployment.NewSmartContractAddress,
		deployment.GetUniqueID(),
	)
	if err != nil {
		logger.WithError(err).Error("Failed to activate chain")
		return err
	}

	return nil
}

func (k Keeper) txAlreadyProcessedStore(ctx context.Context) storetypes.KVStore {
	s := runtime.KVStoreAdapter(k.storeKey.OpenKVStore(ctx))
	return prefix.NewStore(s, []byte("tx-processed"))
}

func (k Keeper) setTxAsAlreadyProcessed(ctx context.Context, tx *ethtypes.Transaction) {
	kv := k.txAlreadyProcessedStore(ctx)
	kv.Set(tx.Hash().Bytes(), []byte{1})
}

func (k Keeper) isTxProcessed(ctx context.Context, tx *ethtypes.Transaction) bool {
	kv := k.txAlreadyProcessedStore(ctx)
	return kv.Has(tx.Hash().Bytes())
}
