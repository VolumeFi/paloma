package keeper

import (
	"errors"
	"math/big"
	"os"
	"sync"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/VolumeFi/whoops"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethcoretypes "github.com/ethereum/go-ethereum/core/types"
	g "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/palomachain/paloma/util/slice"
	consensusmocks "github.com/palomachain/paloma/x/consensus/keeper/consensus/mocks"
	consensustypes "github.com/palomachain/paloma/x/consensus/types"
	"github.com/palomachain/paloma/x/evm/types"
	evmmocks "github.com/palomachain/paloma/x/evm/types/mocks"
	metrixtypes "github.com/palomachain/paloma/x/metrix/types"
	valsettypes "github.com/palomachain/paloma/x/valset/types"
	"github.com/stretchr/testify/mock"
)

var (
	contractAbi         = string(whoops.Must(os.ReadFile("testdata/sample-abi.json")))
	contractBytecodeStr = string(whoops.Must(os.ReadFile("testdata/sample-bytecode.out")))

	sampleTx1RawBytes = common.FromHex(string(whoops.Must(os.ReadFile("testdata/sample-tx-raw.hex"))))

	sampleTx1 = func() *ethcoretypes.Transaction {
		tx := new(ethcoretypes.Transaction)
		whoops.Assert(tx.UnmarshalBinary(sampleTx1RawBytes))
		return tx
	}()

	slcPayload = string(whoops.Must(os.ReadFile("testdata/slc-payload.hex")))
	slcTx1     = ethcoretypes.NewTx(&ethcoretypes.DynamicFeeTx{
		Data: common.FromHex(string(whoops.Must(os.ReadFile("testdata/slc-tx-data.hex")))),
	})

	valsetTx1 = ethcoretypes.NewTx(&ethcoretypes.DynamicFeeTx{
		Data: common.FromHex(string(whoops.Must(os.ReadFile("testdata/valset-tx-data.hex")))),
	})
)

type record struct {
	denom string
	erc20 string
	chain string
}

func (r record) GetDenom() string            { return r.denom }
func (r record) GetErc20() string            { return r.erc20 }
func (r record) GetChainReferenceId() string { return r.chain }

func TestKeeperGinkgo(t *testing.T) {
	RegisterFailHandler(g.Fail)
	g.RunSpecs(t, "Metrix")
}

func createSnapshot(chain *types.AddChainProposal) *valsettypes.Snapshot {
	type valpower struct {
		valAddr       sdk.ValAddress
		power         int64
		externalChain []*valsettypes.ExternalChainInfo
	}

	totalPower := int64(20)
	valpowers := []valpower{
		{
			valAddr: sdk.ValAddress("validator-1"),
			power:   15,
			externalChain: []*valsettypes.ExternalChainInfo{
				{
					ChainType:        "evm",
					ChainReferenceID: chain.GetChainReferenceID(),
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
					ChainReferenceID: chain.GetChainReferenceID(),
					Address:          "addr2",
					Pubkey:           []byte("1"),
				},
			},
		},
	}

	return &valsettypes.Snapshot{
		Validators: slice.Map(valpowers, func(p valpower) valsettypes.Validator {
			return valsettypes.Validator{
				ShareCount:         sdkmath.NewInt(p.power),
				Address:            p.valAddr,
				ExternalChainInfos: p.externalChain,
			}
		}),
		TotalShares: sdkmath.NewInt(totalPower),
	}
}

var _ = g.Describe("attest router", func() {
	var k Keeper
	var ctx sdk.Context
	var q *consensusmocks.Queuer
	var v *evmmocks.ValsetKeeper
	var gk *evmmocks.SkywayKeeper
	var consensukeeper *evmmocks.ConsensusKeeper
	var mk *evmmocks.MetrixKeeper
	var tk *evmmocks.TreasuryKeeper
	var msg *consensustypes.QueuedSignedMessage
	var consensusMsg *types.Message
	var evidence []*consensustypes.Evidence
	var isGoodcase bool
	var isTxProcessed bool
	var execTx *ethcoretypes.Transaction
	newChain := &types.AddChainProposal{
		ChainReferenceID:  "eth-main",
		Title:             "bla",
		Description:       "bla",
		BlockHeight:       uint64(123),
		BlockHashAtHeight: "0x1234",
	}

	type valpower struct {
		valAddr       sdk.ValAddress
		power         int64
		externalChain []*valsettypes.ExternalChainInfo
	}
	var valpowers []valpower

	var totalPower int64

	g.BeforeEach(func() {
		t := g.GinkgoT()
		kpr, ms, _ctx := NewEvmKeeper(t)
		ctx = _ctx
		k = *kpr
		v = ms.ValsetKeeper
		gk = ms.SkywayKeeper
		tk = ms.TreasuryKeeper
		mk = ms.MetrixKeeper
		consensukeeper = ms.ConsensusKeeper
		q = consensusmocks.NewQueuer(t)
		isGoodcase = true
		isTxProcessed = true
		ms.SkywayKeeper.On("GetLastObservedSkywayNonce", mock.Anything, mock.Anything).
			Return(uint64(100), nil).Maybe()
		ms.ValsetKeeper.On("GetLatestSnapshotOnChain", mock.Anything, mock.Anything).
			Return(createSnapshot(newChain), nil).Maybe()
		ms.ValsetKeeper.On("FindSnapshotByID", mock.Anything, mock.Anything).
			Return(createSnapshot(newChain), nil).Maybe()
	})

	g.BeforeEach(func() {
		consensusMsg = &types.Message{}
	})

	g.JustBeforeEach(func() {
		sig := make([]byte, 100)

		msg = &consensustypes.QueuedSignedMessage{
			Id:       123,
			Msg:      whoops.Must(codectypes.NewAnyWithValue(consensusMsg)),
			Evidence: evidence,
			SignData: []*consensustypes.SignData{{
				ExternalAccountAddress: "addr1",
				Signature:              sig,
			}, {
				ExternalAccountAddress: "addr2",
				Signature:              sig,
			}},
			PublicAccessData: &consensustypes.PublicAccessData{
				ValsetID: 1,
			},
		}
	})

	var subject func() error
	var subjectOnce sync.Once
	var subjectErr error
	g.BeforeEach(func() {
		subjectOnce = sync.Once{}
		subject = func() error {
			subjectOnce.Do(func() {
				subjectErr = k.attestRouter(ctx, q, msg)
			})
			return subjectErr
		}
	})

	g.When("snapshot returns an error", func() {
		retErr := whoops.String("random error")
		g.BeforeEach(func() {
			v.On("GetCurrentSnapshot", mock.Anything).Return(
				nil,
				retErr,
			)
		})
		g.BeforeEach(func() {
			evidence = []*consensustypes.Evidence{
				{
					Proof: whoops.Must(codectypes.NewAnyWithValue(&types.SmartContractExecutionErrorProof{ErrorMessage: "doesn't matter"})),
				},
				{
					Proof: whoops.Must(codectypes.NewAnyWithValue(&types.SmartContractExecutionErrorProof{ErrorMessage: "just need something to exist"})),
				},
			}
		})

		g.It("returns error back", func() {
			Expect(subject()).To(MatchError(retErr))
		})
	})

	g.When("snapshot returns an actual snapshot", func() {
		g.JustBeforeEach(func() {
			v.On("GetCurrentSnapshot", mock.Anything).Return(
				&valsettypes.Snapshot{
					Validators: slice.Map(valpowers, func(p valpower) valsettypes.Validator {
						return valsettypes.Validator{
							ShareCount:         sdkmath.NewInt(p.power),
							Address:            p.valAddr,
							ExternalChainInfos: p.externalChain,
						}
					}),
					TotalShares: sdkmath.NewInt(totalPower),
				},
				nil,
			)
		})

		g.When("there is not enough power to reach a consensus", func() {
			g.BeforeEach(func() {
				totalPower = 20
				valpowers = []valpower{
					{
						valAddr: sdk.ValAddress("validator-1"),
						power:   5,
					},
					{
						valAddr: sdk.ValAddress("validator-2"),
						power:   5,
					},
				}
			})

			g.It("returns nil for error", func() {
				Expect(subject()).To(BeNil())
			})
		})

		g.When("there is enough power to reach a consensus", func() {
			setupChainSupport := func() {
				consensukeeper.On("PutMessageInQueue", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(uint64(10), nil)
				mk.On("Validators", mock.Anything, mock.Anything).Return(&metrixtypes.QueryValidatorsResponse{
					ValMetrics: getMetrics(3),
				}, nil)
				tk.On("GetRelayerFeesByChainReferenceID", mock.Anything, mock.Anything).Return(getFees(3), nil)

				err := k.AddSupportForNewChain(
					ctx,
					newChain.GetChainReferenceID(),
					newChain.GetChainID(),
					newChain.GetBlockHeight(),
					newChain.GetBlockHashAtHeight(),
					big.NewInt(55),
				)
				Expect(err).To(BeNil())

				sc, err := k.SaveNewSmartContract(ctx, contractAbi, common.FromHex(contractBytecodeStr))
				Expect(err).To(BeNil())

				err = k.SetAsCompassContract(ctx, sc)
				Expect(err).To(BeNil())

				dep, _ := k.getSmartContractDeploymentByContractID(ctx, sc.GetId(), newChain.GetChainReferenceID())
				Expect(dep).NotTo(BeNil())
			}

			g.BeforeEach(func() {
				totalPower = 20
				valpowers = []valpower{
					{
						valAddr: sdk.ValAddress("validator-1"),
						power:   10,
						externalChain: []*valsettypes.ExternalChainInfo{
							{
								ChainType:        "evm",
								ChainReferenceID: newChain.GetChainReferenceID(),
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
								ChainReferenceID: newChain.GetChainReferenceID(),
								Address:          "addr2",
								Pubkey:           []byte("2"),
							},
						},
					},
					{
						valAddr: sdk.ValAddress("validator-3"),
						power:   5,
						externalChain: []*valsettypes.ExternalChainInfo{
							{
								ChainType:        "evm",
								ChainReferenceID: newChain.GetChainReferenceID(),
								Address:          "addr3",
								Pubkey:           []byte("3"),
							},
						},
					},
				}
			})

			g.Context("with a valid evidence", func() {
				g.BeforeEach(func() {
					q.On("ChainInfo").Return("", "eth-main")
					q.On("Remove", mock.Anything, uint64(123)).Return(nil)
				})

				successfulProcess := func() {
					g.It("processes it successfully", g.Offset(1), func() {
						setupChainSupport()
						Expect(subject()).To(BeNil())
					})
				}

				g.JustBeforeEach(func() {
					Expect(k.isTxProcessed(ctx, execTx)).To(BeFalse())
				})

				g.When("message is SubmitLogicCall", func() {
					g.BeforeEach(func() {
						execTx = slcTx1
						proof := whoops.Must(codectypes.NewAnyWithValue(&types.TxExecutedProof{SerializedTX: whoops.Must(execTx.MarshalBinary())}))
						evidence = []*consensustypes.Evidence{
							{
								ValAddress: sdk.ValAddress("validator-1"),
								Proof:      proof,
							},
							{
								ValAddress: sdk.ValAddress("validator-2"),
								Proof:      proof,
							},
							{
								ValAddress: sdk.ValAddress("validator-3"),
								Proof:      proof,
							},
						}
						consensusMsg.Action = &types.Message_SubmitLogicCall{
							SubmitLogicCall: &types.SubmitLogicCall{
								HexContractAddress: "0x51eca2efb15afacc612278c71f5edb35986f172f",
								Abi:                []byte(contractAbi),
								Payload:            common.FromHex(slcPayload),
							},
						}
					})
					successfulProcess()

					g.When("message is attesting to successful erc20 relink", func() {
						g.When("no more pending messages after this", func() {
							g.It("remove the deployment and activate the chain", func() {
								setupChainSupport()
								dep, _ := k.getSmartContractDeploymentByContractID(ctx, uint64(1), newChain.GetChainReferenceID())
								Expect(dep).ToNot(BeNil())
								dep.Status = types.SmartContractDeployment_WAITING_FOR_ERC20_OWNERSHIP_TRANSFER
								dep.Erc20Transfers = []types.SmartContractDeployment_ERC20Transfer{
									{
										Denom:  "denom",
										Erc20:  "address",
										MsgID:  123,
										Status: types.SmartContractDeployment_ERC20Transfer_PENDING,
									},
								}
								err := k.updateSmartContractDeployment(ctx, uint64(1), newChain.ChainReferenceID, dep)
								Expect(err).To(BeNil())
								k.deploymentCache.Add(ctx, newChain.ChainReferenceID, uint64(1), 123)
								Expect(subject()).To(BeNil())
								res, _ := k.getSmartContractDeploymentByContractID(ctx, uint64(1), newChain.GetChainReferenceID())
								Expect(res).To(BeNil())
								info, err := k.GetChainInfo(ctx, newChain.ChainReferenceID)
								Expect(err).To(BeNil())
								Expect(info.ActiveSmartContractID).To(Equal(uint64(1)))
							})
						})
						g.When("more pending messages after this", func() {
							g.It("updates the deployment", func() {
								setupChainSupport()
								dep, _ := k.getSmartContractDeploymentByContractID(ctx, uint64(1), newChain.GetChainReferenceID())
								Expect(dep).ToNot(BeNil())
								dep.Status = types.SmartContractDeployment_WAITING_FOR_ERC20_OWNERSHIP_TRANSFER
								dep.Erc20Transfers = []types.SmartContractDeployment_ERC20Transfer{
									{
										Denom:  "denom",
										Erc20:  "address",
										MsgID:  123,
										Status: types.SmartContractDeployment_ERC20Transfer_PENDING,
									},
									{
										Denom:  "denom2",
										Erc20:  "address2",
										MsgID:  1234,
										Status: types.SmartContractDeployment_ERC20Transfer_PENDING,
									},
								}
								err := k.updateSmartContractDeployment(ctx, uint64(1), newChain.ChainReferenceID, dep)
								Expect(err).To(BeNil())
								k.deploymentCache.Add(ctx, newChain.ChainReferenceID, uint64(1), 123)
								Expect(subject()).To(BeNil())
								res, _ := k.getSmartContractDeploymentByContractID(ctx, uint64(1), newChain.GetChainReferenceID())
								Expect(res.Erc20Transfers).To(HaveLen(2))
								Expect(res.Erc20Transfers[0].Status).To(Equal(types.SmartContractDeployment_ERC20Transfer_OK))
								Expect(res.Erc20Transfers[1].Status).To(Equal(types.SmartContractDeployment_ERC20Transfer_PENDING))
								info, err := k.GetChainInfo(ctx, newChain.ChainReferenceID)
								Expect(err).To(BeNil())
								Expect(info.ActiveSmartContractID).To(Equal(uint64(0)))
							})
						})
					})

					g.Context("there is error proof", func() {
						g.BeforeEach(func() {
							// We're not expecting an error, but the tx won't be
							// processed either
							isTxProcessed = false
							proof, _ := codectypes.NewAnyWithValue(&types.SmartContractExecutionErrorProof{ErrorMessage: "doesn't matter"})
							evidence = []*consensustypes.Evidence{
								{
									ValAddress: sdk.ValAddress("validator-1"),
									Proof:      proof,
								},
								{
									ValAddress: sdk.ValAddress("validator-2"),
									Proof:      proof,
								},
							}
						})

						g.When("message has not been retried", func() {
							g.BeforeEach(func() {
								consensusMsg.Action = &types.Message_SubmitLogicCall{
									SubmitLogicCall: &types.SubmitLogicCall{
										Retries: uint32(0),
									},
								}
								consensukeeper.On("PutMessageInQueue", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(uint64(10), nil).Once()
							})

							g.It("should attempt to retry", func() {
								setupChainSupport()
								Expect(subject()).To(BeNil())
								consensukeeper.AssertNumberOfCalls(g.GinkgoT(), "PutMessageInQueue", 2)
							})
						})

						g.When("message has been retried too many times", func() {
							g.BeforeEach(func() {
								consensusMsg.Action = &types.Message_SubmitLogicCall{
									SubmitLogicCall: &types.SubmitLogicCall{
										Retries: uint32(2),
									},
								}
							})

							g.It("should not attempt to retry", func() {
								setupChainSupport()
								Expect(subject()).To(BeNil())
								consensukeeper.AssertNumberOfCalls(g.GinkgoT(), "PutMessageInQueue", 1)
							})
						})
					})
				})

				g.When("message is UpdateValset", func() {
					g.BeforeEach(func() {
						execTx = valsetTx1
						proof := whoops.Must(codectypes.NewAnyWithValue(&types.TxExecutedProof{SerializedTX: whoops.Must(execTx.MarshalBinary())}))
						evidence = []*consensustypes.Evidence{
							{
								ValAddress: sdk.ValAddress("validator-1"),
								Proof:      proof,
							},
							{
								ValAddress: sdk.ValAddress("validator-2"),
								Proof:      proof,
							},
							{
								ValAddress: sdk.ValAddress("validator-3"),
								Proof:      proof,
							},
						}
						consensusMsg.Action = &types.Message_UpdateValset{
							UpdateValset: &types.UpdateValset{
								Valset: &types.Valset{
									ValsetID:   1,
									Validators: []string{"addr1", "addr2"},
									Powers:     []uint64{15, 5},
								},
							},
						}
					})

					g.BeforeEach(func() {
						q.On("GetAll", mock.Anything).Return(nil, nil)
					})

					g.When("successfully sets valset for chain", func() {
						g.BeforeEach(func() {
							v.On("SetSnapshotOnChain", mock.Anything, mock.Anything, mock.Anything).Return(nil)
						})
						successfulProcess()
					})

					g.When("unsuccessfully sets valset for chain", func() {
						// We still process successfully even if we get an error here
						g.BeforeEach(func() {
							v.On("SetSnapshotOnChain", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("example error"))
						})
						successfulProcess()
					})
				})

				g.When("message is UploadSmartContract", func() {
					g.BeforeEach(func() {
						execTx = sampleTx1
						proof := whoops.Must(codectypes.NewAnyWithValue(&types.TxExecutedProof{SerializedTX: whoops.Must(execTx.MarshalBinary())}))
						evidence = []*consensustypes.Evidence{
							{
								ValAddress: sdk.ValAddress("validator-1"),
								Proof:      proof,
							},
							{
								ValAddress: sdk.ValAddress("validator-2"),
								Proof:      proof,
							},
							{
								ValAddress: sdk.ValAddress("validator-3"),
								Proof:      proof,
							},
						}
						consensusMsg.Action = &types.Message_UploadSmartContract{
							UploadSmartContract: &types.UploadSmartContract{
								Id:       1,
								Abi:      contractAbi,
								Bytecode: common.FromHex(contractBytecodeStr),
							},
						}
						address, err := sdk.ValAddressFromBech32("cosmosvaloper1pzf9apnk8yw7pjw3v9vtmxvn6guhkslanh8r07")
						Expect(err).To(BeNil())
						consensusMsg.Assignee = address.String()
					})

					g.When("target chain has no deployed ERC20 tokens", func() {
						g.BeforeEach(func() {
							gk.On("CastChainERC20ToDenoms", mock.Anything, mock.Anything).Return(nil, nil)
						})
						g.It("removes deployment", func() {
							setupChainSupport()
							Expect(subject()).To(BeNil())
							v, key := k.getSmartContractDeploymentByContractID(ctx, uint64(1), newChain.GetChainReferenceID())
							Expect(key).To(BeNil())
							Expect(v).To(BeNil())
						})
						g.It("sets chain as active", func() {
							setupChainSupport()
							Expect(subject()).To(BeNil())
							v, err := k.GetChainInfo(ctx, newChain.GetChainReferenceID())
							Expect(err).To(BeNil())
							Expect(v.GetActiveSmartContractID()).To(BeEquivalentTo(uint64(1)))
						})
					})

					g.When("target chain has active ERC20 tokens deployed", func() {
						g.BeforeEach(func() {
							gk.On("CastChainERC20ToDenoms", mock.Anything, newChain.ChainReferenceID).Return([]types.ERC20Record{
								record{"denom", "address1", newChain.ChainReferenceID},
								record{"denom2", "address2", newChain.ChainReferenceID},
							}, nil)
						})
						g.It("updates deployment", func() {
							setupChainSupport()
							Expect(subject()).To(BeNil())
							v, key := k.getSmartContractDeploymentByContractID(ctx, uint64(1), newChain.GetChainReferenceID())
							Expect(key).ToNot(BeNil())
							Expect(v).ToNot(BeNil())
							Expect(v.GetStatus()).To(BeEquivalentTo(types.SmartContractDeployment_WAITING_FOR_ERC20_OWNERSHIP_TRANSFER))
							Expect(v.GetErc20Transfers()).To(BeEquivalentTo([]types.SmartContractDeployment_ERC20Transfer{
								{
									Denom:  "denom",
									Erc20:  "address1",
									MsgID:  10,
									Status: types.SmartContractDeployment_ERC20Transfer_PENDING,
								},
								{
									Denom:  "denom2",
									Erc20:  "address2",
									MsgID:  10,
									Status: types.SmartContractDeployment_ERC20Transfer_PENDING,
								},
							}))
						})
						g.It("doesn't set the chain as active", func() {
							setupChainSupport()
							Expect(subject()).To(BeNil())
							v, err := k.GetChainInfo(ctx, newChain.GetChainReferenceID())
							Expect(err).To(BeNil())
							Expect(v.GetActiveSmartContractID()).To(BeEquivalentTo(uint64(0)))
						})
					})
				})

				g.JustAfterEach(func() {
					if isGoodcase {
						g.By("there is no error when processing evidence")
						Expect(subject()).To(BeNil())
					} else {
						g.By("there is an error when processing evidence")
						Expect(subject()).To(Not(BeNil()))
					}
				})

				g.JustAfterEach(func() {
					Expect(k.isTxProcessed(ctx, execTx)).To(Equal(isTxProcessed))
				})
			})
		})
	})
})
