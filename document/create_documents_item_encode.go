package document // nolint: dupl, revive

import (
	"github.com/pkg/errors"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/util/encoder"
)

func (it *CreateDocumentsItemImpl) unpack(
	enc encoder.Encoder,
	// sdt string,
	bdd []byte,
	scid string,
) error {
	// it.doctype = hint.Type(sdt)

	// unpack documentdata
	if hinter, err := enc.Decode(bdd); err != nil {
		return err
	} else if i, ok := hinter.(DocumentData); !ok {
		return errors.Errorf("not Document: %T", hinter)
	} else {
		it.doc = i
	}

	it.cid = currency.CurrencyID(scid)

	return nil
}
