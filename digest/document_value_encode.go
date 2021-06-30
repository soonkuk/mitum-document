package digest

import (
	"github.com/soonkuk/mitum-data/blocksign"
	"github.com/soonkuk/mitum-data/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util/encoder"
)

func (dv *DocumentValue) unpack(enc encoder.Encoder, bac []byte, bfd []byte, height, previousHeight base.Height) error {
	if bac != nil {
		i, err := currency.DecodeAccount(bac, enc)
		if err != nil {
			return err
		}
		dv.ac = i
	}

	if bfd != nil {
		i, err := blocksign.DecodeFileData(bfd, enc)
		if err != nil {
			return err
		}
		dv.filedata = i
	}

	dv.height = height
	dv.previousHeight = previousHeight

	return nil
}
