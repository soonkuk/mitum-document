package blockcity

import (
	"sync"

	"github.com/pkg/errors"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/base/state"
	"github.com/spikeekips/mitum/util/valuehash"
)

var UpdateDocumentsItemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(UpdateDocumentsItemProcessor)
	},
}

var UpdateDocumentsProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(UpdateDocumentsProcessor)
	},
}

func (op UpdateDocuments) Process(
	func(key string) (state.State, bool, error),
	func(valuehash.Hash, ...state.State) error,
) error {
	return nil
}

type UpdateDocumentsItemProcessor struct {
	cp     *currency.CurrencyPool
	h      valuehash.Hash
	sender base.Address
	item   UpdateDocumentsItem
	nds    state.State       // new document data state (key = document nickname)
	dinv   DocumentInventory // document inventory
	ndinvs state.State       // document inventory state (key = owner address)

}

func (opp *UpdateDocumentsItemProcessor) PreProcess(
	getState func(key string) (state.State, bool, error),
	_ func(valuehash.Hash, ...state.State) error,
) error {

	if err := opp.item.IsValid(nil); err != nil {
		return operation.NewBaseReasonError(err.Error())
	}

	// check existence of owner account
	if _, found, err := getState(currency.StateKeyAccount(opp.item.Doc().Owner())); err != nil {
		return err
	} else if !found {
		return operation.NewBaseReasonError("owner does not exist, %q", opp.item.Doc().Owner())
	}

	// check existence of document inventory state with owner address
	switch st, found, err := getState(StateKeyDocuments(opp.item.Doc().Owner())); {
	case err != nil:
		return err
	case !found:
		return operation.NewBaseReasonError("Owner has no document inventory, %v", opp.item.Doc().Owner())

		// get document inventory of owner
	default:
		dinv, err := StateDocumentsValue(st)
		if err != nil {
			return err
		}
		opp.dinv = dinv
		opp.ndinvs = st
	}

	// check existence of document state with documentid
	switch st, found, err := getState(StateKeyDocumentData(opp.item.DocumentId())); {
	case err != nil:
		return err
	case !found:
		return operation.NewBaseReasonError("document not registered with documentid, %q", opp.item.DocumentId())
	default:
		opp.nds = st
	}

	dd, err := StateDocumentDataValue(opp.nds)
	if err != nil {
		return err
	}

	if !dd.Owner().Equal(opp.item.Doc().Owner()) {
		return operation.NewBaseReasonError("item's Owner not matched with Owner in document, %v", opp.item.Doc().Owner())
	}

	if dd.DocumentType() != opp.item.DocType() {
		return operation.NewBaseReasonError("item's Document type not matched with document type in document, %v", opp.item.DocType())
	}

	// update document data state
	st, err := SetStateDocumentDataValue(opp.nds, opp.item.Doc())
	if err != nil {
		return err
	}
	opp.nds = st

	return nil
}

func (opp *UpdateDocumentsItemProcessor) Process(
	_ func(key string) (state.State, bool, error),
	_ func(valuehash.Hash, ...state.State) error,
) ([]state.State, error) {

	sts := make([]state.State, 1)
	sts[0] = opp.nds

	return sts, nil
}

func (opp *UpdateDocumentsItemProcessor) Close() error {
	opp.cp = nil
	opp.h = nil
	opp.sender = nil
	opp.item = nil
	opp.nds = nil
	opp.dinv = DocumentInventory{}
	opp.ndinvs = nil

	UpdateDocumentsItemProcessorPool.Put(opp)

	return nil
}

type UpdateDocumentsProcessor struct {
	cp *currency.CurrencyPool
	UpdateDocuments
	sb       map[currency.CurrencyID]currency.AmountState // sender StateBalance
	ns       []*UpdateDocumentsItemProcessor              // ItemProcessor
	required map[currency.CurrencyID][2]currency.Big      // Fee
}

func NewUpdateDocumentsProcessor(cp *currency.CurrencyPool) currency.GetNewProcessor {
	return func(op state.Processor) (state.Processor, error) {
		i, ok := op.(UpdateDocuments)
		if !ok {
			return nil, operation.NewBaseReasonError("not UpdateDocuments, %T", op)
		}
		return &UpdateDocumentsProcessor{
			cp:              cp,
			UpdateDocuments: i,
			sb:              nil,
			ns:              nil,
			required:        nil,
		}, nil
	}
}

func (opp *UpdateDocumentsProcessor) PreProcess(
	getState func(key string) (state.State, bool, error),
	setState func(valuehash.Hash, ...state.State) error,
) (state.Processor, error) {
	fact := opp.Fact().(UpdateDocumentsFact)

	// check sender account state existence
	if err := checkExistsState(currency.StateKeyAccount(fact.sender), getState); err != nil {
		return nil, err
	}

	// prepare sender balance state
	if required, err := opp.calculateItemsFee(); err != nil {
		return nil, operation.NewBaseReasonError("failed to calculate fee: %w", err)
	} else if sb, err := CheckDocumentOwnerEnoughBalance(fact.sender, required, getState); err != nil {
		return nil, err
	} else {
		opp.required = required
		opp.sb = sb
	}

	// prepare item processor for each items
	ns := make([]*UpdateDocumentsItemProcessor, len(fact.items))
	for i := range fact.items {
		c := UpdateDocumentsItemProcessorPool.Get().(*UpdateDocumentsItemProcessor)
		c.cp = opp.cp
		c.h = opp.Hash()
		c.sender = fact.sender
		c.item = fact.items[i]

		if err := c.PreProcess(getState, setState); err != nil {
			return nil, err
		}

		ns[i] = c
	}

	// check fact sign
	if err := checkFactSignsByState(fact.sender, opp.Signs(), getState); err != nil {
		return nil, errors.Wrap(err, "invalid signing")
	}

	opp.ns = ns

	return opp, nil
}

func (opp *UpdateDocumentsProcessor) Process( // nolint:dupl
	getState func(key string) (state.State, bool, error),
	setState func(valuehash.Hash, ...state.State) error,
) error {
	// get fact
	fact := opp.Fact().(UpdateDocumentsFact)

	var sts []state.State // nolint:prealloc

	// append document data state and add doc info to owner document inventory
	for i := range opp.ns {
		if s, err := opp.ns[i].Process(getState, setState); err != nil {
			return operation.NewBaseReasonError("failed to process update document item: %w", err)
		} else {
			sts = append(sts, s...)
		}
	}

	// append sender balance state
	for k := range opp.required {
		rq := opp.required[k]
		sts = append(sts, opp.sb[k].Sub(rq[0]).AddFee(rq[1]))
	}

	return setState(fact.Hash(), sts...)
}

func (opp *UpdateDocumentsProcessor) Close() error {
	for i := range opp.ns {
		_ = opp.ns[i].Close()
	}

	opp.cp = nil
	opp.UpdateDocuments = UpdateDocuments{}
	opp.sb = nil
	opp.required = nil

	UpdateDocumentsProcessorPool.Put(opp)

	return nil
}

func (opp *UpdateDocumentsProcessor) calculateItemsFee() (map[currency.CurrencyID][2]currency.Big, error) {
	fact := opp.Fact().(UpdateDocumentsFact)
	items := make([]UpdateDocumentsItem, len(fact.items))
	for i := range fact.items {
		items[i] = fact.items[i]
	}

	return CalculateUpdateDocumentItemsFee(opp.cp, items)
}

func CalculateUpdateDocumentItemsFee(cp *currency.CurrencyPool, items []UpdateDocumentsItem) (map[currency.CurrencyID][2]currency.Big, error) {
	required := map[currency.CurrencyID][2]currency.Big{}

	for i := range items {
		it := items[i]

		rq := [2]currency.Big{currency.ZeroBig, currency.ZeroBig}

		if k, found := required[it.Currency()]; found {
			rq = k
		}

		if cp == nil {
			required[it.Currency()] = [2]currency.Big{rq[0], rq[1]}

			continue
		}

		feeer, found := cp.Feeer(it.Currency())
		if !found {
			return nil, operation.NewBaseReasonError("unknown currency id found, %q", it.Currency())
		}
		switch k, err := feeer.Fee(currency.ZeroBig); {
		case err != nil:
			return nil, err
		case !k.OverZero():
			required[it.Currency()] = [2]currency.Big{rq[0], rq[1]}
		default:
			required[it.Currency()] = [2]currency.Big{rq[0].Add(k), rq[1].Add(k)}
		}

	}

	return required, nil
}
