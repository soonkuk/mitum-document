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
				"filehash": it.fileHash,
				"signers":  it.signers,
				"currency": it.cid,
			}),
	)
}

type CreateDocumentsItemBSONUnpacker struct {
	FH string                `bson:"filehash"`
	SG []base.AddressDecoder `bson:"signers"`
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

	return it.unpack(enc, ht.H, ucd.FH, ucd.SG, ucd.CI)
}
