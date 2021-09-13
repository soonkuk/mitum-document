package blocksign

import (
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
)

func DecodeDocumentData(b []byte, enc encoder.Encoder) (DocumentData, error) {
	if i, err := enc.Decode(b); err != nil {
		return DocumentData{}, err
	} else if i == nil {
		return DocumentData{}, nil
	} else if v, ok := i.(DocumentData); !ok {
		return DocumentData{}, util.WrongTypeError.Errorf("not DocumentData; type=%T", i)
	} else {
		return v, nil
	}
}
