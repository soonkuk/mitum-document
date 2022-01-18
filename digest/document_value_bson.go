package digest

import (
	"github.com/spikeekips/mitum/base"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
	"go.mongodb.org/mongo-driver/bson"
)

func (dv BlocksignDocumentValue) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bsonenc.MergeBSONM(
		bsonenc.NewHintedDoc(dv.Hint()),
		bson.M{
			"document": dv.doc,
			"height":   dv.height,
		},
	))
}

type BlocksignDocumentValueBSONUnpacker struct {
	DM bson.Raw    `bson:"document"`
	HT base.Height `bson:"height"`
}

func (dv *BlocksignDocumentValue) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var uva BlocksignDocumentValueBSONUnpacker
	if err := enc.Unmarshal(b, &uva); err != nil {
		return err
	}

	return dv.unpack(enc, uva.DM, uva.HT)
}

func (dv BlockcityDocumentValue) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bsonenc.MergeBSONM(
		bsonenc.NewHintedDoc(dv.Hint()),
		bson.M{
			"document": dv.doc,
			"height":   dv.height,
		},
	))
}

type BlockcityDocumentValueBSONUnpacker struct {
	DM bson.Raw    `bson:"document"`
	HT base.Height `bson:"height"`
}

func (dv *BlockcityDocumentValue) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var uva BlocksignDocumentValueBSONUnpacker
	if err := enc.Unmarshal(b, &uva); err != nil {
		return err
	}

	return dv.unpack(enc, uva.DM, uva.HT)
}
