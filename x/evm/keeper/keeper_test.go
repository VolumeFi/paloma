package keeper

import (
	"errors"
	"fmt"
	"math/big"
	"testing"
	"time"

	"cosmossdk.io/math"
	sdkmath "cosmossdk.io/math"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/palomachain/paloma/v2/util/slice"
	consensustypes "github.com/palomachain/paloma/v2/x/consensus/types"
	"github.com/palomachain/paloma/v2/x/evm/types"
	evmtypes "github.com/palomachain/paloma/v2/x/evm/types"
	"github.com/palomachain/paloma/v2/x/evm/types/mocks"
	metrixtypes "github.com/palomachain/paloma/v2/x/metrix/types"
	schedulertypes "github.com/palomachain/paloma/v2/x/scheduler/types"
	valsettypes "github.com/palomachain/paloma/v2/x/valset/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type validatorChainInfo struct {
	chainType        string
	chainReferenceID string
}

func getFees(num int) map[string]sdkmath.LegacyDec {
	fees := make(map[string]sdkmath.LegacyDec)
	for i := 0; i < num; i++ {
		fees[sdk.ValAddress(fmt.Sprintf("validator-%d", i)).String()] = sdkmath.LegacyMustNewDecFromStr("0.1")
	}
	return fees
}

func getMetrics(num int) []metrixtypes.ValidatorMetrics {
	metrics := make([]metrixtypes.ValidatorMetrics, num)
	for i := 0; i < num; i++ {
		metrics[i] = metrixtypes.ValidatorMetrics{
			ValAddress:    sdk.ValAddress(fmt.Sprintf("validator-%d", i)).String(),
			Uptime:        sdkmath.LegacyMustNewDecFromStr("1.0"),
			SuccessRate:   sdkmath.LegacyMustNewDecFromStr("0.9"),
			ExecutionTime: sdkmath.NewInt(3),
			FeatureSet:    sdkmath.LegacyMustNewDecFromStr("1.0"),
		}
	}
	return metrics
}

func getValidators(num int, chains []validatorChainInfo) []valsettypes.Validator {
	validators := make([]valsettypes.Validator, num)
	for i := 0; i < num; i++ {
		chainInfos := make([]*valsettypes.ExternalChainInfo, len(chains))
		for i, chain := range chains {
			chainInfos[i] = &valsettypes.ExternalChainInfo{
				ChainType:        chain.chainType,
				ChainReferenceID: chain.chainReferenceID,
				Address:          fmt.Sprintf("0x%02d", i),
			}
		}
		validators[i] = valsettypes.Validator{
			Address:            sdk.ValAddress(fmt.Sprintf("validator-%d", i)),
			State:              valsettypes.ValidatorState_ACTIVE,
			ShareCount:         sdkmath.NewInt(25000),
			ExternalChainInfos: chainInfos,
		}
	}
	return validators
}

func buildKeeper(t *testing.T) (*Keeper, sdk.Context, mockedServices) {
	k, mockServices, ctx := NewEvmKeeper(t)

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
	mockServices.ConsensusKeeper.On("PutMessageInQueue", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(uint64(0), nil)
	mockServices.MetrixKeeper.On("Validators", mock.Anything, mock.Anything).Return(&metrixtypes.QueryValidatorsResponse{
		ValMetrics: getMetrics(3),
	}, nil)
	mockServices.TreasuryKeeper.On("GetRelayerFeesByChainReferenceID", mock.Anything, mock.Anything).Return(getFees(3), nil)

	// invalid-test-chain mocks
	mockServices.ValsetKeeper.On("GetCurrentSnapshot", mock.Anything).Return(unpublishedSnapshot, nil)
	mockServices.ConsensusKeeper.On("PutMessageInQueue", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(uint64(0), nil)
	mockServices.MetrixKeeper.On("Validators", mock.Anything, mock.Anything).Return(&metrixtypes.QueryValidatorsResponse{
		ValMetrics: getMetrics(3),
	}, nil)
	mockServices.TreasuryKeeper.On("GetRelayerFeesByChainReferenceID", mock.Anything, mock.Anything).Return(getFees(3), nil)

	// Add 2 new chains for our tests to use
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

	err = k.SetSmartContractDeployer(ctx, "test-chain", cDummySmartContractDeployer)
	require.NoError(t, err)

	sc, err := k.SaveNewSmartContract(ctx, contractAbi, common.FromHex(contractBytecodeStr))
	require.NoError(t, err)
	err = k.SetAsCompassContract(ctx, sc)
	require.NoError(t, err)

	dep, _ := k.getSmartContractDeploymentByContractID(ctx, sc.GetId(), "test-chain")
	require.NotNil(t, dep)

	err = k.ActivateChainReferenceID(
		ctx,
		"test-chain",
		sc,
		"0x1234",
		dep.GetUniqueID(),
	)
	require.NoError(t, err)

	err = k.AddSupportForNewChain(
		ctx,
		"inactive-test-chain",
		2,
		uint64(123),
		"",
		big.NewInt(55),
	)
	require.NoError(t, err)

	err = k.SetFeeManagerAddress(ctx, "inactive-test-chain", cDummyFeeMgrAddress)
	require.NoError(t, err)

	err = k.SetSmartContractDeployer(ctx, "inactive-test-chain", cDummySmartContractDeployer)
	require.NoError(t, err)

	sc, err = k.SaveNewSmartContract(ctx, contractAbi, common.FromHex(contractBytecodeStr))
	require.NoError(t, err)
	err = k.SetAsCompassContract(ctx, sc)
	require.NoError(t, err)

	dep, _ = k.getSmartContractDeploymentByContractID(ctx, sc.GetId(), "test-chain")
	require.NotNil(t, dep)

	err = k.ActivateChainReferenceID(
		ctx,
		"test-chain",
		sc,
		"0x1234",
		dep.GetUniqueID(),
	)
	require.NoError(t, err)

	return k, ctx, mockServices
}

func TestKeeper_PreJobExecution(t *testing.T) {
	testcases := []struct {
		name             string
		chainReferenceID string
		setupMocks       func(sdk.Context, *Keeper)
		expectedError    error
	}{
		{
			name:             "publishes most recent valset",
			chainReferenceID: "test-chain",
			setupMocks: func(ctx sdk.Context, k *Keeper) {
				valsetKeeperMock := mocks.NewValsetKeeper(t)
				msgSenderMock := mocks.NewMsgSender(t)
				skywayKeeperMock := mocks.NewSkywayKeeper(t)

				unpublishedSnapshot := &valsettypes.Snapshot{
					Id:          1,
					TotalShares: sdkmath.NewInt(75000),
					Validators: []valsettypes.Validator{
						{
							Address:    sdk.ValAddress("validator-0"),
							State:      valsettypes.ValidatorState_ACTIVE,
							ShareCount: sdkmath.NewInt(25000),
							ExternalChainInfos: []*valsettypes.ExternalChainInfo{
								{
									ChainType:        "evm",
									ChainReferenceID: "test-chain",
									Address:          "0x00",
								},
							},
						},
						{
							Address:    sdk.ValAddress("validator-1"),
							State:      valsettypes.ValidatorState_ACTIVE,
							ShareCount: sdkmath.NewInt(25000),
							ExternalChainInfos: []*valsettypes.ExternalChainInfo{
								{
									ChainType:        "evm",
									ChainReferenceID: "test-chain",
									Address:          "0x01",
								},
							},
						},
						{
							Address:    sdk.ValAddress("validator-2"),
							State:      valsettypes.ValidatorState_ACTIVE,
							ShareCount: sdkmath.NewInt(25000),
							ExternalChainInfos: []*valsettypes.ExternalChainInfo{
								{
									ChainType:        "evm",
									ChainReferenceID: "test-chain",
									Address:          "0x02",
								},
							},
						},
					},
				}
				valsetKeeperMock.On("GetCurrentSnapshot", mock.Anything).Return(unpublishedSnapshot, nil)

				publishedSnapshot := &valsettypes.Snapshot{
					Id:     3,
					Chains: []string{"test-chain"},
				}
				valsetKeeperMock.On("GetLatestSnapshotOnChain", mock.Anything, mock.Anything).Return(publishedSnapshot, nil)

				msgSenderMock.On(
					"SendValsetMsgForChain",
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).Return(nil)

				k.Valset = valsetKeeperMock
				k.msgSender = msgSenderMock
				k.Skyway = skywayKeeperMock
			},
			expectedError: nil,
		},
		{
			name:             "no snapshot exists yet, return an error",
			chainReferenceID: "test-chain",
			setupMocks: func(ctx sdk.Context, k *Keeper) {
				valsetKeeperMock := mocks.NewValsetKeeper(t)
				valsetKeeperMock.On("GetCurrentSnapshot", mock.Anything).Return(nil, nil)
				k.Valset = valsetKeeperMock
			},
			expectedError: errors.New("nil, nil returned from Valset.GetCurrentSnapshot"),
		},
		{
			name:             "already using most recent published snapshot.  do nothing",
			chainReferenceID: "test-chain",
			setupMocks: func(ctx sdk.Context, k *Keeper) {
				valsetKeeperMock := mocks.NewValsetKeeper(t)
				skywayKeeperMock := mocks.NewSkywayKeeper(t)

				publishedSnapshot := &valsettypes.Snapshot{
					Id:          1,
					Chains:      []string{"test-chain"},
					TotalShares: sdkmath.NewInt(75000),
					Validators: []valsettypes.Validator{
						{
							State:      valsettypes.ValidatorState_ACTIVE,
							ShareCount: sdkmath.NewInt(25000),
							ExternalChainInfos: []*valsettypes.ExternalChainInfo{
								{
									ChainType:        "evm",
									ChainReferenceID: "test-chain",
								},
							},
						},
						{
							State:      valsettypes.ValidatorState_ACTIVE,
							ShareCount: sdkmath.NewInt(25000),
							ExternalChainInfos: []*valsettypes.ExternalChainInfo{
								{
									ChainType:        "evm",
									ChainReferenceID: "test-chain",
								},
							},
						},
						{
							State:      valsettypes.ValidatorState_ACTIVE,
							ShareCount: sdkmath.NewInt(25000),
							ExternalChainInfos: []*valsettypes.ExternalChainInfo{
								{
									ChainType:        "evm",
									ChainReferenceID: "test-chain",
								},
							},
						},
					},
				}
				valsetKeeperMock.On("GetCurrentSnapshot", mock.Anything).Return(publishedSnapshot, nil)

				valsetKeeperMock.On("GetLatestSnapshotOnChain", mock.Anything, mock.Anything).Return(publishedSnapshot, nil)

				k.Valset = valsetKeeperMock
				k.Skyway = skywayKeeperMock
				// Success is indicated by returning nil before calling msgSender.SendValsetMsgForChain
			},
			expectedError: nil,
		},
		{
			name:             "inactive chain.  do nothing",
			chainReferenceID: "inactive-test-chain",
			setupMocks: func(ctx sdk.Context, k *Keeper) {
				valsetKeeperMock := mocks.NewValsetKeeper(t)
				skywayKeeperMock := mocks.NewSkywayKeeper(t)

				unpublishedSnapshot := &valsettypes.Snapshot{
					Id:          1,
					TotalShares: sdkmath.NewInt(75000),
					Validators: []valsettypes.Validator{
						{
							State:      valsettypes.ValidatorState_ACTIVE,
							ShareCount: sdkmath.NewInt(25000),
							ExternalChainInfos: []*valsettypes.ExternalChainInfo{
								{
									ChainType:        "evm",
									ChainReferenceID: "inactive-test-chain",
								},
							},
						},
						{
							State:      valsettypes.ValidatorState_ACTIVE,
							ShareCount: sdkmath.NewInt(25000),
							ExternalChainInfos: []*valsettypes.ExternalChainInfo{
								{
									ChainType:        "evm",
									ChainReferenceID: "inactive-test-chain",
								},
							},
						},
						{
							State:      valsettypes.ValidatorState_ACTIVE,
							ShareCount: sdkmath.NewInt(25000),
							ExternalChainInfos: []*valsettypes.ExternalChainInfo{
								{
									ChainType:        "evm",
									ChainReferenceID: "inactive-test-chain",
								},
							},
						},
					},
				}
				valsetKeeperMock.On("GetCurrentSnapshot", mock.Anything).Return(unpublishedSnapshot, nil)

				publishedSnapshot := &valsettypes.Snapshot{
					Id:     3,
					Chains: []string{"inactive-test-chain"},
				}
				valsetKeeperMock.On("GetLatestSnapshotOnChain", mock.Anything, mock.Anything).Return(publishedSnapshot, nil)

				k.Valset = valsetKeeperMock
				k.Skyway = skywayKeeperMock
				// Success is indicated by returning nil before calling msgSender.SendValsetMsgForChain
			},
			expectedError: nil,
		},
	}

	asserter := assert.New(t)
	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			k, ctx, _ := buildKeeper(t)
			tt.setupMocks(ctx, k)
			job := &schedulertypes.Job{
				ID: "test_job_1",
				Routing: schedulertypes.Routing{
					ChainType:        "evm",
					ChainReferenceID: tt.chainReferenceID,
				},
				IsPayloadModifiable: false,
			}

			actualErr := k.PreJobExecution(ctx, job)

			// The real assertions we're making is that the mocks call the correct functions for the path
			asserter.Equal(tt.expectedError, actualErr)
		})
	}
}

func TestKeeper_MissingChains(t *testing.T) {
	testcases := []struct {
		name                   string
		inputChainReferenceIDs []string
		setup                  func(sdk.Context, *Keeper)
		expected               []string
		expectedError          error
	}{
		{
			name: "Returns a list of chains that are missing - 2 chain missing, inactive chain ignored",
			inputChainReferenceIDs: []string{
				"test-chain",
			},
			setup: func(ctx sdk.Context, k *Keeper) {
				for i, chainId := range []string{"test-chain-2", "test-chain-3"} {
					err := k.AddSupportForNewChain(
						ctx,
						chainId,
						uint64(i+3), // 2 chains already set up by keeper
						uint64(i+100),
						"",
						big.NewInt(55),
					)
					require.NoError(t, err)

					// Activate chain
					chainInfo, err := k.GetChainInfo(ctx, chainId)
					require.NoError(t, err)

					chainInfo.SmartContractAddr = "0x1234"

					err = k.updateChainInfo(ctx, chainInfo)
					require.NoError(t, err)
				}
			},
			expected: []string{
				"test-chain-2",
				"test-chain-3",
			},
		},
		{
			name: "Returns a list of chains that are missing - extra chain in input ignored",
			inputChainReferenceIDs: []string{
				"test-chain",
				"extra-chain",
				"test-chain-2",
			},
			setup: func(ctx sdk.Context, k *Keeper) {
				for i, chainId := range []string{"test-chain-2", "test-chain-3"} {
					err := k.AddSupportForNewChain(
						ctx,
						chainId,
						uint64(i+3), // 2 chains already set up by keeper
						uint64(i+100),
						"",
						big.NewInt(55),
					)
					require.NoError(t, err)

					// Activate chain
					chainInfo, err := k.GetChainInfo(ctx, chainId)
					require.NoError(t, err)

					chainInfo.SmartContractAddr = "0x1234"

					err = k.updateChainInfo(ctx, chainInfo)
					require.NoError(t, err)
				}
			},
			expected: []string{
				"test-chain-3",
			},
		},
		{
			name: "Returns a list of chains that are missing - nil slice when matching",
			inputChainReferenceIDs: []string{
				"test-chain",
				"test-chain-2",
				"test-chain-3",
			},
			setup: func(ctx sdk.Context, k *Keeper) {
				for i, chainId := range []string{"test-chain-2", "test-chain-3"} {
					err := k.AddSupportForNewChain(
						ctx,
						chainId,
						uint64(i+3), // 2 chains already set up by keeper
						uint64(i+100),
						"",
						big.NewInt(55),
					)
					require.NoError(t, err)

					// Activate chain
					chainInfo, err := k.GetChainInfo(ctx, chainId)
					require.NoError(t, err)

					chainInfo.SmartContractAddr = "0x1234"

					err = k.updateChainInfo(ctx, chainInfo)
					require.NoError(t, err)
				}
			},
			expected: []string(nil),
		},
	}

	asserter := assert.New(t)
	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			k, ctx, _ := buildKeeper(t)
			tt.setup(ctx, k)

			actual, actualErr := k.MissingChains(ctx, tt.inputChainReferenceIDs)
			asserter.Equal(tt.expected, actual)
			asserter.Equal(len(tt.expected), len(actual))
			asserter.Equal(tt.expectedError, actualErr)
		})
	}
}

func TestKeeper_PublishSnapshotToAllChains(t *testing.T) {
	testcases := []struct {
		name          string
		setup         func(sdk.Context, *Keeper, mockedServices)
		forcePublish  bool
		expectedError error
	}{
		{
			name: "Publishes when no previous snapshot on the chain",
			setup: func(ctx sdk.Context, k *Keeper, ms mockedServices) {
				ms.ValsetKeeper.On("GetLatestSnapshotOnChain", mock.Anything, mock.Anything).Return(nil, nil)
				// SendValsetMsgForChain indicates a publish
				ms.MsgSender.On("SendValsetMsgForChain", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
		},
		{
			name: "Doesn't publish when latest valset not over a month old and not forcePublish",
			setup: func(ctx sdk.Context, k *Keeper, ms mockedServices) {
				validators := getValidators(
					3,
					[]validatorChainInfo{
						{
							chainType:        "evm",
							chainReferenceID: "test-chain",
						},
						{
							chainType:        "evm",
							chainReferenceID: "test-chain",
						},
					},
				)
				publishedSnapshot := &valsettypes.Snapshot{
					Id:          1,
					Chains:      []string{"test-chain"},
					TotalShares: sdkmath.NewInt(75000),
					Validators:  validators,
					CreatedAt:   time.Now().Add(time.Duration(-28*24) * time.Hour), // 28 days ago
				}

				ms.ValsetKeeper.On("GetLatestSnapshotOnChain", mock.Anything, mock.Anything).Return(publishedSnapshot, nil)
				// Lack of a call to SendValsetMsgForChain indicates no publish
			},
		},
		{
			name: "Publishes regardless of age when force publish requested",
			setup: func(ctx sdk.Context, k *Keeper, ms mockedServices) {
				validators := getValidators(
					3,
					[]validatorChainInfo{
						{
							chainType:        "evm",
							chainReferenceID: "test-chain",
						},
						{
							chainType:        "evm",
							chainReferenceID: "test-chain",
						},
					},
				)
				publishedSnapshot := &valsettypes.Snapshot{
					Id:          1,
					Chains:      []string{"test-chain"},
					TotalShares: sdkmath.NewInt(75000),
					Validators:  validators,
					CreatedAt:   time.Now().Add(time.Duration(-28*24) * time.Hour), // 28 days ago
				}

				ms.ValsetKeeper.On("GetLatestSnapshotOnChain", mock.Anything, mock.Anything).Return(publishedSnapshot, nil)
				// SendValsetMsgForChain indicates a publish
				ms.MsgSender.On("SendValsetMsgForChain", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			forcePublish: true,
		},
	}

	asserter := assert.New(t)
	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			k, ctx, mockServices := buildKeeper(t)
			tt.setup(ctx, k, mockServices)

			ctx = ctx.WithBlockTime(time.Now())
			newSnapshot := &valsettypes.Snapshot{
				Id:          2,
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

			actualErr := k.PublishSnapshotToAllChains(ctx, newSnapshot, tt.forcePublish)
			asserter.Equal(tt.expectedError, actualErr)
		})
	}
}

func TestKeeper_SendValsetMsgForChain(t *testing.T) {
	k, ms, ctx := NewEvmKeeper(t)
	mSender := msgSender{
		ConsensusKeeper: k.ConsensusKeeper,
		cdc:             k.cdc,
	}
	valset := types.Valset{
		ValsetID:   2,
		Validators: []string{"addr1", "addr2"},
		Powers:     []uint64{15, 5},
	}
	chainInfo := &types.ChainInfo{
		ChainID:               100,
		ChainReferenceID:      "test-chain",
		ReferenceBlockHeight:  1000,
		ReferenceBlockHash:    "0x00",
		MinOnChainBalance:     "100",
		SmartContractUniqueID: []byte("abc"),
		RelayWeights: &types.RelayWeights{
			Fee:           "1.0",
			Uptime:        "1.0",
			SuccessRate:   "1.0",
			ExecutionTime: "1.0",
			FeatureSet:    "1.0",
		},
	}

	t.Run("Should do nothing if valset update is already scheduled", func(t *testing.T) {
		qMsg, _ := codectypes.NewAnyWithValue(&evmtypes.Message{
			TurnstoneID:      "abc",
			ChainReferenceID: "test-chain",
			Assignee:         "addr4",
			Action: &evmtypes.Message_UpdateValset{
				UpdateValset: &evmtypes.UpdateValset{
					Valset: &evmtypes.Valset{
						ValsetID: 2,
					},
				},
			},
		})

		msgs := []consensustypes.QueuedSignedMessageI{
			&consensustypes.QueuedSignedMessage{
				Id:  1,
				Msg: qMsg,
			},
		}

		ms.ConsensusKeeper.On("GetMessagesFromQueue", mock.Anything, mock.Anything, mock.Anything).
			Return(msgs, nil).
			Once()

		err := mSender.SendValsetMsgForChain(ctx, chainInfo, valset, "addr3", "0x03")
		assert.NoError(t, err)
	})

	t.Run("Should add valset update to queue", func(t *testing.T) {
		qMsg, _ := codectypes.NewAnyWithValue(&evmtypes.Message{
			TurnstoneID:      "abc",
			ChainReferenceID: "test-chain",
			Assignee:         "addr4",
			Action: &evmtypes.Message_UpdateValset{
				UpdateValset: &evmtypes.UpdateValset{
					Valset: &evmtypes.Valset{
						ValsetID: 1,
					},
				},
			},
		})

		msgs := []consensustypes.QueuedSignedMessageI{
			&consensustypes.QueuedSignedMessage{
				Id:  1,
				Msg: qMsg,
			},
		}

		ms.ConsensusKeeper.On("GetMessagesFromQueue", mock.Anything, mock.Anything, mock.Anything).
			Return(msgs, nil).
			Once()
		ms.ConsensusKeeper.On("DeleteJob", mock.Anything, mock.Anything, mock.Anything).
			Return(nil).
			Once()
		ms.ConsensusKeeper.On("PutMessageInQueue", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(uint64(1), nil).
			Once()

		err := mSender.SendValsetMsgForChain(ctx, chainInfo, valset, "addr3", "0x03")
		assert.NoError(t, err)
	})
}

func TestKeeper_AddJustInTimeValsetUpdates(t *testing.T) {
	chainInfo := &types.ChainInfo{
		ChainID:               100,
		ChainReferenceID:      "test-chain",
		ReferenceBlockHeight:  1000,
		ReferenceBlockHash:    "0x00",
		MinOnChainBalance:     "100",
		SmartContractUniqueID: []byte("abc"),
		SmartContractAddr:     "0x01",
		RelayWeights: &types.RelayWeights{
			Fee:           "1.0",
			Uptime:        "1.0",
			SuccessRate:   "1.0",
			ExecutionTime: "1.0",
			FeatureSet:    "1.0",
		},
	}

	type valpower struct {
		valAddr       sdk.ValAddress
		power         int64
		externalChain []*valsettypes.ExternalChainInfo
	}

	var totalPower int64 = 20
	valpowers := []valpower{
		{
			valAddr: sdk.ValAddress("validator-1"),
			power:   15,
			externalChain: []*valsettypes.ExternalChainInfo{
				{
					ChainType:        "evm",
					ChainReferenceID: chainInfo.GetChainReferenceID(),
					Address:          "addr1",
					Pubkey:           []byte("1"),
				},
			},
		},
		{
			valAddr: sdk.ValAddress("validator-2"),
			power:   5,
			externalChain: []*valsettypes.ExternalChainInfo{
				{
					ChainType:        "evm",
					ChainReferenceID: chainInfo.GetChainReferenceID(),
					Address:          "addr1",
					Pubkey:           []byte("1"),
				},
			},
		},
	}

	currentSnapshot := &valsettypes.Snapshot{
		Id: 5,
		Validators: slice.Map(valpowers, func(p valpower) valsettypes.Validator {
			return valsettypes.Validator{
				ShareCount:         sdkmath.NewInt(p.power),
				Address:            p.valAddr,
				ExternalChainInfos: p.externalChain,
			}
		}),
		TotalShares: sdkmath.NewInt(totalPower),
	}

	validatorMetrics := &metrixtypes.QueryValidatorsResponse{
		ValMetrics: []metrixtypes.ValidatorMetrics{
			{
				ValAddress:    sdk.ValAddress("validator-1").String(),
				Uptime:        sdkmath.LegacyOneDec(),
				SuccessRate:   sdkmath.LegacyOneDec(),
				ExecutionTime: sdkmath.NewInt(0),
				Fee:           sdkmath.NewInt(0),
				FeatureSet:    sdkmath.LegacyOneDec(),
			},
		},
	}

	fee, _ := sdkmath.LegacyNewDecFromStr("1.1")
	relayerFees := map[string]math.LegacyDec{
		sdk.ValAddress("validator-1").String(): fee,
	}

	t.Run("Should add valset update with SLC", func(t *testing.T) {
		k, ms, ctx := NewEvmKeeper(t)
		err := k.updateChainInfo(ctx, chainInfo)
		require.NoError(t, err)

		qMsg, _ := codectypes.NewAnyWithValue(&evmtypes.Message{
			TurnstoneID:      "abc",
			ChainReferenceID: "test-chain",
			Assignee:         "addr4",
			Action: &evmtypes.Message_SubmitLogicCall{
				SubmitLogicCall: &evmtypes.SubmitLogicCall{
					SenderAddress: sdk.ValAddress("sender"),
				},
			},
		})

		msgs := []consensustypes.QueuedSignedMessageI{
			&consensustypes.QueuedSignedMessage{
				Id:  1,
				Msg: qMsg,
			},
		}

		ms.ConsensusKeeper.On("GetMessagesFromQueue", mock.Anything,
			mock.Anything, mock.Anything).
			Return(msgs, nil).
			Once()

		ms.ValsetKeeper.On("GetCurrentSnapshot", mock.Anything).
			Return(currentSnapshot, nil)

		ms.ValsetKeeper.On("GetLatestSnapshotOnChain", mock.Anything, mock.Anything).
			Return(&valsettypes.Snapshot{Id: 1}, nil)

		ms.MetrixKeeper.On("Validators", mock.Anything, mock.Anything).
			Return(validatorMetrics, nil)

		ms.TreasuryKeeper.On("GetRelayerFeesByChainReferenceID", mock.Anything,
			chainInfo.ChainReferenceID).
			Return(relayerFees, nil)

		ms.MsgSender.On("SendValsetMsgForChain", mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything).
			Return(nil)

		k.AddJustInTimeValsetUpdates(ctx)
	})

	t.Run("Should do nothing if no SLC in queue", func(t *testing.T) {
		k, ms, ctx := NewEvmKeeper(t)
		err := k.updateChainInfo(ctx, chainInfo)
		require.NoError(t, err)

		qMsg, _ := codectypes.NewAnyWithValue(&evmtypes.Message{
			TurnstoneID:      "abc",
			ChainReferenceID: "test-chain",
			Assignee:         "addr4",
			Action: &evmtypes.Message_UploadSmartContract{
				UploadSmartContract: &evmtypes.UploadSmartContract{},
			},
		})

		msgs := []consensustypes.QueuedSignedMessageI{
			&consensustypes.QueuedSignedMessage{
				Id:  1,
				Msg: qMsg,
			},
		}

		ms.ConsensusKeeper.On("GetMessagesFromQueue", mock.Anything, mock.Anything, mock.Anything).
			Return(msgs, nil).
			Once()

		k.AddJustInTimeValsetUpdates(ctx)
	})

	t.Run("Should do nothing if valset update is already scheduled", func(t *testing.T) {
		k, ms, ctx := NewEvmKeeper(t)
		err := k.updateChainInfo(ctx, chainInfo)
		require.NoError(t, err)

		qMsg, _ := codectypes.NewAnyWithValue(&evmtypes.Message{
			TurnstoneID:      "abc",
			ChainReferenceID: "test-chain",
			Assignee:         "addr4",
			Action: &evmtypes.Message_UpdateValset{
				UpdateValset: &evmtypes.UpdateValset{
					Valset: &evmtypes.Valset{
						ValsetID: 2,
					},
				},
			},
		})

		slcMsg, _ := codectypes.NewAnyWithValue(&evmtypes.Message{
			TurnstoneID:      "abc",
			ChainReferenceID: "test-chain",
			Assignee:         "addr4",
			Action: &evmtypes.Message_SubmitLogicCall{
				SubmitLogicCall: &evmtypes.SubmitLogicCall{
					SenderAddress: sdk.ValAddress("sender"),
				},
			},
		})

		msgs := []consensustypes.QueuedSignedMessageI{
			&consensustypes.QueuedSignedMessage{
				Id:  1,
				Msg: qMsg,
			},
			&consensustypes.QueuedSignedMessage{
				Id:  2,
				Msg: slcMsg,
			},
		}

		ms.ConsensusKeeper.On("GetMessagesFromQueue", mock.Anything, mock.Anything, mock.Anything).
			Return(msgs, nil).
			Once()

		k.AddJustInTimeValsetUpdates(ctx)
	})

	t.Run("Should add valset update with UploadUserSmartContract", func(t *testing.T) {
		k, ms, ctx := NewEvmKeeper(t)
		err := k.updateChainInfo(ctx, chainInfo)
		require.NoError(t, err)

		qMsg, _ := codectypes.NewAnyWithValue(&evmtypes.Message{
			TurnstoneID:      "abc",
			ChainReferenceID: "test-chain",
			Assignee:         "addr4",
			Action: &evmtypes.Message_UploadUserSmartContract{
				UploadUserSmartContract: &evmtypes.UploadUserSmartContract{},
			},
		})

		msgs := []consensustypes.QueuedSignedMessageI{
			&consensustypes.QueuedSignedMessage{
				Id:  1,
				Msg: qMsg,
			},
		}

		ms.ConsensusKeeper.On("GetMessagesFromQueue", mock.Anything,
			mock.Anything, mock.Anything).
			Return(msgs, nil).
			Once()

		ms.ValsetKeeper.On("GetCurrentSnapshot", mock.Anything).
			Return(currentSnapshot, nil)

		ms.ValsetKeeper.On("GetLatestSnapshotOnChain", mock.Anything, mock.Anything).
			Return(&valsettypes.Snapshot{Id: 1}, nil)

		ms.MetrixKeeper.On("Validators", mock.Anything, mock.Anything).
			Return(validatorMetrics, nil)

		ms.TreasuryKeeper.On("GetRelayerFeesByChainReferenceID", mock.Anything,
			chainInfo.ChainReferenceID).
			Return(relayerFees, nil)

		ms.MsgSender.On("SendValsetMsgForChain", mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything).
			Return(nil)

		k.AddJustInTimeValsetUpdates(ctx)
	})
}
