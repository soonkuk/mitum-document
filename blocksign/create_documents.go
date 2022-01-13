package blocksign

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
	CreateDocumentsFactType   = hint.Type("mitum-blocksign-create-documents-operation-fact")
	CreateDocumentsFactHint   = hint.NewHint(CreateDocumentsFactType, "v0.0.1")
	CreateDocumentsFactHinter = CreateDocumentsFact{BaseHinter: hint.NewBaseHinter(CreateDocumentsFactHint)}
	CreateDocumentsType       = hint.Type("mitum-blocksign-create-documents-operation")
	CreateDocumentsHint       = hint.NewHint(CreateDocumentsType, "v0.0.1")
	CreateDocumentsHinter     = CreateDocuments{BaseOperation: operationHinter(CreateDocumentsHint)}
)

var MaxCreateDocumentsItems uint = 10

type CreateDocumentsItem interface {
	hint.Hinter
	isvalid.IsValider
	Bytes() []byte
	FileHash() FileHash
	DocumentId() currency.Big
	Signcode() string
	Title() string
	Size() currency.Big
	Signers() []base.Address
	Signcodes() []string
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

func (fact CreateDocumentsFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := currency.IsValidOperationFact(fact, b); err != nil {
		return err
	}
	if len(fact.token) < 1 {
		return errors.Errorf("empty token for CreateDocumentsFact")
	} else if n := len(fact.items); n < 1 {
		return errors.Errorf("empty items")
	} else if n > int(MaxCreateDocumentsItems) {
		return errors.Errorf("items, %d over max, %d", n, MaxCreateDocumentsItems)
	}

	if err := isvalid.Check(nil, false, fact.sender); err != nil {
		return err
	}

	fhmap := map[string]bool{}
	for i := range fact.items {
		if err := isvalid.Check(nil, false, fact.items[i]); err != nil {
			return err
		}

		it := fact.items[i]
		k := it.FileHash().String()
		if _, found := fhmap[k]; found {
			return errors.Errorf("duplicated filehash, %s", k)
		}
		fhmap[fact.items[i].FileHash().String()] = true
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

func (fact CreateDocumentsFact) Signers() []base.Address {
	var as []base.Address
	for i := range fact.items {
		a := fact.items[i].Signers()
		if len(a) > 0 {
			as = append(as, a...)
		}
	}

	return as
}

func (fact CreateDocumentsFact) Addresses() ([]base.Address, error) {
	var as []base.Address

	signers := fact.Signers()
	if len(signers) > 0 {
		copy(as, signers)
	}

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
	} else {

		return CreateDocuments{BaseOperation: bo}, nil
	}
}
