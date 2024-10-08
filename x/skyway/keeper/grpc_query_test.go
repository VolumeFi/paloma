package keeper

import (
	"bytes"
	"context"
	"math/big"
	"testing"
	"time"

	"cosmossdk.io/math"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/palomachain/paloma/v2/x/skyway/types"
	vtypes "github.com/palomachain/paloma/v2/x/valset/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// nolint: exhaustruct
func TestLastPendingBatchRequest(t *testing.T) {
	specs := map[string]struct {
		expResp types.QueryLastPendingBatchRequestByAddrResponse
	}{
		"find batch": {
			expResp: types.QueryLastPendingBatchRequestByAddrResponse{
				Batch: []types.OutgoingTxBatch{
					{
						BatchNonce: 1,
						Transactions: []types.OutgoingTransferTx{
							{
								Id:          4,
								Sender:      "paloma1qyqszqgpqyqszqgpqyqszqgpqyqszqgp2kvale",
								DestAddress: "0x320915BD0F1bad11cBf06e85D5199DBcAC4E9934",
								Erc20Token: types.ERC20Token{
									Amount:           math.NewInt(103),
									Contract:         testERC20Address,
									ChainReferenceId: "test-chain",
								},
								BridgeTaxAmount: math.ZeroInt(),
							},
							{
								Id:          3,
								Sender:      "paloma1qyqszqgpqyqszqgpqyqszqgpqyqszqgp2kvale",
								DestAddress: "0x320915BD0F1bad11cBf06e85D5199DBcAC4E9934",
								Erc20Token: types.ERC20Token{
									Amount:           math.NewInt(102),
									Contract:         testERC20Address,
									ChainReferenceId: "test-chain",
								},
								BridgeTaxAmount: math.ZeroInt(),
							},
						},
						TokenContract:      testERC20Address,
						PalomaBlockCreated: 1235067,
						ChainReferenceId:   "test-chain",
					},
				},
			},
		},
	}
	// any lower than this and a validator won't be created
	const minStake = 1000000
	input, _ := SetupTestChain(t, []uint64{minStake, minStake, minStake, minStake, minStake})
	sdkCtx := sdk.UnwrapSDKContext(input.Context)

	defer func() { sdkCtx.Logger().Info("Asserting invariants at test end"); input.AssertInvariants() }()

	ctx := sdk.WrapSDKContext(sdkCtx)
	var valAddr sdk.AccAddress = bytes.Repeat([]byte{byte(1)}, 20)
	createTestBatch(t, input, 2)
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			req := new(types.QueryLastPendingBatchRequestByAddrRequest)
			req.Address = valAddr.String()
			got, err := input.SkywayKeeper.LastPendingBatchRequestByAddr(ctx, req)
			require.NoError(t, err)

			// Don't bother comparing some computed values
			got.Batch[0].BatchTimeout = 0
			got.Batch[0].BytesToSign = nil
			got.Batch[0].Assignee = ""
			got.Batch[0].AssigneeRemoteAddress = nil

			assert.Equal(t, &spec.expResp, got, got)
		})
	}
}

// nolint: exhaustruct
func createTestBatch(t *testing.T, input TestInput, maxTxElements uint) {
	sdkCtx := sdk.UnwrapSDKContext(input.Context)
	var (
		mySender   = bytes.Repeat([]byte{1}, 20)
		myReceiver = "0x320915BD0F1bad11cBf06e85D5199DBcAC4E9934"
		now        = time.Now().UTC()
	)
	receiver, err := types.NewEthAddress(myReceiver)
	require.NoError(t, err)
	tokenContract, err := types.NewEthAddress(testERC20Address)
	require.NoError(t, err)
	// mint some voucher first
	token, err := types.NewInternalERC20Token(math.NewInt(99999), testERC20Address, "test-chain")
	require.NoError(t, err)
	allVouchers := sdk.Coins{sdk.NewCoin(testDenom, token.Amount)}
	err = input.BankKeeper.MintCoins(input.Context, types.ModuleName, allVouchers)
	require.NoError(t, err)

	// set senders balance
	input.AccountKeeper.NewAccountWithAddress(input.Context, mySender)
	err = input.BankKeeper.SendCoinsFromModuleToAccount(input.Context, types.ModuleName, mySender, allVouchers)
	require.NoError(t, err)

	// add some TX to the pool
	for i := 0; i < 4; i++ {
		amountToken, err := types.NewInternalERC20Token(math.NewInt(int64(i+100)), testERC20Address, "test-chain")
		require.NoError(t, err)
		amount := sdk.NewCoin(testDenom, amountToken.Amount)
		_, err = input.SkywayKeeper.AddToOutgoingPool(input.Context, mySender, *receiver, amount, "test-chain")
		require.NoError(t, err)
		// Should create:
		// 1: amount 100
		// 2: amount 101
		// 3: amount 102
		// 4: amount 103
	}
	// when
	input.Context = sdkCtx.WithBlockTime(now)

	// tx batch size is 2, so that some of them stay behind
	_, err = input.SkywayKeeper.BuildOutgoingTXBatch(input.Context, "test-chain", *tokenContract, maxTxElements)
	require.NoError(t, err)
	// Should have 2 and 3 from above
	// 1 and 4 should be unbatched
}

// nolint: exhaustruct
func TestQueryAllBatchConfirms(t *testing.T) {
	input, ctx := SetupFiveValChain(t)
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	defer func() { sdkCtx.Logger().Info("Asserting invariants at test end"); input.AssertInvariants() }()

	k := input.SkywayKeeper

	var (
		tokenContract      = "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"
		validatorAddr, err = sdk.AccAddressFromBech32("paloma1mgamdcs9dah0vn0gqupl05up7pedg2mvyy5e9j")
	)
	require.NoError(t, err)

	_, err = input.SkywayKeeper.SetBatchConfirm(sdkCtx, &types.MsgConfirmBatch{
		Nonce:         1,
		TokenContract: tokenContract,
		EthSigner:     "0xf35e2cc8e6523d683ed44870f5b7cc785051a77d",
		Orchestrator:  validatorAddr.String(),
		Signature:     "d34db33f",
		Metadata: vtypes.MsgMetadata{
			Creator: validatorAddr.String(),
			Signers: []string{validatorAddr.String()},
		},
	})
	require.NoError(t, err)

	batchConfirms, err := k.BatchConfirms(ctx, &types.QueryBatchConfirmsRequest{Nonce: 1, ContractAddress: tokenContract})
	require.NoError(t, err)

	expectedRes := types.QueryBatchConfirmsResponse{
		Confirms: []types.MsgConfirmBatch{
			{
				Nonce:         1,
				TokenContract: "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B",
				EthSigner:     "0xf35e2cc8e6523d683ed44870f5b7cc785051a77d",
				Orchestrator:  "paloma1mgamdcs9dah0vn0gqupl05up7pedg2mvyy5e9j",
				Signature:     "d34db33f",
				Metadata: vtypes.MsgMetadata{
					Creator: "paloma1mgamdcs9dah0vn0gqupl05up7pedg2mvyy5e9j",
					Signers: []string{"paloma1mgamdcs9dah0vn0gqupl05up7pedg2mvyy5e9j"},
				},
			},
		},
	}

	assert.Equal(t, &expectedRes, batchConfirms, "json is equal")
}

// nolint: exhaustruct
// Check with multiple nonces and tokenContracts
func TestQueryBatch(t *testing.T) {
	input, ctx := SetupFiveValChain(t)
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	defer func() { sdkCtx.Logger().Info("Asserting invariants at test end"); input.AssertInvariants() }()

	k := input.SkywayKeeper

	createTestBatch(t, input, 2)

	batch, err := k.BatchRequestByNonce(ctx, &types.QueryBatchRequestByNonceRequest{Nonce: 1, ContractAddress: testERC20Address})
	require.NoError(t, err)

	expectedRes := types.QueryBatchRequestByNonceResponse{
		Batch: types.OutgoingTxBatch{
			BatchTimeout: 0,
			Transactions: []types.OutgoingTransferTx{
				{
					DestAddress: "0x320915BD0F1bad11cBf06e85D5199DBcAC4E9934",
					Erc20Token: types.ERC20Token{
						Amount:           math.NewInt(103),
						Contract:         testERC20Address,
						ChainReferenceId: "test-chain",
					},
					Sender:          "paloma1qyqszqgpqyqszqgpqyqszqgpqyqszqgp2kvale",
					Id:              4,
					BridgeTaxAmount: math.ZeroInt(),
				},
				{
					DestAddress: "0x320915BD0F1bad11cBf06e85D5199DBcAC4E9934",
					Erc20Token: types.ERC20Token{
						Amount:           math.NewInt(102),
						Contract:         testERC20Address,
						ChainReferenceId: "test-chain",
					},
					Sender:          "paloma1qyqszqgpqyqszqgpqyqszqgpqyqszqgp2kvale",
					Id:              3,
					BridgeTaxAmount: math.ZeroInt(),
				},
			},
			BatchNonce:         1,
			PalomaBlockCreated: 1234567,
			TokenContract:      testERC20Address,
			ChainReferenceId:   "test-chain",
		},
	}

	// Don't bother comparing some computed values
	batch.Batch.BatchTimeout = 0
	batch.Batch.BytesToSign = nil
	batch.Batch.Assignee = ""
	batch.Batch.AssigneeRemoteAddress = nil

	assert.Equal(t, &expectedRes, batch, batch)
}

// nolint: exhaustruct
func TestLastBatchesRequest(t *testing.T) {
	input, ctx := SetupFiveValChain(t)
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	defer func() { sdkCtx.Logger().Info("Asserting invariants at test end"); input.AssertInvariants() }()

	// evm/test-chain/evm-turnstone-message
	// evm/test-chain/evm-turnstone-message
	input.EvmKeeper.AddSupportForNewChain(ctx, "test-chain", 42, 100, "0x123", big.NewInt(0))

	k := input.SkywayKeeper

	createTestBatch(t, input, 2)
	createTestBatch(t, input, 3)

	lastBatches, err := k.OutgoingTxBatches(ctx, &types.QueryOutgoingTxBatchesRequest{
		ChainReferenceId: "test-chain",
	})
	require.NoError(t, err)

	// Should have no batches, since no gas estimates are set
	require.Len(t, lastBatches.Batches, 0)

	// Make sure gas is set, otherwise we don't hand it out for relaying
	oldCheckpoints := make(map[uint64][]byte)
	k.IterateOutgoingTxBatches(ctx, func(key []byte, batch types.InternalOutgoingTxBatch) bool {
		oldCheckpoints[batch.BatchNonce] = batch.BytesToSign
		k.UpdateBatchGasEstimate(ctx, batch, 21_000)
		return false
	})
	// Make sure the bytes to sign are changed after updating the gas estimate
	k.IterateOutgoingTxBatches(ctx, func(key []byte, batch types.InternalOutgoingTxBatch) bool {
		require.NotEqual(t, oldCheckpoints[batch.BatchNonce], batch.BytesToSign, "should have changed the bytes to sign")
		return false
	})

	lastBatches, err = k.OutgoingTxBatches(ctx, &types.QueryOutgoingTxBatchesRequest{
		ChainReferenceId: "test-chain",
	})
	require.NoError(t, err)

	expectedRes := types.QueryOutgoingTxBatchesResponse{
		Batches: []types.OutgoingTxBatch{
			{
				GasEstimate: 21_000,
				Transactions: []types.OutgoingTransferTx{
					{
						DestAddress: "0x320915BD0F1bad11cBf06e85D5199DBcAC4E9934",
						Erc20Token: types.ERC20Token{
							Amount:           math.NewInt(103),
							Contract:         testERC20Address,
							ChainReferenceId: "test-chain",
						},
						Sender:          "paloma1qyqszqgpqyqszqgpqyqszqgpqyqszqgp2kvale",
						Id:              8,
						BridgeTaxAmount: math.ZeroInt(),
					},
					{
						DestAddress: "0x320915BD0F1bad11cBf06e85D5199DBcAC4E9934",
						Erc20Token: types.ERC20Token{
							Amount:           math.NewInt(102),
							Contract:         testERC20Address,
							ChainReferenceId: "test-chain",
						},
						Sender:          "paloma1qyqszqgpqyqszqgpqyqszqgpqyqszqgp2kvale",
						Id:              7,
						BridgeTaxAmount: math.ZeroInt(),
					},
					{
						DestAddress: "0x320915BD0F1bad11cBf06e85D5199DBcAC4E9934",
						Erc20Token: types.ERC20Token{
							Amount:           math.NewInt(101),
							Contract:         testERC20Address,
							ChainReferenceId: "test-chain",
						},
						Sender:          "paloma1qyqszqgpqyqszqgpqyqszqgpqyqszqgp2kvale",
						Id:              6,
						BridgeTaxAmount: math.ZeroInt(),
					},
				},
				BatchNonce:         2,
				PalomaBlockCreated: 1234567,
				TokenContract:      testERC20Address,
				ChainReferenceId:   "test-chain",
			},
			{
				GasEstimate: 21_000,
				Transactions: []types.OutgoingTransferTx{
					{
						DestAddress: "0x320915BD0F1bad11cBf06e85D5199DBcAC4E9934",
						Erc20Token: types.ERC20Token{
							Amount:           math.NewInt(103),
							Contract:         testERC20Address,
							ChainReferenceId: "test-chain",
						},
						Sender:          "paloma1qyqszqgpqyqszqgpqyqszqgpqyqszqgp2kvale",
						Id:              4,
						BridgeTaxAmount: math.ZeroInt(),
					},
					{
						DestAddress: "0x320915BD0F1bad11cBf06e85D5199DBcAC4E9934",
						Erc20Token: types.ERC20Token{
							Amount:           math.NewInt(102),
							Contract:         testERC20Address,
							ChainReferenceId: "test-chain",
						},
						Sender:          "paloma1qyqszqgpqyqszqgpqyqszqgpqyqszqgp2kvale",
						Id:              3,
						BridgeTaxAmount: math.ZeroInt(),
					},
				},
				BatchNonce:         1,
				PalomaBlockCreated: 1234567,
				TokenContract:      testERC20Address,
				ChainReferenceId:   "test-chain",
			},
		},
	}

	// Don't bother comparing some computed values
	lastBatches.Batches[0].BatchTimeout = 0
	lastBatches.Batches[0].BytesToSign = nil
	lastBatches.Batches[0].Assignee = ""
	lastBatches.Batches[0].AssigneeRemoteAddress = nil
	lastBatches.Batches[1].BatchTimeout = 0
	lastBatches.Batches[1].BytesToSign = nil
	lastBatches.Batches[1].Assignee = ""
	lastBatches.Batches[1].AssigneeRemoteAddress = nil

	assert.Equal(t, &expectedRes, lastBatches, "json is equal")
}

// nolint: exhaustruct
func TestQueryERC20ToDenom(t *testing.T) {
	var (
		chainReferenceID = "test-chain"
		erc20, err       = types.NewEthAddress(testERC20Address)
	)
	require.NoError(t, err)
	response := types.QueryERC20ToDenomResponse{
		Denom: testDenom,
	}
	input := CreateTestEnv(t)
	sdkCtx := sdk.UnwrapSDKContext(input.Context)
	defer func() { sdkCtx.Logger().Info("Asserting invariants at test end"); input.AssertInvariants() }()

	ctx := sdk.UnwrapSDKContext(sdkCtx)
	k := input.SkywayKeeper
	err = input.SkywayKeeper.setDenomToERC20(sdkCtx, chainReferenceID, testDenom, *erc20)
	require.NoError(t, err)

	queriedDenom, err := k.ERC20ToDenom(ctx, &types.QueryERC20ToDenomRequest{
		Erc20:            erc20.GetAddress().Hex(),
		ChainReferenceId: chainReferenceID,
	})
	require.NoError(t, err)

	assert.Equal(t, &response, queriedDenom)
}

// nolint: exhaustruct
func TestQueryDenomToERC20(t *testing.T) {
	var (
		chainReferenceID = "test-chain"
		erc20, err       = types.NewEthAddress(testERC20Address)
	)
	require.NoError(t, err)
	response := types.QueryDenomToERC20Response{
		Erc20: erc20.GetAddress().Hex(),
	}
	input := CreateTestEnv(t)
	defer func() {
		sdk.UnwrapSDKContext(input.Context).Logger().Info("Asserting invariants at test end")
		input.AssertInvariants()
	}()

	sdkCtx := input.Context

	k := input.SkywayKeeper

	err = input.SkywayKeeper.setDenomToERC20(sdkCtx, chainReferenceID, testDenom, *erc20)
	require.NoError(t, err)

	queriedERC20, err := k.DenomToERC20(input.Context, &types.QueryDenomToERC20Request{
		Denom:            testDenom,
		ChainReferenceId: chainReferenceID,
	})
	require.NoError(t, err)

	assert.Equal(t, &response, queriedERC20)
}

// nolint: exhaustruct
func TestQueryPendingSendToRemote(t *testing.T) {
	input, ctx := SetupFiveValChain(t)
	sdkCtx := sdk.UnwrapSDKContext(input.Context)
	defer func() { sdkCtx.Logger().Info("Asserting invariants at test end"); input.AssertInvariants() }()
	k := input.SkywayKeeper
	var (
		now            = time.Now().UTC()
		mySender, err1 = sdk.AccAddressFromBech32("paloma1ahx7f8wyertuus9r20284ej0asrs085c945jyk")
		myReceiver     = "0xd041c41EA1bf0F006ADBb6d2c9ef9D425dE5eaD7"
		token, err2    = types.NewInternalERC20Token(math.NewInt(99999), testERC20Address, "test-chain")
		allVouchers    = sdk.NewCoins(sdk.NewCoin(testDenom, token.Amount))
	)
	require.NoError(t, err1)
	require.NoError(t, err2)
	receiver, err := types.NewEthAddress(myReceiver)
	require.NoError(t, err)
	tokenContract, err := types.NewEthAddress(testERC20Address)
	require.NoError(t, err)

	// mint some voucher first
	require.NoError(t, input.BankKeeper.MintCoins(sdkCtx, types.ModuleName, allVouchers))
	// set senders balance
	input.AccountKeeper.NewAccountWithAddress(sdkCtx, mySender)
	require.NoError(t, input.BankKeeper.SendCoinsFromModuleToAccount(sdkCtx, types.ModuleName, mySender, allVouchers))

	// CREATE FIRST BATCH
	// ==================

	// add some TX to the pool
	for i := 0; i < 4; i++ {
		amountToken, err := types.NewInternalERC20Token(math.NewInt(int64(i+100)), testERC20Address, "test-chain")
		require.NoError(t, err)
		amount := sdk.NewCoin(testDenom, amountToken.Amount)
		_, err = input.SkywayKeeper.AddToOutgoingPool(sdkCtx, mySender, *receiver, amount, "test-chain")
		require.NoError(t, err)
		// Should create:
		// 1: amount 100
		// 2: amount 101
		// 3: amount 102
		// 4: amount 104
	}

	// when
	sdkCtx = sdkCtx.WithBlockTime(now)

	// tx batch size is 2, so that some of them stay behind
	// Should contain 2 and 3 from above
	_, err = input.SkywayKeeper.BuildOutgoingTXBatch(sdkCtx, "test-chain", *tokenContract, 2)
	require.NoError(t, err)

	// Should receive 1 and 4 unbatched, 2 and 3 batched in response
	response, err := k.GetPendingSendToRemote(ctx, &types.QueryPendingSendToRemote{SenderAddress: mySender.String()})
	require.NoError(t, err)
	expectedRes := types.QueryPendingSendToRemoteResponse{
		TransfersInBatches: []types.OutgoingTransferTx{
			{
				Id:          4,
				Sender:      "paloma1ahx7f8wyertuus9r20284ej0asrs085c945jyk",
				DestAddress: "0xd041c41EA1bf0F006ADBb6d2c9ef9D425dE5eaD7",
				Erc20Token: types.ERC20Token{
					Contract:         testERC20Address,
					Amount:           math.NewInt(103),
					ChainReferenceId: "test-chain",
				},
				BridgeTaxAmount: math.ZeroInt(),
			},
			{
				Id:          3,
				Sender:      "paloma1ahx7f8wyertuus9r20284ej0asrs085c945jyk",
				DestAddress: "0xd041c41EA1bf0F006ADBb6d2c9ef9D425dE5eaD7",
				Erc20Token: types.ERC20Token{
					Contract:         testERC20Address,
					Amount:           math.NewInt(102),
					ChainReferenceId: "test-chain",
				},
				BridgeTaxAmount: math.ZeroInt(),
			},
		},

		UnbatchedTransfers: []types.OutgoingTransferTx{
			{
				Id:          2,
				Sender:      "paloma1ahx7f8wyertuus9r20284ej0asrs085c945jyk",
				DestAddress: "0xd041c41EA1bf0F006ADBb6d2c9ef9D425dE5eaD7",
				Erc20Token: types.ERC20Token{
					Contract:         testERC20Address,
					Amount:           math.NewInt(101),
					ChainReferenceId: "test-chain",
				},
				BridgeTaxAmount: math.ZeroInt(),
			},
			{
				Id:          1,
				Sender:      "paloma1ahx7f8wyertuus9r20284ej0asrs085c945jyk",
				DestAddress: "0xd041c41EA1bf0F006ADBb6d2c9ef9D425dE5eaD7",
				Erc20Token: types.ERC20Token{
					Contract:         testERC20Address,
					Amount:           math.NewInt(100),
					ChainReferenceId: "test-chain",
				},
				BridgeTaxAmount: math.ZeroInt(),
			},
		},
	}

	assert.Equal(t, &expectedRes, response, "json is equal")
}

func TestGetUnobservedBlocksByAddr(t *testing.T) {
	chainReferenceID := "test-chain"
	address := common.BytesToAddress(bytes.Repeat([]byte{0x1}, 20)).String()

	tests := []struct {
		name     string
		setup    func(context.Context, Keeper)
		expected []uint64
	}{
		{
			name: "three unobserved messages for current compass",
			setup: func(ctx context.Context, k Keeper) {
				sdkCtx := sdktypes.UnwrapSDKContext(ctx)

				for i := 0; i < 3; i++ {
					msg := types.MsgLightNodeSaleClaim{
						SkywayNonce:    uint64(i + 1),
						EthBlockHeight: uint64(i + 1),
					}
					claim, err := codectypes.NewAnyWithValue(&msg)
					require.NoError(t, err)

					hash, err := msg.ClaimHash()
					require.NoError(t, err)

					att := &types.Attestation{
						Observed: false,
						Votes:    []string{},
						Height:   uint64(sdkCtx.BlockHeight()),
						Claim:    claim,
					}

					k.SetAttestation(ctx, chainReferenceID, uint64(i+1), hash, att)
				}
			},
			expected: []uint64{1, 2, 3},
		},
		{
			name: "three observed messages for current compass",
			setup: func(ctx context.Context, k Keeper) {
				sdkCtx := sdktypes.UnwrapSDKContext(ctx)

				for i := 0; i < 3; i++ {
					msg := types.MsgLightNodeSaleClaim{
						SkywayNonce:    uint64(i + 1),
						EthBlockHeight: uint64(i + 1),
					}
					claim, err := codectypes.NewAnyWithValue(&msg)
					require.NoError(t, err)

					hash, err := msg.ClaimHash()
					require.NoError(t, err)

					att := &types.Attestation{
						Observed: true,
						Votes:    []string{},
						Height:   uint64(sdkCtx.BlockHeight()),
						Claim:    claim,
					}

					k.SetAttestation(ctx, chainReferenceID, uint64(i+1), hash, att)
				}
			},
			expected: nil,
		},
		{
			name: "unobserved messages already signed by validator",
			setup: func(ctx context.Context, k Keeper) {
				sdkCtx := sdktypes.UnwrapSDKContext(ctx)

				for i := 0; i < 3; i++ {
					msg := types.MsgLightNodeSaleClaim{
						SkywayNonce:    uint64(i + 1),
						EthBlockHeight: uint64(i + 1),
					}
					claim, err := codectypes.NewAnyWithValue(&msg)
					require.NoError(t, err)

					hash, err := msg.ClaimHash()
					require.NoError(t, err)

					att := &types.Attestation{
						Observed: false,
						Votes:    []string{address},
						Height:   uint64(sdkCtx.BlockHeight()),
						Claim:    claim,
					}

					k.SetAttestation(ctx, chainReferenceID, uint64(i+1), hash, att)
				}
			},
			expected: nil,
		},
		{
			name: "unobserved messages behind last observed nonce",
			setup: func(ctx context.Context, k Keeper) {
				sdkCtx := sdktypes.UnwrapSDKContext(ctx)

				err := k.setLastObservedSkywayNonce(ctx, chainReferenceID, 2)
				require.NoError(t, err)

				for i := 0; i < 3; i++ {
					msg := types.MsgLightNodeSaleClaim{
						SkywayNonce:    uint64(i + 1),
						EthBlockHeight: uint64(i + 1),
					}
					claim, err := codectypes.NewAnyWithValue(&msg)
					require.NoError(t, err)

					hash, err := msg.ClaimHash()
					require.NoError(t, err)

					att := &types.Attestation{
						Observed: false,
						Votes:    []string{},
						Height:   uint64(sdkCtx.BlockHeight()),
						Claim:    claim,
					}

					k.SetAttestation(ctx, chainReferenceID, uint64(i+1), hash, att)
				}
			},
			expected: []uint64{3},
		},
		{
			name: "unobserved messages for current and porevious compass",
			setup: func(ctx context.Context, k Keeper) {
				sdkCtx := sdktypes.UnwrapSDKContext(ctx)

				k.setLatestCompassID(ctx, chainReferenceID, "2")

				msg := types.MsgLightNodeSaleClaim{
					SkywayNonce:    2,
					EthBlockHeight: 2,
					CompassId:      "1",
				}
				claim, err := codectypes.NewAnyWithValue(&msg)
				require.NoError(t, err)

				hash, err := msg.ClaimHash()
				require.NoError(t, err)

				att := &types.Attestation{
					Observed: false,
					Votes:    []string{},
					Height:   uint64(sdkCtx.BlockHeight()),
					Claim:    claim,
				}

				k.SetAttestation(ctx, chainReferenceID, 2, hash, att)

				msg = types.MsgLightNodeSaleClaim{
					SkywayNonce:    1,
					EthBlockHeight: 1,
					CompassId:      "2",
				}
				claim, err = codectypes.NewAnyWithValue(&msg)
				require.NoError(t, err)

				hash, err = msg.ClaimHash()
				require.NoError(t, err)

				att = &types.Attestation{
					Observed: false,
					Votes:    []string{},
					Height:   uint64(sdkCtx.BlockHeight()),
					Claim:    claim,
				}

				k.SetAttestation(ctx, chainReferenceID, 1, hash, att)
			},
			expected: []uint64{1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := CreateTestEnv(t)

			k := input.SkywayKeeper
			ctx := input.Context

			tt.setup(ctx, k)

			req := &types.QueryUnobservedBlocksByAddrRequest{
				ChainReferenceId: chainReferenceID,
				Address:          address,
			}

			res, err := k.GetUnobservedBlocksByAddr(ctx, req)
			require.NoError(t, err)
			require.Equal(t, tt.expected, res.Blocks)
		})
	}
}
