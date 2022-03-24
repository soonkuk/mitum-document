package digest

import (
	"github.com/protoconNet/mitum-document/document"
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
	var addresses = make([]string, len(doc.Accounts()))
	for i := range doc.Accounts() {
		addresses[i] = doc.Accounts()[i].String()
	}
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

func (doc DocumentDoc) DocumentID() string {
	return doc.va.doc.DocumentID()
}

func (doc DocumentDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}

	m["owner"] = doc.va.Document().Owner()
	m["documentid"] = doc.va.Document().DocumentID()
	m["docid"] = doc.va.Document().DocumentID()[:len(doc.va.Document().DocumentID())-3]
	m["doctype"] = doc.va.Document().DocumentID()[len(doc.va.Document().DocumentID())-3:]
	m["addresses"] = doc.addresses
	m["height"] = doc.height

	return bsonenc.Marshal(m)
}
