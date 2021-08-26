package blocksign

import (
	"github.com/pkg/errors"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util/encoder"
)

func (doc *DocumentData) unpack(
	enc encoder.Encoder,
	di []byte,
	cr base.AddressDecoder, // creator
	ow base.AddressDecoder, // owner
	bsg []byte, // signers
) error {

	// unpack document info
	if hinter, err := enc.Decode(di); err != nil {
		return err
	} else if i, ok := hinter.(DocInfo); !ok {
		return errors.Errorf("not DocInfo: %T", hinter)
	} else {
		doc.info = i
	}

	// unpack creator
	if a, err := cr.Encode(enc); err != nil {
		return err
	} else {
		doc.creator = a
	}

	// unpack owner
	if a, err := ow.Encode(enc); err != nil {
		return err
	} else {
		doc.owner = a
	}

	hits, err := enc.DecodeSlice(bsg)
	if err != nil {
		return err
	}
	// unpack signers
	signers := make([]DocSign, len(hits))

	for i := range hits {
		s, ok := hits[i].(DocSign)
		if !ok {
			return errors.Errorf("not DocSign : %T", s)
		}

		signers[i] = s
	}
	doc.signers = signers

	return nil
}

func (ds *DocSign) unpack(
	enc encoder.Encoder,
	ad base.AddressDecoder, // address
	sg bool, // signed
) error {

	a, err := ad.Encode(enc)
	if err != nil {
		return err
	}
	ds.address = a
	ds.signed = sg

	return nil
}
