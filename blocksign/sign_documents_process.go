package blocksign

import (
	"github.com/pkg/errors"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/base/state"
	"github.com/spikeekips/mitum/util/valuehash"
)

func (op SignDocuments) Process(
	func(key string) (state.State, bool, error),
	func(valuehash.Hash, ...state.State) error,
) error {
	return nil
}

type SignDocumentsItemProcessor struct {
	cp     *currency.CurrencyPool
	sender base.Address
	h      valuehash.Hash
	item   SignDocumentItem
	nds    state.State       // new document data state (key = document filehash)
	dinv   DocumentInventory // document inventory
	ndinvs state.State       // document inventory state (key = owner address)

}

func (opp *SignDocumentsItemProcessor) PreProcess(
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
		return errors.Errorf("owner does not exist, %q", opp.item.Owner())
	}

	// check existence of document inventory state with owner address
	switch st, found, err := getState(StateKeyDocuments(opp.item.Owner())); {
	case err != nil:
		return err
	case !found:
		return errors.Errorf("Owner has no document inventory, %v", opp.item.Owner())
	default:
		dinv, err := StateDocumentsValue(st)
		if err != nil {
			return err
		}
		opp.dinv = dinv
		opp.ndinvs = st
	}

	docinfo, err := opp.dinv.Get(opp.item.DocumentId())
	if err != nil {
		return err
	}

	// check existence of new document state with filehash
	switch st, found, err := getState(StateKeyDocumentData(docinfo.FileHash())); {
	case err != nil:
		return err
	case !found:
		return errors.Errorf("document not registered with filehash, %q", docinfo.FileHash())
	default:
		opp.nds = st
	}

	dd, err := StateDocumentDataValue(opp.nds)
	if err != nil {
		return err
	}

	if !dd.Owner().Equal(opp.item.Owner()) {
		return errors.Errorf("Owner not matched with owner in document, %v", opp.item.Owner())
	}

	if len(dd.Signers()) < 1 {
		return errors.Errorf("sender not found in document Signers, %v", opp.sender)
	}
	// check signer exist in document data signers
	for i := range dd.Signers() {
		if dd.Signers()[i].Address().Equal(opp.sender) {
			dd.Signers()[i].SetSigned()
			break
		}
		if i == (len(dd.Signers()) - 1) {
			return errors.Errorf("sender not found in document Signers, %v", opp.sender)
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

type SignDocumentsProcessor struct {
	cp *currency.CurrencyPool
	SignDocuments
	sb       map[currency.CurrencyID]currency.AmountState // sender StateBalance
	ns       []*SignDocumentsItemProcessor                // ItemProcessor
	required map[currency.CurrencyID][2]currency.Big      // Fee
}

func NewSignDocumentsProcessor(cp *currency.CurrencyPool) currency.GetNewProcessor {
	return func(op state.Processor) (state.Processor, error) {
		if i, ok := op.(SignDocuments); !ok {
			return nil, errors.Errorf("not SignDocuments, %T", op)
		} else {
			return &SignDocumentsProcessor{
				cp:            cp,
				SignDocuments: i,
			}, nil
		}
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

func (opp *SignDocumentsProcessor) calculateItemsFee() (map[currency.CurrencyID][2]currency.Big, error) {
	fact := opp.Fact().(SignDocumentsFact)
	items := make([]SignDocumentItem, len(fact.items))
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
