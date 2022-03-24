//go:build test
// +build test

package document

import (
	// "encoding/json"
	"encoding/json"
	"testing"

	"github.com/spikeekips/mitum/util/encoder"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/xerrors"
)

type testDocumentData struct {
	suite.Suite
}

func (t *testDocumentData) TestNewBSDocData() {
	d, oacc, sacc := newBSDocData("filehash", "1sdi", account{})
	t.Equal(d.Info(), MustNewDocInfo("1sdi", BSDocDataType))
	t.Equal(d.DocumentID(), "1sdi")
	t.Equal(d.DocumentType(), BSDocDataType)
	t.Equal(d.fileHash, FileHash("filehash"))
	t.Equal(d.Creator(), MustNewDocSign(oacc.Address, "signcode0", true))
	t.Equal(d.Owner(), oacc.Address)
	t.Equal(d.Signers(), []DocSign{MustNewDocSign(sacc.Address, "signcode1", false)})
}

func (t *testDocumentData) TestNewBCUserData() {
	d, oacc, stat := newBCUserData("1cui", account{})
	t.Equal(d.Owner(), oacc.Address)
	t.Equal(d.Info(), MustNewDocInfo("1cui", BCUserDataType))
	t.Equal(d.DocumentID(), "1cui")
	t.Equal(d.DocumentType(), BCUserDataType)
	t.Equal(d.gold, uint(10))
	t.Equal(d.bankgold, uint(10))
	t.Equal(d.statistics, stat)
}

func (t *testDocumentData) TestNewBCLandData() {
	d, oacc, racc := newBCLandData("1cli", account{})
	t.Equal(d.Owner(), oacc.Address)
	t.Equal(d.Info(), MustNewDocInfo("1cli", BCLandDataType))
	t.Equal(d.DocumentID(), "1cli")
	t.Equal(d.DocumentType(), BCLandDataType)
	t.Equal(d.address, "address")
	t.Equal(d.area, "area")
	t.Equal(d.renter, "renter")
	t.Equal(d.account, racc.Address)
	t.Equal(d.rentdate, "rentdate")
	t.Equal(d.periodday, uint(10))
}

func (t *testDocumentData) TestNewBCVotingData() {
	d, oacc, bacc := newBCVotingData("1cvi", account{})
	t.Equal(d.Owner(), oacc.Address)
	t.Equal(d.Info(), MustNewDocInfo("1cvi", BCVotingDataType))
	t.Equal(d.DocumentID(), "1cvi")
	t.Equal(d.DocumentType(), BCVotingDataType)
	t.Equal(d.round, uint(10))
	t.Equal(d.endVoteTime, "endVoteTime")
	t.Equal(d.candidates, []VotingCandidate{
		MustNewVotingCandidate(bacc.Address, "nickname", "manifest", 10),
	})
	t.Equal(d.bossname, "bossname")
	t.Equal(d.account, bacc.Address)
	t.Equal(d.termofoffice, "termofoffice")
}

func (t *testDocumentData) TestNewBCHistoryData() {
	d, oacc, bacc := newBCHistoryData("1chi", account{})
	t.Equal(d.Owner(), oacc.Address)
	t.Equal(d.Info(), MustNewDocInfo("1chi", BCHistoryDataType))
	t.Equal(d.DocumentID(), "1chi")
	t.Equal(d.DocumentType(), BCHistoryDataType)
	t.Equal(d.name, "name")
	t.Equal(d.account, bacc.Address)
	t.Equal(d.date, "date")
	t.Equal(d.usage, "usage")
	t.Equal(d.application, "application")
}

func TestDocumentData(t *testing.T) {
	suite.Run(t, new(testDocumentData))
}

func testEncodeFunc() func(encoder.Encoder, interface{}) ([]byte, error) {
	return func(enc encoder.Encoder, v interface{}) ([]byte, error) {
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
}

func testDecodeFunc() func(encoder.Encoder, []byte) (interface{}, error) {
	return func(enc encoder.Encoder, b []byte) (interface{}, error) {
		return DecodeDocumentData(b, enc)
	}
}

func testBSDocDataEncode(enc encoder.Encoder) suite.TestingSuite {
	t := new(baseTestEncode)
	t.enc = enc
	t.newObject = func() interface{} {
		a, _, _ := newBSDocData("filehash", "1sdi", account{})
		t.NoError(a.IsValid(nil))

		return *a
	}

	t.encode = testEncodeFunc()
	t.decode = testDecodeFunc()
	t.compare = func(a, b interface{}) {
		ca := a.(BSDocData)
		cb := b.(BSDocData)

		t.True(ca.Info().Equal(cb.Info()))
		t.True(ca.Owner().Equal(cb.Owner()))
		t.Equal(ca.creator, cb.creator)
		t.Equal(ca.fileHash, cb.fileHash)
		signers := ca.signers
		for i := range signers {
			t.True(signers[i].Equal(cb.signers[i]))
		}
	}

	return t
}

func testBCUserDataEncode(enc encoder.Encoder) suite.TestingSuite {
	t := new(baseTestEncode)
	t.enc = enc
	t.newObject = func() interface{} {
		a, _, _ := newBCUserData("1cui", account{})
		t.NoError(a.IsValid(nil))

		return *a
	}

	t.encode = testEncodeFunc()
	t.decode = testDecodeFunc()
	t.compare = func(a, b interface{}) {
		ca := a.(BCUserData)
		cb := b.(BCUserData)

		t.True(ca.Info().Equal(cb.Info()))
		t.True(ca.Owner().Equal(cb.Owner()))
		t.Equal(ca.gold, cb.gold)
		t.Equal(ca.bankgold, cb.bankgold)
		t.Equal(ca.statistics, cb.statistics)
	}

	return t
}

func testBCLandDataEncode(enc encoder.Encoder) suite.TestingSuite {
	t := new(baseTestEncode)

	t.enc = enc
	t.newObject = func() interface{} {
		a, _, _ := newBCLandData("1cli", account{})
		t.NoError(a.IsValid(nil))

		return *a
	}

	t.encode = testEncodeFunc()
	t.decode = testDecodeFunc()
	t.compare = func(a, b interface{}) {
		ca := a.(BCLandData)
		cb := b.(BCLandData)

		t.True(ca.Info().Equal(cb.Info()))
		t.True(ca.Owner().Equal(cb.Owner()))
		t.Equal(ca.account, cb.account)
		t.Equal(ca.address, cb.address)
		t.Equal(ca.area, cb.area)
		t.Equal(ca.periodday, cb.periodday)
		t.Equal(ca.rentdate, cb.rentdate)
		t.Equal(ca.renter, cb.renter)
	}

	return t
}

func testBCVotingDataEncode(enc encoder.Encoder) suite.TestingSuite {
	t := new(baseTestEncode)

	t.enc = enc
	t.newObject = func() interface{} {
		a, _, _ := newBCVotingData("1cvi", account{})
		t.NoError(a.IsValid(nil))

		return *a
	}

	t.encode = testEncodeFunc()
	t.decode = testDecodeFunc()
	t.compare = func(a, b interface{}) {
		ca := a.(BCVotingData)
		cb := b.(BCVotingData)

		t.True(ca.Info().Equal(cb.Info()))
		t.True(ca.Owner().Equal(cb.Owner()))
		t.Equal(ca.account, cb.account)
		t.Equal(ca.bossname, cb.bossname)
		t.Equal(ca.candidates, cb.candidates)
		t.Equal(ca.endVoteTime, cb.endVoteTime)
		t.Equal(ca.round, cb.round)
		t.Equal(ca.termofoffice, cb.termofoffice)
	}

	return t
}

func testBCHistoryDataEncode(enc encoder.Encoder) suite.TestingSuite {
	t := new(baseTestEncode)

	t.enc = enc
	t.newObject = func() interface{} {
		a, _, _ := newBCHistoryData("1chi", account{})
		t.NoError(a.IsValid(nil))

		return *a
	}

	t.encode = testEncodeFunc()
	t.decode = testDecodeFunc()
	t.compare = func(a, b interface{}) {
		ca := a.(BCHistoryData)
		cb := b.(BCHistoryData)

		t.True(ca.Info().Equal(cb.Info()))
		t.True(ca.Owner().Equal(cb.Owner()))
		t.Equal(ca.account, cb.account)
		t.Equal(ca.application, cb.application)
		t.Equal(ca.date, cb.date)
		t.Equal(ca.name, cb.name)
		t.Equal(ca.usage, cb.usage)
	}

	return t
}

func TestDocumentDataEncodeJSON(t *testing.T) {
	suite.Run(t, testBSDocDataEncode(jsonenc.NewEncoder()))
	suite.Run(t, testBCUserDataEncode(jsonenc.NewEncoder()))
	suite.Run(t, testBCLandDataEncode(jsonenc.NewEncoder()))
	suite.Run(t, testBCVotingDataEncode(jsonenc.NewEncoder()))
	suite.Run(t, testBCHistoryDataEncode(jsonenc.NewEncoder()))
}

func TestDocumentDataEncodeBSON(t *testing.T) {
	suite.Run(t, testBSDocDataEncode(bsonenc.NewEncoder()))
	suite.Run(t, testBCUserDataEncode(bsonenc.NewEncoder()))
	suite.Run(t, testBCLandDataEncode(bsonenc.NewEncoder()))
	suite.Run(t, testBCVotingDataEncode(bsonenc.NewEncoder()))
	suite.Run(t, testBCHistoryDataEncode(bsonenc.NewEncoder()))
}
