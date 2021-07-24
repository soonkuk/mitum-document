package blocksign

import (
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util/encoder"
	"github.com/spikeekips/mitum/util/valuehash"
)

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
