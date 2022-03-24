package document // nolint:dupl

import (
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
)

var (
	UpdateDocumentsItemImplType   = hint.Type("mitum-update-documents-item")
	UpdateDocumentsItemImplHint   = hint.NewHint(UpdateDocumentsItemImplType, "v0.0.1")
	UpdateDocumentsItemImplHinter = UpdateDocumentsItemImpl{BaseHinter: hint.NewBaseHinter(UpdateDocumentsItemImplHint)}
)

type UpdateDocumentsItemImpl struct {
	hint.BaseHinter
	doc DocumentData
	cid currency.CurrencyID
}

func NewUpdateDocumentsItemImpl(
	doc DocumentData,
	cid currency.CurrencyID) UpdateDocumentsItemImpl {
	return UpdateDocumentsItemImpl{
		BaseHinter: hint.NewBaseHinter(UpdateDocumentsItemImplHint),
		doc:        doc,
		cid:        cid,
	}
}

func (it UpdateDocumentsItemImpl) Bytes() []byte {
	bs := make([][]byte, 2)
	bs[0] = it.doc.Bytes()
	bs[1] = it.cid.Bytes()

	return util.ConcatBytesSlice(bs...)
}

func (it UpdateDocumentsItemImpl) IsValid([]byte) error {
	if err := isvalid.Check(
		nil, false,
		it.BaseHinter,
		it.doc,
		it.cid,
	); err != nil {
		return isvalid.InvalidError.Errorf("invalid UpdateDocumentsItem: %w", err)
	}
	return nil
}

func (it UpdateDocumentsItemImpl) DocumentID() string {
	return it.doc.DocumentID()
}

func (it UpdateDocumentsItemImpl) Doc() DocumentData {
	return it.doc
}

func (it UpdateDocumentsItemImpl) Currency() currency.CurrencyID {
	return it.cid
}

func (it UpdateDocumentsItemImpl) Rebuild() UpdateDocumentsItem {
	return it
}
