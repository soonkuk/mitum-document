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

func DecodeCreateDocumentsItem(enc encoder.Encoder, b []byte) (CreateDocumentsItem, error) {
	if i, err := enc.Decode(b); err != nil {
		return nil, err
	} else if i == nil {
		return nil, nil
	} else if v, ok := i.(CreateDocumentsItem); !ok {
		return nil, util.WrongTypeError.Errorf("not CreateDocumentsItem; type=%T", i)
	} else {
		return v, nil
	}
}

func DecodeTransferDocumentsItem(enc encoder.Encoder, b []byte) (TransferDocumentsItem, error) {
	if i, err := enc.Decode(b); err != nil {
		return nil, err
	} else if i == nil {
		return nil, nil
	} else if v, ok := i.(TransferDocumentsItem); !ok {
		return nil, util.WrongTypeError.Errorf("not TransferDocumentsItem; type=%T", i)
	} else {
		return v, nil
	}
}
