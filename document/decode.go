package document

import (
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
)

func DecodeDocument(b []byte, enc encoder.Encoder) (DocumentData, error) {
	if i, err := enc.Decode(b); err != nil {
		return nil, err
	} else if i == nil {
		return nil, nil
	} else if v, ok := i.(DocumentData); !ok {
		return nil, util.WrongTypeError.Errorf("not Document; type=%T", i)
	} else {
		return v, nil
	}
}
