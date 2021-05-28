package digest

import (
	"github.com/spikeekips/mitum/base"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
	"go.mongodb.org/mongo-driver/bson"
)

func (dv DocumentValue) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bsonenc.MergeBSONM(
		bsonenc.NewHintedDoc(dv.Hint()),
		bson.M{
			"ac":              dv.ac,
			"filedata":        dv.filedata,
			"height":          dv.height,
			"previous_height": dv.previousHeight,
		},
	))
}

type DocumentValueBSONUnpacker struct {
	AC bson.Raw    `bson:"ac"`
	FD bson.Raw    `bson:"filedata"`
	HT base.Height `bson:"height"`
	PT base.Height `bson:"previous_height"`
}

func (dv *DocumentValue) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var uva DocumentValueBSONUnpacker
	if err := enc.Unmarshal(b, &uva); err != nil {
		return err
	}

	return dv.unpack(enc, uva.AC, uva.FD, uva.HT, uva.PT)
}
