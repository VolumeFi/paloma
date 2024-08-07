package consensus

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
	keeperutil "github.com/palomachain/paloma/util/keeper"
	"github.com/palomachain/paloma/util/liblog"
	"github.com/palomachain/paloma/x/consensus/types"
)

var _ Queuer = Queue{}

type ConsensusMsg = types.ConsensusMsg

// Queue is a database storing messages that need to be signed.
type Queue struct {
	qo QueueOptions
}

type QueueOptions struct {
	Batched               bool
	QueueTypeName         string
	Sg                    keeperutil.StoreGetter
	Ider                  keeperutil.IDGenerator
	Cdc                   codec.Codec
	TypeCheck             types.TypeChecker
	BytesToSignCalculator types.BytesToSignFunc
	VerifySignature       types.VerifySignatureFunc
	ChainType             types.ChainType
	Attestator            types.Attestator
	ChainReferenceID      string
}

type OptFnc func(*QueueOptions)

func WithQueueTypeName(val string) OptFnc {
	return func(opt *QueueOptions) {
		opt.QueueTypeName = val
	}
}

func WithStaticTypeCheck(val any) OptFnc {
	return func(opt *QueueOptions) {
		opt.TypeCheck = types.StaticTypeChecker(val)
	}
}

func WithBytesToSignCalc(val types.BytesToSignFunc) OptFnc {
	return func(opt *QueueOptions) {
		opt.BytesToSignCalculator = val
	}
}

func WithVerifySignature(val types.VerifySignatureFunc) OptFnc {
	return func(opt *QueueOptions) {
		opt.VerifySignature = val
	}
}

func WithChainInfo(chainType, chainReferenceID string) OptFnc {
	return func(opt *QueueOptions) {
		opt.ChainReferenceID = chainReferenceID
		opt.ChainType = types.ChainType(chainType)
	}
}

func WithBatch(batch bool) OptFnc {
	return func(opt *QueueOptions) {
		opt.Batched = batch
	}
}

func WithAttestator(att types.Attestator) OptFnc {
	return func(opt *QueueOptions) {
		opt.Attestator = att
	}
}

func ApplyOpts(opts *QueueOptions, fncs ...OptFnc) *QueueOptions {
	if opts == nil {
		opts = &QueueOptions{}
	}
	for _, fnc := range fncs {
		fnc(opts)
	}
	return opts
}

func NewQueue(qo QueueOptions) (Queue, error) {
	if qo.TypeCheck == nil {
		return Queue{}, ErrNilTypeCheck
	}

	if qo.BytesToSignCalculator == nil {
		return Queue{}, ErrNilBytesToSignCalculator
	}

	if qo.VerifySignature == nil {
		return Queue{}, ErrNilVerifySignature
	}

	if len(qo.ChainType) == 0 {
		return Queue{}, ErrEmptyChainType
	}

	if len(qo.ChainReferenceID) == 0 {
		return Queue{}, ErrEmptyChainReferenceID
	}

	if len(qo.QueueTypeName) == 0 {
		return Queue{}, ErrEmptyQueueTypeName
	}

	return Queue{
		qo: qo,
	}, nil
}

// Put puts raw message into a signing queue.
func (c Queue) Put(ctx context.Context, msg ConsensusMsg, opts *PutOptions) (uint64, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	requireSignatures := true
	var publicAccessData *types.PublicAccessData

	if opts != nil {
		requireSignatures = opts.RequireSignatures
		publicAccessData = &types.PublicAccessData{
			ValAddress: nil,
			Data:       opts.PublicAccessData,
		}
	}

	if !c.qo.TypeCheck(msg) {
		return 0, ErrIncorrectMessageType.Format(msg)
	}
	newID := c.qo.Ider.IncrementNextID(sdkCtx, consensusQueueIDCounterKey)
	// just so it's clear that nonce is an actual ID
	nonce := newID
	anyMsg, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return 0, err
	}
	queuedMsg := &types.QueuedSignedMessage{
		Id:                 newID,
		Msg:                anyMsg,
		SignData:           []*types.SignData{},
		AddedAtBlockHeight: sdkCtx.BlockHeight(),
		AddedAt:            sdkCtx.BlockTime(),
		RequireSignatures:  requireSignatures,
		PublicAccessData:   publicAccessData,
		BytesToSign: c.qo.BytesToSignCalculator(msg, types.Salt{
			Nonce: nonce,
		}),
	}
	if err := c.save(sdkCtx, queuedMsg); err != nil {
		return 0, err
	}
	return newID, nil
}

// getAll returns all messages from a signing queu
func (c Queue) GetAll(ctx context.Context) ([]types.QueuedSignedMessageI, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	var msgs []types.QueuedSignedMessageI

	queue, err := c.queue(sdkCtx)
	if err != nil {
		return nil, err
	}

	iterator := queue.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		iterData := iterator.Value()

		var sm types.QueuedSignedMessageI
		if err := c.qo.Cdc.UnmarshalInterface(iterData, &sm); err != nil {
			return nil, err
		}

		msgs = append(msgs, sm)
	}

	return msgs, nil
}

func (c Queue) AddEvidence(ctx context.Context, msgID uint64, evidence *types.Evidence) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	msg, err := c.GetMsgByID(sdkCtx, msgID)
	if err != nil {
		return err
	}

	msg.AddEvidence(*evidence)

	return c.save(sdkCtx, msg)
}

type assignableMessage interface {
	proto.Message
	SetAssignee(ctx sdk.Context, val string)
}

func (c Queue) ReassignValidator(ctx sdk.Context, msgID uint64, val string) error {
	imsg, err := c.GetMsgByID(ctx, msgID)
	if err != nil {
		return err
	}

	var m types.ConsensusMsg
	if err := c.qo.Cdc.UnpackAny(imsg.GetMsg(), &m); err != nil {
		return err
	}
	assignable, ok := m.(assignableMessage)
	if !ok {
		return fmt.Errorf("message does not support setting assignee")
	}
	assignable.SetAssignee(ctx, val)
	var a proto.Message
	a = assignable
	msg := imsg.(*types.QueuedSignedMessage)
	anyMsg, err := codectypes.NewAnyWithValue(a)
	if err != nil {
		return err
	}

	msg.Msg = anyMsg

	return c.save(ctx, msg)
}

// SetPublicAccessData sets data that should be visible publicly so that other can provide proofs.
func (c Queue) SetPublicAccessData(ctx context.Context, msgID uint64, data *types.PublicAccessData) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	msg, err := c.GetMsgByID(sdkCtx, msgID)
	if err != nil {
		return err
	}

	if msg.GetPublicAccessData() != nil {
		return nil
	}

	msg.SetPublicAccessData(data)
	msg.SetHandledAtBlockHeight(math.NewInt(sdkCtx.BlockHeight()))

	return c.save(sdkCtx, msg)
}

func (c Queue) GetPublicAccessData(ctx context.Context, msgID uint64) (*types.PublicAccessData, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	msg, err := c.GetMsgByID(sdkCtx, msgID)
	if err != nil {
		return nil, err
	}

	return msg.GetPublicAccessData(), nil
}

func (c Queue) SetErrorData(ctx context.Context, msgID uint64, data *types.ErrorData) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	msg, err := c.GetMsgByID(sdkCtx, msgID)
	if err != nil {
		return err
	}

	if msg.GetErrorData() != nil || msg.GetPublicAccessData() != nil {
		return nil
	}

	msg.SetErrorData(data)
	msg.SetHandledAtBlockHeight(math.NewInt(sdkCtx.BlockHeight()))

	return c.save(sdkCtx, msg)
}

func (c Queue) GetErrorData(ctx context.Context, msgID uint64) (*types.ErrorData, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	msg, err := c.GetMsgByID(sdkCtx, msgID)
	if err != nil {
		return nil, err
	}

	return msg.GetErrorData(), nil
}

// AddSignature adds a signature to the message and checks if the signature is valid.
func (c Queue) AddSignature(ctx context.Context, msgID uint64, signData *types.SignData) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	msg, err := c.GetMsgByID(sdkCtx, msgID)
	if err != nil {
		return err
	}

	// check if the same public key already signed this message or of this
	// validator has already signed this message
	for _, existingSigData := range msg.GetSignData() {
		if bytes.Equal(existingSigData.PublicKey, signData.PublicKey) {
			return ErrAlreadySignedWithKey.Format(msgID, c.qo.QueueTypeName, existingSigData.PublicKey)
		}
		if signData.ValAddress.Equals(existingSigData.ValAddress) {
			return ErrValidatorAlreadySigned.Format(signData.ValAddress)
		}
	}

	if !c.qo.VerifySignature(msg.GetBytesToSign(), signData.Signature, signData.PublicKey) {
		return ErrInvalidSignature
	}

	msg.AddSignData(signData)

	return c.save(sdkCtx, msg)
}

// remove removes the message from the queue.
func (c Queue) Remove(ctx context.Context, msgID uint64) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	msg, err := c.GetMsgByID(sdkCtx, msgID)
	if err != nil {
		return err
	}

	queue, err := c.queue(sdkCtx)
	if err != nil {
		return err
	}

	queue.Delete(sdk.Uint64ToBigEndian(msgID))

	keeperutil.EmitEvent(keeperutil.ModuleNameFunc(types.ModuleName), sdkCtx, types.ItemRemovedEventKey,
		types.ItemRemovedEventID.With(fmt.Sprintf("%d", msgID)),
		types.ItemRemovedChainReferenceID.With(c.qo.ChainReferenceID),
	)

	logger := liblog.FromSDKLogger(sdkCtx.Logger()).WithFields("msg-id", msgID)

	msgWithSigs, err := ToMessageWithSignatures(msg, c.qo.Cdc)
	if err != nil {
		logger.WithError(err).Error("Failed to convert message with signatures")
		return nil
	}

	jsonMsg, err := c.qo.Cdc.MarshalJSON(&msgWithSigs)
	if err != nil {
		logger.WithError(err).Error("Failed to marshal message as json")
		return nil
	}

	logger.WithFields("msg", json.RawMessage(jsonMsg),
		"block_height", sdkCtx.BlockHeight()).
		Info("Removed message from queue")

	return nil
}

// getMsgByID given a message ID, it returns the message
func (c Queue) GetMsgByID(ctx context.Context, id uint64) (types.QueuedSignedMessageI, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	queue, err := c.queue(sdkCtx)
	if err != nil {
		return nil, err
	}

	data := queue.Get(sdk.Uint64ToBigEndian(id))

	if data == nil {
		return nil, ErrMessageDoesNotExist.Format(id)
	}

	var sm types.QueuedSignedMessageI
	if err := c.qo.Cdc.UnmarshalInterface(data, &sm); err != nil {
		return nil, err
	}

	return sm, nil
}

// save saves the message into the queue
func (c Queue) save(ctx context.Context, msg types.QueuedSignedMessageI) error {
	if msg.GetId() == 0 {
		return ErrUnableToSaveMessageWithoutID
	}
	data, err := c.qo.Cdc.MarshalInterface(msg)
	if err != nil {
		return err
	}

	q, err := c.queue(ctx)
	if err != nil {
		return err
	}

	q.Set(sdk.Uint64ToBigEndian(msg.GetId()), data)
	return nil
}

// queue is a simple helper function to return the queue store
func (c Queue) queue(ctx context.Context) (prefix.Store, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := c.qo.Sg.Store(sdkCtx)

	key, err := c.signingQueueKey()
	if err != nil {
		return prefix.Store{}, err
	}

	return prefix.NewStore(store, []byte(key)), nil
}

// signingQueueKey builds a key for the store where are we going to store those.
func (c Queue) signingQueueKey() (string, error) {
	if c.qo.QueueTypeName == "" {
		return "", ErrEmptyQueueTypeName
	}

	return fmt.Sprintf("%s-%s", consensusQueueSigningKey, c.qo.QueueTypeName), nil
}

func (c Queue) ChainInfo() (types.ChainType, string) {
	return c.qo.ChainType, c.qo.ChainReferenceID
}

func RemoveQueueCompletely(ctx context.Context, cq Queuer) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	var store storetypes.KVStore
	var err error

	switch typ := cq.(type) {
	case Queue:
		store, err = typ.queue(sdkCtx)
	case BatchQueue:
		store, err = typ.batchQueue(sdkCtx)
	default:
		err = errors.New("cq type is unknown!")
	}

	if err != nil {
		liblog.FromSDKLogger(sdkCtx.Logger()).WithError(err).Error("Error while removing the queue")
		return
	}

	iterator := store.Iterator(nil, nil)
	deleteKeys := [][]byte{}
	for ; iterator.Valid(); iterator.Valid() {
		deleteKeys = append(deleteKeys, iterator.Key())
	}

	for _, key := range deleteKeys {
		store.Delete(key)
	}
}

func ToMessageWithSignatures(msg types.QueuedSignedMessageI, cdc codec.BinaryCodec) (types.MessageWithSignatures, error) {
	origMsg, err := msg.ConsensusMsg(cdc)
	if err != nil {
		return types.MessageWithSignatures{}, err
	}
	anyMsg, err := codectypes.NewAnyWithValue(origMsg)
	if err != nil {
		return types.MessageWithSignatures{}, err
	}

	var publicAccessData []byte
	var valsetID uint64

	if msg.GetPublicAccessData() != nil {
		publicAccessData = msg.GetPublicAccessData().GetData()
		valsetID = msg.GetPublicAccessData().GetValsetID()
	}

	var errorData []byte

	if msg.GetErrorData() != nil {
		errorData = msg.GetErrorData().GetData()
	}

	approvedMessage := types.MessageWithSignatures{
		Nonce:            msg.Nonce(),
		Id:               msg.GetId(),
		Msg:              anyMsg,
		BytesToSign:      msg.GetBytesToSign(),
		SignData:         []*types.ValidatorSignature{},
		PublicAccessData: publicAccessData,
		ValsetID:         valsetID,
		ErrorData:        errorData,
	}
	for _, signData := range msg.GetSignData() {
		approvedMessage.SignData = append(approvedMessage.SignData, &types.ValidatorSignature{
			ValAddress:             signData.GetValAddress(),
			Signature:              signData.GetSignature(),
			ExternalAccountAddress: signData.GetExternalAccountAddress(),
			PublicKey:              signData.GetPublicKey(),
		})
	}

	return approvedMessage, nil
}
