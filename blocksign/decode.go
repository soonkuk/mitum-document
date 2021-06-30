package blocksign

import (
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
)

func DecodeFileData(b []byte, enc encoder.Encoder) (FileData, error) {
	if i, err := enc.Decode(b); err != nil {
		return FileData{}, err
	} else if i == nil {
		return FileData{}, nil
	} else if v, ok := i.(FileData); !ok {
		return FileData{}, util.WrongTypeError.Errorf("not FileData; type=%T", i)
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
