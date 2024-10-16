package keeper

import (
	"context"
	"fmt"
	"testing"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/palomachain/paloma/v2/util/libmsg"
	"github.com/palomachain/paloma/v2/x/consensus/keeper/consensus"
	"github.com/palomachain/paloma/v2/x/consensus/types"
	evmtypes "github.com/palomachain/paloma/v2/x/evm/types"
	valsettypes "github.com/palomachain/paloma/v2/x/valset/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type treasyreKeeperMock struct{}

func (tk *treasyreKeeperMock) GetCombinedFeesForRelay(ctx context.Context, valAddress sdk.ValAddress, chainReferenceID string) (*types.MessageFeeSettings, error) {
	return &types.MessageFeeSettings{
		RelayerFee:   math.LegacyMustNewDecFromStr("1.25"),
		CommunityFee: math.LegacyMustNewDecFromStr("0.3"),
		SecurityFee:  math.LegacyMustNewDecFromStr("0.01"),
	}, nil
}

func Test_CheckAndProcessEstimatedMessages(t *testing.T) {
	k, ms, ctx := newConsensusKeeper(t)
	queue := types.Queue(defaultQueueName, chainType, chainReferenceID)

	k.feeProvider = &treasyreKeeperMock{}

	uvType := &evmtypes.Message{}

	types.RegisterInterfaces(types.ModuleCdc.InterfaceRegistry())
	evmtypes.RegisterInterfaces(types.ModuleCdc.InterfaceRegistry())

	types.ModuleCdc.InterfaceRegistry().RegisterImplementations((*types.ConsensusMsg)(nil), &evmtypes.Message{})
	types.ModuleCdc.InterfaceRegistry().RegisterImplementations((*evmtypes.TurnstoneMsg)(nil), &evmtypes.Message{})

	k.registry.Add(
		queueSupporter{
			opt: consensus.ApplyOpts(nil,
				consensus.WithQueueTypeName(queue),
				consensus.WithStaticTypeCheck(uvType),
				consensus.WithChainInfo(chainType, chainReferenceID),
				consensus.WithVerifySignature(func([]byte, []byte, []byte) bool {
					return true
				}),
			),
		},
	)

	validators := []valsettypes.Validator{}
	for i := range 10 {
		validators = append(validators, valsettypes.Validator{
			Address:    sdk.ValAddress(fmt.Sprintf("validator-%d", i)),
			ShareCount: math.NewInt(500),
		})
	}
	// Estimates for the first votes which won't reach consensus yet
	estimates := []uint64{
		21_000,
		19_000,
		11_000,
		25_000,
		21_000,
		25_000,
	}

	ms.ValsetKeeper.On("GetCurrentSnapshot", mock.Anything).Return(&valsettypes.Snapshot{
		Id:          1,
		Height:      500,
		Validators:  validators,
		TotalShares: math.NewInt(10 * 500),
		CreatedAt:   time.Time{},
		Chains:      []string{},
	}, nil)

	tt := []struct {
		name     string
		msg      *evmtypes.Message
		slcCheck func(*evmtypes.Message, *require.Assertions, bool) bool
	}{
		{
			name: "regular old message",
			msg: &evmtypes.Message{
				TurnstoneID:      "abc",
				ChainReferenceID: chainReferenceID,
				Assignee:         validators[0].Address.String(),
			},
			slcCheck: func(_ *evmtypes.Message, _ *require.Assertions, expected bool) bool {
				return expected
			},
		},
		{
			name: "SLC message",
			msg: &evmtypes.Message{
				TurnstoneID:      "abc",
				ChainReferenceID: chainReferenceID,
				Assignee:         validators[0].Address.String(),
				Action: &evmtypes.Message_SubmitLogicCall{
					SubmitLogicCall: &evmtypes.SubmitLogicCall{},
				},
			},
			slcCheck: func(m *evmtypes.Message, r *require.Assertions, expected bool) bool {
				slc := m.GetSubmitLogicCall()
				if slc == nil {
					return expected
				}
				if !expected {
					return slc.Fees == nil
				}

				r.NotNil(slc.Fees)
				r.Equal(uint64(26250), slc.Fees.RelayerFee, "relayer fee: got %d", slc.Fees.RelayerFee)
				r.Equal(uint64(7875), slc.Fees.CommunityFee, "community fee: got %d", slc.Fees.CommunityFee)
				r.Equal(uint64(263), slc.Fees.SecurityFee, "security fee: got %d", slc.Fees.SecurityFee)
				return true
			},
		},
		{
			name: "Upload user smart contract message",
			msg: &evmtypes.Message{
				TurnstoneID:      "abc",
				ChainReferenceID: chainReferenceID,
				Assignee:         validators[0].Address.String(),
				Action: &evmtypes.Message_UploadUserSmartContract{
					UploadUserSmartContract: &evmtypes.UploadUserSmartContract{},
				},
			},
			slcCheck: func(m *evmtypes.Message, r *require.Assertions, expected bool) bool {
				usc := m.GetUploadUserSmartContract()
				if usc == nil {
					return expected
				}
				if !expected {
					return usc.Fees == nil
				}

				r.NotNil(usc.Fees)
				r.Equal(uint64(26250), usc.Fees.RelayerFee, "relayer fee: got %d", usc.Fees.RelayerFee)
				r.Equal(uint64(7875), usc.Fees.CommunityFee, "community fee: got %d", usc.Fees.CommunityFee)
				r.Equal(uint64(263), usc.Fees.SecurityFee, "security fee: got %d", usc.Fees.SecurityFee)
				return true
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r := require.New(t)
			q, err := k.getConsensusQueue(ctx, queue)
			r.NoError(err)

			getMsg := func(mid uint64) *evmtypes.Message {
				m, err := q.GetMsgByID(ctx, mid)
				r.NoError(err)
				em, err := libmsg.ToEvmMessage(m, k.cdc)
				r.NoError(err)
				return em
			}

			// message with no need for estimation
			mid, err := k.PutMessageInQueue(ctx, queue, tc.msg, &consensus.PutOptions{RequireGasEstimation: true})
			r.NoError(err)
			r.NoError(k.CheckAndProcessEstimatedMessages(ctx))
			tc.slcCheck(getMsg(mid), r, false)

			// Start processing the first votes
			// Not enough validators will have vast their vote yet
			for i, estimate := range estimates {
				err = k.AddMessageGasEstimates(ctx, validators[i].Address, []*types.MsgAddMessageGasEstimates_GasEstimate{
					{
						MsgId:              mid,
						QueueTypeName:      queue,
						Value:              estimate,
						EstimatedByAddress: validators[i].Address.String(),
					},
				})
				r.NoError(err)
				r.NoError(k.CheckAndProcessEstimatedMessages(ctx))
				m, err := q.GetMsgByID(ctx, mid)
				r.NoError(err)
				r.Len(m.GetGasEstimates(), i+1)
				r.Equal(uint64(0), m.GetGasEstimate())
				tc.slcCheck(getMsg(mid), r, false)
			}

			// We now get a new estimate that will push us above the threshold
			// Paloma should build an expected gas cost from the median of all submitted estimates
			err = k.AddMessageGasEstimates(ctx, validators[6].Address, []*types.MsgAddMessageGasEstimates_GasEstimate{
				{
					MsgId:              mid,
					QueueTypeName:      queue,
					Value:              18_000,
					EstimatedByAddress: validators[6].Address.String(),
				},
			})
			r.NoError(err)
			r.NoError(k.CheckAndProcessEstimatedMessages(ctx))
			m, err := q.GetMsgByID(ctx, mid)
			r.NoError(err)
			r.Len(m.GetGasEstimates(), 7)
			r.Equal(uint64(21000), m.GetGasEstimate())
			tc.slcCheck(getMsg(mid), r, true)
		})
	}
}
