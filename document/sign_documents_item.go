package document // nolint: dupl

import (
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
)

type BaseSignDocumentsItem struct {
	hint.BaseHinter
	id    string
	owner base.Address
	cid   currency.CurrencyID
}

func NewBaseSignDocumentsItem(
	ht hint.Hint, id string, owner base.Address, cid currency.CurrencyID,
) BaseSignDocumentsItem {
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
	_, docidtype, err := ParseDocID(it.id)
	if err != nil {
		return err
	}
	if docidtype != BSDocIDType {
		return isvalid.InvalidError.Errorf("invalid docID type: %v", docidtype)
	}

	err = it.owner.IsValid(nil)
	if err != nil {
		return err
	}

	return it.cid.IsValid(nil)
}

func (it BaseSignDocumentsItem) DocumentID() string {
	return it.id
}

func (it BaseSignDocumentsItem) Owner() base.Address {
	return it.owner
}

func (it BaseSignDocumentsItem) Currency() currency.CurrencyID {
	return it.cid
}

func (it BaseSignDocumentsItem) Rebuild() SignDocumentsItem {
	return it
}
