package blocksign

import (
	"github.com/soonkuk/mitum-data/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util/encoder"
	"github.com/spikeekips/mitum/util/hint"
)

func (it *BaseTransferDocumentsItem) unpack(
	enc encoder.Encoder,
	ht hint.Hint,
	docId currency.Big,
	owner base.AddressDecoder,
	receiver base.AddressDecoder,
	cid string,
) error {

	it.docId = docId

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
