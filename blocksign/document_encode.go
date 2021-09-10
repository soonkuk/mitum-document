package blocksign

import (
	"github.com/pkg/errors"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util/encoder"
)

func (doc *DocumentData) unpack(
	enc encoder.Encoder,
	di []byte,
	cr []byte, // creator
	tl string,
	sz currency.Big,
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
	if hinter, err := enc.Decode(cr); err != nil {
		return err
	} else if i, ok := hinter.(DocSign); !ok {
		return errors.Errorf("not DocSign: %T", hinter)
	} else {
		doc.creator = i
	}

	doc.title = tl
	doc.size = sz

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
	sc string,
	sg bool, // signed
) error {

	a, err := ad.Encode(enc)
	if err != nil {
		return err
	}
	ds.address = a
	ds.signcode = sc
	ds.signed = sg

	return nil
}
