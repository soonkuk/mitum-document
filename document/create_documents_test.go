//go:build test
// +build test

package document

import (
	"fmt"
	"testing"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
	"github.com/stretchr/testify/suite"
)

type testCreateDocuments struct {
	baseTest
}

func (t *testCreateDocuments) TestNewCreateDocuments() {
	bsDocID := "1sdi"
	bcUserDocID := "1cui"
	bcLandDocID := "1cli"
	bcVotingDocID := "1cvi"
	bcHistoryDocID := "1chi"
	// create BSDocData
	bsDocData, ownerAccount, signerAccount := newBSDocData("filehash", bsDocID, account{})
	// create BCUserData
	bcUserData, _, stat := newBCUserData(bcUserDocID, *ownerAccount)
	// create BCLandData
	bcLandData, _, renterAccount := newBCLandData(bcLandDocID, *ownerAccount)
	// create BCVotingData
	bcVotingData, _, bossAccount := newBCVotingData(bcVotingDocID, *ownerAccount)
	// create BCHistoryData
	bcHistoryData, _, cityAdminAccount := newBCHistoryData(bcHistoryDocID, *ownerAccount)
	// sender address is same with owner address
	sender := ownerAccount.Address
	// random token
	token := util.UUID().Bytes()
	// currency id
	cid := currency.CurrencyID("SHOWME")
	// create document item
	items := []CreateDocumentsItem{
		NewCreateDocumentsItemImpl(*bsDocData, cid),
		NewCreateDocumentsItemImpl(*bcUserData, cid),
		NewCreateDocumentsItemImpl(*bcLandData, cid),
		NewCreateDocumentsItemImpl(*bcVotingData, cid),
		NewCreateDocumentsItemImpl(*bcHistoryData, cid),
	}
	// create document fact
	fact := NewCreateDocumentsFact(token, sender, items)
	var fs []base.FactSign
	// generate fact signature
	sig, err := base.NewFactSignature(ownerAccount.Priv, fact, nil)
	t.NoError(err)
	// make fact sign
	fs = append(fs, base.NewBaseFactSign(ownerAccount.Priv.Publickey(), sig))
	// create document with fact and fact sign
	cd, err := NewCreateDocuments(fact, fs, "")
	t.NoError(err)

	t.NoError(cd.IsValid(nil))
	t.Implements((*base.Fact)(nil), cd.Fact())
	t.Implements((*operation.Operation)(nil), cd)
	ufact := cd.Fact().(CreateDocumentsFact)
	uBSDocData, ok := ufact.Items()[0].Doc().(BSDocData)
	t.True(ok)
	uBCUserData, ok := ufact.Items()[1].Doc().(BCUserData)
	t.True(ok)
	uBCLandData, ok := ufact.Items()[2].Doc().(BCLandData)
	t.True(ok)
	uBCVotingData, ok := ufact.Items()[3].Doc().(BCVotingData)
	t.True(ok)
	uBCHistoryData, ok := ufact.Items()[4].Doc().(BCHistoryData)
	t.True(ok)
	// compare filedata from created BSDocData's fact with original filedata
	t.Equal(MustNewBSDocID(bsDocID), uBSDocData.info.id)
	t.Equal(BSDocDataType, uBSDocData.info.docType)
	t.Equal(currency.NewBig(100), uBSDocData.size)
	t.Equal(FileHash("filehash"), uBSDocData.fileHash)
	t.Equal(MustNewDocSign(ownerAccount.Address, "signcode0", true), uBSDocData.creator)
	t.Equal("title", uBSDocData.title)
	t.Equal(MustNewDocSign(signerAccount.Address, "signcode1", false), uBSDocData.signers[0])
	// compare filedata from created BCUserData's fact with original filedata
	t.Equal(MustNewBCUserDocID(bcUserDocID), uBCUserData.info.id)
	t.Equal(BCUserDataType, uBCUserData.info.docType)
	t.Equal(uint(10), uBCUserData.gold)
	t.Equal(uint(10), uBCUserData.bankgold)
	t.Equal(stat, uBCUserData.statistics)
	// compare filedata from created BCLandData's fact with original filedata
	t.Equal(MustNewBCLandDocID(bcLandDocID), uBCLandData.info.id)
	t.Equal(BCLandDataType, uBCLandData.info.docType)
	t.Equal(renterAccount.Address, uBCLandData.account)
	t.Equal("address", uBCLandData.address)
	t.Equal("area", uBCLandData.area)
	t.Equal(uint(10), uBCLandData.periodday)
	t.Equal("rentdate", uBCLandData.rentdate)
	t.Equal("renter", uBCLandData.renter)
	// compare filedata from created BCVotingData's fact with original filedata
	t.Equal(MustNewBCVotingDocID(bcVotingDocID), uBCVotingData.info.id)
	t.Equal(BCVotingDataType, uBCVotingData.info.docType)
	t.Equal(bossAccount.Address, uBCVotingData.account)
	t.Equal("bossname", uBCVotingData.bossname)
	t.Equal([]VotingCandidate{MustNewVotingCandidate(bossAccount.Address, "nickname", "manifest", uint(10))}, uBCVotingData.candidates)
	// compare filedata from created BCHistoryData's fact with original filedata
	t.Equal(MustNewBCHistoryDocID(bcHistoryDocID), uBCHistoryData.info.id)
	t.Equal(BCHistoryDataType, uBCHistoryData.info.docType)
	t.Equal(cityAdminAccount.Address, uBCHistoryData.account)
	t.Equal("application", uBCHistoryData.application)
	t.Equal("date", uBCHistoryData.date)
	t.Equal("name", uBCHistoryData.name)
	t.Equal("usage", uBCHistoryData.usage)
}

func (t *testCreateDocuments) TestEmptyToken() {
	bsDocID := "1sdi"
	// create BSDocData
	bsDocData, ownerAccount, _ := newBSDocData("filehash", bsDocID, account{})
	// sender address is same with owner address
	sender := ownerAccount.Address
	// currency id
	cid := currency.CurrencyID("SHOWME")
	// create document item
	items := []CreateDocumentsItem{
		NewCreateDocumentsItemImpl(*bsDocData, cid),
	}
	// create document fact
	fact := NewCreateDocumentsFact(nil, sender, items)
	var fs []base.FactSign
	// generate fact signature
	sig, err := base.NewFactSignature(ownerAccount.Priv, fact, nil)
	t.NoError(err)
	// make fact sign
	fs = append(fs, base.NewBaseFactSign(ownerAccount.Priv.Publickey(), sig))
	// create document with fact and fact sign
	cd, err := NewCreateDocuments(fact, fs, "")
	err = cd.IsValid(nil)
	t.Contains(err.Error(), "Operation has empty token")
}

func (t *testCreateDocuments) TestEmptyItem() {
	sender := generateAccount()
	token := util.UUID().Bytes()
	// create document item
	items := []CreateDocumentsItem{}
	// create document fact
	fact := NewCreateDocumentsFact(token, sender.Address, items)
	var fs []base.FactSign
	// generate fact signature
	sig, err := base.NewFactSignature(sender.Priv, fact, nil)
	t.NoError(err)
	// make fact sign
	fs = append(fs, base.NewBaseFactSign(sender.Priv.Publickey(), sig))
	// create document with fact and fact sign
	cd, err := NewCreateDocuments(fact, fs, "")
	err = cd.IsValid(nil)
	t.Contains(err.Error(), "empty items")
}

func (t *testCreateDocuments) TestMaxItem() {
	i := uint(0)
	items := make([]CreateDocumentsItem, MaxCreateDocumentsItems+1)
	account := generateAccount()
	for i < (MaxCreateDocumentsItems + 1) {
		bsDocID := fmt.Sprint(i) + "sdi"
		// create BSDocData
		bsDocData, _, _ := newBSDocData("filehash", bsDocID, *account)
		// currency id
		cid := currency.CurrencyID("SHOWME")
		// create document item
		items[i] = NewCreateDocumentsItemImpl(*bsDocData, cid)
		i++
	}
	// token
	token := util.UUID().Bytes()
	// create document fact
	fact := NewCreateDocumentsFact(token, account.Address, items)
	var fs []base.FactSign
	// generate fact signature
	sig, err := base.NewFactSignature(account.Priv, fact, nil)
	t.NoError(err)
	// make fact sign
	fs = append(fs, base.NewBaseFactSign(account.Priv.Publickey(), sig))
	// create document with fact and fact sign
	cd, err := NewCreateDocuments(fact, fs, "")
	err = cd.IsValid(nil)
	t.Contains(err.Error(), "over max")
}

func (t *testCreateDocuments) TestInvalidAddress() {
	bsDocID := "1sdi"
	// create BSDocData
	bsDocData, ownerAccount, _ := newBSDocData("filehash", bsDocID, account{})
	// currency id
	cid := currency.CurrencyID("SHOWME")
	// create document item
	items := []CreateDocumentsItem{
		NewCreateDocumentsItemImpl(*bsDocData, cid),
	}
	// token
	token := util.UUID().Bytes()
	// invalid short sender address
	n := 0
	stringAddress := ""
	for n < (base.MinAddressSize - base.AddressTypeSize - 1) {
		stringAddress = stringAddress + "a"
		n++
	}
	sender := currency.NewAddress(stringAddress)
	// create document fact
	fact := NewCreateDocumentsFact(token, sender, items)
	var fs []base.FactSign
	// generate fact signature
	sig, err := base.NewFactSignature(ownerAccount.Priv, fact, nil)
	t.NoError(err)
	// make fact sign
	fs = append(fs, base.NewBaseFactSign(ownerAccount.Priv.Publickey(), sig))
	// create document with fact and fact sign
	cd, err := NewCreateDocuments(fact, fs, "")
	err = cd.IsValid(nil)
	t.Contains(err.Error(), "too short string address")

	// invalid long sender address
	n = 0
	stringAddress = ""
	for n < (base.MaxAddressSize - base.AddressTypeSize + 1) {
		stringAddress = stringAddress + "a"
		n++
	}
	sender = currency.NewAddress(stringAddress)
	// create document fact
	fact = NewCreateDocumentsFact(token, sender, items)
	fs = []base.FactSign{}
	// generate fact signature
	sig, err = base.NewFactSignature(ownerAccount.Priv, fact, nil)
	t.NoError(err)
	// make fact sign
	fs = append(fs, base.NewBaseFactSign(ownerAccount.Priv.Publickey(), sig))
	// create document with fact and fact sign
	cd, err = NewCreateDocuments(fact, fs, "")
	err = cd.IsValid(nil)
	t.Contains(err.Error(), "too long string address")
}

func (t *testCreateDocuments) TestDuplicatedDocID() {
	bsDocID := "1sdi"
	// create BSDocData
	bsDocData, ownerAccount, _ := newBSDocData("filehash", bsDocID, account{})
	// currency id
	cid := currency.CurrencyID("SHOWME")
	// create document item
	item := NewCreateDocumentsItemImpl(*bsDocData, cid)
	items := []CreateDocumentsItem{
		item,
		item,
	}
	// token
	token := util.UUID().Bytes()
	// create document fact
	fact := NewCreateDocumentsFact(token, ownerAccount.Address, items)
	var fs []base.FactSign
	// generate fact signature
	sig, err := base.NewFactSignature(ownerAccount.Priv, fact, nil)
	t.NoError(err)
	// make fact sign
	fs = append(fs, base.NewBaseFactSign(ownerAccount.Priv.Publickey(), sig))
	// create document with fact and fact sign
	cd, err := NewCreateDocuments(fact, fs, "")
	err = cd.IsValid(nil)
	t.Contains(err.Error(), "duplicated documentID")
}

func TestCreateDocuments(t *testing.T) {
	suite.Run(t, new(testCreateDocuments))
}

func testCreateDocumentsEncode(enc encoder.Encoder) suite.TestingSuite {
	t := new(baseTestOperationEncode)

	t.enc = enc
	t.newObject = func() interface{} {
		bsDocID := "1sdi"
		bcUserDocID := "1cui"
		bcLandDocID := "1cli"
		bcVotingDocID := "1cvi"
		bcHistoryDocID := "1chi"
		// create BSDocData
		bsDocData, ownerAccount, _ := newBSDocData("filehash", bsDocID, account{})
		// create BCUserData
		bcUserData, _, _ := newBCUserData(bcUserDocID, *ownerAccount)
		// create BCLandData
		bcLandData, _, _ := newBCLandData(bcLandDocID, *ownerAccount)
		// create BCVotingData
		bcVotingData, _, _ := newBCVotingData(bcVotingDocID, *ownerAccount)
		// create BCHistoryData
		bcHistoryData, _, _ := newBCHistoryData(bcHistoryDocID, *ownerAccount)
		// sender address is same with owner address
		sender := ownerAccount.Address
		// random token
		token := util.UUID().Bytes()
		// currency id
		cid := currency.CurrencyID("SHOWME")
		// create document item
		items := []CreateDocumentsItem{
			NewCreateDocumentsItemImpl(*bsDocData, cid),
			NewCreateDocumentsItemImpl(*bcUserData, cid),
			NewCreateDocumentsItemImpl(*bcLandData, cid),
			NewCreateDocumentsItemImpl(*bcVotingData, cid),
			NewCreateDocumentsItemImpl(*bcHistoryData, cid),
		}
		// create document fact
		fact := NewCreateDocumentsFact(token, sender, items)
		var fs []base.FactSign
		// generate fact signature
		sig, err := base.NewFactSignature(ownerAccount.Priv, fact, nil)
		t.NoError(err)
		// make fact sign
		fs = append(fs, base.NewBaseFactSign(ownerAccount.Priv.Publickey(), sig))
		// create document with fact and fact sign
		cd, err := NewCreateDocuments(fact, fs, "")
		t.NoError(err)
		return cd
	}

	t.compare = func(a, b interface{}) {
		da := a.(CreateDocuments)
		db := b.(CreateDocuments)

		t.Equal(da.Memo, db.Memo)

		fact := da.Fact().(CreateDocumentsFact)
		ufact := db.Fact().(CreateDocumentsFact)

		t.True(fact.Hint().Equal(ufact.Hint()))
		t.True(fact.sender.Equal(ufact.sender))
		t.Equal(len(fact.Items()), len(ufact.Items()))

		for i := range fact.Items() {
			a := fact.Items()[i]
			b := ufact.Items()[i]
			t.Equal(a.Doc(), b.Doc())
			t.Equal(a.Currency(), b.Currency())
		}
	}

	return t
}

func TestCreateDocumentsEncodeJSON(t *testing.T) {
	suite.Run(t, testCreateDocumentsEncode(jsonenc.NewEncoder()))
}

func TestCreateDocumentsEncodeBSON(t *testing.T) {
	suite.Run(t, testCreateDocumentsEncode(bsonenc.NewEncoder()))
}
