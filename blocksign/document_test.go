package blocksign

import (
	"encoding/json"
	"testing"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/xerrors"

	"github.com/spikeekips/mitum/base/key"
	"github.com/spikeekips/mitum/util/encoder"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
)

type testDocumentData struct {
	suite.Suite
}

func (t *testDocumentData) TestNew() {
	fh := FileHash("ABCD")
	cPkey := key.MustNewBTCPrivatekey()
	cKey, _ := currency.NewKey(cPkey.Publickey(), 100)
	cKeys, _ := currency.NewKeys([]currency.Key{cKey}, 100)
	aCreator, _ := currency.NewAddressFromKeys(cKeys)

	sPkey := key.MustNewBTCPrivatekey()
	sKey, _ := currency.NewKey(sPkey.Publickey(), 100)
	sKeys, _ := currency.NewKeys([]currency.Key{sKey}, 100)
	aSigner, _ := currency.NewAddressFromKeys(sKeys)
	aSigncode0 := "user0"
	aSigncode1 := "user1"

	sDocSigns := []DocSign{MustNewDocSign(aSigner, aSigncode1, true)}
	title := "title"
	size := currency.NewBig(333)

	info := DocInfo{
		idx:      currency.NewBig(0),
		filehash: fh,
	}

	a := MustNewDocumentData(info, aCreator, aSigncode0, title, size, sDocSigns)
	t.Equal(a.FileHash(), FileHash("ABCD"))
	t.Equal(a.Creator(), aCreator)
	t.Equal(a.Signers(), sDocSigns)
}

func TestDocumentData(t *testing.T) {
	suite.Run(t, new(testDocumentData))
}

func testDocumentDataEncode(enc encoder.Encoder) suite.TestingSuite {
	t := new(baseTestEncode)

	t.enc = enc
	t.newObject = func() interface{} {
		fh := FileHash("ABCD")
		cPkey := key.MustNewBTCPrivatekey()
		cKey, _ := currency.NewKey(cPkey.Publickey(), 100)
		cKeys, _ := currency.NewKeys([]currency.Key{cKey}, 100)
		aCreator, _ := currency.NewAddressFromKeys(cKeys)

		sPkey := key.MustNewBTCPrivatekey()
		sKey, _ := currency.NewKey(sPkey.Publickey(), 100)
		sKeys, _ := currency.NewKeys([]currency.Key{sKey}, 100)
		aSigner, _ := currency.NewAddressFromKeys(sKeys)
		aSigncode0 := "user0"
		aSigncode1 := "user1"

		sDocSigns := []DocSign{MustNewDocSign(aSigner, aSigncode1, true)}
		title := "title"
		size := currency.NewBig(333)

		info := DocInfo{
			idx:      currency.NewBig(0),
			filehash: fh,
		}

		a := MustNewDocumentData(info, aCreator, aSigncode0, title, size, sDocSigns)

		t.NoError(a.IsValid(nil))

		return a
	}

	t.encode = func(enc encoder.Encoder, v interface{}) ([]byte, error) {
		b, err := enc.Marshal(struct {
			F DocumentData
		}{F: v.(DocumentData)})
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
		return DecodeDocumentData(b, enc)
	}

	t.compare = func(a, b interface{}) {
		ca := a.(DocumentData)
		cb := b.(DocumentData)

		t.True(ca.FileHash().Equal(cb.FileHash()))
		t.True(ca.Creator().Equal(cb.Creator()))
		signers := ca.Signers()
		for i := range signers {
			t.True(signers[i].Equal(cb.Signers()[i]))
		}
	}

	return t
}

func TestDocumentDataEncodeJSON(t *testing.T) {
	suite.Run(t, testDocumentDataEncode(jsonenc.NewEncoder()))
}

func TestDocumentDataEncodeBSON(t *testing.T) {
	suite.Run(t, testDocumentDataEncode(bsonenc.NewEncoder()))
}
