package blocksign

import (
	"github.com/soonkuk/mitum-data/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
)

type BaseCreateDocumentsItem struct {
	hint hint.Hint
	keys currency.Keys
	doc  DocumentData
	cid  currency.CurrencyID
}

func NewBaseCreateDocumentsItem(ht hint.Hint, keys currency.Keys, doc DocumentData, cid currency.CurrencyID) BaseCreateDocumentsItem {
	return BaseCreateDocumentsItem{
		hint: ht,
		keys: keys,
		doc:  doc,
		cid:  cid,
	}
}

func (it BaseCreateDocumentsItem) Hint() hint.Hint {
	return it.hint
}

func (it BaseCreateDocumentsItem) Bytes() []byte {
	bs := make([][]byte, 3)
	bs[0] = it.keys.Bytes()
	bs[1] = it.doc.Bytes()
	bs[2] = it.cid.Bytes()

	return util.ConcatBytesSlice(bs...)
}

func (it BaseCreateDocumentsItem) IsValid([]byte) error {

	// empty key, duplicated key, threshold check
	if err := it.keys.IsValid(nil); err != nil {
		return err
	}

	if err := it.doc.IsValid(nil); err != nil {
		return err
	}

	return nil
}

// Keys return BaseCreateDocumentsItem's keys.
func (it BaseCreateDocumentsItem) Keys() currency.Keys {
	return it.keys
}

// Address get address from BaseCreateDocumentsItem's keys and return it.
func (it BaseCreateDocumentsItem) Address() (base.Address, error) {
	return currency.NewAddressFromKeys(it.keys)
}

// FileHash return BaseCreateDocumetsItem's owner address.
func (it BaseCreateDocumentsItem) DocumentData() DocumentData {
	return it.doc
}

// FileData return BaseCreateDocumentsItem's fileData.
func (it BaseCreateDocumentsItem) Currency() currency.CurrencyID {
	return it.cid
}

func (it BaseCreateDocumentsItem) Rebuild() CreateDocumentsItem {

	return it
}
