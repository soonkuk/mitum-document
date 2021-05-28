package currency

import (
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util/encoder"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/valuehash"
)

func (it *BaseTransferDocumentsItem) unpack(
	enc encoder.Encoder,
	ht hint.Hint,
	sdocument base.AddressDecoder,
	sreceiver base.AddressDecoder,
	scid string,
) error {

	a, err := sdocument.Encode(enc)
	if err != nil {
		return err
	}
	it.document = a
	b, err := sreceiver.Encode(enc)
	if err != nil {
		return err
	}
	it.receiver = b

	it.hint = ht
	it.cid = CurrencyID(scid)

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
