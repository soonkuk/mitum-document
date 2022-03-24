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

type testSignDocuments struct {
	baseTest
}

func (t *testSignDocuments) TestNewSignDocuments() {
	bsDocID := "1sdi"
	ownerAccount := generateAccount()
	// sender address is same with owner address
	sender := ownerAccount.Address
	// random token
	token := util.UUID().Bytes()
	// currency id
	cid := currency.CurrencyID("SHOWME")
	// sign document item
	items := []SignDocumentsItem{
		NewSignDocumentsItemSingleFile(bsDocID, ownerAccount.Address, cid),
	}
	// sign document fact
	fact := NewSignDocumentsFact(token, sender, items)
	var fs []base.FactSign
	// generate fact signature
	sig, err := base.NewFactSignature(ownerAccount.Priv, fact, nil)
	t.NoError(err)
	// make fact sign
	fs = append(fs, base.NewBaseFactSign(ownerAccount.Priv.Publickey(), sig))
	// sign document with fact and fact sign
	cd, err := NewSignDocuments(fact, fs, "")
	t.NoError(err)
	t.NoError(cd.IsValid(nil))
	t.Implements((*base.Fact)(nil), cd.Fact())
	t.Implements((*operation.Operation)(nil), cd)
	ufact := cd.Fact().(SignDocumentsFact)
	uBSDocItem := ufact.Items()[0]
	// compare filedata from signed BSDocData's fact with original filedata
	t.Equal(bsDocID, uBSDocItem.DocumentID())
	t.Equal(ownerAccount.Address, uBSDocItem.Owner())
}

func (t *testSignDocuments) TestEmptyToken() {
	bsDocID := "1sdi"
	ownerAccount := generateAccount()
	// sender address is same with owner address
	sender := ownerAccount.Address
	// currency id
	cid := currency.CurrencyID("SHOWME")
	// sign document item
	items := []SignDocumentsItem{
		NewSignDocumentsItemSingleFile(bsDocID, ownerAccount.Address, cid),
	}
	// sign document fact
	fact := NewSignDocumentsFact(nil, sender, items)
	var fs []base.FactSign
	// generate fact signature
	sig, err := base.NewFactSignature(ownerAccount.Priv, fact, nil)
	t.NoError(err)
	// make fact sign
	fs = append(fs, base.NewBaseFactSign(ownerAccount.Priv.Publickey(), sig))
	// sign document with fact and fact sign
	cd, err := NewSignDocuments(fact, fs, "")
	err = cd.IsValid(nil)
	t.Contains(err.Error(), "Operation has empty token")
}

func (t *testSignDocuments) TestEmptyItem() {
	sender := generateAccount()
	token := util.UUID().Bytes()
	// sign document item
	items := []SignDocumentsItem{}
	// sign document fact
	fact := NewSignDocumentsFact(token, sender.Address, items)
	var fs []base.FactSign
	// generate fact signature
	sig, err := base.NewFactSignature(sender.Priv, fact, nil)
	t.NoError(err)
	// make fact sign
	fs = append(fs, base.NewBaseFactSign(sender.Priv.Publickey(), sig))
	// sign document with fact and fact sign
	cd, err := NewSignDocuments(fact, fs, "")
	err = cd.IsValid(nil)
	t.Contains(err.Error(), "empty items")
}

func (t *testSignDocuments) TestMaxItem() {
	i := uint(0)
	items := make([]SignDocumentsItem, MaxSignDocumentsItems+1)
	account := generateAccount()
	cid := currency.CurrencyID("SHOWME")
	for i < (MaxSignDocumentsItems + 1) {
		bsDocID := fmt.Sprint(i) + "sdi"
		// sign document item
		items[i] = NewSignDocumentsItemSingleFile(bsDocID, account.Address, cid)
		i++
	}
	// token
	token := util.UUID().Bytes()
	// sign document fact
	fact := NewSignDocumentsFact(token, account.Address, items)
	var fs []base.FactSign
	// generate fact signature
	sig, err := base.NewFactSignature(account.Priv, fact, nil)
	t.NoError(err)
	// make fact sign
	fs = append(fs, base.NewBaseFactSign(account.Priv.Publickey(), sig))
	// sign document with fact and fact sign
	cd, err := NewSignDocuments(fact, fs, "")
	err = cd.IsValid(nil)
	t.Contains(err.Error(), "over max")
}

func (t *testSignDocuments) TestInvalidAddress() {
	bsDocID := "1sdi"
	// sign BSDocData
	ownerAccount := generateAccount()
	// currency id
	cid := currency.CurrencyID("SHOWME")
	// sign document item
	items := []SignDocumentsItem{
		NewSignDocumentsItemSingleFile(bsDocID, ownerAccount.Address, cid),
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
	// sign document fact
	fact := NewSignDocumentsFact(token, sender, items)
	var fs []base.FactSign
	// generate fact signature
	sig, err := base.NewFactSignature(ownerAccount.Priv, fact, nil)
	t.NoError(err)
	// make fact sign
	fs = append(fs, base.NewBaseFactSign(ownerAccount.Priv.Publickey(), sig))
	// sign document with fact and fact sign
	cd, err := NewSignDocuments(fact, fs, "")
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
	// sign document fact
	fact = NewSignDocumentsFact(token, sender, items)
	fs = []base.FactSign{}
	// generate fact signature
	sig, err = base.NewFactSignature(ownerAccount.Priv, fact, nil)
	t.NoError(err)
	// make fact sign
	fs = append(fs, base.NewBaseFactSign(ownerAccount.Priv.Publickey(), sig))
	// sign document with fact and fact sign
	cd, err = NewSignDocuments(fact, fs, "")
	err = cd.IsValid(nil)
	t.Contains(err.Error(), "too long string address")
}

func (t *testSignDocuments) TestDuplicatedDocID() {
	bsDocID := "1sdi"
	// sign BSDocData
	ownerAccount := generateAccount()
	// currency id
	cid := currency.CurrencyID("SHOWME")
	// sign document item
	item := NewSignDocumentsItemSingleFile(bsDocID, ownerAccount.Address, cid)

	items := []SignDocumentsItem{
		item,
		item,
	}
	// token
	token := util.UUID().Bytes()
	// sign document fact
	fact := NewSignDocumentsFact(token, ownerAccount.Address, items)
	var fs []base.FactSign
	// generate fact signature
	sig, err := base.NewFactSignature(ownerAccount.Priv, fact, nil)
	t.NoError(err)
	// make fact sign
	fs = append(fs, base.NewBaseFactSign(ownerAccount.Priv.Publickey(), sig))
	// sign document with fact and fact sign
	cd, err := NewSignDocuments(fact, fs, "")
	err = cd.IsValid(nil)
	t.Contains(err.Error(), "duplicated documentID")
}

func TestSignDocuments(t *testing.T) {
	suite.Run(t, new(testSignDocuments))
}

func testSignDocumentsEncode(enc encoder.Encoder) suite.TestingSuite {
	t := new(baseTestOperationEncode)

	t.enc = enc
	t.newObject = func() interface{} {
		bsDocID := "1sdi"
		ownerAccount := generateAccount()
		cid := currency.CurrencyID("SHOWME")
		sender := ownerAccount.Address
		token := util.UUID().Bytes()

		// sign document item
		items := []SignDocumentsItem{
			NewSignDocumentsItemSingleFile(bsDocID, ownerAccount.Address, cid),
		}
		// sign document fact
		fact := NewSignDocumentsFact(token, sender, items)
		var fs []base.FactSign
		// generate fact signature
		sig, err := base.NewFactSignature(ownerAccount.Priv, fact, nil)
		t.NoError(err)
		// make fact sign
		fs = append(fs, base.NewBaseFactSign(ownerAccount.Priv.Publickey(), sig))
		// sign document with fact and fact sign
		sd, err := NewSignDocuments(fact, fs, "")
		t.NoError(err)

		return sd
	}

	t.compare = func(a, b interface{}) {
		da := a.(SignDocuments)
		db := b.(SignDocuments)

		t.Equal(da.Memo, db.Memo)

		fact := da.Fact().(SignDocumentsFact)
		ufact := db.Fact().(SignDocumentsFact)

		t.True(fact.Hint().Equal(ufact.Hint()))
		t.True(fact.sender.Equal(ufact.sender))
		t.Equal(len(fact.Items()), len(ufact.Items()))

		for i := range fact.Items() {
			a := fact.Items()[i]
			b := ufact.Items()[i]
			t.Equal(a.DocumentID(), b.DocumentID())
			t.Equal(a.Currency(), b.Currency())
		}
	}

	return t
}

func TestSignDocumentsEncodeJSON(t *testing.T) {
	suite.Run(t, testSignDocumentsEncode(jsonenc.NewEncoder()))
}

func TestSignDocumentsEncodeBSON(t *testing.T) {
	suite.Run(t, testSignDocumentsEncode(bsonenc.NewEncoder()))
}
