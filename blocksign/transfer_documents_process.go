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
	cp    *currency.CurrencyPool
	h     valuehash.Hash
	item  TransferDocumentsItem
	sas   state.State // sender account state
	das   state.State // document account state
	dds   state.State // document filedata state
	dd    DocumentData
	keys  currency.Keys // document owner keys
	skeys currency.Keys // sender keys
}

func (opp *TransferDocumentsItemProcessor) PreProcess(
	getState func(key string) (state.State, bool, error),
	_ func(valuehash.Hash, ...state.State) error,
) error {
	if st, err := currency.ExistsState(currency.StateKeyDocument(opp.item.Sender()), "sender", getState); err != nil {
		return err
	} else {
		opp.sas = st
	}

	if skeys, err := currency.StateKeysValue(opp.sas); err != nil {
		return err
	} else {
		opp.skeys = skeys
	}

	if st, err := currency.ExistsState(currency.StateKeyDocument(opp.item.Document()), "document", getState); err != nil {
		return err
	} else {
		opp.das = st
	}

	dkeys, err := currency.StateKeysValue(opp.das)
	if err != nil {
		return err
	}
	if !opp.skeys.Equal(dkeys) {
		return err
	}

	if st, err := currency.ExistsState(StateKeyDocumentData(opp.item.Document()), "document data", getState); err != nil {
		return err
	} else {
		opp.dds = st
	}
	if dd, err := StateDocumentDataValue(opp.dds); err != nil {
		return err
	} else {
		opp.dd = dd
	}

	// get account address of receiver
	receiver := opp.item.Receiver()
	// change owner of documentData as receiver
	opp.dd = opp.dd.WithData(opp.dd.FileHash(), opp.dd.Creator(), opp.item.Receiver(), opp.dd.Signers())

	// check existence of reciever account and get receiver keys
	if st, err := currency.ExistsState(currency.StateKeyAccount(receiver), "receiver", getState); err != nil {
		return err
	} else if ks, err := currency.StateKeysValue(st); err != nil {
		return err
	} else {
		opp.keys = ks
	}

	return nil
}

func (opp *TransferDocumentsItemProcessor) Process(
	_ func(key string) (state.State, bool, error),
	_ func(valuehash.Hash, ...state.State) error,
) ([]state.State, error) {
	sts := make([]state.State, 2)

	// replace Document acccount key as receivers(new owner) key
	if st, err := currency.SetStateKeysValue(opp.das, opp.keys); err != nil {
		return nil, err
	} else {
		sts[0] = st
	}
	// replace Document data
	if st, err := SetStateDocumentDataValue(opp.dds, opp.dd); err != nil {
		return nil, err
	} else {
		sts[1] = st
	}

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
		if !fact.sender.Equal(fact.items[i].Sender()) {
			return nil, xerrors.Errorf("item sender is not same with fact sender, %q", fact.items[i].Sender())
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
