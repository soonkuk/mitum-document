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
	TransferDocumentsFactType = hint.Type("mitum-blocksign-transfer-documents-operation-fact")
	TransferDocumentsFactHint = hint.NewHint(TransferDocumentsFactType, "v0.0.1")
	TransferDocumentsType     = hint.Type("mitum-blocksign-transfer-documents-operation")
	TransferDocumentsHint     = hint.NewHint(TransferDocumentsType, "v0.0.1")
)

var MaxTransferDocumentsItems uint = 10

type TransferDocumentsItem interface {
	hint.Hinter
	isvalid.IsValider
	Bytes() []byte
	DocumentId() currency.Big
	Owner() base.Address
	Receiver() base.Address
	Currency() currency.CurrencyID
	Rebuild() TransferDocumentsItem
}

type TransferDocumentsFact struct {
	h      valuehash.Hash
	token  []byte
	sender base.Address
	items  []TransferDocumentsItem
}

func NewTransferDocumentsFact(token []byte, sender base.Address, items []TransferDocumentsItem) TransferDocumentsFact {
	fact := TransferDocumentsFact{
		token:  token,
		sender: sender,
		items:  items,
	}
	fact.h = fact.GenerateHash()

	return fact
}

func (fact TransferDocumentsFact) Hint() hint.Hint {
	return TransferDocumentsFactHint
}

func (fact TransferDocumentsFact) Hash() valuehash.Hash {
	return fact.h
}

func (fact TransferDocumentsFact) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact TransferDocumentsFact) Bytes() []byte {
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

func (fact TransferDocumentsFact) IsValid([]byte) error {
	if len(fact.token) < 1 {
		return xerrors.Errorf("empty token for TransferDocumentsFact")
	} else if n := len(fact.items); n < 1 {
		return xerrors.Errorf("empty items")
	} else if n > int(MaxTransferDocumentsItems) {
		return xerrors.Errorf("items, %d over max, %d", n, MaxTransferDocumentsItems)
	}

	if err := isvalid.Check([]isvalid.IsValider{
		fact.h,
		fact.sender,
	}, nil, false); err != nil {
		return err
	}

	// check receiver same with sender and duplicated document id
	foundDocId := map[string]bool{}
	for i := range fact.items {
		it := fact.items[i]
		if err := it.IsValid(nil); err != nil {
			return err
		}
		r := it.Receiver().String()
		if r == fact.sender.String() {
			return xerrors.Errorf("receiver is same with sender, %q", fact.sender)
		}
		k := it.DocumentId().String()
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

func (fact TransferDocumentsFact) Token() []byte {
	return fact.token
}

func (fact TransferDocumentsFact) Sender() base.Address {
	return fact.sender
}

func (fact TransferDocumentsFact) Items() []TransferDocumentsItem {
	return fact.items
}

func (fact TransferDocumentsFact) Targets() ([]base.Address, error) {
	as := make([]base.Address, len(fact.items))
	for i := range fact.items {
		if a, err := currency.NewAddress(fact.items[i].Receiver().String()); err != nil {
			return nil, err
		} else {
			as[i] = a
		}
	}

	return as, nil
}

func (fact TransferDocumentsFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, len(fact.items)+1)

	if tas, err := fact.Targets(); err != nil {
		return nil, err
	} else {
		copy(as, tas)
	}

	as[len(fact.items)] = fact.Sender()

	return as, nil
}

func (fact TransferDocumentsFact) Rebuild() TransferDocumentsFact {
	items := make([]TransferDocumentsItem, len(fact.items))
	for i := range fact.items {
		it := fact.items[i]
		items[i] = it.Rebuild()
	}

	fact.items = items
	fact.h = fact.GenerateHash()

	return fact
}

type TransferDocuments struct {
	operation.BaseOperation
	Memo string
}

func NewTransferDocuments(fact TransferDocumentsFact, fs []operation.FactSign, memo string) (TransferDocuments, error) {
	if bo, err := operation.NewBaseOperationFromFact(TransferDocumentsHint, fact, fs); err != nil {
		return TransferDocuments{}, err
	} else {
		op := TransferDocuments{BaseOperation: bo, Memo: memo}

		op.BaseOperation = bo.SetHash(op.GenerateHash())

		return op, nil
	}
}

func (op TransferDocuments) Hint() hint.Hint {
	return TransferDocumentsHint
}

func (op TransferDocuments) IsValid(networkID []byte) error {
	if err := currency.IsValidMemo(op.Memo); err != nil {
		return err
	}

	return operation.IsValidOperation(op, networkID)
}

func (op TransferDocuments) GenerateHash() valuehash.Hash {
	bs := make([][]byte, len(op.Signs())+1)
	for i := range op.Signs() {
		bs[i] = op.Signs()[i].Bytes()
	}

	bs[len(bs)-1] = []byte(op.Memo)

	e := util.ConcatBytesSlice(op.Fact().Hash().Bytes(), util.ConcatBytesSlice(bs...))

	return valuehash.NewSHA256(e)
}

func (op TransferDocuments) AddFactSigns(fs ...operation.FactSign) (operation.FactSignUpdater, error) {
	if o, err := op.BaseOperation.AddFactSigns(fs...); err != nil {
		return nil, err
	} else {
		op.BaseOperation = o.(operation.BaseOperation)
	}

	op.BaseOperation = op.SetHash(op.GenerateHash())

	return op, nil
}
