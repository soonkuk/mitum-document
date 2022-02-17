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
			"document": dv.doc,
			"height":   dv.height,
		},
	))
}

type DocumentValueBSONUnpacker struct {
	DM bson.Raw    `bson:"document"`
	HT base.Height `bson:"height"`
}

func (dv *DocumentValue) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var uva DocumentValueBSONUnpacker
	if err := enc.Unmarshal(b, &uva); err != nil {
		return err
	}

	return dv.unpack(enc, uva.DM, uva.HT)
}
