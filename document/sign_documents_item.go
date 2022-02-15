package document

import (
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
)

type BaseSignDocumentsItem struct {
	hint.BaseHinter
	id    string
	owner base.Address
	cid   currency.CurrencyID
}

func NewBaseSignDocumentsItem(ht hint.Hint, id string, owner base.Address, cid currency.CurrencyID) BaseSignDocumentsItem {
	return BaseSignDocumentsItem{
		BaseHinter: hint.NewBaseHinter(ht),
		id:         id,
		owner:      owner,
		cid:        cid,
	}
}

func (it BaseSignDocumentsItem) Bytes() []byte {
	bs := make([][]byte, 3)
	bs[0] = []byte(it.id)
	bs[1] = it.owner.Bytes()
	bs[2] = it.cid.Bytes()

	return util.ConcatBytesSlice(bs...)
}

func (it BaseSignDocumentsItem) IsValid([]byte) error {

	if err := it.owner.IsValid(nil); err != nil {
		return err
	}

	if err := it.cid.IsValid(nil); err != nil {
		return err
	}

	return nil
}

// FileHash return BaseCreateDocumetsItem's owner address.
func (it BaseSignDocumentsItem) DocumentId() string {
	return it.id
}

func (it BaseSignDocumentsItem) Owner() base.Address {
	return it.owner
}

// FileData return BaseCreateDocumentsItem's fileData.
func (it BaseSignDocumentsItem) Currency() currency.CurrencyID {
	return it.cid
}

func (it BaseSignDocumentsItem) Rebuild() SignDocumentItem {
	return it
}
