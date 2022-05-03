package extension // nolint:dupl

import (
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
	"go.mongodb.org/mongo-driver/bson"
)

func (it BaseCreateContractAccountsItem) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bsonenc.MergeBSONM(bsonenc.NewHintedDoc(it.Hint()),
			bson.M{
				"keys":    it.keys,
				"amounts": it.amounts,
			}),
	)
}

type CreateContractAccountsItemBSONUnpacker struct {
	KS bson.Raw `bson:"keys"`
	AM bson.Raw `bson:"amounts"`
}

func (it *BaseCreateContractAccountsItem) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var uca CreateContractAccountsItemBSONUnpacker
	if err := bson.Unmarshal(b, &uca); err != nil {
		return err
	}

	return it.unpack(enc, uca.KS, uca.AM)
}
