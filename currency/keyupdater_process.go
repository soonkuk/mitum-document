package currency

import (
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/base/state"
	"github.com/spikeekips/mitum/util/valuehash"
	"golang.org/x/xerrors"
)

func (KeyUpdater) Process(
	func(key string) (state.State, bool, error),
	func(valuehash.Hash, ...state.State) error,
) error {
	return nil
}

type KeyUpdaterProcessor struct {
	cp *CurrencyPool
	KeyUpdater
	sa  state.State // target의 state account
	sb  AmountState // sender의 state balance
	fee Big         // operation
}

func NewKeyUpdaterProcessor(cp *CurrencyPool) GetNewProcessor {
	return func(op state.Processor) (state.Processor, error) {
		i, ok := op.(KeyUpdater)
		if !ok {
			return nil, xerrors.Errorf("not KeyUpdater, %T", op)
		}
		return &KeyUpdaterProcessor{
			cp:         cp,
			KeyUpdater: i,
		}, nil
	}
}

func (op *KeyUpdaterProcessor) PreProcess(
	getState func(string) (state.State, bool, error),
	_ func(valuehash.Hash, ...state.State) error,
) (state.Processor, error) {
	fact := op.Fact().(KeyUpdaterFact)

	// key에 관련된 state 가져오기
	st, err := existsState(StateKeyAccount(fact.target), "target keys", getState)
	if err != nil {
		return nil, err
	}
	op.sa = st

	// 업데이트하려는 키가 같은 키인지 확인
	if ks, err := StateKeysValue(op.sa); err != nil {
		return nil, operation.NewBaseReasonErrorFromError(err)
	} else if ks.Equal(fact.Keys()) {
		return nil, operation.NewBaseReasonError("same Keys with the existing")
	}

	// balance에 대한 state 가져오기 없으면 empty state 받아옴.
	st, err = existsState(StateKeyBalance(fact.target, fact.currency), "balance of target", getState)
	if err != nil {
		return nil, err
	}
	op.sb = NewAmountState(st, fact.currency)

	// fact sign 확인
	if err := checkFactSignsByState(fact.target, op.Signs(), getState); err != nil {
		return nil, operation.NewBaseReasonError("invalid signing: %w", err)
	}
	// feeere 확인
	feeer, found := op.cp.Feeer(fact.currency)
	if !found {
		return nil, operation.NewBaseReasonError("currency, %q not found of KeyUpdater", fact.currency)
	}
	// sender의 balance에서 fee를 감당할 정도로 충분한지 확인
	fee, err := feeer.Fee(ZeroBig)
	if err != nil {
		return nil, operation.NewBaseReasonErrorFromError(err)
	}
	switch b, err := StateBalanceValue(op.sb); {
	case err != nil:
		return nil, operation.NewBaseReasonErrorFromError(err)
	case b.Big().Compare(fee) < 0:
		return nil, operation.NewBaseReasonError("insufficient balance with fee")
	default:
		op.fee = fee
	}

	return op, nil
}

func (op *KeyUpdaterProcessor) Process(
	_ func(key string) (state.State, bool, error),
	setState func(valuehash.Hash, ...state.State) error,
) error {
	fact := op.Fact().(KeyUpdaterFact)
	// sender의 balance state(amount state)를 업데이트한다.
	op.sb = op.sb.Sub(op.fee).AddFee(op.fee)
	// key state를 저장하기
	st, err := SetStateKeysValue(op.sa, fact.keys)
	if err != nil {
		return err
	}
	return setState(fact.Hash(), st, op.sb)
}
