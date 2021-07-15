package blocksign

import (
	"github.com/soonkuk/mitum-data/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util/encoder"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/valuehash"
	"golang.org/x/xerrors"
)

func (it *BaseTransferDocumentsItem) unpack(
	enc encoder.Encoder,
	ht hint.Hint,
	document []byte,
	owner base.AddressDecoder,
	receiver base.AddressDecoder,
	cid string,
) error {

	if hinter, err := enc.Decode(document); err != nil {
		return err
	} else if d, ok := hinter.(DocId); !ok {
		return xerrors.Errorf("not Document Id : %T", d)
	} else {
		it.documentId = d
	}

	o, err := owner.Encode(enc)
	if err != nil {
		return err
	}
	it.owner = o

	r, err := receiver.Encode(enc)
	if err != nil {
		return err
	}
	it.receiver = r

	it.hint = ht
	it.cid = currency.CurrencyID(cid)

	return nil
}

func (fact *TransferDocumentsFact) unpack(
	enc encoder.Encoder,
	h valuehash.Hash,
	token []byte,
	bSender base.AddressDecoder,
	bitems [][]byte,
) error {
	a, err := bSender.Encode(enc)
	if err != nil {
		return err
	}
	fact.sender = a

	items := make([]TransferDocumentsItem, len(bitems))
	for i := range bitems {
		if j, err := DecodeTransferDocumentsItem(enc, bitems[i]); err != nil {
			return err
		} else {
			items[i] = j
		}
	}

	fact.h = h
	fact.token = token
	fact.items = items

	return nil
}
