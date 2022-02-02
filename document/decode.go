package document

import (
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
)

func DecodeDocument(b []byte, enc encoder.Encoder) (Document, error) {
	if i, err := enc.Decode(b); err != nil {
		return Document{}, err
	} else if i == nil {
		return Document{}, nil
	} else if v, ok := i.(Document); !ok {
		return Document{}, util.WrongTypeError.Errorf("not blockcity Document; type=%T", i)
	} else {
		return v, nil
	}
}
