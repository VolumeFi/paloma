package keeper_test

import (
	"fmt"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"cosmossdk.io/core/header"
	"cosmossdk.io/math"
	"github.com/VolumeFi/whoops"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/palomachain/paloma/v2/tests/integration/helper"
	"github.com/palomachain/paloma/v2/testutil"
	"github.com/palomachain/paloma/v2/testutil/rand"
	"github.com/palomachain/paloma/v2/testutil/sample"
	utilkeeper "github.com/palomachain/paloma/v2/util/keeper"
	consensustypes "github.com/palomachain/paloma/v2/x/consensus/types"
	"github.com/palomachain/paloma/v2/x/evm/keeper"
	"github.com/palomachain/paloma/v2/x/evm/types"
	treasurytypes "github.com/palomachain/paloma/v2/x/treasury/types"
	valsettypes "github.com/palomachain/paloma/v2/x/valset/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const cDummyFeeMgrAddress = "0xb794f5ea0ba39494ce839613fffba74279579268"

var (
	contractAbi         = string(whoops.Must(os.ReadFile("testdata/sample-abi.json")))
	contractBytecodeStr = string(whoops.Must(os.ReadFile("testdata/sample-bytecode.out")))
)

func genValidators(numValidators, totalConsPower int) []stakingtypes.Validator {
	return testutil.GenValidators(numValidators, totalConsPower)
}

func TestEndToEndForEvmArbitraryCall(t *testing.T) {
	chainType, chainReferenceID := consensustypes.ChainTypeEVM, "eth-main"
	t1 := GinkgoT()
	f := helper.InitFixture(t1)
	ctx := f.Ctx.WithBlockHeight(5)

	newChain := &types.AddChainProposal{
		ChainReferenceID:  "eth-main",
		Title:             "bla",
		Description:       "bla",
		BlockHeight:       uint64(123),
		BlockHashAtHeight: "0x1234",
	}

	err := f.EvmKeeper.AddSupportForNewChain(
		ctx,
		newChain.GetChainReferenceID(),
		newChain.GetChainID(),
		newChain.GetBlockHeight(),
		newChain.GetBlockHashAtHeight(),
		big.NewInt(55),
	)
	require.NoError(t, err)
	err = f.EvmKeeper.SetFeeManagerAddress(ctx, newChain.GetChainReferenceID(), cDummyFeeMgrAddress)
	require.NoError(t, err)

	err = f.EvmKeeper.ActivateChainReferenceID(ctx, newChain.ChainReferenceID, &types.SmartContract{Id: 123}, "addr", []byte("abc"))
	require.NoError(t, err)

	validators := genValidators(25, 25000)
	for _, val := range validators {
		f.StakingKeeper.SetValidator(ctx, val)
	}

	for _, validator := range validators {
		valAddr, err := validator.GetConsAddr()
		require.NoError(t, err)
		pubKey, err := validator.ConsPubKey()
		require.NoError(t, err)
		operator, err := utilkeeper.ValAddressFromBech32(f.EvmKeeper.AddressCodec, validator.GetOperator())
		require.NoError(t, err)
		valAddress, err := f.EvmKeeper.AddressCodec.BytesToString(valAddr)
		require.NoError(t, err)
		err = f.ValsetKeeper.AddExternalChainInfo(ctx, operator, []*valsettypes.ExternalChainInfo{
			{
				ChainType:        "evm",
				ChainReferenceID: newChain.GetChainReferenceID(),
				Address:          valAddress,
				Pubkey:           pubKey.Bytes(),
			},
		})
		require.NoError(t, err)
		err = f.TreasuryKeeper.SetRelayerFee(ctx, sdk.ValAddress(operator), &treasurytypes.RelayerFeeSetting{
			ValAddress: sdk.ValAddress(operator).String(),
			Fees: []treasurytypes.RelayerFeeSetting_FeeSetting{
				{
					Multiplicator:    math.LegacyMustNewDecFromStr("1.1"),
					ChainReferenceId: chainReferenceID,
				},
				{
					Multiplicator:    math.LegacyMustNewDecFromStr("1.1"),
					ChainReferenceId: newChain.ChainReferenceID,
				},
			},
		})
		require.NoError(t, err)
	}

	_, err = f.ValsetKeeper.TriggerSnapshotBuild(ctx)
	require.NoError(t, err)
	f.MetrixKeeper.UpdateUptime(ctx)

	smartContractAddr := common.BytesToAddress(rand.Bytes(5))
	_, err = f.EvmKeeper.AddSmartContractExecutionToConsensus(
		ctx,
		chainReferenceID,
		"",
		&types.SubmitLogicCall{
			Payload: func() []byte {
				evm := whoops.Must(abi.JSON(strings.NewReader(sample.SimpleABI)))
				return whoops.Must(evm.Pack("store", big.NewInt(1337)))
			}(),
			HexContractAddress: smartContractAddr.Hex(),
			Abi:                []byte(sample.SimpleABI),
			Deadline:           1337,
		},
	)

	require.NoError(t, err)

	private, err := crypto.GenerateKey()
	require.NoError(t, err)

	accAddr := crypto.PubkeyToAddress(private.PublicKey)
	operator, err := utilkeeper.ValAddressFromBech32(f.EvmKeeper.AddressCodec, validators[0].GetOperator())
	require.NoError(t, err)
	err = f.ValsetKeeper.AddExternalChainInfo(ctx, operator, []*valsettypes.ExternalChainInfo{
		{
			ChainType:        chainType,
			ChainReferenceID: chainReferenceID,
			Address:          accAddr.Hex(),
			Pubkey:           accAddr[:],
		},
	})

	require.NoError(t, err)
	queue := consensustypes.Queue(types.ConsensusTurnstoneMessage, chainType, chainReferenceID)
	msgs, err := f.ConsensusKeeper.GetMessagesForSigning(ctx, queue, operator)
	require.NoError(t, err)

	for _, msg := range msgs {
		bytesToSign, err := msg.GetBytesToSign(f.Codec)
		require.NoError(t, err)

		sigbz, err := crypto.Sign(
			crypto.Keccak256(
				[]byte(keeper.SignaturePrefix),
				bytesToSign,
			),
			private,
		)
		require.NoError(t, err)
		err = f.ConsensusKeeper.AddMessageSignature(
			ctx,
			operator,
			[]*consensustypes.ConsensusMessageSignature{
				{
					Id:              msg.GetId(),
					QueueTypeName:   queue,
					Signature:       sigbz,
					SignedByAddress: accAddr.Hex(),
				},
			},
		)
		require.NoError(t, err)
	}
}

func TestFirstSnapshot_OnSnapshotBuilt(t *testing.T) {
	t1 := GinkgoT()
	f := helper.InitFixture(t1)
	ctx := f.Ctx.WithHeaderInfo(header.Info{
		Height: 5,
		Time:   time.Now(),
	})

	newChain := &types.AddChainProposal{
		ChainReferenceID:  "bob",
		Title:             "bla",
		Description:       "bla",
		BlockHeight:       uint64(123),
		BlockHashAtHeight: "0x1234",
	}
	err := f.EvmKeeper.AddSupportForNewChain(
		ctx,
		newChain.GetChainReferenceID(),
		newChain.GetChainID(),
		newChain.GetBlockHeight(),
		newChain.GetBlockHashAtHeight(),
		big.NewInt(55),
	)
	require.NoError(t, err)
	err = f.EvmKeeper.SetFeeManagerAddress(ctx, newChain.GetChainReferenceID(), cDummyFeeMgrAddress)
	require.NoError(t, err)
	err = f.EvmKeeper.ActivateChainReferenceID(
		ctx,
		newChain.ChainReferenceID,
		&types.SmartContract{
			Id: 123,
		},
		"addr",
		[]byte("abc"),
	)
	require.NoError(t, err)

	validators := genValidators(25, 25000)
	for _, val := range validators {
		err := f.StakingKeeper.SetValidator(ctx, val)
		require.NoError(t, err)
		operator, err := utilkeeper.ValAddressFromBech32(f.EvmKeeper.AddressCodec, val.GetOperator())
		require.NoError(t, err)
		err = f.ValsetKeeper.AddExternalChainInfo(ctx, operator, []*valsettypes.ExternalChainInfo{
			{
				ChainType:        "evm",
				ChainReferenceID: "bob",
				Address:          rand.ETHAddress().Hex(),
				Pubkey:           []byte("pk" + rand.ETHAddress().Hex()),
			},
		})
		require.NoError(t, err)
	}

	queue := fmt.Sprintf("evm/%s/%s", newChain.GetChainReferenceID(), types.ConsensusTurnstoneMessage)
	msgs, err := f.ConsensusKeeper.GetMessagesFromQueue(ctx, queue, 100)
	require.NoError(t, err)
	require.Empty(t, msgs)
	_, err = f.ConsensusKeeper.PutMessageInQueue(ctx, queue, &types.Message{
		TurnstoneID:      "abc",
		ChainReferenceID: "new-chain",
		Action: &types.Message_UpdateValset{
			UpdateValset: &types.UpdateValset{
				Valset: &types.Valset{
					ValsetID: 777,
				},
			},
		},
	}, nil)
	require.NoError(t, err)
	_, err = f.ValsetKeeper.TriggerSnapshotBuild(ctx)
	require.NoError(t, err)

	msgs, err = f.ConsensusKeeper.GetMessagesFromQueue(ctx, queue, 100)
	require.NoError(t, err)
	require.Len(t, msgs, 1)
}

func TestRecentPublishedSnapshot_OnSnapshotBuilt(t *testing.T) {
	t1 := GinkgoT()
	f := helper.InitFixture(t1)
	ctx := f.Ctx
	newChain := &types.AddChainProposal{
		ChainReferenceID:  "bob",
		Title:             "bla",
		Description:       "bla",
		BlockHeight:       uint64(123),
		BlockHashAtHeight: "0x1234",
	}
	err := f.EvmKeeper.AddSupportForNewChain(
		ctx,
		newChain.GetChainReferenceID(),
		newChain.GetChainID(),
		newChain.GetBlockHeight(),
		newChain.GetBlockHashAtHeight(),
		big.NewInt(55),
	)
	err = f.EvmKeeper.SetFeeManagerAddress(ctx, newChain.GetChainReferenceID(), cDummyFeeMgrAddress)
	require.NoError(t, err)
	require.NoError(t, err)
	err = f.EvmKeeper.ActivateChainReferenceID(
		ctx,
		newChain.ChainReferenceID,
		&types.SmartContract{
			Id: 123,
		},
		"addr",
		[]byte("abc"),
	)
	require.NoError(t, err)

	validators := genValidators(25, 25000)
	for _, val := range validators {
		f.StakingKeeper.SetValidator(ctx, val)
		operator, err := utilkeeper.ValAddressFromBech32(f.EvmKeeper.AddressCodec, val.GetOperator())
		require.NoError(t, err)
		err = f.ValsetKeeper.AddExternalChainInfo(ctx, operator, []*valsettypes.ExternalChainInfo{
			{
				ChainType:        "evm",
				ChainReferenceID: "bob",
				Address:          rand.ETHAddress().Hex(),
				Pubkey:           []byte("pk" + rand.ETHAddress().Hex()),
			},
		})
		require.NoError(t, err)
	}

	queue := fmt.Sprintf("evm/%s/%s", newChain.GetChainReferenceID(), types.ConsensusTurnstoneMessage)

	msgs, err := f.ConsensusKeeper.GetMessagesFromQueue(ctx, queue, 1)
	require.NoError(t, err)
	require.Empty(t, msgs)

	// Remove the listeners to set current state
	snapshotListeners := f.ValsetKeeper.SnapshotListeners
	f.ValsetKeeper.SnapshotListeners = []valsettypes.OnSnapshotBuiltListener{}

	_, err = f.ValsetKeeper.TriggerSnapshotBuild(ctx)
	require.NoError(t, err)

	msgs, err = f.ConsensusKeeper.GetMessagesFromQueue(ctx, queue, 100)
	require.NoError(t, err)
	require.Len(t, msgs, 0)

	latestSnapshot, err := f.ValsetKeeper.GetCurrentSnapshot(ctx)
	require.NoError(t, err)

	latestSnapshot.Chains = []string{"bob"}
	err = f.ValsetKeeper.SaveModifiedSnapshot(ctx, latestSnapshot)
	require.NoError(t, err)

	// Add two validators to make this new snapshot worthy
	validators = genValidators(2, 25000)
	for _, val := range validators {
		err := f.StakingKeeper.SetValidator(ctx, val)
		require.NoError(t, err)
		operator, err := utilkeeper.ValAddressFromBech32(f.EvmKeeper.AddressCodec, val.GetOperator())
		require.NoError(t, err)
		err = f.ValsetKeeper.AddExternalChainInfo(ctx, operator, []*valsettypes.ExternalChainInfo{
			{
				ChainType:        "evm",
				ChainReferenceID: "bob",
				Address:          rand.ETHAddress().Hex(),
				Pubkey:           []byte("pk" + rand.ETHAddress().Hex()),
			},
		})
		require.NoError(t, err)
	}

	// Add the listeners back on for the test
	f.ValsetKeeper.SnapshotListeners = snapshotListeners

	_, err = f.ValsetKeeper.TriggerSnapshotBuild(ctx)
	require.NoError(t, err)

	msgs, err = f.ConsensusKeeper.GetMessagesFromQueue(ctx, queue, 100)
	require.NoError(t, err)
	require.Len(t, msgs, 0) // We don't expect a message because there is already a recent snapshot for the chain
}

func TestOldPublishedSnapshot_OnSnapshotBuilt(t *testing.T) {
	var f *helper.Fixture
	var ctx sdk.Context
	t1 := GinkgoT()
	f = helper.InitFixture(t1)
	ctx = f.Ctx.WithHeaderInfo(header.Info{
		Height: 5,
		Time:   time.Now(),
	})
	newChain := &types.AddChainProposal{
		ChainReferenceID:  "bob",
		Title:             "bla",
		Description:       "bla",
		BlockHeight:       uint64(123),
		BlockHashAtHeight: "0x1234",
	}

	err := f.EvmKeeper.AddSupportForNewChain(
		ctx,
		newChain.GetChainReferenceID(),
		newChain.GetChainID(),
		newChain.GetBlockHeight(),
		newChain.GetBlockHashAtHeight(),
		big.NewInt(55),
	)
	require.NoError(t, err)
	err = f.EvmKeeper.SetFeeManagerAddress(ctx, newChain.GetChainReferenceID(), cDummyFeeMgrAddress)
	require.NoError(t, err)
	err = f.EvmKeeper.ActivateChainReferenceID(
		ctx,
		newChain.ChainReferenceID,
		&types.SmartContract{
			Id: 123,
		},
		"addr",
		[]byte("abc"),
	)
	require.NoError(t, err)

	validators := genValidators(25, 25000)
	for _, val := range validators {
		err := f.StakingKeeper.SetValidator(ctx, val)
		require.NoError(t, err)
		operator, err := utilkeeper.ValAddressFromBech32(f.EvmKeeper.AddressCodec, val.GetOperator())
		require.NoError(t, err)
		err = f.ValsetKeeper.AddExternalChainInfo(ctx, operator, []*valsettypes.ExternalChainInfo{
			{
				ChainType:        "evm",
				ChainReferenceID: "bob",
				Address:          rand.ETHAddress().Hex(),
				Pubkey:           []byte("pk" + rand.ETHAddress().Hex()),
			},
		})
		require.NoError(t, err)
	}

	queue := fmt.Sprintf("evm/%s/%s", newChain.GetChainReferenceID(), types.ConsensusTurnstoneMessage)

	msgs, err := f.ConsensusKeeper.GetMessagesFromQueue(ctx, queue, 1)
	require.NoError(t, err)
	require.Empty(t, msgs)

	// Remove the listeners to set current state
	snapshotListeners := f.ValsetKeeper.SnapshotListeners
	f.ValsetKeeper.SnapshotListeners = []valsettypes.OnSnapshotBuiltListener{}

	_, err = f.ValsetKeeper.TriggerSnapshotBuild(ctx)
	require.NoError(t, err)

	msgs, err = f.ConsensusKeeper.GetMessagesFromQueue(ctx, queue, 100)
	require.NoError(t, err)
	require.Len(t, msgs, 0)

	latestSnapshot, err := f.ValsetKeeper.GetCurrentSnapshot(ctx)
	require.NoError(t, err)

	// Age the latest snapshot by 30 days, 1 minute, set as active on chain
	latestSnapshot.Chains = []string{"bob"}

	latestSnapshot.CreatedAt = ctx.HeaderInfo().Time.Add(-((30 * 24 * time.Hour) + time.Minute))
	err = f.ValsetKeeper.SaveModifiedSnapshot(ctx, latestSnapshot)
	require.NoError(t, err)

	// Add two validators to make this new snapshot worthy
	validators = genValidators(2, 25000)
	for _, val := range validators {
		err := f.StakingKeeper.SetValidator(ctx, val)
		require.NoError(t, err)
		operator, err := utilkeeper.ValAddressFromBech32(f.EvmKeeper.AddressCodec, val.GetOperator())
		require.NoError(t, err)
		err = f.ValsetKeeper.AddExternalChainInfo(ctx, operator, []*valsettypes.ExternalChainInfo{
			{
				ChainType:        "evm",
				ChainReferenceID: "bob",
				Address:          rand.ETHAddress().Hex(),
				Pubkey:           []byte("pk" + rand.ETHAddress().Hex()),
			},
		})
		require.NoError(t, err)
	}

	// Add the listeners back on for the test
	f.ValsetKeeper.SnapshotListeners = snapshotListeners

	f.ConsensusKeeper.PutMessageInQueue(ctx, queue, &types.Message{
		TurnstoneID:      "abc",
		ChainReferenceID: "new-chain",
		Action: &types.Message_UpdateValset{
			UpdateValset: &types.UpdateValset{
				Valset: &types.Valset{
					ValsetID: 777,
				},
			},
		},
	}, nil)
	_, err = f.ValsetKeeper.TriggerSnapshotBuild(ctx)
	require.NoError(t, err)

	msgs, err = f.ConsensusKeeper.GetMessagesFromQueue(ctx, queue, 100)
	require.NoError(t, err)
	require.Len(t, msgs, 1) // We expect a new message because the previous one is a week old
}

func TestInactiveChain_OnSnapshotBuilt(t *testing.T) {
	t1 := GinkgoT()
	f := helper.InitFixture(t1)
	ctx := f.Ctx.WithBlockHeight(5)

	validators := genValidators(25, 25000)
	for _, val := range validators {
		f.StakingKeeper.SetValidator(ctx, val)
	}

	queue := fmt.Sprintf("evm/%s/%s", "bob", types.ConsensusTurnstoneMessage)

	_, err := f.ValsetKeeper.TriggerSnapshotBuild(ctx)
	require.NoError(t, err)

	_, err = f.ConsensusKeeper.GetMessagesFromQueue(ctx, queue, 100)
	require.Error(t, err) // We expect an error from this
}

func TestAddingSupportForNewChain(t *testing.T) {
	t1 := GinkgoT()
	f := helper.InitFixture(t1)
	ctx := f.Ctx.WithBlockHeight(5)

	t.Run("with happy path there are no errors", func(t *testing.T) {
		newChain := &types.AddChainProposal{
			ChainReferenceID:  "bob",
			Title:             "bla",
			Description:       "bla",
			BlockHeight:       uint64(123),
			BlockHashAtHeight: "0x1234",
		}
		err := f.EvmKeeper.AddSupportForNewChain(
			ctx,
			newChain.GetChainReferenceID(),
			newChain.GetChainID(),
			newChain.GetBlockHeight(),
			newChain.GetBlockHashAtHeight(),
			big.NewInt(55),
		)
		require.NoError(t, err)
		err = f.EvmKeeper.SetFeeManagerAddress(ctx, newChain.GetChainReferenceID(), cDummyFeeMgrAddress)
		require.NoError(t, err)

		gotChainInfo, err := f.EvmKeeper.GetChainInfo(ctx, newChain.GetChainReferenceID())
		require.NoError(t, err)

		require.Equal(t, newChain.GetChainReferenceID(), gotChainInfo.GetChainReferenceID())
		require.Equal(t, newChain.GetBlockHashAtHeight(), gotChainInfo.GetReferenceBlockHash())
		require.Equal(t, newChain.GetBlockHeight(), gotChainInfo.GetReferenceBlockHeight())
		t.Run("it returns an error if we try to add a chian whose chainID already exists", func(t *testing.T) {
			newChain.ChainReferenceID = "something_new"
			err := f.EvmKeeper.AddSupportForNewChain(
				ctx,
				newChain.GetChainReferenceID(),
				newChain.GetChainID(),
				newChain.GetBlockHeight(),
				newChain.GetBlockHashAtHeight(),
				big.NewInt(55),
			)
			require.ErrorIs(t, err, keeper.ErrCannotAddSupportForChainThatExists)
		})
	})

	t.Run("when chainReferenceID already exists then it returns an error", func(t *testing.T) {
		newChain := &types.AddChainProposal{
			ChainReferenceID:  "bob",
			Title:             "bla",
			Description:       "bla",
			BlockHeight:       uint64(123),
			BlockHashAtHeight: "0x1234",
		}
		err := f.EvmKeeper.AddSupportForNewChain(
			ctx,
			newChain.GetChainReferenceID(),
			newChain.GetChainID(),
			newChain.GetBlockHeight(),
			newChain.GetBlockHashAtHeight(),

			big.NewInt(55),
		)
		require.Error(t, err)
	})

	t.Run("activating chain", func(t *testing.T) {
		t.Run("if the chain does not exist it returns the error", func(t *testing.T) {
			err := f.EvmKeeper.ActivateChainReferenceID(ctx, "i don't exist", &types.SmartContract{}, "", []byte{})
			require.Error(t, err)
		})
		t.Run("works when chain exists", func(t *testing.T) {
			err := f.EvmKeeper.ActivateChainReferenceID(ctx, "bob", &types.SmartContract{Id: 123}, "addr", []byte("unique id"))
			require.NoError(t, err)
			gotChainInfo, err := f.EvmKeeper.GetChainInfo(ctx, "bob")
			require.NoError(t, err)

			require.Equal(t, "addr", gotChainInfo.GetSmartContractAddr())
			require.Equal(t, []byte("unique id"), gotChainInfo.GetSmartContractUniqueID())
		})
	})

	t.Run("removing chain", func(t *testing.T) {
		t.Run("if the chain does not exist it returns the error", func(t *testing.T) {
			err := f.EvmKeeper.RemoveSupportForChain(ctx, &types.RemoveChainProposal{
				ChainReferenceID: "i don't exist",
			})
			require.Error(t, err)
		})
		t.Run("works when chain exists", func(t *testing.T) {
			err := f.EvmKeeper.RemoveSupportForChain(ctx, &types.RemoveChainProposal{
				ChainReferenceID: "bob",
			})
			require.NoError(t, err)
			_, err = f.EvmKeeper.GetChainInfo(ctx, "bob")
			require.Error(t, keeper.ErrChainNotFound)
		})
	})
}

func TestKeeper_ValidatorSupportsAllChains(t *testing.T) {
	var f *helper.Fixture

	testcases := []struct {
		name     string
		setup    func(sdk.Context, *helper.Fixture) sdk.ValAddress
		expected bool
	}{
		{
			name: "returns true when all chains supported",
			setup: func(ctx sdk.Context, a *helper.Fixture) sdk.ValAddress {
				for i, chainId := range []string{"chain-1", "chain-2"} {
					newChain := &types.AddChainProposal{
						ChainReferenceID:  chainId,
						ChainID:           uint64(i),
						Title:             "bla",
						Description:       "bla",
						BlockHeight:       uint64(123),
						BlockHashAtHeight: "0x1234",
					}

					err := a.EvmKeeper.AddSupportForNewChain(
						ctx,
						newChain.GetChainReferenceID(),
						newChain.GetChainID(),
						newChain.GetBlockHeight(),
						newChain.GetBlockHashAtHeight(),
						big.NewInt(55),
					)
					require.NoError(t, err)
					err = f.EvmKeeper.SetFeeManagerAddress(ctx, newChain.GetChainReferenceID(), cDummyFeeMgrAddress)
					require.NoError(t, err)

					err = a.EvmKeeper.ActivateChainReferenceID(ctx, newChain.ChainReferenceID, &types.SmartContract{Id: 123}, fmt.Sprintf("addr%d", i), []byte("abc"))
					require.NoError(t, err)
				}

				validator := genValidators(1, 1000)[0]
				err := a.StakingKeeper.SetValidator(ctx, validator)
				require.NoError(t, err)
				private, err := crypto.GenerateKey()
				require.NoError(t, err)

				accAddr := crypto.PubkeyToAddress(private.PublicKey)

				// Add support for both chains created
				externalChains := make([]*valsettypes.ExternalChainInfo, 2)
				for i, chainId := range []string{"chain-1", "chain-2"} {
					externalChains[i] = &valsettypes.ExternalChainInfo{
						ChainType:        "evm",
						ChainReferenceID: chainId,
						Address:          accAddr.Hex(),
						Pubkey:           accAddr[:],
					}
				}
				operator, err := utilkeeper.ValAddressFromBech32(f.EvmKeeper.AddressCodec, validator.GetOperator())
				require.NoError(t, err)
				err = f.ValsetKeeper.AddExternalChainInfo(ctx, operator, externalChains)
				require.NoError(t, err)

				return operator
			},
			expected: true,
		},
		{
			name: "returns false when a chain is not supported",
			setup: func(ctx sdk.Context, a *helper.Fixture) sdk.ValAddress {
				for i, chainId := range []string{"chain-1", "chain-2"} {
					newChain := &types.AddChainProposal{
						ChainReferenceID:  chainId,
						ChainID:           uint64(i),
						Title:             "bla",
						Description:       "bla",
						BlockHeight:       uint64(123),
						BlockHashAtHeight: "0x1234",
					}

					err := a.EvmKeeper.AddSupportForNewChain(
						ctx,
						newChain.GetChainReferenceID(),
						newChain.GetChainID(),
						newChain.GetBlockHeight(),
						newChain.GetBlockHashAtHeight(),
						big.NewInt(55),
					)
					require.NoError(t, err)
					err = f.EvmKeeper.SetFeeManagerAddress(ctx, newChain.GetChainReferenceID(), cDummyFeeMgrAddress)
					require.NoError(t, err)

					err = a.EvmKeeper.ActivateChainReferenceID(ctx, newChain.ChainReferenceID, &types.SmartContract{Id: 123}, fmt.Sprintf("addr%d", i), []byte("abc"))
					require.NoError(t, err)
				}

				validator := genValidators(1, 1000)[0]
				a.StakingKeeper.SetValidator(ctx, validator)

				private, err := crypto.GenerateKey()
				require.NoError(t, err)

				accAddr := crypto.PubkeyToAddress(private.PublicKey)

				// Only add support for one of two chains created
				operator, err := utilkeeper.ValAddressFromBech32(f.EvmKeeper.AddressCodec, validator.GetOperator())
				require.NoError(t, err)
				err = f.ValsetKeeper.AddExternalChainInfo(
					ctx,
					operator,
					[]*valsettypes.ExternalChainInfo{
						{
							ChainType:        "evm",
							ChainReferenceID: "chain-1",
							Address:          accAddr.Hex(),
							Pubkey:           accAddr[:],
						},
					},
				)
				require.NoError(t, err)

				return operator
			},
			expected: false,
		},
	}

	asserter := assert.New(t)
	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			t1 := GinkgoT()
			f = helper.InitFixture(t1)
			ctx := f.Ctx.WithBlockHeight(5)

			validatorAddress := tt.setup(ctx, f)

			actual := f.ValsetKeeper.ValidatorSupportsAllChains(ctx, validatorAddress)
			asserter.Equal(tt.expected, actual)
		})
	}
}

func TestWithGinkgo(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "EVM keeper")
}

var _ = Describe("evm", func() {
	// smartContractAddr := common.BytesToAddress(rand.Bytes(5))
	// chainType, chainReferenceID := consensustypes.ChainTypeEVM, "eth-main"
	t := GinkgoT()
	var a *helper.Fixture
	var ctx sdk.Context
	var validators []stakingtypes.Validator
	newChain := &types.AddChainProposal{
		ChainReferenceID:  "eth-main",
		Title:             "bla",
		Description:       "bla",
		BlockHeight:       uint64(123),
		BlockHashAtHeight: "0x1234",
	}
	smartContract := &types.SmartContract{
		Id:       1,
		AbiJSON:  contractAbi,
		Bytecode: common.FromHex(contractBytecodeStr),
	}
	smartContract2 := &types.SmartContract{
		Id:       2,
		AbiJSON:  contractAbi,
		Bytecode: common.FromHex(contractBytecodeStr),
	}

	BeforeEach(func() {
		a = helper.InitFixture(t) // app.NewTestApp(GinkgoT(), false)
		ctx = a.Ctx.WithHeaderInfo(header.Info{Height: 5})
	})

	Context("multiple chains and smart contracts", func() {
		Describe("trying to add support for the same chain twice", func() {
			It("returns an error", func() {
				err := a.EvmKeeper.AddSupportForNewChain(
					ctx,
					newChain.GetChainReferenceID(),
					newChain.GetChainID(),
					newChain.GetBlockHeight(),
					newChain.GetBlockHashAtHeight(),
					big.NewInt(55),
				)
				Expect(err).To(BeNil())

				err = a.EvmKeeper.AddSupportForNewChain(
					ctx,
					newChain.GetChainReferenceID(),
					newChain.GetChainID(),
					newChain.GetBlockHeight(),
					newChain.GetBlockHashAtHeight(),
					big.NewInt(55),
				)
				Expect(err).To(MatchError(keeper.ErrCannotAddSupportForChainThatExists))
			})
		})

		Describe("ensuring that there can be two chains at the same time", func() {
			chain1 := &types.AddChainProposal{
				ChainReferenceID:  "chain1",
				Title:             "bla",
				Description:       "bla",
				BlockHeight:       uint64(456),
				BlockHashAtHeight: "0x1234",
				ChainID:           1,
			}
			chain2 := &types.AddChainProposal{
				ChainReferenceID:  "chain2",
				Title:             "bla",
				Description:       "bla",
				BlockHeight:       uint64(123),
				BlockHashAtHeight: "0x5678",
				ChainID:           2,
			}
			BeforeEach(func() {
				validators = genValidators(25, 25000)
				for _, val := range validators {
					a.StakingKeeper.SetValidator(ctx, val)
				}
			})

			JustBeforeEach(func() {
				for _, val := range validators {
					private1, err := crypto.GenerateKey()
					Expect(err).To(BeNil())
					private2, err := crypto.GenerateKey()
					Expect(err).To(BeNil())
					accAddr1 := crypto.PubkeyToAddress(private1.PublicKey)
					accAddr2 := crypto.PubkeyToAddress(private2.PublicKey)
					valAddr, err := utilkeeper.ValAddressFromBech32(a.ValsetKeeper.AddressCodec, val.GetOperator())
					Expect(err).To(BeNil())
					err = a.ValsetKeeper.AddExternalChainInfo(ctx, valAddr, []*valsettypes.ExternalChainInfo{
						{
							ChainType:        "evm",
							ChainReferenceID: chain1.ChainReferenceID,
							Address:          accAddr1.Hex(),
							Pubkey:           []byte("pub key 1" + accAddr1.Hex()),
						},
						{
							ChainType:        "evm",
							ChainReferenceID: chain2.ChainReferenceID,
							Address:          accAddr2.Hex(),
							Pubkey:           []byte("pub key 2" + accAddr2.Hex()),
						},
					})
					Expect(err).To(BeNil())
					err = a.TreasuryKeeper.SetRelayerFee(ctx, valAddr, &treasurytypes.RelayerFeeSetting{
						ValAddress: sdk.ValAddress(valAddr).String(),
						Fees: []treasurytypes.RelayerFeeSetting_FeeSetting{
							{
								Multiplicator:    math.LegacyMustNewDecFromStr("1.1"),
								ChainReferenceId: chain1.ChainReferenceID,
							},
							{
								Multiplicator:    math.LegacyMustNewDecFromStr("1.1"),
								ChainReferenceId: chain2.ChainReferenceID,
							},
						},
					})
					require.NoError(t, err)
				}
				_, err := a.ValsetKeeper.TriggerSnapshotBuild(ctx)
				a.MetrixKeeper.UpdateUptime(ctx)
				Expect(err).To(BeNil())
			})

			BeforeEach(func() {
				By("adding chain1 works")
				err := a.EvmKeeper.AddSupportForNewChain(
					ctx,
					chain1.GetChainReferenceID(),
					chain1.GetChainID(),
					chain1.GetBlockHeight(),
					chain1.GetBlockHashAtHeight(),
					big.NewInt(55),
				)
				Expect(err).To(BeNil())
				err = a.EvmKeeper.SetFeeManagerAddress(ctx, chain1.GetChainReferenceID(), cDummyFeeMgrAddress)
				require.NoError(t, err)

				By("adding chain2 works")
				err = a.EvmKeeper.AddSupportForNewChain(
					ctx,
					chain2.GetChainReferenceID(),
					chain2.GetChainID(),
					chain2.GetBlockHeight(),
					chain2.GetBlockHashAtHeight(),
					big.NewInt(55),
				)
				Expect(err).To(BeNil())
				err = a.EvmKeeper.SetFeeManagerAddress(ctx, chain2.GetChainReferenceID(), cDummyFeeMgrAddress)
				require.NoError(t, err)
			})

			Context("adding smart contract", func() {
				It("adds a new smart contract deployment", func() {
					By("simple assertion that two smart contracts share different ids", func() {
						Expect(smartContract.GetId()).NotTo(Equal(smartContract2.GetId()))
					})
					By("saving a new smart contract", func() {
						Expect(
							a.EvmKeeper.HasAnySmartContractDeployment(ctx, chain1.GetChainReferenceID()),
						).To(BeFalse())
						Expect(
							a.EvmKeeper.HasAnySmartContractDeployment(ctx, chain2.GetChainReferenceID()),
						).To(BeFalse())

						sc, err := a.EvmKeeper.SaveNewSmartContract(ctx, smartContract.GetAbiJSON(), smartContract.GetBytecode())
						Expect(err).To(BeNil())

						err = a.EvmKeeper.SetAsCompassContract(ctx, sc)
						Expect(err).To(BeNil())

						Expect(
							a.EvmKeeper.HasAnySmartContractDeployment(ctx, chain1.GetChainReferenceID()),
						).To(BeTrue())
						Expect(
							a.EvmKeeper.HasAnySmartContractDeployment(ctx, chain2.GetChainReferenceID()),
						).To(BeTrue())
					})

					By("removing a smart deployment for chain1 - it means that it was successfully uploaded", func() {
						a.EvmKeeper.DeleteSmartContractDeploymentByContractID(ctx, smartContract.GetId(), chain1.GetChainReferenceID())
						Expect(
							a.EvmKeeper.HasAnySmartContractDeployment(ctx, chain1.GetChainReferenceID()),
						).To(BeFalse())
						Expect(
							a.EvmKeeper.HasAnySmartContractDeployment(ctx, chain2.GetChainReferenceID()),
						).To(BeTrue())
					})

					By("activating a new smart contract it removes a deployment for chain1 but it doesn't for chain2", func() {
						err := a.EvmKeeper.ActivateChainReferenceID(ctx, chain1.GetChainReferenceID(), smartContract, "addr1", []byte("id1"))
						Expect(err).To(BeNil())
						Expect(
							a.EvmKeeper.HasAnySmartContractDeployment(ctx, chain1.GetChainReferenceID()),
						).To(BeFalse())
						Expect(
							a.EvmKeeper.HasAnySmartContractDeployment(ctx, chain2.GetChainReferenceID()),
						).To(BeTrue())

						By("verify that the chain's smart contract id has been deployed", func() {
							ci, err := a.EvmKeeper.GetChainInfo(ctx, chain1.GetChainReferenceID())
							Expect(err).To(BeNil())
							Expect(ci.GetActiveSmartContractID()).To(Equal(smartContract.GetId()))
						})
					})

					By("adding a new smart contract deployment deploys it to chain1 only", func() {
						sc, err := a.EvmKeeper.SaveNewSmartContract(ctx, smartContract2.GetAbiJSON(), smartContract2.GetBytecode())
						Expect(err).To(BeNil())
						err = a.EvmKeeper.SetAsCompassContract(ctx, sc)
						Expect(err).To(BeNil())
						Expect(
							a.EvmKeeper.HasAnySmartContractDeployment(ctx, chain1.GetChainReferenceID()),
						).To(BeTrue())
					})

					By("activating a new-new smart contract it deploys it to chain 1", func() {
						err := a.EvmKeeper.ActivateChainReferenceID(ctx, chain1.GetChainReferenceID(), smartContract2, "addr2", []byte("id2"))
						Expect(err).To(BeNil())
						Expect(
							a.EvmKeeper.HasAnySmartContractDeployment(ctx, chain2.GetChainReferenceID()),
						).To(BeTrue())
						By("verify that the chain's smart contract id has been deployed", func() {
							ci, err := a.EvmKeeper.GetChainInfo(ctx, chain1.GetChainReferenceID())
							Expect(err).To(BeNil())
							Expect(ci.GetActiveSmartContractID()).To(Equal(smartContract2.GetId()))
						})
					})
				})
			})
		})
	})

	Describe("on snapshot build", func() {
		var snapshot *valsettypes.Snapshot
		When("validator set is valid", func() {
			BeforeEach(func() {
				validators = genValidators(25, 25000)
				for _, val := range validators {
					a.StakingKeeper.SetValidator(ctx, val)
				}
			})

			When("evm chain and smart contract both exist", func() {
				BeforeEach(func() {
					for _, val := range validators {
						private, err := crypto.GenerateKey()
						Expect(err).To(BeNil())
						accAddr := crypto.PubkeyToAddress(private.PublicKey)
						valAddr, err := utilkeeper.ValAddressFromBech32(a.ValsetKeeper.AddressCodec, val.GetOperator())
						Expect(err).To(BeNil())
						err = a.ValsetKeeper.AddExternalChainInfo(ctx, valAddr, []*valsettypes.ExternalChainInfo{
							{
								ChainType:        "evm",
								ChainReferenceID: newChain.ChainReferenceID,
								Address:          accAddr.Hex(),
								Pubkey:           []byte("pub key" + accAddr.Hex()),
							},
							{
								ChainType:        "evm",
								ChainReferenceID: "new-chain",
								Address:          accAddr.Hex(),
								Pubkey:           []byte("pub key" + accAddr.Hex()),
							},
						})
						Expect(err).To(BeNil())
						err = a.TreasuryKeeper.SetRelayerFee(ctx, valAddr, &treasurytypes.RelayerFeeSetting{
							ValAddress: sdk.ValAddress(valAddr).String(),
							Fees: []treasurytypes.RelayerFeeSetting_FeeSetting{
								{
									Multiplicator:    math.LegacyMustNewDecFromStr("1.1"),
									ChainReferenceId: newChain.ChainReferenceID,
								},
								{
									Multiplicator:    math.LegacyMustNewDecFromStr("1.1"),
									ChainReferenceId: "new-chain",
								},
							},
						})
						require.NoError(t, err)
					}
					var err error
					snapshot, err = a.ValsetKeeper.TriggerSnapshotBuild(ctx)
					Expect(err).To(BeNil())
					a.MetrixKeeper.UpdateUptime(ctx)
				})

				BeforeEach(func() {
					err := a.EvmKeeper.AddSupportForNewChain(
						ctx,
						newChain.GetChainReferenceID(),
						newChain.GetChainID(),
						newChain.GetBlockHeight(),
						newChain.GetBlockHashAtHeight(),
						big.NewInt(55),
					)
					Expect(err).To(BeNil())
					err = a.EvmKeeper.SetFeeManagerAddress(ctx, newChain.GetChainReferenceID(), cDummyFeeMgrAddress)
					require.NoError(t, err)

					sc, err := a.EvmKeeper.SaveNewSmartContract(ctx, smartContract.GetAbiJSON(), smartContract.GetBytecode())
					Expect(err).To(BeNil())
					err = a.EvmKeeper.SetAsCompassContract(ctx, sc)
					Expect(err).To(BeNil())

					err = a.EvmKeeper.ActivateChainReferenceID(ctx, newChain.ChainReferenceID, smartContract, "addr", []byte("abc"))
					Expect(err).To(BeNil())

					By("it should have upload smart contract message", func() {
						msgs, err := a.ConsensusKeeper.GetMessagesFromQueue(ctx, "evm/eth-main/evm-turnstone-message", 5)

						Expect(err).To(BeNil())
						Expect(len(msgs)).To(Equal(1))

						con, err := msgs[0].ConsensusMsg(a.Codec)
						Expect(err).To(BeNil())

						evmMsg, ok := con.(*types.Message)
						Expect(ok).To(BeTrue())

						_, ok = evmMsg.GetAction().(*types.Message_UploadSmartContract)
						Expect(ok).To(BeTrue())

						a.ConsensusKeeper.DeleteJob(ctx, "evm/eth-main/evm-turnstone-message", msgs[0].GetId())
					})
				})

				It("expects update valset message to exist", func() {
					a.EvmKeeper.OnSnapshotBuilt(ctx, snapshot)
					msgs, err := a.ConsensusKeeper.GetMessagesFromQueue(ctx, "evm/eth-main/evm-turnstone-message", 5)

					Expect(err).To(BeNil())
					Expect(len(msgs)).To(Equal(1))

					con, err := msgs[0].ConsensusMsg(a.Codec)
					Expect(err).To(BeNil())

					evmMsg, ok := con.(*types.Message)
					Expect(ok).To(BeTrue())

					_, ok = evmMsg.GetAction().(*types.Message_UpdateValset)
					Expect(ok).To(BeTrue())
				})

				When("adding another chain which is not yet active", func() {
					BeforeEach(func() {
						err := a.EvmKeeper.AddSupportForNewChain(
							ctx,
							"new-chain",
							123,
							uint64(123),
							"0x1234",
							big.NewInt(55),
						)
						Expect(err).To(BeNil())
						err = a.EvmKeeper.SetFeeManagerAddress(ctx, "new-chain", cDummyFeeMgrAddress)
						require.NoError(t, err)
					})

					It("tries to deploy a smart contract to it", func() {
						a.EvmKeeper.OnSnapshotBuilt(ctx, snapshot)
						msgs, err := a.ConsensusKeeper.GetMessagesFromQueue(ctx, "evm/new-chain/evm-turnstone-message", 5)
						Expect(err).To(BeNil())
						Expect(len(msgs)).To(Equal(1))

						con, err := msgs[0].ConsensusMsg(a.Codec)
						Expect(err).To(BeNil())

						evmMsg, ok := con.(*types.Message)
						Expect(ok).To(BeTrue())

						_, ok = evmMsg.GetAction().(*types.Message_UploadSmartContract)
						Expect(ok).To(BeTrue())
					})
				})

				When("there is another upload valset already in", func() {
					BeforeEach(func() {
						err := a.EvmKeeper.AddSupportForNewChain(
							ctx,
							"new-chain",
							123,
							uint64(123),
							"0x1234",
							big.NewInt(55),
						)
						Expect(err).To(BeNil())
						err = a.EvmKeeper.SetFeeManagerAddress(ctx, "new-chain", cDummyFeeMgrAddress)
						require.NoError(t, err)
						err = a.EvmKeeper.ActivateChainReferenceID(ctx, "new-chain", &types.SmartContract{Id: 123}, "addr", []byte("abc"))
						Expect(err).To(BeNil())
						for _, val := range validators {
							private, err := crypto.GenerateKey()
							Expect(err).To(BeNil())
							accAddr := crypto.PubkeyToAddress(private.PublicKey)
							valAddr, err := utilkeeper.ValAddressFromBech32(a.ValsetKeeper.AddressCodec, val.GetOperator())
							Expect(err).To(BeNil())
							err = a.ValsetKeeper.AddExternalChainInfo(ctx, valAddr, []*valsettypes.ExternalChainInfo{
								{
									ChainType:        "evm",
									ChainReferenceID: "new-chain",
									Address:          accAddr.Hex(),
									Pubkey:           []byte("pub key" + accAddr.Hex()),
								},
							})
							Expect(err).To(BeNil())
							err = a.TreasuryKeeper.SetRelayerFee(ctx, valAddr, &treasurytypes.RelayerFeeSetting{
								ValAddress: sdk.ValAddress(valAddr).String(),
								Fees: []treasurytypes.RelayerFeeSetting_FeeSetting{
									{
										Multiplicator:    math.LegacyMustNewDecFromStr("1.1"),
										ChainReferenceId: newChain.ChainReferenceID,
									},
									{
										Multiplicator:    math.LegacyMustNewDecFromStr("1.1"),
										ChainReferenceId: "new-chain",
									},
								},
							})
							require.NoError(t, err)
						}
						a.MetrixKeeper.UpdateUptime(ctx)
					})
					BeforeEach(func() {
						msgs, err := a.ConsensusKeeper.GetMessagesFromQueue(ctx, "evm/new-chain/evm-turnstone-message", 5)
						Expect(err).To(BeNil())
						for _, msg := range msgs {
							// we are now clearing the deploy smart contract from the queue as we don't need it
							a.ConsensusKeeper.DeleteJob(ctx, "evm/new-chain/evm-turnstone-message", msg.GetId())
						}
						a.ConsensusKeeper.PutMessageInQueue(ctx, "evm/new-chain/evm-turnstone-message", &types.Message{
							TurnstoneID:      "abc",
							ChainReferenceID: "new-chain",
							Action: &types.Message_UpdateValset{
								UpdateValset: &types.UpdateValset{
									Valset: &types.Valset{
										ValsetID: 777,
									},
								},
							},
						}, nil)
					})
					It("deletes the old smart deployment", func() {
						a.EvmKeeper.OnSnapshotBuilt(ctx, snapshot)
						msgs, err := a.ConsensusKeeper.GetMessagesFromQueue(ctx, "evm/new-chain/evm-turnstone-message", 5)
						Expect(err).To(BeNil())
						Expect(len(msgs)).To(Equal(1))

						con, err := msgs[0].ConsensusMsg(a.Codec)
						Expect(err).To(BeNil())

						evmMsg, ok := con.(*types.Message)
						Expect(ok).To(BeTrue())

						vset, ok := evmMsg.GetAction().(*types.Message_UpdateValset)
						Expect(ok).To(BeTrue())
						Expect(vset.UpdateValset.GetValset().GetValsetID()).NotTo(Equal(uint64(777)))
						Expect(len(vset.UpdateValset.GetValset().GetValidators())).NotTo(BeZero())
					})
				})
			})
		})

		When("validator set is too tiny", func() {
			BeforeEach(func() {
				validators = genValidators(25, 25000)[:5]
				for _, val := range validators {
					a.StakingKeeper.SetValidator(ctx, val)
				}
				_, err := a.ValsetKeeper.TriggerSnapshotBuild(ctx)
				Expect(err).To(BeNil())
			})

			Context("evm chain and smart contract both exist", func() {
				BeforeEach(func() {
					err := a.EvmKeeper.AddSupportForNewChain(
						ctx,
						newChain.GetChainReferenceID(),
						newChain.GetChainID(),
						newChain.GetBlockHeight(),
						newChain.GetBlockHashAtHeight(),
						big.NewInt(55),
					)
					Expect(err).To(BeNil())
					err = a.EvmKeeper.SetFeeManagerAddress(ctx, newChain.ChainReferenceID, cDummyFeeMgrAddress)
					require.NoError(t, err)
					sc, err := a.EvmKeeper.SaveNewSmartContract(ctx, smartContract.GetAbiJSON(), smartContract.GetBytecode())
					Expect(err).To(BeNil())
					err = a.EvmKeeper.SetAsCompassContract(ctx, sc)
					Expect(err).To(BeNil())
				})

				It("doesn't put any message into a queue", func() {
					msgs, err := a.ConsensusKeeper.GetMessagesFromQueue(ctx, "evm/eth-main/evm-turnstone-message", 5)
					Expect(err).To(BeNil())
					Expect(msgs).To(BeZero())
				})
			})
		})
	})
})

var _ = Describe("change min on chain balance", func() {
	var a *helper.Fixture
	t := GinkgoT()
	var ctx sdk.Context
	newChain := &types.AddChainProposal{
		ChainReferenceID:  "eth-main",
		Title:             "bla",
		Description:       "bla",
		BlockHeight:       uint64(123),
		BlockHashAtHeight: "0x1234",
	}

	BeforeEach(func() {
		a = helper.InitFixture(t)
		ctx = a.Ctx.WithHeaderInfo(header.Info{Height: 5})
	})

	When("chain info exists", func() {
		BeforeEach(func() {
			err := a.EvmKeeper.AddSupportForNewChain(ctx, newChain.GetChainReferenceID(), newChain.GetChainID(), 1, "a", big.NewInt(55))
			Expect(err).To(BeNil())
			err = a.EvmKeeper.SetFeeManagerAddress(ctx, newChain.ChainReferenceID, cDummyFeeMgrAddress)
			require.NoError(t, err)
		})

		BeforeEach(func() {
			ci, err := a.EvmKeeper.GetChainInfo(ctx, newChain.GetChainReferenceID())
			Expect(err).To(BeNil())
			balance, err := ci.GetMinOnChainBalanceBigInt()
			Expect(err).To(BeNil())
			Expect(balance.Text(10)).To(Equal(big.NewInt(55).Text(10)))
		})

		It("changes the on chain balance", func() {
			err := a.EvmKeeper.ChangeMinOnChainBalance(ctx, newChain.GetChainReferenceID(), big.NewInt(888))
			Expect(err).To(BeNil())

			ci, err := a.EvmKeeper.GetChainInfo(ctx, newChain.GetChainReferenceID())
			Expect(err).To(BeNil())
			balance, err := ci.GetMinOnChainBalanceBigInt()
			Expect(err).To(BeNil())
			Expect(balance.Text(10)).To(Equal(big.NewInt(888).Text(10)))
		})
	})

	When("chain info does not exists", func() {
		It("returns an error", func() {
			err := a.EvmKeeper.ChangeMinOnChainBalance(ctx, newChain.GetChainReferenceID(), big.NewInt(888))
			Expect(err).To(MatchError(keeper.ErrChainNotFound))
		})
	})
})

var _ = Describe("change relay weights", func() {
	var a *helper.Fixture
	t := GinkgoT()
	var ctx sdk.Context
	newChain := &types.AddChainProposal{
		ChainReferenceID:  "eth-main",
		Title:             "bla",
		Description:       "bla",
		BlockHeight:       uint64(123),
		BlockHashAtHeight: "0x1234",
	}

	BeforeEach(func() {
		a = helper.InitFixture(t)
		ctx = a.Ctx.WithHeaderInfo(header.Info{Height: 5})
	})

	When("chain info exists", func() {
		BeforeEach(func() {
			err := a.EvmKeeper.AddSupportForNewChain(ctx, newChain.GetChainReferenceID(), newChain.GetChainID(), 1, "a", big.NewInt(55))
			Expect(err).To(BeNil())
			err = a.EvmKeeper.SetFeeManagerAddress(ctx, newChain.ChainReferenceID, cDummyFeeMgrAddress)
			require.NoError(t, err)
		})

		BeforeEach(func() {
			ci, err := a.EvmKeeper.GetChainInfo(ctx, newChain.GetChainReferenceID())
			Expect(err).To(BeNil())
			weights := ci.GetRelayWeights()
			Expect(weights).To(Equal(&types.RelayWeights{
				Fee:           "1.0",
				Uptime:        "1.0",
				SuccessRate:   "1.0",
				ExecutionTime: "1.0",
				FeatureSet:    "1.0",
			}))
		})

		It("changes the relay weights", func() {
			newWeights := &types.RelayWeights{
				Fee:           "0.12",
				Uptime:        "0.34",
				SuccessRate:   "0.56",
				ExecutionTime: "0.78",
				FeatureSet:    "0.99",
			}
			err := a.EvmKeeper.SetRelayWeights(ctx, newChain.GetChainReferenceID(), newWeights)
			Expect(err).To(BeNil())

			ci, err := a.EvmKeeper.GetChainInfo(ctx, newChain.GetChainReferenceID())
			Expect(err).To(BeNil())
			weights := ci.GetRelayWeights()
			Expect(weights).To(Equal(newWeights))
		})
	})

	When("chain info does not exists", func() {
		It("returns an error", func() {
			err := a.EvmKeeper.SetRelayWeights(ctx, newChain.GetChainReferenceID(), &types.RelayWeights{
				Fee:           "0.12",
				Uptime:        "0.34",
				SuccessRate:   "0.56",
				ExecutionTime: "0.78",
				FeatureSet:    "0.99",
			})
			Expect(err).To(MatchError(keeper.ErrChainNotFound))
		})
	})
})
