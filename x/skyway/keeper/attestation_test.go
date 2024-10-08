package keeper

import (
	"bytes"
	"testing"

	"cosmossdk.io/math"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/palomachain/paloma/v2/x/skyway/types"
	vtypes "github.com/palomachain/paloma/v2/x/valset/types"
	"github.com/stretchr/testify/require"
)

func TestGetAndDeleteAttestation(t *testing.T) {
	input := CreateTestEnv(t)
	k := input.SkywayKeeper
	ctx := input.Context
	sdkCtx := sdktypes.UnwrapSDKContext(ctx)
	chainReferenceID := "test-chain"

	length := 10
	_, _, hashes := createAttestations(t, chainReferenceID, length, k, sdkCtx)

	// Get created attestations
	for i := 0; i < length; i++ {
		nonce := uint64(1 + i)
		att := k.GetAttestation(ctx, chainReferenceID, nonce, hashes[i])
		require.NotNil(t, att)
		att = k.GetAttestation(ctx, "fake-chain", nonce, hashes[i])
		require.Nil(t, att)
	}

	recentAttestations, err := k.GetMostRecentAttestations(ctx, "wrong-chain", uint64(length))
	require.NoError(t, err)
	require.Empty(t, recentAttestations)
	recentAttestations, err = k.GetMostRecentAttestations(ctx, chainReferenceID, uint64(length))
	require.NoError(t, err)
	require.True(t, len(recentAttestations) == length)

	// Delete last 3 attestations
	var nilAtt *types.Attestation
	for i := 7; i < length; i++ {
		nonce := uint64(1 + i)
		att := k.GetAttestation(ctx, chainReferenceID, nonce, hashes[i])
		err = k.DeleteAttestation(ctx, chainReferenceID, *att)
		require.NoError(t, err)

		att = k.GetAttestation(ctx, chainReferenceID, nonce, hashes[i])
		require.Equal(t, nilAtt, att)
	}
	recentAttestations, err = k.GetMostRecentAttestations(ctx, chainReferenceID, uint64(10))
	require.NoError(t, err)
	require.True(t, len(recentAttestations) == 7)

	// Check all attestations again
	for i := 0; i < 7; i++ {
		nonce := uint64(1 + i)
		att := k.GetAttestation(ctx, chainReferenceID, nonce, hashes[i])
		require.NotNil(t, att)
	}
	for i := 7; i < length; i++ {
		nonce := uint64(1 + i)
		att := k.GetAttestation(ctx, chainReferenceID, nonce, hashes[i])
		require.Equal(t, nilAtt, att)
	}
}

// Sets up 10 attestations and checks that they are returned in the correct order
func TestGetMostRecentAttestations(t *testing.T) {
	input := CreateTestEnv(t)

	defer func() {
		sdktypes.UnwrapSDKContext(input.Context).Logger().Info("Asserting invariants at test end")
		input.AssertInvariants()
	}()

	k := input.SkywayKeeper
	ctx := input.Context
	chainReferenceID := "test-chain"

	length := 10
	msgs, anys, _ := createAttestations(t, chainReferenceID, length, k, sdktypes.UnwrapSDKContext(ctx))

	recentAttestations, err := k.GetMostRecentAttestations(ctx, chainReferenceID, uint64(length))
	require.NoError(t, err)
	require.True(t, len(recentAttestations) == length,
		"recentAttestations should have len %v but instead has %v", length, len(recentAttestations))
	for n, attest := range recentAttestations {
		require.Equal(t, attest.Claim.GetCachedValue(), anys[n].GetCachedValue(),
			"The %vth claim does not match our message: claim %v\n message %v", n, attest.Claim, msgs[n])
	}
}

func createAttestations(t *testing.T, chainReferenceID string, length int, k Keeper, ctx sdktypes.Context) ([]types.MsgSendToPalomaClaim, []codectypes.Any, [][]byte) {
	msgs := make([]types.MsgSendToPalomaClaim, 0, length)
	anys := make([]codectypes.Any, 0, length)
	hashes := make([][]byte, 0, length)
	for i := 0; i < length; i++ {
		nonce := uint64(1 + i)

		contract := common.BytesToAddress(bytes.Repeat([]byte{0x1}, 20)).String()
		sender := common.BytesToAddress(bytes.Repeat([]byte{0x2}, 20)).String()
		orch := sdktypes.AccAddress(bytes.Repeat([]byte{0x3}, 20)).String()
		receiver := sdktypes.AccAddress(bytes.Repeat([]byte{0x4}, 20)).String()
		msg := types.MsgSendToPalomaClaim{
			SkywayNonce:    nonce,
			EventNonce:     nonce,
			EthBlockHeight: 1,
			TokenContract:  contract,
			Amount:         math.NewInt(10000000000 + int64(i)),
			EthereumSender: sender,
			PalomaReceiver: receiver,
			Orchestrator:   orch,
			Metadata: vtypes.MsgMetadata{
				Creator: receiver,
				Signers: []string{receiver},
			},
		}
		msgs = append(msgs, msg)

		any, err := codectypes.NewAnyWithValue(&msg)
		require.NoError(t, err)
		anys = append(anys, *any)
		att := &types.Attestation{
			Observed: false,
			Votes:    []string{},
			Height:   uint64(ctx.BlockHeight()),
			Claim:    any,
		}
		unpackedClaim, err := k.UnpackAttestationClaim(att)
		require.NoError(t, err)
		err = unpackedClaim.ValidateBasic()
		require.NoError(t, err)
		hash, err := msg.ClaimHash()
		hashes = append(hashes, hash)
		require.NoError(t, err)
		k.SetAttestation(ctx, chainReferenceID, nonce, hash, att)
	}

	return msgs, anys, hashes
}

func TestGetSetLastObservedEthereumBlockHeight(t *testing.T) {
	input := CreateTestEnv(t)
	k := input.SkywayKeeper
	ctx := input.Context
	chainReferenceID := "test-chain"

	ethereumHeight := uint64(7654321)

	err := k.SetLastObservedEthereumBlockHeight(ctx, chainReferenceID, ethereumHeight)
	require.NoError(t, err)

	ethHeight := k.GetLastObservedEthereumBlockHeight(ctx, chainReferenceID)
	require.Equal(t, uint64(sdktypes.UnwrapSDKContext(ctx).BlockHeight()), ethHeight.PalomaBlockHeight)
	require.Equal(t, ethereumHeight, ethHeight.EthereumBlockHeight)
}

func TestGetSetLastEventNonceByValidator(t *testing.T) {
	input, ctx := SetupFiveValChain(t)
	k := input.SkywayKeeper
	chainReferenceID := "test-chain"

	valAddrString := "paloma1ahx7f8wyertuus9r20284ej0asrs085c945jyk"
	valAccAddress, err := sdktypes.AccAddressFromBech32(valAddrString)
	require.NoError(t, err)
	valAccount := k.accountKeeper.NewAccountWithAddress(ctx, valAccAddress)
	require.NotNil(t, valAccount)

	nonce := uint64(1234)
	addrInBytes := valAccount.GetAddress().Bytes()

	// In case this is first time validator is submiting claim, nonce is expected to be LastObservedNonce-1
	err = k.setLastObservedSkywayNonce(ctx, chainReferenceID, nonce)
	require.NoError(t, err)
	skywayNonce, err := k.GetLastSkywayNonceByValidator(ctx, addrInBytes, chainReferenceID)
	require.NoError(t, err)

	require.Equal(t, nonce-1, skywayNonce)

	err = k.SetLastSkywayNonceByValidator(ctx, addrInBytes, chainReferenceID, nonce)
	require.NoError(t, err)

	skywayNonce, err = k.GetLastSkywayNonceByValidator(ctx, addrInBytes, chainReferenceID)
	require.NoError(t, err)
	require.Equal(t, nonce, skywayNonce)
}

func TestInvalidHeight(t *testing.T) {
	input, ctx := SetupFiveValChain(t)
	sdkCtx := sdktypes.UnwrapSDKContext(ctx)
	chainReferenceID := "test-chain"
	defer func() { sdkCtx.Logger().Info("Asserting invariants at test end"); input.AssertInvariants() }()
	pk := input.SkywayKeeper
	msgServer := NewMsgServerImpl(pk)
	log := sdkCtx.Logger()

	val0 := ValAddrs[0]
	sender := AccAddrs[0]
	receiver := EthAddrs[0]

	lastNonce, err := pk.GetLastObservedSkywayNonce(ctx, chainReferenceID)
	require.NoError(t, err)

	lastEthHeight := pk.GetLastObservedEthereumBlockHeight(ctx, chainReferenceID).EthereumBlockHeight
	lastBatchNonce := 0
	tokenContract := "0x2260FAC5E5542a773Aa44fBCfeDf7C193bc2C599"
	goodHeight := lastEthHeight + 1
	batchTimeout := lastEthHeight + 100
	badHeight := batchTimeout

	// Setup a batch with a timeout
	batch := types.OutgoingTxBatch{
		BatchNonce:   uint64(lastBatchNonce + 1),
		BatchTimeout: batchTimeout,
		Transactions: []types.OutgoingTransferTx{{
			Id:          0,
			Sender:      sender.String(),
			DestAddress: receiver.String(),
			Erc20Token: types.ERC20Token{
				Contract: tokenContract,
				Amount:   math.NewInt(1),
			},
		}},
		TokenContract:      tokenContract,
		PalomaBlockCreated: 0,
		ChainReferenceId:   "test-chain",
	}
	b, err := batch.ToInternal()
	require.NoError(t, err)
	err = pk.StoreBatch(ctx, *b)
	require.NoError(t, err)

	// Submit a bad claim with EthBlockHeight >= timeout

	bad := types.MsgBatchSendToRemoteClaim{
		SkywayNonce:      lastNonce + 1,
		EventNonce:       lastNonce + 15,
		EthBlockHeight:   badHeight,
		BatchNonce:       uint64(lastBatchNonce + 1),
		TokenContract:    tokenContract,
		ChainReferenceId: "test-chain",
		Orchestrator:     sender.String(),
		Metadata: vtypes.MsgMetadata{
			Creator: sender.String(),
			Signers: []string{sender.String()},
		},
	}
	context := sdktypes.UnwrapSDKContext(ctx)
	log.Info("Submitting bad eth claim from orchestrator 0", "sender", sender.String(), "val", val0.String())

	_, err = msgServer.BatchSendToRemoteClaim(context, &bad)
	require.Error(t, err)

	// Assert that there is no attestation since the above failed
	badHash, err := bad.ClaimHash()
	require.NoError(t, err)
	att := pk.GetAttestation(ctx, chainReferenceID, bad.GetSkywayNonce(), badHash)
	require.Nil(t, att)

	// Attest the actual batch, and assert the votes are correct
	for i, orch := range AccAddrs[1:] {
		log.Info("Submitting good eth claim from orchestrators", "orch", orch.String())
		good := types.MsgBatchSendToRemoteClaim{
			SkywayNonce:      lastNonce + 1,
			EventNonce:       lastNonce + 15,
			EthBlockHeight:   goodHeight,
			BatchNonce:       uint64(lastBatchNonce + 1),
			TokenContract:    tokenContract,
			ChainReferenceId: "test-chain",
			Orchestrator:     orch.String(),
			Metadata: vtypes.MsgMetadata{
				Creator: orch.String(),
				Signers: []string{orch.String()},
			},
		}
		_, err := msgServer.BatchSendToRemoteClaim(context, &good)
		require.NoError(t, err)

		goodHash, err := good.ClaimHash()
		require.NoError(t, err)

		att := pk.GetAttestation(ctx, chainReferenceID, good.GetSkywayNonce(), goodHash)
		require.NotNil(t, att)
		log.Info("Asserting that the bad attestation only has one claimer", "attVotes", att.Votes)
		require.Equal(t, len(att.Votes), i+1) // Only these good orchestrators votes should be counted
	}
}
