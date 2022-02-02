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
	UpdateDocumentsFactType   = hint.Type("mitum-blockcity-update-documents-operation-fact")
	UpdateDocumentsFactHint   = hint.NewHint(UpdateDocumentsFactType, "v0.0.1")
	UpdateDocumentsFactHinter = UpdateDocumentsFact{BaseHinter: hint.NewBaseHinter(UpdateDocumentsFactHint)}
	UpdateDocumentsType       = hint.Type("mitum-blockcity-update-documents-operation")
	UpdateDocumentsHint       = hint.NewHint(UpdateDocumentsType, "v0.0.1")
	UpdateDocumentsHinter     = UpdateDocuments{BaseOperation: operationHinter(UpdateDocumentsHint)}
)

var MaxUpdateDocumentsItems uint = 10

type UpdateDocumentsItem interface {
	hint.Hinter
	isvalid.IsValider
	Bytes() []byte
	DocumentId() string
	DocType() hint.Type
	Doc() Document
	Currency() currency.CurrencyID
	Rebuild() UpdateDocumentsItem
}

type UpdateDocumentsFact struct {
	hint.BaseHinter
	h      valuehash.Hash
	token  []byte
	sender base.Address
	items  []UpdateDocumentsItem
}

func NewUpdateDocumentsFact(token []byte, sender base.Address, items []UpdateDocumentsItem) UpdateDocumentsFact {
	fact := UpdateDocumentsFact{
		BaseHinter: hint.NewBaseHinter(UpdateDocumentsFactHint),
		token:      token,
		sender:     sender,
		items:      items,
	}
	fact.h = fact.GenerateHash()

	return fact
}

func (fact UpdateDocumentsFact) Hash() valuehash.Hash {
	return fact.h
}

func (fact UpdateDocumentsFact) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact UpdateDocumentsFact) Bytes() []byte {
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

func (fact UpdateDocumentsFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := currency.IsValidOperationFact(fact, b); err != nil {
		return err
	}
	if len(fact.token) < 1 {
		return errors.Errorf("empty token for UpdateDocumentsFact")
	} else if n := len(fact.items); n < 1 {
		return errors.Errorf("empty items")
	} else if n > int(MaxUpdateDocumentsItems) {
		return errors.Errorf("items, %d over max, %d", n, MaxUpdateDocumentsItems)
	}

	if err := isvalid.Check(nil, false, fact.sender); err != nil {
		return err
	}

	docIdMap := map[string]bool{}
	for i := range fact.items {
		if err := isvalid.Check(nil, false, fact.items[i]); err != nil {
			return err
		}

		it := fact.items[i]
		k := it.Doc().DocumentId()
		if _, found := docIdMap[k]; found {
			return errors.Errorf("duplicated document user Id, %s", k)
		}
		docIdMap[k] = true
	}

	if !fact.h.Equal(fact.GenerateHash()) {
		return isvalid.InvalidError.Errorf("wrong Fact hash")
	}

	return nil
}

func (fact UpdateDocumentsFact) Token() []byte {
	return fact.token
}

func (fact UpdateDocumentsFact) Sender() base.Address {
	return fact.sender
}

func (fact UpdateDocumentsFact) Items() []UpdateDocumentsItem {
	return fact.items
}

func (fact UpdateDocumentsFact) Rebuild() UpdateDocumentsFact {
	items := make([]UpdateDocumentsItem, len(fact.items))
	for i := range fact.items {
		it := fact.items[i]
		items[i] = it.Rebuild()
	}

	fact.items = items
	fact.h = fact.GenerateHash()

	return fact
}

type UpdateDocuments struct {
	currency.BaseOperation
}

func NewUpdateDocuments(fact UpdateDocumentsFact, fs []base.FactSign, memo string) (UpdateDocuments, error) {
	bo, err := currency.NewBaseOperationFromFact(UpdateDocumentsHint, fact, fs, memo)
	if err != nil {
		return UpdateDocuments{}, err
	} else {

		return UpdateDocuments{BaseOperation: bo}, nil
	}
}
