package digest

import (
	"github.com/soonkuk/mitum-data/blocksign"
	"github.com/soonkuk/mitum-data/currency"
	"github.com/spikeekips/mitum/base"
	mongodbstorage "github.com/spikeekips/mitum/storage/mongodb"
	"github.com/spikeekips/mitum/util/encoder"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
)

type DocumentDoc struct {
	mongodbstorage.BaseDoc
	va     DocumentValue
	height base.Height
}

func NewDocumentDoc(
	enc encoder.Encoder,
	doc blocksign.DocumentData,
	height base.Height,
) (DocumentDoc, error) {

	va := NewDocumentValue(doc, height)
	b, err := mongodbstorage.NewBaseDoc(nil, va, enc)
	if err != nil {
		return DocumentDoc{}, err
	}

	return DocumentDoc{
		BaseDoc: b,
		va:      va,
		height:  height,
	}, nil
}

func (doc DocumentDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}

	m["filehash"] = doc.va.Document().FileHash()
	m["documentid"] = doc.va.Document().Info().Index()
	m["owner"] = currency.StateAddressKeyPrefix(doc.va.Document().Owner())
	m["height"] = doc.height

	return bsonenc.Marshal(m)
}

/*
type DocumentDoc struct {
	mongodbstorage.BaseDoc
	st state.State
	fh blocksign.FileHash
	ow string
	di blocksign.DocInfo
}

// NewDocumentDoc gets the State of DocumentData
func NewDocumentDoc(st state.State, enc encoder.Encoder) (DocumentDoc, error) {

	var doc blocksign.DocumentData
	if i, err := blocksign.StateDocumentDataValue(st); err != nil {
		return DocumentDoc{}, xerrors.Errorf("DocumentDoc needs DocumentData state: %w", err)
	} else {
		doc = i
	}

	b, err := mongodbstorage.NewBaseDoc(nil, st, enc)
	if err != nil {
		return DocumentDoc{}, err
	}
	return DocumentDoc{
		BaseDoc: b,
		st:      st,
		fh:      doc.FileHash(),
		ow:      currency.StateAddressKeyPrefix(doc.Owner()),
		di:      doc.Info(),
	}, nil
}

func (doc DocumentDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}
	filehash := doc.st.Key()[:len(doc.st.Key())-len(blocksign.StateKeyDocumentDataSuffix)-1]
	m["filehash"] = filehash
	m["owner"] = doc.ow
	m["documentid"] = doc.di.Index()
	m["height"] = doc.st.Height()

	return bsonenc.Marshal(m)
}
*/
