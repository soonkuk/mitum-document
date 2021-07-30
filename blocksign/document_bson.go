package blocksign

import (
	"github.com/spikeekips/mitum/base"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
	"go.mongodb.org/mongo-driver/bson"
)

type DocumentBSONPacker struct {
	FH FileHash            `bson:"filehash"`
	CR base.AddressDecoder `bson:"creator"`
	CD bson.Raw            `bson:"createdby"`
	OW base.AddressDecoder `bson:"owner"`
	SG []DocSign           `bson:"signers"`
}

func (doc DocumentData) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bsonenc.MergeBSONM(
		bsonenc.NewHintedDoc(doc.Hint()),
		bson.M{
			"filehash":     doc.fileHash,
			"documentinfo": doc.info,
			"creator":      doc.creator,
			"owner":        doc.owner,
			"signers":      doc.signers,
		}),
	)
}

type DocumentBSONUnpacker struct {
	FH string              `bson:"filehash"`
	DI bson.Raw            `bson:"documentinfo"`
	CR base.AddressDecoder `bson:"creator"`
	OW base.AddressDecoder `bson:"owner"`
	SG bson.Raw            `bson:"signers"`
}

func (doc *DocumentData) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var udoc DocumentBSONUnpacker
	if err := enc.Unmarshal(b, &udoc); err != nil {
		return err
	}

	return doc.unpack(enc, udoc.FH, udoc.DI, udoc.CR, udoc.OW, udoc.SG)
}
