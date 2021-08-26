package blocksign

import (
	"github.com/pkg/errors"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/base/state"
	"github.com/spikeekips/mitum/util/valuehash"
)

func (op TransferDocuments) Process(
	func(key string) (state.State, bool, error),
	func(valuehash.Hash, ...state.State) error,
) error {
	// NOTE Process is nil func
	return nil
}

type TransferDocumentsItemProcessor struct {
	cp     *currency.CurrencyPool
	h      valuehash.Hash
	item   TransferDocumentsItem
	rdinvs state.State // receiver document inventory state
	dds    state.State // document data state
	di     DocInfo
	dd     DocumentData
}

func (opp *TransferDocumentsItemProcessor) PreProcess(
	getState func(key string) (state.State, bool, error),
	_ func(valuehash.Hash, ...state.State) error,
) error {
	// check existence of reciever account
	if _, found, err := getState(currency.StateKeyAccount(opp.item.Receiver())); err != nil {
		return err
	} else if !found {
		return errors.Errorf("receiver account not found, %q", opp.item.Receiver())
	}

	// check existence of owner account
	if _, found, err := getState(currency.StateKeyAccount(opp.item.Owner())); err != nil {
		return err
	} else if !found {
		return errors.Errorf("owner account not found, %q", opp.item.Owner())
	}

	// check document id is lesser than or equal to global last document id
	switch st, found, err := getState(StateKeyLastDocumentId); {
	case err != nil:
		return err
	case !found:
		return errors.Errorf("Document is not registered")
	default:
		v, err := StateLastDocumentIdValue(st)
		if err != nil {
			return err
		}
		if opp.item.DocumentId().Sub(v.idx).OverZero() {
			return errors.Errorf("Document Id is greater than last registered document id")
		}
	}

	// get owner document inventory by document id and get filehash
	if st, err := existsState(StateKeyDocuments(opp.item.Owner()), "document inventory", getState); err != nil {
		return err
	} else {
		dinv, err := StateDocumentsValue(st)
		if err != nil {
			return err
		}
		if !dinv.Exists(opp.item.DocumentId()) {
			return errors.Errorf("document id not registered to account, %v", opp.item.DocumentId())
		}
		docInfo, err := dinv.Get(opp.item.DocumentId())
		if err != nil {
			return err
		}
		opp.di = docInfo

		// sopp.odinvs = st
	}

	// get document data state by filehash
	if st, err := existsState(StateKeyDocumentData(opp.di.filehash), "document data", getState); err != nil {
		return errors.Errorf("document data of filehash not found, %v", opp.di.filehash)
	} else {
		opp.dds = st
	}

	// replace document data owner with receiver
	if dd, err := StateDocumentDataValue(opp.dds); err != nil {
		return err
	} else if !dd.Owner().Equal(opp.item.Owner()) {
		return err
	} else {
		ndd := dd.WithData(dd.Info(), dd.Creator(), opp.item.Receiver(), dd.Signers())
		opp.dd = ndd
	}

	// get receiver document inventory state.
	dst, found, err := getState(StateKeyDocuments(opp.item.Receiver()))
	if err != nil {
		return err
	} else if !found {
		dinv := NewDocumentInventory(nil)
		ndst, err := SetStateDocumentsValue(dst, dinv)
		if err != nil {
			return err
		}
		dst = ndst
	} else {
		dinv, err := StateDocumentsValue(dst)
		if err != nil {
			return err
		}
		if dinv.Exists(opp.di.Index()) {
			return errors.Errorf("Document id already registered in receiver's document inventory, %v", opp.di.idx)
		}
	}
	opp.rdinvs = dst

	return nil
}

func (opp *TransferDocumentsItemProcessor) Process(
	getState func(key string) (state.State, bool, error),
	_ func(valuehash.Hash, ...state.State) error,
) ([]state.State, error) {
	sts := make([]state.State, 2)

	sts[0] = opp.rdinvs

	ndst, err := SetStateDocumentDataValue(opp.dds, opp.dd)
	if err != nil {
		return nil, err
	}
	sts[1] = ndst

	return sts, nil
}

type TransferDocumentsProcessor struct {
	cp *currency.CurrencyPool
	TransferDocuments
	odinv    DocumentInventory
	odinvs   state.State
	sb       map[currency.CurrencyID]currency.AmountState
	ns       []*TransferDocumentsItemProcessor
	required map[currency.CurrencyID][2]currency.Big
}

func NewTransferDocumentsProcessor(cp *currency.CurrencyPool) currency.GetNewProcessor {
	return func(op state.Processor) (state.Processor, error) {
		if i, ok := op.(TransferDocuments); !ok {
			return nil, errors.Errorf("not TransferDocuments, %T", op)
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
	if err := checkExistsState(currency.StateKeyAccount(fact.sender), getState); err != nil {
		return nil, err
	}
	// check sender is owner
	for i := range fact.items {
		if !fact.sender.Equal(fact.items[i].Owner()) {
			return nil, errors.Errorf("item Owner is not same with fact sender, %q", fact.items[i].Owner())
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

	if err := checkFactSignsByState(fact.sender, opp.Signs(), getState); err != nil {
		return nil, errors.Errorf("invalid signing: %w", err)
	}

	opp.ns = ns

	return opp, nil
}

func (opp *TransferDocumentsProcessor) Process( // nolint:dupl
	getState func(key string) (state.State, bool, error),
	setState func(valuehash.Hash, ...state.State) error,
) error {
	fact := opp.Fact().(TransferDocumentsFact)

	// get owner document inventory state and document inventory value
	if st, err := existsState(StateKeyDocuments(fact.Sender()), "document inventory", getState); err != nil {
		return err
	} else {
		dinv, err := StateDocumentsValue(st)
		if err != nil {
			return err
		}
		opp.odinv = dinv
		opp.odinvs = st
	}

	rdinvs := map[string]state.State{}
	var sts []state.State // nolint:prealloc
	for i := range opp.ns {
		if s, err := opp.ns[i].Process(getState, setState); err != nil {
			return operation.NewBaseReasonError("failed to process transfer document item: %w", err)
		} else {
			sts = append(sts, s[1])
			doc, err := StateDocumentDataValue(s[1])
			if err != nil {
				return err
			}
			opp.odinv.Romove(doc.Info())

			ra := s[0].Key()[:len(s[0].Key())-len(StateKeyDocumentsSuffix)]
			if _, found := rdinvs[ra]; !found {
				rdinvs[ra] = s[0]
			}

			rdinv, err := StateDocumentsValue(rdinvs[ra])
			if err != nil {
				return err
			}

			rdinv.Append(doc.Info())
			rdinv.Sort(true)
			nst, err := SetStateDocumentsValue(rdinvs[ra], rdinv)
			if err != nil {
				return err
			}
			rdinvs[ra] = nst
		}
	}

	opp.odinv.Sort(true)
	// set owner document inventory state value
	onst, err := SetStateDocumentsValue(opp.odinvs, opp.odinv)
	if err != nil {
		return err
	}
	sts = append(sts, onst)

	for k := range rdinvs {
		rnst := rdinvs[k]
		sts = append(sts, rnst)
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
			return nil, errors.Errorf("unknown currency id found, %q", it.Currency())
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
