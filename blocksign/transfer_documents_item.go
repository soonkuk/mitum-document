package blocksign

import (
	"github.com/soonkuk/mitum-data/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
)

type BaseTransferDocumentsItem struct {
	hint       hint.Hint
	documentId DocId        // document address
	owner      base.Address // document owner
	receiver   base.Address // document receiver
	cid        currency.CurrencyID
}

func NewBaseTransferDocumentsItem(ht hint.Hint, documentId DocId, owner base.Address, receiver base.Address, cid currency.CurrencyID) BaseTransferDocumentsItem {
	return BaseTransferDocumentsItem{
		hint:       ht,
		documentId: documentId,
		owner:      owner,
		receiver:   receiver,
		cid:        cid,
	}
}

func (it BaseTransferDocumentsItem) Hint() hint.Hint {
	return it.hint
}

func (it BaseTransferDocumentsItem) Bytes() []byte {
	bs := make([][]byte, 4)
	bs[0] = it.documentId.Bytes()
	bs[1] = it.owner.Bytes()
	bs[2] = it.receiver.Bytes()
	bs[3] = it.cid.Bytes()

	return util.ConcatBytesSlice(bs...)
}

func (it BaseTransferDocumentsItem) IsValid([]byte) error {
	if err := it.documentId.IsValid(nil); err != nil {
		return err
	} else if err := it.owner.IsValid(nil); err != nil {
		return err
	} else if err := it.receiver.IsValid(nil); err != nil {
		return err
	} else if err := it.cid.IsValid(nil); err != nil {
		return err
	}

	// TODO : empty check
	/*
		if n := len(it.amounts); n == 0 {
			return xerrors.Errorf("empty amounts")
		}
	*/

	return nil
}

func (it BaseTransferDocumentsItem) DocumentId() DocId {
	return it.documentId
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
