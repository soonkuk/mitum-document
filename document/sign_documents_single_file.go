package document

import (
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util/hint"
)

var (
	SignItemSingleDocumentType   = hint.Type("mitum-blocksign-sign-item-single-document")
	SignItemSingleDocumentHint   = hint.NewHint(SignItemSingleDocumentType, "v0.0.1")
	SignItemSingleDocumentHinter = SignDocumentsItemSingleDocument{
		BaseSignDocumentsItem{
			BaseHinter: hint.NewBaseHinter(SignItemSingleDocumentHint),
		},
	}
)

type SignDocumentsItemSingleDocument struct {
	BaseSignDocumentsItem
}

func NewSignDocumentsItemSingleFile(
	docID string, owner base.Address, cid currency.CurrencyID,
) SignDocumentsItemSingleDocument {
	return SignDocumentsItemSingleDocument{
		BaseSignDocumentsItem: NewBaseSignDocumentsItem(SignItemSingleDocumentHint, docID, owner, cid),
	}
}

func (it SignDocumentsItemSingleDocument) IsValid([]byte) error {
	return it.BaseSignDocumentsItem.IsValid(nil)
}

func (it SignDocumentsItemSingleDocument) Rebuild() SignDocumentsItem {
	it.BaseSignDocumentsItem = it.BaseSignDocumentsItem.Rebuild().(BaseSignDocumentsItem)

	return it
}
