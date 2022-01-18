package digest

import (
	"github.com/soonkuk/mitum-blocksign/blockcity"
	"github.com/soonkuk/mitum-blocksign/blocksign"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util/encoder"
)

func (dv *BlocksignDocumentValue) unpack(enc encoder.Encoder, bdm []byte, height base.Height) error {

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

func (dv *BlockcityDocumentValue) unpack(enc encoder.Encoder, bdm []byte, height base.Height) error {

	if bdm != nil {
		i, err := blockcity.DecodeDocument(bdm, enc)
		if err != nil {
			return err
		}
		dv.doc = i
	}

	dv.height = height

	return nil
}
