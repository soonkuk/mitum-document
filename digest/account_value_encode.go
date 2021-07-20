package digest

import (
	"github.com/soonkuk/mitum-data/blocksign"
	"github.com/soonkuk/mitum-data/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
)

func (va *AccountValue) unpack(enc encoder.Encoder, bac []byte, bl []byte, dm []byte, height, previousHeight base.Height) error {
	if bac != nil {
		i, err := currency.DecodeAccount(bac, enc)
		if err != nil {
			return err
		}
		va.ac = i
	}

	hbl, err := enc.DecodeSlice(bl)
	if err != nil {
		return err
	}

	balance := make([]currency.Amount, len(hbl))
	for i := range hbl {
		j, ok := hbl[i].(currency.Amount)
		if !ok {
			return util.WrongTypeError.Errorf("expected currency.Amount, not %T", hbl[i])
		}
		balance[i] = j
	}

	va.balance = balance

	hdm, err := enc.DecodeSlice(dm)
	if err != nil {
		return err
	}

	document := make([]blocksign.DocumentData, len(dm))
	for i := range hdm {
		j, ok := hdm[i].(blocksign.DocumentData)
		if !ok {
			return util.WrongTypeError.Errorf("expected blocksign.DocumentData, not %T", hbl[i])
		}
		document[i] = j
	}

	va.document = document

	va.height = height
	va.previousHeight = previousHeight

	return nil
}
