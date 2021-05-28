package currency

import (
	"github.com/spikeekips/mitum/base"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
	"go.mongodb.org/mongo-driver/bson"
)

type FileDataBSONPacker struct {
	US SignCode            `bson:"signcode"`
	OW base.AddressDecoder `bson:"owner"`
}

func (fd FileData) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bsonenc.MergeBSONM(
		bsonenc.NewHintedDoc(fd.Hint()),
		bson.M{
			"signcode": fd.signcode,
			"owner":    fd.owner,
		}),
	)
}

type FileDataBSONUnpacker struct {
	US string              `bson:"signcode"`
	OW base.AddressDecoder `bson:"owner"`
}

func (fd *FileData) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var ufd FileDataBSONUnpacker
	if err := enc.Unmarshal(b, &ufd); err != nil {
		return err
	}

	return fd.unpack(enc, ufd.US, ufd.OW)
}
