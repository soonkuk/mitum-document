package document

import (
	"sync"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/base/state"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/valuehash"
)

var CreateDocumentsItemProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(CreateDocumentsItemProcessor)
	},
}

var CreateDocumentsProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(CreateDocumentsProcessor)
	},
}

func (op CreateDocuments) Process(
	func(key string) (state.State, bool, error),
	func(valuehash.Hash, ...state.State) error,
) error {
	return nil
}

type CreateDocumentsItemProcessor struct {
	cp      *currency.CurrencyPool
	h       valuehash.Hash
	sender  base.Address
	item    CreateDocumentsItem
	nds     state.State // new document data state (key = document id)
	docInfo DocInfo     // new document info

}

func (opp *CreateDocumentsItemProcessor) PreProcess(
	getState func(key string) (state.State, bool, error),
	_ func(valuehash.Hash, ...state.State) error,
) error {

	if err := opp.item.IsValid(nil); err != nil {
		return operation.NewBaseReasonError(err.Error())
	}

	// check existence of new document state with documentid and get document state
	switch st, found, err := getState(StateKeyDocumentData(opp.item.DocumentId())); {
	case err != nil:
		return err
	case found:
		return operation.NewBaseReasonError("documentid already registered, %q", opp.item.DocumentId())
	default:
		opp.nds = st
	}

	id := NewDocId(opp.item.DocumentId())

	// check existence of DocumentData related accounts
	for i := range opp.item.Doc().Accounts() {
		switch _, found, err := getState(currency.StateKeyAccount(opp.item.Doc().Accounts()[i])); {
		case err != nil:
			return err
		case !found:
			return operation.NewBaseReasonError("DocumentData related accounts not found, document type : %q, address : %q", opp.item.DocType(), opp.item.Doc().Accounts()[i])
		}
	}

	// prepare doccInfo
	opp.docInfo = DocInfo{
		BaseHinter: hint.NewBaseHinter(DocInfoHint),
		id:         id,
		docType:    id.Hint().Type(),
	}

	return nil
}

func (opp *CreateDocumentsItemProcessor) Process(
	_ func(key string) (state.State, bool, error),
	_ func(valuehash.Hash, ...state.State) error,
) ([]state.State, error) {

	sts := make([]state.State, 1)

	// set new document state
	if dst, err := SetStateDocumentDataValue(opp.nds, opp.item.Doc()); err != nil {
		return nil, err
	} else {
		sts[0] = dst
	}

	return sts, nil
}

func (opp *CreateDocumentsItemProcessor) Close() error {
	opp.cp = nil
	opp.h = nil
	opp.sender = nil
	opp.item = nil
	opp.nds = nil
	opp.docInfo = DocInfo{}

	CreateDocumentsItemProcessorPool.Put(opp)

	return nil
}

type CreateDocumentsProcessor struct {
	cp *currency.CurrencyPool
	CreateDocuments
	dinv     DocumentInventory
	ndinvs   state.State
	sb       map[currency.CurrencyID]currency.AmountState // sender StateBalance
	ns       []*CreateDocumentsItemProcessor              // ItemProcessor
	required map[currency.CurrencyID][2]currency.Big      // Fee
}

func NewCreateDocumentsProcessor(cp *currency.CurrencyPool) currency.GetNewProcessor {
	return func(op state.Processor) (state.Processor, error) {
		i, ok := op.(CreateDocuments)
		if !ok {
			return nil, operation.NewBaseReasonError("not CreateDocuments, %T", op)
		}

		opp := CreateDocumentsProcessorPool.Get().(*CreateDocumentsProcessor)

		opp.cp = cp
		opp.CreateDocuments = i
		opp.dinv = DocumentInventory{}
		opp.ndinvs = nil
		opp.sb = nil
		opp.ns = nil
		opp.required = nil

		return opp, nil

	}
}

func (opp *CreateDocumentsProcessor) PreProcess(
	getState func(key string) (state.State, bool, error),
	setState func(valuehash.Hash, ...state.State) error,
) (state.Processor, error) {
	fact := opp.Fact().(CreateDocumentsFact)

	// check sender account state existence
	if err := checkExistsState(currency.StateKeyAccount(fact.sender), getState); err != nil {
		return nil, err
	}

	// check existence of document inventory state with address and get document inventory state
	switch st, found, err := getState(StateKeyDocuments(fact.sender)); {
	case err != nil:
		return nil, err
	case !found:
		opp.dinv = NewDocumentInventory(nil)
		opp.ndinvs = st
	default:
		dinv, err := StateDocumentsValue(st)
		if err != nil {
			return nil, err
		}
		opp.dinv = dinv
		opp.ndinvs = st
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
	ns := make([]*CreateDocumentsItemProcessor, len(fact.items))
	for i := range fact.items {

		c := CreateDocumentsItemProcessorPool.Get().(*CreateDocumentsItemProcessor)
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
		return nil, operation.NewBaseReasonError("invalid signing: %w", err)
	}

	opp.ns = ns

	return opp, nil
}

func (opp *CreateDocumentsProcessor) Process( // nolint:dupl
	getState func(key string) (state.State, bool, error),
	setState func(valuehash.Hash, ...state.State) error,
) error {
	// get fact
	fact := opp.Fact().(CreateDocumentsFact)

	var sts []state.State // nolint:prealloc

	// append document data state and add doc info to owner document inventory
	for i := range opp.ns {
		if s, err := opp.ns[i].Process(getState, setState); err != nil {
			return operation.NewBaseReasonError("failed to process create document item: %w", err)
		} else {
			sts = append(sts, s...)
			doc, err := StateDocumentDataValue(s[0])
			if err != nil {
				return err
			}

			if err := opp.dinv.Append(doc.Info()); err != nil {
				return err
			}
		}
	}

	opp.dinv.Sort(true)

	// prepare document inventory state and append it
	if dinvs, err := SetStateDocumentsValue(opp.ndinvs, opp.dinv); err != nil {
		return err
	} else {
		sts = append(sts, dinvs)
	}

	// append sender balance state
	for k := range opp.required {
		rq := opp.required[k]
		sts = append(sts, opp.sb[k].Sub(rq[0]).AddFee(rq[1]))
	}

	return setState(fact.Hash(), sts...)
}

func (opp *CreateDocumentsProcessor) Close() error {
	for i := range opp.ns {
		_ = opp.ns[i].Close()
	}

	opp.cp = nil
	opp.CreateDocuments = CreateDocuments{}
	opp.dinv = DocumentInventory{}
	opp.sb = nil
	opp.ndinvs = nil
	opp.required = nil

	CreateDocumentsProcessorPool.Put(opp)

	return nil
}

func (opp *CreateDocumentsProcessor) calculateItemsFee() (map[currency.CurrencyID][2]currency.Big, error) {
	fact := opp.Fact().(CreateDocumentsFact)

	items := make([]CreateDocumentsItem, len(fact.items))
	for i := range fact.items {
		items[i] = fact.items[i]
	}

	return CalculateDocumentItemsFee(opp.cp, items)
}

func CalculateDocumentItemsFee(cp *currency.CurrencyPool, items []CreateDocumentsItem) (map[currency.CurrencyID][2]currency.Big, error) {
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

func CheckDocumentOwnerEnoughBalance(
	holder base.Address,
	required map[currency.CurrencyID][2]currency.Big,
	getState func(key string) (state.State, bool, error),
) (map[currency.CurrencyID]currency.AmountState, error) {
	sb := map[currency.CurrencyID]currency.AmountState{}

	for cid := range required {
		rq := required[cid]

		st, err := existsState(currency.StateKeyBalance(holder, cid), "currency of holder", getState)
		if err != nil {
			return nil, err
		}

		am, err := currency.StateBalanceValue(st)
		if err != nil {
			return nil, operation.NewBaseReasonError("insufficient balance of sender: %w", err)
		}

		if am.Big().Compare(rq[0]) < 0 {
			return nil, operation.NewBaseReasonError(
				"insufficient balance of sender, %s; %d !> %d", holder.String(), am.Big(), rq[0])
		} else {
			sb[cid] = currency.NewAmountState(st, cid)
		}
	}

	return sb, nil
}
