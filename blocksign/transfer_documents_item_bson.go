package blocksign // nolint:dupl

import (
	"go.mongodb.org/mongo-driver/bson"

	"github.com/spikeekips/mitum/base"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
)

func (it BaseTransferDocumentsItem) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bsonenc.MergeBSONM(bsonenc.NewHintedDoc(it.Hint()),
			bson.M{
				"document": it.document,
				"receiver": it.receiver,
				"currency": it.cid,
			}),
	)
}

type BaseTransferDocumentsItemBSONUnpacker struct {
	DM base.AddressDecoder `bson:"document"`
	RC base.AddressDecoder `bson:"receiver"`
	CI string              `bson:"currency"`
}

func (it *BaseTransferDocumentsItem) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var ht bsonenc.HintedHead
	if err := enc.Unmarshal(b, &ht); err != nil {
		return err
	}

	var uit BaseTransferDocumentsItemBSONUnpacker
	if err := enc.Unmarshal(b, &uit); err != nil {
		return err
	}

	return it.unpack(enc, ht.H, uit.DM, uit.RC, uit.CI)
}
