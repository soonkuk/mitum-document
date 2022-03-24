package document // nolint: dupl, revive

import (
	"sync"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/base/state"
	"github.com/spikeekips/mitum/util/valuehash"
)

var SignDocumentsItemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(SignDocumentsItemProcessor)
	},
}

var SignDocumentsProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(SignDocumentsProcessor)
	},
}

func (SignDocuments) Process(
	func(key string) (state.State, bool, error),
	func(valuehash.Hash, ...state.State) error,
) error {
	return nil
}

type SignDocumentsItemProcessor struct {
	cp     *currency.CurrencyPool
	h      valuehash.Hash
	sender base.Address
	item   SignDocumentsItem
	nds    state.State       // new document data state (key = document filehash)
	dinv   DocumentInventory // document inventory
	ndinvs state.State       // document inventory state (key = owner address)

}

func (opp *SignDocumentsItemProcessor) PreProcess( // nolint:revive
	getState func(key string) (state.State, bool, error),
	_ func(valuehash.Hash, ...state.State) error,
) error {
	if err := opp.item.IsValid(nil); err != nil {
		return err
	}

	// check existence of owner account
	if _, found, err := getState(currency.StateKeyAccount(opp.item.Owner())); err != nil {
		return err
	} else if !found {
		return operation.NewBaseReasonError("owner does not exist, %q", opp.item.Owner())
	}

	// check existence of document inventory state with owner address
	switch st, found, err := getState(StateKeyDocuments(opp.item.Owner())); {
	case err != nil:
		return err
	case !found:
		return operation.NewBaseReasonError("Owner has no document inventory, %v", opp.item.Owner())
	default:
		dinv, err := StateDocumentsValue(st)
		if err != nil {
			return err
		}
		opp.dinv = dinv
		opp.ndinvs = st
	}

	if _, err := opp.dinv.Get(opp.item.DocumentID()); err != nil {
		return err
	}

	// check existence of new document state with documentid
	switch st, found, err := getState(StateKeyDocumentData(opp.item.DocumentID())); {
	case err != nil:
		return err
	case !found:
		return operation.NewBaseReasonError("document not registered with documentid, %q", opp.item.DocumentID())
	default:
		opp.nds = st
	}

	dd, err := StateDocumentDataValue(opp.nds)
	if err != nil {
		return err
	}

	v, ok := dd.(BSDocData)
	if !ok {
		return operation.NewBaseReasonError("Document is not Blocksign Document, %v", opp.item.DocumentID())
	}
	if !v.Creator().Address().Equal(opp.item.Owner()) {
		return operation.NewBaseReasonError("Owner not matched with creator in document, %v", opp.item.Owner())
	}

	if len(v.Signers()) < 1 {
		return operation.NewBaseReasonError("sender not found in document Signers, %v", opp.sender)
	}
	// check signer exist in document data signers
	for i := range v.Signers() {
		if v.Signers()[i].Address().Equal(opp.sender) {
			v.Signers()[i].SetSigned()
			break
		}
		if i == (len(v.Signers()) - 1) {
			return operation.NewBaseReasonError("sender not found in document Signers, %v", opp.sender)
		}
	}

	// update document data state
	st, err := SetStateDocumentDataValue(opp.nds, dd)
	if err != nil {
		return err
	}
	opp.nds = st

	return nil
}

func (opp *SignDocumentsItemProcessor) Process(
	_ func(key string) (state.State, bool, error),
	_ func(valuehash.Hash, ...state.State) error,
) ([]state.State, error) {
	sts := make([]state.State, 1)
	sts[0] = opp.nds

	return sts, nil
}

func (opp *SignDocumentsItemProcessor) Close() error {
	opp.cp = nil
	opp.h = nil
	opp.sender = nil
	opp.item = nil
	opp.nds = nil
	opp.dinv = DocumentInventory{}
	opp.ndinvs = nil

	CreateDocumentsItemProcessorPool.Put(opp)

	return nil
}

type SignDocumentsProcessor struct {
	cp *currency.CurrencyPool
	SignDocuments
	sb       map[currency.CurrencyID]currency.AmountState // sender StateBalance
	ns       []*SignDocumentsItemProcessor                // ItemProcessor
	required map[currency.CurrencyID][2]currency.Big      // Fee
}

func NewSignDocumentsProcessor(cp *currency.CurrencyPool) currency.GetNewProcessor {
	return func(op state.Processor) (state.Processor, error) {
		i, ok := op.(SignDocuments)
		if !ok {
			return nil, operation.NewBaseReasonError("not SignDocuments, %T", op)
		}
		return &SignDocumentsProcessor{
			cp:            cp,
			SignDocuments: i,
		}, nil
	}
}

func (opp *SignDocumentsProcessor) PreProcess(
	getState func(key string) (state.State, bool, error),
	setState func(valuehash.Hash, ...state.State) error,
) (state.Processor, error) {
	fact := opp.Fact().(SignDocumentsFact)

	// check sender account state existence
	if err := checkExistsState(currency.StateKeyAccount(fact.sender), getState); err != nil {
		return nil, err
	}

	if required, err := opp.calculateItemsFee(); err != nil {
		return nil, operation.NewBaseReasonError("failed to calculate fee: %w", err)
	} else if sb, err := CheckDocumentOwnerEnoughBalance(fact.sender, required, getState); err != nil {
		return nil, err
	} else {
		opp.required = required
		opp.sb = sb
	}

	ns := make([]*SignDocumentsItemProcessor, len(fact.items))
	for i := range fact.items {
		c := &SignDocumentsItemProcessor{cp: opp.cp, sender: fact.sender, h: opp.Hash(), item: fact.items[i]}
		if err := c.PreProcess(getState, setState); err != nil {
			return nil, operation.NewBaseReasonErrorFromError(err)
		}

		ns[i] = c
	}

	// check fact sign
	if err := checkFactSignsByState(fact.sender, opp.Signs(), getState); err != nil {
		return nil, operation.NewBaseReasonError("invalid signing: %w", err)
	}

	opp.ns = ns

	return opp, nil
}

func (opp *SignDocumentsProcessor) Process( // nolint:dupl
	getState func(key string) (state.State, bool, error),
	setState func(valuehash.Hash, ...state.State) error,
) error {
	fact := opp.Fact().(SignDocumentsFact)

	var sts []state.State // nolint:prealloc

	for i := range opp.ns {
		s, err := opp.ns[i].Process(getState, setState)
		if err != nil {
			return operation.NewBaseReasonError("failed to process create document item: %w", err)
		}
		sts = append(sts, s...)
	}

	for k := range opp.required {
		rq := opp.required[k]
		sts = append(sts, opp.sb[k].Sub(rq[0]).AddFee(rq[1]))
	}

	return setState(fact.Hash(), sts...)
}

func (opp *SignDocumentsProcessor) Close() error {
	for i := range opp.ns {
		_ = opp.ns[i].Close()
	}

	opp.cp = nil
	opp.SignDocuments = SignDocuments{}
	opp.sb = nil
	opp.required = nil

	CreateDocumentsProcessorPool.Put(opp)

	return nil
}

func (opp *SignDocumentsProcessor) calculateItemsFee() (map[currency.CurrencyID][2]currency.Big, error) {
	fact := opp.Fact().(SignDocumentsFact)
	items := make([]SignDocumentsItem, len(fact.items))
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
			return nil, operation.NewBaseReasonError("unknown currency id found, %q", it.Currency())
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
