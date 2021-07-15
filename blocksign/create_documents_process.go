package blocksign

import (
	"github.com/soonkuk/mitum-data/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/base/state"
	"github.com/spikeekips/mitum/util/valuehash"
	"golang.org/x/xerrors"
)

func (op CreateDocuments) Process(
	func(key string) (state.State, bool, error),
	func(valuehash.Hash, ...state.State) error,
) error {
	return nil
}

type CreateDocumentsItemProcessor struct {
	cp     *currency.CurrencyPool
	sender base.Address
	h      valuehash.Hash
	item   CreateDocumentsItem
	// nas  state.State   // new document account state
	nds state.State // new document data state
	// keys  currency.Keys // creator keys
	docId DocId // new document id
}

func (opp *CreateDocumentsItemProcessor) PreProcess(
	getState func(key string) (state.State, bool, error),
	_ func(valuehash.Hash, ...state.State) error,
) error {

	if err := opp.item.IsValid(nil); err != nil {
		return err
	}

	// get global last document id
	switch st, found, err := getState(StateKeyLastDocumentId); {
	case err != nil:
		return err
	case !found:
		opp.docId = NewDocId(0)
	default:
		v, err := StateLastDocumentIdValue(st)
		if err != nil {
			return err
		}
		v.idx.Add(currency.NewBig(1))
		opp.docId = v
	}

	// check existence of new document state with account address
	if _, found, err := getState(StateKeyDocumentDataExcl(opp.item.FileHash())); err != nil {
		return err
	} else if found {
		return xerrors.Errorf("already registered, %q", opp.item.FileHash())
	}

	// check existence of new document state with account address
	switch st, found, err := getState(StateKeyDocumentData(opp.sender, opp.docId)); {
	case err != nil:
		return err
	case found:
		return xerrors.Errorf("document id already registered, %q", opp.docId)
	default:
		opp.nds = st
	}

	// check sigenrs account existence
	signers := opp.item.Signers()
	for i := range signers {
		switch _, found, err := getState(currency.StateKeyAccount(signers[i])); {
		case err != nil:
			return err
		case !found:
			return xerrors.Errorf("signer account not found, %q", signers[i])
		}
	}

	return nil
}

func (opp *CreateDocumentsItemProcessor) Process(
	_ func(key string) (state.State, bool, error),
	_ func(valuehash.Hash, ...state.State) error,
) ([]state.State, error) {

	sts := make([]state.State, 1)

	signers := make([]DocSign, len(opp.item.Signers()))
	for i := range opp.item.Signers() {
		signers[i] = NewDocSign(opp.item.Signers()[i], false)
	}

	// document state with new document id
	docData := DocumentData{
		fileHash: opp.item.FileHash(),
		id:       opp.docId,
		creator:  opp.sender,
		signers:  signers,
	}

	if d, err := SetStateDocumentDataValue(opp.nds, docData); err != nil {
		return nil, err
	} else {
		sts[0] = d
	}

	return sts, nil
}

type CreateDocumentsProcessor struct {
	cp *currency.CurrencyPool
	CreateDocuments
	sb       map[currency.CurrencyID]currency.AmountState // sender StateBalance
	ns       []*CreateDocumentsItemProcessor              // ItemProcessor
	required map[currency.CurrencyID][2]currency.Big      // Fee
}

func NewCreateDocumentsProcessor(cp *currency.CurrencyPool) currency.GetNewProcessor {
	return func(op state.Processor) (state.Processor, error) {
		if i, ok := op.(CreateDocuments); !ok {
			return nil, xerrors.Errorf("not CreateDocuments, %T", op)
		} else {
			return &CreateDocumentsProcessor{
				cp:              cp,
				CreateDocuments: i,
			}, nil
		}
	}
}

func (opp *CreateDocumentsProcessor) PreProcess(
	getState func(key string) (state.State, bool, error),
	setState func(valuehash.Hash, ...state.State) error,
) (state.Processor, error) {
	fact := opp.Fact().(CreateDocumentsFact)

	// check sender account state existence
	if err := currency.CheckExistsState(currency.StateKeyAccount(fact.sender), getState); err != nil {
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

	ns := make([]*CreateDocumentsItemProcessor, len(fact.items))
	for i := range fact.items {

		c := &CreateDocumentsItemProcessor{cp: opp.cp, sender: fact.sender, h: opp.Hash(), item: fact.items[i]}
		if err := c.PreProcess(getState, setState); err != nil {
			return nil, operation.NewBaseReasonErrorFromError(err)
		}

		ns[i] = c
	}

	// check fact sign
	if err := currency.CheckFactSignsByState(fact.sender, opp.Signs(), getState); err != nil {
		return nil, operation.NewBaseReasonError("invalid signing: %w", err)
	}

	opp.ns = ns

	return opp, nil
}

func (opp *CreateDocumentsProcessor) Process( // nolint:dupl
	getState func(key string) (state.State, bool, error),
	setState func(valuehash.Hash, ...state.State) error,
) error {
	fact := opp.Fact().(CreateDocumentsFact)

	var sts []state.State // nolint:prealloc
	for i := range opp.ns {
		if s, err := opp.ns[i].Process(getState, setState); err != nil {
			return operation.NewBaseReasonError("failed to process create document item: %w", err)
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
			return nil, xerrors.Errorf("unknown currency id found, %q", it.Currency())
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

		st, err := currency.ExistsState(currency.StateKeyBalance(holder, cid), "currency of holder", getState)
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
