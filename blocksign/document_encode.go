package blocksign

import (
	"golang.org/x/xerrors"

	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util/encoder"
)

func (doc *DocumentData) unpack(
	enc encoder.Encoder,
	filehash string, // filehash
	cr []byte, // creator
	ow base.AddressDecoder, // owner
	bsg []byte, // signers
) error {

	// unpack filehash
	doc.fileHash = FileHash(filehash)

	// unpack creator
	if hinter, err := enc.Decode(cr); err != nil {
		return err
	} else if c, ok := hinter.(DocSign); !ok {
		return xerrors.Errorf("not DocSign : %T", hinter)
	} else {
		doc.creator = c
	}

	// unpack owner
	o, err := ow.Encode(enc)
	if err != nil {
		return err
	}
	doc.owner = o

	hits, err := enc.DecodeSlice(bsg)
	if err != nil {
		return err
	}
	// unpack signers
	signers := make([]DocSign, len(bsg))

	for i := range hits {
		s, ok := hits[i].(DocSign)
		if !ok {
			xerrors.Errorf("not DocSign : %T", s)
		}

		signers[i] = s
	}
	doc.signers = signers

	return nil
}
