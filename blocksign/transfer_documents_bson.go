package blocksign // nolint: dupl

import (
	"go.mongodb.org/mongo-driver/bson"

	"github.com/soonkuk/mitum-data/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/operation"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
	"github.com/spikeekips/mitum/util/valuehash"
)

func (fact TransferDocumentsFact) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bsonenc.MergeBSONM(bsonenc.NewHintedDoc(fact.Hint()),
			bson.M{
				"hash":   fact.h,
				"token":  fact.token,
				"sender": fact.sender,
				"items":  fact.items,
			}))
}

type TransferDocumentsFactBSONUnpacker struct {
	H  valuehash.Bytes     `bson:"hash"`
	TK []byte              `bson:"token"`
	SD base.AddressDecoder `bson:"sender"`
	IT []bson.Raw          `bson:"items"`
}

func (fact *TransferDocumentsFact) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var ufact TransferDocumentsFactBSONUnpacker
	if err := enc.Unmarshal(b, &ufact); err != nil {
		return err
	}

	its := make([][]byte, len(ufact.IT))
	for i := range ufact.IT {
		its[i] = ufact.IT[i]
	}

	return fact.unpack(enc, ufact.H, ufact.TK, ufact.SD, its)
}

func (op TransferDocuments) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bsonenc.MergeBSONM(
			op.BaseOperation.BSONM(),
			bson.M{"memo": op.Memo},
		))
}

func (op *TransferDocuments) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var ubo operation.BaseOperation
	if err := ubo.UnpackBSON(b, enc); err != nil {
		return err
	}

	*op = TransferDocuments{BaseOperation: ubo}

	var um currency.MemoBSONUnpacker
	if err := enc.Unmarshal(b, &um); err != nil {
		return err
	} else {
		op.Memo = um.Memo
	}

	return nil
}
