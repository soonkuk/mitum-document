package blocksign // nolint:dupl

import (
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
	"go.mongodb.org/mongo-driver/bson"
)

func (it BaseSignDocumentsItem) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bsonenc.MergeBSONM(bsonenc.NewHintedDoc(it.Hint()),
			bson.M{
				"documentid": it.id,
				"owner":      it.owner,
				"currency":   it.cid,
			}),
	)
}

type SignDocumentsItemBSONUnpacker struct {
	DI currency.Big        `bson:"documentid"`
	OW base.AddressDecoder `bson:"owner"`
	CI string              `bson:"currency"`
}

func (it *BaseSignDocumentsItem) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var ucd SignDocumentsItemBSONUnpacker
	if err := bson.Unmarshal(b, &ucd); err != nil {
		return err
	}

	return it.unpack(enc, ucd.DI, ucd.OW, ucd.CI)
}
