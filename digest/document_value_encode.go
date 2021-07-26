package digest

import (
	"github.com/soonkuk/mitum-data/blocksign"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util/encoder"
)

func (dv *DocumentValue) unpack(enc encoder.Encoder, bdm []byte, height base.Height) error {

	if bdm != nil {
		i, err := blocksign.DecodeDocumentData(bdm, enc)
		if err != nil {
			return err
		}
		dv.doc = i
	}

	dv.height = height

	return nil
}
