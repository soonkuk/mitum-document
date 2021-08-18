package blocksign

import (
	"golang.org/x/xerrors"

	"github.com/soonkuk/mitum-data/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
	"github.com/spikeekips/mitum/util/valuehash"
)

var (
	SignDocumentsFactType = hint.Type("mitum-blocksign-sign-documents-operation-fact")
	SignDocumentsFactHint = hint.NewHint(SignDocumentsFactType, "v0.0.1")
	SignDocumentsType     = hint.Type("mitum-blocksign-sign-documents-operation")
	SignDocumentsHint     = hint.NewHint(SignDocumentsType, "v0.0.1")
)

var MaxSignDocumentsItems uint = 10

type SignDocumentItem interface {
	hint.Hinter
	isvalid.IsValider
	Bytes() []byte
	DocumentId() currency.Big
	Owner() base.Address
	Currency() currency.CurrencyID
	Rebuild() SignDocumentItem
}

type SignDocumentsFact struct {
	h      valuehash.Hash
	token  []byte
	sender base.Address
	items  []SignDocumentItem
}

func NewSignDocumentsFact(token []byte, sender base.Address, items []SignDocumentItem) SignDocumentsFact {
	fact := SignDocumentsFact{
		token:  token,
		sender: sender,
		items:  items,
	}
	fact.h = fact.GenerateHash()

	return fact
}

func (fact SignDocumentsFact) Hint() hint.Hint {
	return SignDocumentsFactHint
}

func (fact SignDocumentsFact) Hash() valuehash.Hash {
	return fact.h
}

func (fact SignDocumentsFact) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact SignDocumentsFact) Bytes() []byte {
	is := make([][]byte, len(fact.items))
	for i := range fact.items {
		is[i] = fact.items[i].Bytes()
	}

	return util.ConcatBytesSlice(
		fact.token,
		fact.sender.Bytes(),
		util.ConcatBytesSlice(is...),
	)
}

func (fact SignDocumentsFact) IsValid([]byte) error {
	if len(fact.token) < 1 {
		return xerrors.Errorf("empty token for SignDocumentsFact")
	} else if n := len(fact.items); n < 1 {
		return xerrors.Errorf("empty items")
	} else if n > int(MaxCreateDocumentsItems) {
		return xerrors.Errorf("items, %d over max, %d", n, MaxSignDocumentsItems)
	}

	if err := isvalid.Check([]isvalid.IsValider{
		fact.h,
		fact.sender,
	}, nil, false); err != nil {
		return err
	}

	// check duplicated document
	foundDocId := map[string]bool{}
	for i := range fact.items {
		if err := fact.items[i].IsValid(nil); err != nil {
			return err
		}
		k := fact.items[i].DocumentId().String()
		if _, found := foundDocId[k]; found {
			return xerrors.Errorf("duplicated document found, %s", k)
		}
		foundDocId[k] = true
	}

	if !fact.h.Equal(fact.GenerateHash()) {
		return isvalid.InvalidError.Errorf("wrong Fact hash")
	}

	return nil
}

func (fact SignDocumentsFact) Token() []byte {
	return fact.token
}

func (fact SignDocumentsFact) Sender() base.Address {
	return fact.sender
}

func (fact SignDocumentsFact) Items() []SignDocumentItem {
	return fact.items
}

func (fact SignDocumentsFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 1)

	as[0] = fact.Sender()

	return as, nil
}

func (fact SignDocumentsFact) Rebuild() SignDocumentsFact {
	items := make([]SignDocumentItem, len(fact.items))
	for i := range fact.items {
		it := fact.items[i]
		items[i] = it.Rebuild()
	}

	fact.items = items
	fact.h = fact.GenerateHash()

	return fact
}

type SignDocuments struct {
	operation.BaseOperation
	Memo string
}

func NewSignDocuments(fact SignDocumentsFact, fs []operation.FactSign, memo string) (SignDocuments, error) {
	if bo, err := operation.NewBaseOperationFromFact(SignDocumentsHint, fact, fs); err != nil {
		return SignDocuments{}, err
	} else {
		op := SignDocuments{BaseOperation: bo, Memo: memo}

		op.BaseOperation = bo.SetHash(op.GenerateHash())

		return op, nil
	}
}

func (op SignDocuments) Hint() hint.Hint {
	return SignDocumentsHint
}

func (op SignDocuments) IsValid(networkID []byte) error {
	if err := currency.IsValidMemo(op.Memo); err != nil {
		return err
	}

	return operation.IsValidOperation(op, networkID)
}

func (op SignDocuments) GenerateHash() valuehash.Hash {
	bs := make([][]byte, len(op.Signs())+1)
	for i := range op.Signs() {
		bs[i] = op.Signs()[i].Bytes()
	}

	bs[len(bs)-1] = []byte(op.Memo)

	e := util.ConcatBytesSlice(op.Fact().Hash().Bytes(), util.ConcatBytesSlice(bs...))

	return valuehash.NewSHA256(e)
}

func (op SignDocuments) AddFactSigns(fs ...operation.FactSign) (operation.FactSignUpdater, error) {
	if o, err := op.BaseOperation.AddFactSigns(fs...); err != nil {
		return nil, err
	} else {
		op.BaseOperation = o.(operation.BaseOperation)
	}

	op.BaseOperation = op.SetHash(op.GenerateHash())

	return op, nil
}
