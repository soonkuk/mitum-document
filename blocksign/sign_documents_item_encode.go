package blocksign

import (
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util/encoder"
	"github.com/spikeekips/mitum/util/hint"
)

func (it *BaseSignDocumentsItem) unpack(
	enc encoder.Encoder,
	ht hint.Hint,
	di currency.Big,
	ow base.AddressDecoder,
	scid string,

) error {
	it.hint = ht

	it.id = di

	a, err := ow.Encode(enc)
	if err != nil {
		return err
	}
	it.owner = a
	it.cid = currency.CurrencyID(scid)

	return nil
}
