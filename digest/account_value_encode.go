package digest

import (
	"github.com/pkg/errors"
	"github.com/soonkuk/mitum-blocksign/blocksign"
	"github.com/soonkuk/mitum-blocksign/document"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
)

func (va *AccountValue) unpack(enc encoder.Encoder, bac []byte, bl []byte, sd []byte, cd []byte, height, previousHeight base.Height) error {
	if err := encoder.Decode(bac, enc, &va.ac); err != nil {
		return err
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

	if hinter, err := enc.Decode(sd); err != nil {
		return err
	} else if k, ok := hinter.(blocksign.DocumentInventory); !ok {
		return errors.Errorf("not Blocksign DocumentInventory: %T", hinter)
	} else {
		va.bsDocument = k
	}

	if hinter, err := enc.Decode(cd); err != nil {
		return err
	} else if l, ok := hinter.(document.DocumentInventory); !ok {
		return errors.Errorf("not Blockcity DocumentInventory: %T", hinter)
	} else {
		va.bcDocument = l
	}

	va.height = height
	va.previousHeight = previousHeight

	return nil
}
