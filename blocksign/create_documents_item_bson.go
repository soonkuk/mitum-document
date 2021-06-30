package blocksign // nolint:dupl

import (
	"github.com/spikeekips/mitum/base"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
	"go.mongodb.org/mongo-driver/bson"
)

func (it BaseCreateDocumentsItem) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bsonenc.MergeBSONM(bsonenc.NewHintedDoc(it.Hint()),
			bson.M{
				"keys":     it.keys,
				"signcode": it.sc,
				"owner":    it.owner,
				"currency": it.cid,
			}),
	)
}

type CreateDocumentsItemBSONUnpacker struct {
	KS bson.Raw            `bson:"keys"`
	SC string              `bson:"signcode"`
	OW base.AddressDecoder `bson:"owner"`
	CI string              `bson:"currency"`
}

func (it *BaseCreateDocumentsItem) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var ht bsonenc.HintedHead
	if err := enc.Unmarshal(b, &ht); err != nil {
		return err
	}

	var ucd CreateDocumentsItemBSONUnpacker
	if err := bson.Unmarshal(b, &ucd); err != nil {
		return err
	}

	return it.unpack(enc, ht.H, ucd.KS, ucd.SC, ucd.OW, ucd.CI)
}
