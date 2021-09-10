package blocksign // nolint:dupl

import (
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
	"go.mongodb.org/mongo-driver/bson"
)

func (it BaseCreateDocumentsItem) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bsonenc.MergeBSONM(bsonenc.NewHintedDoc(it.Hint()),
			bson.M{
				"filehash":   it.fileHash,
				"documentid": it.documentid,
				"signcode":   it.signcode,
				"title":      it.title,
				"size":       it.size,
				"signers":    it.signers,
				"signcodes":  it.signcodes,
				"currency":   it.cid,
			}),
	)
}

type CreateDocumentsItemBSONUnpacker struct {
	FH string                `bson:"filehash"`
	DI currency.Big          `bson:"documentid"`
	SC string                `bson:"signcode"`
	TL string                `bson:"title"`
	SZ currency.Big          `bson:"size"`
	SG []base.AddressDecoder `bson:"signers"`
	SD []string              `bson:"signcodes"`
	CI string                `bson:"currency"`
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

	return it.unpack(enc, ht.H, ucd.FH, ucd.DI, ucd.SC, ucd.TL, ucd.SZ, ucd.SG, ucd.SD, ucd.CI)
}
