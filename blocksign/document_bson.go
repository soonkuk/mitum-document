package blocksign

import (
	"github.com/spikeekips/mitum/base"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
	"go.mongodb.org/mongo-driver/bson"
)

/*
type DocumentBSONPacker struct {
	FH FileHash              `bson:"filehash"`
	CR base.AddressDecoder   `bson:"creator"`
	CD bson.Raw              `bson:"createdby"`
	OW base.AddressDecoder   `bson:"owner"`
	SG []base.AddressDecoder `bson:"signers"`
	SD bson.Raw              `bson:"signedby"`
}
*/

func (doc DocumentData) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bsonenc.MergeBSONM(
		bsonenc.NewHintedDoc(doc.Hint()),
		bson.M{
			"filehash": doc.fileHash,
			"creator":  doc.creator,
			"owner":    doc.owner,
			"signers":  doc.signers,
		}),
	)
}

type DocumentBSONUnpacker struct {
	FH string              `bson:"filehash"`
	CR []byte              `bson:"creator"`
	OW base.AddressDecoder `bson:"owner"`
	SG []byte              `bson:"signers"`
}

func (doc *DocumentData) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var udoc DocumentBSONUnpacker
	if err := enc.Unmarshal(b, &udoc); err != nil {
		return err
	}

	return doc.unpack(enc, udoc.FH, udoc.CR, udoc.OW, udoc.SG)
}
