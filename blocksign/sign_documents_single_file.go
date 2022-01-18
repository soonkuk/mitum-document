package blocksign

import (
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util/hint"
)

var (
	SignItemSingleDocumentType   = hint.Type("mitum-blocksign-sign-item-single-document")
	SignItemSingleDocumentHint   = hint.NewHint(SignItemSingleDocumentType, "v0.0.1")
	SignItemSingleDocumentHinter = BaseSignDocumentsItem{BaseHinter: hint.NewBaseHinter(SignItemSingleDocumentHint)}
)

type SignDocumentsItemSingleFile struct {
	BaseSignDocumentsItem
}

func NewSignDocumentsItemSingleFile(docId currency.Big, owner base.Address, cid currency.CurrencyID) SignDocumentsItemSingleFile {
	return SignDocumentsItemSingleFile{
		BaseSignDocumentsItem: NewBaseSignDocumentsItem(SignItemSingleDocumentHint, docId, owner, cid),
	}
}

func (it SignDocumentsItemSingleFile) IsValid([]byte) error {
	if err := it.BaseSignDocumentsItem.IsValid(nil); err != nil {
		return err
	}

	return nil
}

func (it SignDocumentsItemSingleFile) Rebuild() SignDocumentItem {
	it.BaseSignDocumentsItem = it.BaseSignDocumentsItem.Rebuild().(BaseSignDocumentsItem)

	return it
}
