package currency

import (
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/base/state"
	"github.com/spikeekips/mitum/util/valuehash"
	"golang.org/x/xerrors"
)

func (op TransferDocuments) Process(
	func(key string) (state.State, bool, error),
	func(valuehash.Hash, ...state.State) error,
) error {
	// NOTE Process is nil func
	return nil
}

type TransferDocumentsItemProcessor struct {
	cp   *CurrencyPool
	h    valuehash.Hash
	item TransferDocumentsItem
	das  state.State // document account state
	dfs  state.State // document filedata state
	fd   FileData
	keys Keys
}

// PrePrecess는 document의 StateAccount를 fetch하고 Receiver의 existence확인 후 Keys를 fetch한다.
func (opp *TransferDocumentsItemProcessor) PreProcess(
	getState func(key string) (state.State, bool, error),
	_ func(valuehash.Hash, ...state.State) error,
) error {
	if st, err := existsState(StateKeyDocument(opp.item.Document()), "document", getState); err != nil {
		return err
	} else {
		opp.das = st
	}
	if st, err := existsState(StateKeyFileData(opp.item.Document()), "filedata", getState); err != nil {
		return err
	} else {
		opp.dfs = st
	}
	if fd, err := StateFileDataValue(opp.dfs); err != nil {
		return err
	} else {
		opp.fd = fd
	}

	// get account address of receiver
	receiver := opp.item.Receiver()
	// check existence of reciever account and get key State
	opp.fd = opp.fd.WithData(opp.fd.signcode, opp.item.Receiver())

	if st, err := existsState(StateKeyAccount(receiver), "receiver", getState); err != nil {
		return err
	} else if ks, err := StateKeysValue(st); err != nil {
		return err
	} else {
		opp.keys = ks
	}

	return nil
}

// Process는 Document의 StateAccount에 대한 value를 Receiver의 Keys 값으로 업데이트한다.
func (opp *TransferDocumentsItemProcessor) Process(
	_ func(key string) (state.State, bool, error),
	_ func(valuehash.Hash, ...state.State) error,
) ([]state.State, error) {
	sts := make([]state.State, 2)

	// document의 StateAccount에서 Keys값을 업데이트
	if st, err := SetStateKeysValue(opp.das, opp.keys); err != nil {
		return nil, err
	} else {
		sts[0] = st
	}
	// document의 StateFileData에서 Owner값을 업데이트
	if st, err := SetStateFileDataValue(opp.dfs, opp.fd); err != nil {
		return nil, err
	} else {
		sts[1] = st
	}

	return sts, nil
}

type TransferDocumentsProcessor struct {
	cp *CurrencyPool
	TransferDocuments
	sb       map[CurrencyID]AmountState
	ns       []*TransferDocumentsItemProcessor
	required map[CurrencyID][2]Big
}

func NewTransferDocumentsProcessor(cp *CurrencyPool) GetNewProcessor {
	return func(op state.Processor) (state.Processor, error) {
		if i, ok := op.(TransferDocuments); !ok {
			return nil, xerrors.Errorf("not TransferDocuments, %T", op)
		} else {
			return &TransferDocumentsProcessor{
				cp:                cp,
				TransferDocuments: i,
			}, nil
		}
	}
}

func (opp *TransferDocumentsProcessor) PreProcess(
	getState func(key string) (state.State, bool, error),
	setState func(valuehash.Hash, ...state.State) error,
) (state.Processor, error) {
	fact := opp.Fact().(TransferDocumentsFact)

	// fetch sender StateAccount
	if err := checkExistsState(StateKeyAccount(fact.sender), getState); err != nil {
		return nil, err
	}

	if required, err := opp.calculateItemsFee(); err != nil {
		return nil, operation.NewBaseReasonErrorFromError(err)
	} else if sb, err := CheckDocumentOwnerEnoughBalance(fact.sender, required, getState); err != nil {
		return nil, err
	} else {
		opp.required = required
		// DocumentOwner의 현재 StateBalance
		opp.sb = sb
	}

	ns := make([]*TransferDocumentsItemProcessor, len(fact.items))
	for i := range fact.items {
		t := &TransferDocumentsItemProcessor{cp: opp.cp, h: opp.Hash(), item: fact.items[i]}
		if err := t.PreProcess(getState, setState); err != nil {
			return nil, operation.NewBaseReasonErrorFromError(err)
		}

		ns[i] = t
	}

	if err := checkFactSignsByState(fact.sender, opp.Signs(), getState); err != nil {
		return nil, xerrors.Errorf("invalid signing: %w", err)
	}

	opp.ns = ns

	return opp, nil
}

func (opp *TransferDocumentsProcessor) Process( // nolint:dupl
	getState func(key string) (state.State, bool, error),
	setState func(valuehash.Hash, ...state.State) error,
) error {
	fact := opp.Fact().(TransferDocumentsFact)

	var sts []state.State // nolint:prealloc
	for i := range opp.ns {
		if s, err := opp.ns[i].Process(getState, setState); err != nil {
			return operation.NewBaseReasonError("failed to process transfer document item: %w", err)
		} else {
			sts = append(sts, s...)
		}
	}

	// sender의 currencyID별 stateBalance를 변경하여 배열에 모은다.
	for k := range opp.required {
		rq := opp.required[k]
		sts = append(sts, opp.sb[k].Sub(rq[0]).AddFee(rq[1]))
	}

	return setState(fact.Hash(), sts...)
}

func (opp *TransferDocumentsProcessor) calculateItemsFee() (map[CurrencyID][2]Big, error) {
	fact := opp.Fact().(TransferDocumentsFact)
	items := make([]TransferDocumentsItem, len(fact.items))
	for i := range fact.items {
		items[i] = fact.items[i]
	}

	required := map[CurrencyID][2]Big{}

	for i := range items {
		it := items[i]

		// currencyPool에 currency가 없으니 fee 계산을 안하고 그냥 넘어간다.
		// currencyPool에 currency가 없으면 에러인 것 같은데 이런 경우도 있을까?
		rq := [2]Big{ZeroBig, ZeroBig}

		if k, found := required[it.Currency()]; found {
			rq = k
		}
		if opp.cp == nil {
			required[it.Currency()] = rq

			continue
		}

		// currencyID로 fee 계산
		feeer, found := opp.cp.Feeer(it.Currency())
		if !found {
			return nil, xerrors.Errorf("unknown currency id found, %q", it.Currency())
		}
		switch k, err := feeer.Fee(ZeroBig); {
		case err != nil:
			return nil, err
		case !k.OverZero():
			required[it.Currency()] = rq
		// Fee가 있으면 required amount 누적하고 required fee도 누적한다.
		default:
			required[it.Currency()] = [2]Big{rq[0].Add(k), rq[1].Add(k)}
		}

	}

	return required, nil
}
