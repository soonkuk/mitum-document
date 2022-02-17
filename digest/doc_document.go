package digest

import (
	"github.com/soonkuk/mitum-blocksign/document"
	"github.com/spikeekips/mitum/base"
	mongodbstorage "github.com/spikeekips/mitum/storage/mongodb"
	"github.com/spikeekips/mitum/util/encoder"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
)

type DocumentDoc struct {
	mongodbstorage.BaseDoc
	va        DocumentValue
	addresses []string
	height    base.Height
}

func NewDocumentDoc(
	enc encoder.Encoder,
	doc document.DocumentData,
	height base.Height,
) (DocumentDoc, error) {

	var addresses = make([]string, 1)
	addresses[0] = doc.Owner().String()
	va := NewDocumentValue(doc, height)
	b, err := mongodbstorage.NewBaseDoc(nil, va, enc)
	if err != nil {
		return DocumentDoc{}, err
	}

	return DocumentDoc{
		BaseDoc:   b,
		va:        va,
		addresses: addresses,
		height:    height,
	}, nil
}

func (doc DocumentDoc) DocumentId() string {
	return doc.va.doc.DocumentId()
}

func (doc DocumentDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}

	m["owner"] = doc.va.Document().Owner()
	m["documentid"] = doc.va.Document().DocumentId()
	m["docid"] = doc.va.Document().DocumentId()[:len(doc.va.Document().DocumentId())-3]
	m["doctype"] = doc.va.Document().DocumentId()[len(doc.va.Document().DocumentId())-3:]
	m["addresses"] = doc.addresses
	m["height"] = doc.height

	return bsonenc.Marshal(m)
}
