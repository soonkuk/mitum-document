package blocksign

import (
	"golang.org/x/xerrors"

	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util/encoder"
)

func (doc *DocumentData) unpack(
	enc encoder.Encoder,
	filehash string, // filehash
	id []byte,
	cr base.AddressDecoder, // creator
	bsg []byte, // signers
) error {

	// unpack filehash
	doc.fileHash = FileHash(filehash)

	// unpack document id
	if hinter, err := enc.Decode(id); err != nil {
		return err
	} else if i, ok := hinter.(DocId); !ok {
		return xerrors.Errorf("not DocId: %T", hinter)
	} else {
		doc.id = i
	}

	// unpack creator
	if a, err := cr.Encode(enc); err != nil {
		return err
	} else {
		doc.creator = a
	}

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
