package blocksign

import (
	"github.com/soonkuk/mitum-data/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util/hint"
)

var (
	TransfersItemSingleDocumentType   = hint.Type("mitum-blocksign-transfer-item-single-document")
	TransfersItemSingleDocumentHint   = hint.NewHint(TransfersItemSingleDocumentType, "v0.0.1")
	TransfersItemSingleDocumentHinter = BaseTransferDocumentsItem{hint: TransfersItemSingleDocumentHint}
)

type TransferDocumentsItemSingleFile struct {
	BaseTransferDocumentsItem
}

func NewTransferDocumentsItemSingleFile(docId currency.Big, owner base.Address, receiver base.Address, cid currency.CurrencyID) TransferDocumentsItemSingleFile {
	return TransferDocumentsItemSingleFile{
		BaseTransferDocumentsItem: NewBaseTransferDocumentsItem(TransfersItemSingleDocumentHint, docId, owner, receiver, cid),
	}
}

func (it TransferDocumentsItemSingleFile) IsValid([]byte) error {
	if err := it.BaseTransferDocumentsItem.IsValid(nil); err != nil {
		return err
	}

	return nil
}

func (it TransferDocumentsItemSingleFile) Rebuild() TransferDocumentsItem {
	it.BaseTransferDocumentsItem = it.BaseTransferDocumentsItem.Rebuild().(BaseTransferDocumentsItem)

	return it
}
