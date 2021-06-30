package blocksign

import (
	"github.com/soonkuk/mitum-data/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util/hint"
)

var (
	CreateDocumentsItemSingleFileType   = hint.Type("mitum-currency-create-documents-single-file")
	CreateDocumentsItemSingleFileHint   = hint.NewHint(CreateDocumentsItemSingleFileType, "v0.0.1")
	CreateDocumentsItemSingleFileHinter = BaseCreateDocumentsItem{hint: CreateDocumentsItemSingleFileHint}
)

type CreateDocumentsItemSingleFile struct {
	BaseCreateDocumentsItem
}

func NewCreateDocumentsItemSingleFile(
	keys currency.Keys,
	sc SignCode,
	owner base.Address,
	cid currency.CurrencyID,
) CreateDocumentsItemSingleFile {
	return CreateDocumentsItemSingleFile{
		BaseCreateDocumentsItem: NewBaseCreateDocumentsItem(
			CreateDocumentsItemSingleFileHint,
			keys,
			sc,
			owner,
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
