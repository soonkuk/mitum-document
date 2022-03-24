package document

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
	SignDocumentsFactType   = hint.Type("mitum-blocksign-sign-documents-operation-fact")
	SignDocumentsFactHint   = hint.NewHint(SignDocumentsFactType, "v0.0.1")
	SignDocumentsFactHinter = SignDocumentsFact{BaseHinter: hint.NewBaseHinter(SignDocumentsFactHint)}
	SignDocumentsType       = hint.Type("mitum-blocksign-sign-documents-operation")
	SignDocumentsHint       = hint.NewHint(SignDocumentsType, "v0.0.1")
	SignDocumentsHinter     = SignDocuments{BaseOperation: operationHinter(SignDocumentsHint)}
)

var MaxSignDocumentsItems uint = 10

type SignDocumentsItem interface {
	hint.Hinter
	isvalid.IsValider
	Bytes() []byte
	DocumentID() string
	Owner() base.Address
	Currency() currency.CurrencyID
	Rebuild() SignDocumentsItem
}

type SignDocumentsFact struct {
	hint.BaseHinter
	h      valuehash.Hash
	token  []byte
	sender base.Address
	items  []SignDocumentsItem
}

func NewSignDocumentsFact(token []byte, sender base.Address, items []SignDocumentsItem) SignDocumentsFact {
	fact := SignDocumentsFact{
		BaseHinter: hint.NewBaseHinter(SignDocumentsFactHint),
		token:      token,
		sender:     sender,
		items:      items,
	}
	fact.h = fact.GenerateHash()

	return fact
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

func (fact SignDocumentsFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := currency.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if n := len(fact.items); n < 1 {
		return errors.Errorf("empty items")
	} else if n > int(MaxCreateDocumentsItems) {
		return errors.Errorf("items, %d over max, %d", n, MaxSignDocumentsItems)
	}

	if err := isvalid.Check(
		nil, false, fact.h,
		fact.sender); err != nil {
		return err
	}

	// check duplicated document
	foundDocID := map[string]bool{}
	for i := range fact.items {
		if err := fact.items[i].IsValid(nil); err != nil {
			return err
		}
		k := fact.items[i].DocumentID()
		if _, found := foundDocID[k]; found {
			return errors.Errorf("duplicated documentID, %s", k)
		}
		foundDocID[k] = true
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

func (fact SignDocumentsFact) Items() []SignDocumentsItem {
	return fact.items
}

func (fact SignDocumentsFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, 1)

	as[0] = fact.Sender()

	return as, nil
}

func (fact SignDocumentsFact) Rebuild() SignDocumentsFact {
	items := make([]SignDocumentsItem, len(fact.items))
	for i := range fact.items {
		it := fact.items[i]
		items[i] = it.Rebuild()
	}

	fact.items = items
	fact.h = fact.GenerateHash()

	return fact
}

type SignDocuments struct {
	currency.BaseOperation
}

func NewSignDocuments(fact SignDocumentsFact, fs []base.FactSign, memo string) (SignDocuments, error) {
	bo, err := currency.NewBaseOperationFromFact(SignDocumentsHint, fact, fs, memo)
	if err != nil {
		return SignDocuments{}, err
	}
	return SignDocuments{BaseOperation: bo}, nil
}
