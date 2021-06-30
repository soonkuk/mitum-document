package blocksign

import (
	"golang.org/x/xerrors"

	"github.com/soonkuk/mitum-data/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util/encoder"
	"github.com/spikeekips/mitum/util/hint"
)

func (it *BaseCreateDocumentsItem) unpack(
	enc encoder.Encoder,
	ht hint.Hint,
	bks []byte,
	bsc string,
	bOwner base.AddressDecoder,
	scid string,

) error {
	it.hint = ht

	if hinter, err := enc.Decode(bks); err != nil {
		return err
	} else if k, ok := hinter.(currency.Keys); !ok {
		return xerrors.Errorf("not Keys: %T", hinter)
	} else {
		it.keys = k
	}
	a, err := bOwner.Encode(enc)
	if err != nil {
		return err
	}
	it.owner = a
	it.sc = SignCode(bsc)
	it.cid = currency.CurrencyID(scid)

	return nil
}
