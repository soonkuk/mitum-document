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

type SignDocumentItem interface {
	hint.Hinter
	isvalid.IsValider
	Bytes() []byte
	DocumentId() string
	Owner() base.Address
	Currency() currency.CurrencyID
	Rebuild() SignDocumentItem
}

type SignDocumentsFact struct {
	hint.BaseHinter
	h      valuehash.Hash
	token  []byte
	sender base.Address
	items  []SignDocumentItem
}

func NewSignDocumentsFact(token []byte, sender base.Address, items []SignDocumentItem) SignDocumentsFact {
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

	if len(fact.token) < 1 {
		return errors.Errorf("empty token for SignDocumentsFact")
	} else if n := len(fact.items); n < 1 {
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
	foundDocId := map[string]bool{}
	for i := range fact.items {
		if err := fact.items[i].IsValid(nil); err != nil {
			return err
		}
		k := fact.items[i].DocumentId()
		if _, found := foundDocId[k]; found {
			return errors.Errorf("duplicated document found, %s", k)
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
	currency.BaseOperation
}

func NewSignDocuments(fact SignDocumentsFact, fs []base.FactSign, memo string) (SignDocuments, error) {
	bo, err := currency.NewBaseOperationFromFact(SignDocumentsHint, fact, fs, memo)
	if err != nil {
		return SignDocuments{}, err
	}
	return SignDocuments{BaseOperation: bo}, nil
}
