package blocksign

import (
	"github.com/soonkuk/mitum-data/currency"
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
	cp   *currency.CurrencyPool
	h    valuehash.Hash
	item TransferDocumentsItem
	dds  state.State // document filedata state
	dd   DocumentData
}

func (opp *TransferDocumentsItemProcessor) PreProcess(
	getState func(key string) (state.State, bool, error),
	_ func(valuehash.Hash, ...state.State) error,
) error {
	// check existence of reciever account and get receiver keys
	if _, found, err := getState(currency.StateKeyAccount(opp.item.Receiver())); err != nil {
		return err
	} else if !found {
		return xerrors.Errorf("receiver account not found, %q", opp.item.Receiver())
	}

	if st, err := currency.ExistsState(StateKeyDocumentData(opp.item.Owner(), opp.item.DocumentId()), "document data", getState); err != nil {
		return err
	} else {
		opp.dds = st
	}

	if dd, err := StateDocumentDataValue(opp.dds); err != nil {
		return err
	} else {
		opp.dd = dd
	}

	return nil
}

func (opp *TransferDocumentsItemProcessor) Process(
	getState func(key string) (state.State, bool, error),
	_ func(valuehash.Hash, ...state.State) error,
) ([]state.State, error) {
	sts := make([]state.State, 1)

	dst, found, err := getState(StateKeyDocumentData(opp.item.Receiver(), opp.dd.DocumentId()))
	if err != nil {
		return nil, err
	} else if found {
		return nil, xerrors.Errorf("receiver already own document, %q", opp.dd.DocumentId())
	}

	st, err := SetStateDocumentDataValue(dst, opp.dd)
	if err != nil {
		return nil, err
	}

	sts[0] = st
	return sts, nil
}

type TransferDocumentsProcessor struct {
	cp *currency.CurrencyPool
	TransferDocuments
	sb       map[currency.CurrencyID]currency.AmountState
	ns       []*TransferDocumentsItemProcessor
	required map[currency.CurrencyID][2]currency.Big
}

func NewTransferDocumentsProcessor(cp *currency.CurrencyPool) currency.GetNewProcessor {
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
	if err := currency.CheckExistsState(currency.StateKeyAccount(fact.sender), getState); err != nil {
		return nil, err
	}
	// check sender is owner
	for i := range fact.items {
		if !fact.sender.Equal(fact.items[i].Owner()) {
			return nil, xerrors.Errorf("item Owner is not same with fact sender, %q", fact.items[i].Owner())
		}
	}

	if required, err := opp.calculateItemsFee(); err != nil {
		return nil, operation.NewBaseReasonErrorFromError(err)
	} else if sb, err := CheckDocumentOwnerEnoughBalance(fact.sender, required, getState); err != nil {
		return nil, err
	} else {
		opp.required = required
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

	if err := currency.CheckFactSignsByState(fact.sender, opp.Signs(), getState); err != nil {
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

	for k := range opp.required {
		rq := opp.required[k]
		sts = append(sts, opp.sb[k].Sub(rq[0]).AddFee(rq[1]))
	}

	return setState(fact.Hash(), sts...)
}

func (opp *TransferDocumentsProcessor) calculateItemsFee() (map[currency.CurrencyID][2]currency.Big, error) {
	fact := opp.Fact().(TransferDocumentsFact)
	items := make([]TransferDocumentsItem, len(fact.items))
	for i := range fact.items {
		items[i] = fact.items[i]
	}

	required := map[currency.CurrencyID][2]currency.Big{}

	for i := range items {
		it := items[i]

		rq := [2]currency.Big{currency.ZeroBig, currency.ZeroBig}

		if k, found := required[it.Currency()]; found {
			rq = k
		}
		if opp.cp == nil {
			required[it.Currency()] = rq

			continue
		}

		feeer, found := opp.cp.Feeer(it.Currency())
		if !found {
			return nil, xerrors.Errorf("unknown currency id found, %q", it.Currency())
		}
		switch k, err := feeer.Fee(currency.ZeroBig); {
		case err != nil:
			return nil, err
		case !k.OverZero():
			required[it.Currency()] = rq
		default:
			required[it.Currency()] = [2]currency.Big{rq[0].Add(k), rq[1].Add(k)}
		}

	}

	return required, nil
}
