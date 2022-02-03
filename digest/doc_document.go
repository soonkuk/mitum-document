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

type BlocksignDocumentDoc struct {
	mongodbstorage.BaseDoc
	va        BlocksignDocumentValue
	addresses []string
	height    base.Height
}

func NewBlocksignDocumentDoc(
	enc encoder.Encoder,
	doc blocksign.DocumentData,
	height base.Height,
) (BlocksignDocumentDoc, error) {

	var addresses []string
	ads, err := doc.Addresses()
	if err != nil {
		return BlocksignDocumentDoc{}, err
	}
	addresses = make([]string, len(ads))
	for i := range ads {
		addresses[i] = ads[i].String()
	}
	va := NewBlocksignDocumentValue(doc, height)
	b, err := mongodbstorage.NewBaseDoc(nil, va, enc)
	if err != nil {
		return BlocksignDocumentDoc{}, err
	}

	return BlocksignDocumentDoc{
		BaseDoc:   b,
		va:        va,
		addresses: addresses,
		height:    height,
	}, nil
}

func (doc BlocksignDocumentDoc) DocumentId() currency.Big {
	return doc.va.doc.Info().Index()
}

func (doc BlocksignDocumentDoc) MarshalBSON() ([]byte, error) {
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

type BlockcityDocumentDoc struct {
	mongodbstorage.BaseDoc
	va        BlockcityDocumentValue
	addresses []string
	height    base.Height
}

func NewBlockcityDocumentDoc(
	enc encoder.Encoder,
	doc document.DocumentData,
	height base.Height,
) (BlockcityDocumentDoc, error) {

	var addresses = make([]string, 1)
	addresses[0] = doc.Owner().String()
	va := NewBlockcityDocumentValue(doc, height)
	b, err := mongodbstorage.NewBaseDoc(nil, va, enc)
	if err != nil {
		return BlockcityDocumentDoc{}, err
	}

	return BlockcityDocumentDoc{
		BaseDoc:   b,
		va:        va,
		addresses: addresses,
		height:    height,
	}, nil
}

func (doc BlockcityDocumentDoc) DocumentId() string {
	return doc.va.doc.DocumentId()
}

func (doc BlockcityDocumentDoc) MarshalBSON() ([]byte, error) {
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
