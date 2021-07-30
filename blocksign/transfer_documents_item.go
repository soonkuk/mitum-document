package blocksign

import (
	"github.com/soonkuk/mitum-data/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
)

type BaseTransferDocumentsItem struct {
	hint     hint.Hint
	docId    currency.Big // document id
	owner    base.Address // document owner
	receiver base.Address // document receiver
	cid      currency.CurrencyID
}

func NewBaseTransferDocumentsItem(ht hint.Hint, docId currency.Big, owner base.Address, receiver base.Address, cid currency.CurrencyID) BaseTransferDocumentsItem {
	return BaseTransferDocumentsItem{
		hint:     ht,
		docId:    docId,
		owner:    owner,
		receiver: receiver,
		cid:      cid,
	}
}

func (it BaseTransferDocumentsItem) Hint() hint.Hint {
	return it.hint
}

func (it BaseTransferDocumentsItem) Bytes() []byte {
	bs := make([][]byte, 4)
	bs[0] = it.docId.Bytes()
	bs[1] = it.owner.Bytes()
	bs[2] = it.receiver.Bytes()
	bs[3] = it.cid.Bytes()

	return util.ConcatBytesSlice(bs...)
}

func (it BaseTransferDocumentsItem) IsValid([]byte) error {
	if err := it.docId.IsValid(nil); err != nil {
		return err
	} else if err := it.owner.IsValid(nil); err != nil {
		return err
	} else if err := it.receiver.IsValid(nil); err != nil {
		return err
	} else if err := it.cid.IsValid(nil); err != nil {
		return err
	}

	return nil
}

func (it BaseTransferDocumentsItem) DocumentId() currency.Big {
	return it.docId
}

func (it BaseTransferDocumentsItem) Owner() base.Address {
	return it.owner
}

func (it BaseTransferDocumentsItem) Receiver() base.Address {
	return it.receiver
}

func (it BaseTransferDocumentsItem) Currency() currency.CurrencyID {
	return it.cid
}

func (it BaseTransferDocumentsItem) Rebuild() TransferDocumentsItem {

	return it
}
