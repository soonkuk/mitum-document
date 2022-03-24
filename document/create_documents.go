package document // nolint: dupl, revive

import (
	"github.com/pkg/errors"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
	"github.com/spikeekips/mitum/util/valuehash"
)

var (
	CreateDocumentsFactType   = hint.Type("mitum-create-documents-operation-fact")
	CreateDocumentsFactHint   = hint.NewHint(CreateDocumentsFactType, "v0.0.1")
	CreateDocumentsFactHinter = CreateDocumentsFact{BaseHinter: hint.NewBaseHinter(CreateDocumentsFactHint)}
	CreateDocumentsType       = hint.Type("mitum-create-documents-operation")
	CreateDocumentsHint       = hint.NewHint(CreateDocumentsType, "v0.0.1")
	CreateDocumentsHinter     = CreateDocuments{BaseOperation: operationHinter(CreateDocumentsHint)}
)

var MaxCreateDocumentsItems uint = 10

type CreateDocumentsItem interface {
	hint.Hinter
	isvalid.IsValider
	Bytes() []byte
	DocumentID() string
	// DocType() hint.Type
	Doc() DocumentData
	Currency() currency.CurrencyID
	Rebuild() CreateDocumentsItem
}

type CreateDocumentsFact struct {
	hint.BaseHinter
	h      valuehash.Hash
	token  []byte
	sender base.Address
	items  []CreateDocumentsItem
}

func NewCreateDocumentsFact(token []byte, sender base.Address, items []CreateDocumentsItem) CreateDocumentsFact {
	fact := CreateDocumentsFact{
		BaseHinter: hint.NewBaseHinter(CreateDocumentsFactHint),
		token:      token,
		sender:     sender,
		items:      items,
	}
	fact.h = fact.GenerateHash()

	return fact
}

func (fact CreateDocumentsFact) Hash() valuehash.Hash {
	return fact.h
}

func (fact CreateDocumentsFact) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact CreateDocumentsFact) Bytes() []byte {
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

func (fact CreateDocumentsFact) IsValid(b []byte) error { // nolint:dupl
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := currency.IsValidOperationFact(fact, b); err != nil {
		return err
	}
	if n := len(fact.items); n < 1 {
		return errors.Errorf("empty items")
	} else if n > int(MaxCreateDocumentsItems) {
		return errors.Errorf("items, %d over max, %d", n, MaxCreateDocumentsItems)
	}
	if err := isvalid.Check(nil, false, fact.sender); err != nil {
		return err
	}

	docIDMap := map[string]bool{}
	for i := range fact.items {
		if err := isvalid.Check(nil, false, fact.items[i]); err != nil {
			return err
		}

		it := fact.items[i]
		k := it.Doc().DocumentID()
		if _, found := docIDMap[k]; found {
			return errors.Errorf("duplicated documentID, %s", k)
		}
		docIDMap[k] = true
	}

	if !fact.h.Equal(fact.GenerateHash()) {
		return isvalid.InvalidError.Errorf("wrong Fact hash")
	}

	return nil
}

func (fact CreateDocumentsFact) Token() []byte {
	return fact.token
}

func (fact CreateDocumentsFact) Sender() base.Address {
	return fact.sender
}

func (fact CreateDocumentsFact) Items() []CreateDocumentsItem {
	return fact.items
}

func (fact CreateDocumentsFact) Addresses() ([]base.Address, error) {
	var as []base.Address

	as = append(as, fact.Sender())

	return as, nil
}

func (fact CreateDocumentsFact) Rebuild() CreateDocumentsFact {
	items := make([]CreateDocumentsItem, len(fact.items))
	for i := range fact.items {
		it := fact.items[i]
		items[i] = it.Rebuild()
	}

	fact.items = items
	fact.h = fact.GenerateHash()

	return fact
}

type CreateDocuments struct {
	currency.BaseOperation
}

func NewCreateDocuments(fact CreateDocumentsFact, fs []base.FactSign, memo string) (CreateDocuments, error) {
	bo, err := currency.NewBaseOperationFromFact(CreateDocumentsHint, fact, fs, memo)
	if err != nil {
		return CreateDocuments{}, err
	}
	return CreateDocuments{BaseOperation: bo}, nil
}
