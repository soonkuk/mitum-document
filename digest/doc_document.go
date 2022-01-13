package digest

import (
	"github.com/soonkuk/mitum-blocksign/blocksign"
	"github.com/spikeekips/mitum-currency/currency"
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
	doc blocksign.DocumentData,
	height base.Height,
) (DocumentDoc, error) {

	var addresses []string
	ads, err := doc.Addresses()
	if err != nil {
		return DocumentDoc{}, err
	}
	addresses = make([]string, len(ads))
	for i := range ads {
		addresses[i] = ads[i].String()
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

func (doc DocumentDoc) DocumentId() currency.Big {
	return doc.va.doc.Info().Index()
}

func (doc DocumentDoc) MarshalBSON() ([]byte, error) {
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
