package digest

import (
	"github.com/soonkuk/mitum-blocksign/blocksign"
	"github.com/soonkuk/mitum-blocksign/document"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	mongodbstorage "github.com/spikeekips/mitum/storage/mongodb"
	"github.com/spikeekips/mitum/util/encoder"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
)

type BSDocumentDoc struct {
	mongodbstorage.BaseDoc
	va        BSDocumentValue
	addresses []string
	height    base.Height
}

func NewBSDocumentDoc(
	enc encoder.Encoder,
	doc blocksign.DocumentData,
	height base.Height,
) (BSDocumentDoc, error) {

	var addresses []string
	ads, err := doc.Addresses()
	if err != nil {
		return BSDocumentDoc{}, err
	}
	addresses = make([]string, len(ads))
	for i := range ads {
		addresses[i] = ads[i].String()
	}
	va := NewBSDocumentValue(doc, height)
	b, err := mongodbstorage.NewBaseDoc(nil, va, enc)
	if err != nil {
		return BSDocumentDoc{}, err
	}

	return BSDocumentDoc{
		BaseDoc:   b,
		va:        va,
		addresses: addresses,
		height:    height,
	}, nil
}

func (doc BSDocumentDoc) DocumentId() currency.Big {
	return doc.va.doc.Info().Index()
}

func (doc BSDocumentDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}

	m["filehash"] = doc.va.Document().FileHash()
	m["documentid"] = doc.va.Document().Info().Index()
	m["creator"] = doc.va.Document().Creator().String()
	m["addresses"] = doc.addresses
	m["height"] = doc.height

	return bsonenc.Marshal(m)
}

type BCDocumentDoc struct {
	mongodbstorage.BaseDoc
	va        BCDocumentValue
	addresses []string
	height    base.Height
}

func NewBCDocumentDoc(
	enc encoder.Encoder,
	doc document.DocumentData,
	height base.Height,
) (BCDocumentDoc, error) {

	var addresses = make([]string, 1)
	addresses[0] = doc.Owner().String()
	va := NewBCDocumentValue(doc, height)
	b, err := mongodbstorage.NewBaseDoc(nil, va, enc)
	if err != nil {
		return BCDocumentDoc{}, err
	}

	return BCDocumentDoc{
		BaseDoc:   b,
		va:        va,
		addresses: addresses,
		height:    height,
	}, nil
}

func (doc BCDocumentDoc) DocumentId() string {
	return doc.va.doc.DocumentId()
}

func (doc BCDocumentDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}

	m["owner"] = doc.va.Document().Owner()
	m["documentid"] = doc.va.Document().DocumentId()
	m["addresses"] = doc.addresses
	m["height"] = doc.height

	return bsonenc.Marshal(m)
}
