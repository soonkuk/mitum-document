package blocksign

import (
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util/encoder"
)

func (fd *FileData) unpack(enc encoder.Encoder, us string, ow base.AddressDecoder) error {
	a, err := ow.Encode(enc)
	if err != nil {
		return err
	}
	fd.owner = a
	fd.signcode = SignCode(us)

	return nil
}
