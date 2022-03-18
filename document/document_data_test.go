package document

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
	"github.com/spikeekips/mitum/util/hint"
)

type testDocumentData struct {
	suite.Suite
}

func (t *testDocumentData) TestNewBSDocData() {
	fh := FileHash("ABCD")
	oPkey := key.NewBasePrivatekey()
	oKey, _ := currency.NewBaseAccountKey(oPkey.Publickey(), 100)
	oKeys, _ := currency.NewBaseAccountKeys([]currency.AccountKey{oKey}, 100)
	ownerAddress, _ := currency.NewAddressFromKeys(oKeys)
	ownerDocSign := MustNewDocSign(ownerAddress, "signcode0", true)

	sPkey := key.NewBasePrivatekey()
	sKey, _ := currency.NewBaseAccountKey(sPkey.Publickey(), 100)
	sKeys, _ := currency.NewBaseAccountKeys([]currency.AccountKey{sKey}, 100)
	aSigner, _ := currency.NewAddressFromKeys(sKeys)

	sDocSigns := []DocSign{MustNewDocSign(aSigner, "signcode1", true)}

	info := DocInfo{
		BaseHinter: hint.NewBaseHinter(DocInfoHint),
		id:         MustNewBSDocId("1sdi"),
		docType:    BSDocDataType,
	}

	a := MustNewBSDocData(info, ownerAddress, fh, ownerDocSign, "title", currency.NewBig(100), sDocSigns)
	t.Equal(a.fileHash, FileHash("ABCD"))
	t.Equal(a.Creator(), ownerAddress)
	t.Equal(a.Owner(), ownerAddress)
	t.Equal(a.Signers(), sDocSigns)
}

func (t *testDocumentData) TestNewBCUserData() {
	cPkey := key.NewBasePrivatekey()
	cKey, _ := currency.NewBaseAccountKey(cPkey.Publickey(), 100)
	cKeys, _ := currency.NewBaseAccountKeys([]currency.AccountKey{cKey}, 100)
	aCreator, _ := currency.NewAddressFromKeys(cKeys)

	sPkey := key.NewBasePrivatekey()
	sKey, _ := currency.NewBaseAccountKey(sPkey.Publickey(), 100)
	sKeys, _ := currency.NewBaseAccountKeys([]currency.AccountKey{sKey}, 100)
	aSigner, _ := currency.NewAddressFromKeys(sKeys)

	sDocSigns := []DocSign{MustNewDocSign(aSigner, sc, true)}

	info := DocInfo{
		BaseHinter: hint.NewBaseHinter(DocInfoHint),
		idx:        currency.NewBig(0),
		filehash:   fh,
	}

	a := MustNewDocumentData(info, aCreator, aOwner, sDocSigns)
	t.Equal(a.FileHash(), FileHash("ABCD"))
	t.Equal(a.Creator(), aCreator)
	t.Equal(a.Owner(), aOwner)
	t.Equal(a.Signers(), sDocSigns)
}

func (t *testDocumentData) TestNewBCLandData() {
	fh := FileHash("ABCD")
	sc := "signcode"
	cPkey := key.NewBasePrivatekey()
	cKey, _ := currency.NewBaseAccountKey(cPkey.Publickey(), 100)
	cKeys, _ := currency.NewBaseAccountKeys([]currency.AccountKey{cKey}, 100)
	aCreator, _ := currency.NewAddressFromKeys(cKeys)

	sPkey := key.NewBasePrivatekey()
	sKey, _ := currency.NewBaseAccountKey(sPkey.Publickey(), 100)
	sKeys, _ := currency.NewBaseAccountKeys([]currency.AccountKey{sKey}, 100)
	aSigner, _ := currency.NewAddressFromKeys(sKeys)

	sDocSigns := []DocSign{MustNewDocSign(aSigner, sc, true)}

	info := DocInfo{
		BaseHinter: hint.NewBaseHinter(DocInfoHint),
		idx:        currency.NewBig(0),
		filehash:   fh,
	}

	a := MustNewDocumentData(info, aCreator, aOwner, sDocSigns)
	t.Equal(a.FileHash(), FileHash("ABCD"))
	t.Equal(a.Creator(), aCreator)
	t.Equal(a.Owner(), aOwner)
	t.Equal(a.Signers(), sDocSigns)
}

func (t *testDocumentData) TestNewBCVotingData() {
	fh := FileHash("ABCD")
	sc := "signcode"
	cPkey := key.NewBasePrivatekey()
	cKey, _ := currency.NewBaseAccountKey(cPkey.Publickey(), 100)
	cKeys, _ := currency.NewBaseAccountKeys([]currency.AccountKey{cKey}, 100)
	aCreator, _ := currency.NewAddressFromKeys(cKeys)

	sPkey := key.NewBasePrivatekey()
	sKey, _ := currency.NewBaseAccountKey(sPkey.Publickey(), 100)
	sKeys, _ := currency.NewBaseAccountKeys([]currency.AccountKey{sKey}, 100)
	aSigner, _ := currency.NewAddressFromKeys(sKeys)

	sDocSigns := []DocSign{MustNewDocSign(aSigner, sc, true)}

	info := DocInfo{
		BaseHinter: hint.NewBaseHinter(DocInfoHint),
		idx:        currency.NewBig(0),
		filehash:   fh,
	}

	a := MustNewDocumentData(info, aCreator, aOwner, sDocSigns)
	t.Equal(a.FileHash(), FileHash("ABCD"))
	t.Equal(a.Creator(), aCreator)
	t.Equal(a.Owner(), aOwner)
	t.Equal(a.Signers(), sDocSigns)
}

func (t *testDocumentData) TestNewBCHistoryData() {
	fh := FileHash("ABCD")
	sc := "signcode"
	cPkey := key.NewBasePrivatekey()
	cKey, _ := currency.NewBaseAccountKey(cPkey.Publickey(), 100)
	cKeys, _ := currency.NewBaseAccountKeys([]currency.AccountKey{cKey}, 100)
	aCreator, _ := currency.NewAddressFromKeys(cKeys)

	sPkey := key.NewBasePrivatekey()
	sKey, _ := currency.NewBaseAccountKey(sPkey.Publickey(), 100)
	sKeys, _ := currency.NewBaseAccountKeys([]currency.AccountKey{sKey}, 100)
	aSigner, _ := currency.NewAddressFromKeys(sKeys)

	sDocSigns := []DocSign{MustNewDocSign(aSigner, sc, true)}

	info := DocInfo{
		BaseHinter: hint.NewBaseHinter(DocInfoHint),
		idx:        currency.NewBig(0),
		filehash:   fh,
	}

	a := MustNewDocumentData(info, aCreator, aOwner, sDocSigns)
	t.Equal(a.FileHash(), FileHash("ABCD"))
	t.Equal(a.Creator(), aCreator)
	t.Equal(a.Owner(), aOwner)
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

		oPkey := key.MustNewBTCPrivatekey()
		oKey, _ := currency.NewKey(oPkey.Publickey(), 100)
		oKeys, _ := currency.NewKeys([]currency.Key{oKey}, 100)
		aOwner, _ := currency.NewAddressFromKeys(oKeys)

		sPkey := key.MustNewBTCPrivatekey()
		sKey, _ := currency.NewKey(sPkey.Publickey(), 100)
		sKeys, _ := currency.NewKeys([]currency.Key{sKey}, 100)
		aSigner, _ := currency.NewAddressFromKeys(sKeys)

		sDocSigns := []DocSign{MustNewDocSign(aSigner, true)}

		info := DocInfo{
			idx:      currency.NewBig(0),
			filehash: fh,
		}

		a := MustNewDocumentData(info, aCreator, aOwner, sDocSigns)

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
		t.True(ca.Owner().Equal(cb.Owner()))
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
