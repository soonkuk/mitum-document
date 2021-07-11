package blocksign

import (
	"golang.org/x/xerrors"

	"github.com/soonkuk/mitum-data/currency"
	"github.com/spikeekips/mitum/util/encoder"
	"github.com/spikeekips/mitum/util/hint"
)

func (it *BaseCreateDocumentsItem) unpack(
	enc encoder.Encoder,
	ht hint.Hint,
	bks []byte,
	bDoc []byte,
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

	if hinter, err := enc.Decode(bDoc); err != nil {
		return err
	} else if d, ok := hinter.(DocumentData); !ok {
		return xerrors.Errorf("not DocumetData type : %T", d)
	} else {
		it.doc = d
	}

	it.cid = currency.CurrencyID(scid)

	return nil
}
