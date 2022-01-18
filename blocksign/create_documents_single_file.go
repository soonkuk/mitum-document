package blocksign

import (
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util/hint"
)

var (
	CreateDocumentsItemSingleFileType   = hint.Type("mitum-blocksign-create-documents-single-file")
	CreateDocumentsItemSingleFileHint   = hint.NewHint(CreateDocumentsItemSingleFileType, "v0.0.1")
	CreateDocumentsItemSingleFileHinter = BaseCreateDocumentsItem{BaseHinter: hint.NewBaseHinter(CreateDocumentsItemSingleFileHint)}
)

type CreateDocumentsItemSingleFile struct {
	BaseCreateDocumentsItem
}

func NewCreateDocumentsItemSingleFile(
	fh FileHash,
	documentid currency.Big,
	signcode, title string,
	size currency.Big,
	signers []base.Address,
	signcodes []string,
	cid currency.CurrencyID,
) CreateDocumentsItemSingleFile {
	return CreateDocumentsItemSingleFile{
		BaseCreateDocumentsItem: NewBaseCreateDocumentsItem(
			CreateDocumentsItemSingleFileHint,
			fh,
			documentid,
			signcode,
			title,
			size,
			signers,
			signcodes,
			cid,
		),
	}
}

func (it CreateDocumentsItemSingleFile) IsValid([]byte) error {
	if err := it.BaseCreateDocumentsItem.IsValid(nil); err != nil {
		return err
	}
	return nil
}
