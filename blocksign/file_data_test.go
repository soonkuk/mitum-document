package blocksign

import (
	"encoding/json"
	"testing"

	"github.com/soonkuk/mitum-data/currency"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/xerrors"

	"github.com/spikeekips/mitum/base/key"
	"github.com/spikeekips/mitum/util/encoder"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
)

type testFileData struct {
	suite.Suite
}

func (t *testFileData) TestNew() {
	sc := SignCode("ABCD")
	sKey, _ := currency.NewKey(key.MustNewBTCPrivatekey().Publickey(), 100)
	sKeys, _ := currency.NewKeys([]currency.Key{sKey}, 100)
	sOwner, _ := currency.NewAddressFromKeys(sKeys)

	a := MustNewFileData(sc, sOwner)
	t.Equal(a, a.WithData(sc, sOwner))

	nKey, _ := currency.NewKey(key.MustNewBTCPrivatekey().Publickey(), 100)
	nKeys, _ := currency.NewKeys([]currency.Key{nKey}, 100)
	nOwner, _ := currency.NewAddressFromKeys(nKeys)

	_ = a.WithData(SignCode("EFGH"), nOwner)

	t.Equal(a.SignCode(), SignCode("ABCD"))
	t.Equal(a.Owner(), sOwner)
}

func TestFileData(t *testing.T) {
	suite.Run(t, new(testFileData))
}

func testFileDataEncode(enc encoder.Encoder) suite.TestingSuite {
	t := new(baseTestEncode)

	t.enc = enc
	t.newObject = func() interface{} {
		sc := SignCode("ABCD")
		sKey, _ := currency.NewKey(key.MustNewBTCPrivatekey().Publickey(), 100)
		sKeys, _ := currency.NewKeys([]currency.Key{sKey}, 100)
		sOwner, _ := currency.NewAddressFromKeys(sKeys)

		a := MustNewFileData(sc, sOwner)
		t.NoError(a.IsValid(nil))

		return a
	}

	t.encode = func(enc encoder.Encoder, v interface{}) ([]byte, error) {
		b, err := enc.Marshal(struct {
			F FileData
		}{F: v.(FileData)})
		if err != nil {
			return nil, err
		}

		switch enc.Hint().Type() {
		case jsonenc.JSONEncoderType:
			var D struct {
				F json.RawMessage
			}
			if err := enc.Unmarshal(b, &D); err != nil {
				return nil, err
			} else {
				return []byte(D.F), nil
			}
		case bsonenc.BSONEncoderType:
			var D struct {
				F bson.Raw
			}
			if err := enc.Unmarshal(b, &D); err != nil {
				return nil, err
			} else {
				return []byte(D.F), nil
			}
		default:
			return nil, xerrors.Errorf("unknown encoder, %v", enc)
		}
	}

	t.decode = func(enc encoder.Encoder, b []byte) (interface{}, error) {
		return DecodeFileData(b, enc)
	}

	t.compare = func(a, b interface{}) {
		ca := a.(FileData)
		cb := b.(FileData)

		t.True(ca.SignCode().Equal(cb.SignCode()))
		t.True(ca.Owner().Equal(cb.Owner()))
	}

	return t
}

func TestFileDataEncodeJSON(t *testing.T) {
	suite.Run(t, testFileDataEncode(jsonenc.NewEncoder()))
}

func TestFileDataEncodeBSON(t *testing.T) {
	suite.Run(t, testFileDataEncode(bsonenc.NewEncoder()))
}
