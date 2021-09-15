package blocksign

import (
	"testing"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/key"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
	"github.com/stretchr/testify/suite"
)

type testCreateDocumentsSingleFile struct {
	baseTest
}

func (t *testCreateDocumentsSingleFile) TestNew() {
	// owner private key
	ownerPrvk := key.MustNewBTCPrivatekey()

	// owner key
	ownerKey, err := currency.NewKey(ownerPrvk.Publickey(), 100)
	t.NoError(err)

	// owner keys
	ownerKeys, _ := currency.NewKeys([]currency.Key{ownerKey}, 100)
	// owner address calculated from owner keys
	owner, _ := currency.NewAddressFromKeys(ownerKeys)
	// sender address is same with owner address
	sender := owner

	// signer private key
	signerPrvk := key.MustNewBTCPrivatekey()

	// signer key
	signerKey, err := currency.NewKey(signerPrvk.Publickey(), 100)
	t.NoError(err)

	// signer keys
	signerKeys, _ := currency.NewKeys([]currency.Key{signerKey}, 100)
	// owner address calculated from owner keys
	signer, _ := currency.NewAddressFromKeys(signerKeys)

	// random token
	token := util.UUID().Bytes()

	// currency id
	cid := currency.CurrencyID("SHOWME")

	// uploaderSignCode for document
	fh := FileHash("ABCD")
	documentid := currency.NewBig(0)
	signcode0 := "user0"
	title := "title01"
	size := currency.NewBig(555)
	signcode1 := "user1"

	// create document item
	items := []CreateDocumentsItem{
		NewCreateDocumentsItemSingleFile(
			fh,
			documentid,
			signcode0,
			title,
			size,
			[]base.Address{signer},
			[]string{signcode1},
			cid,
		),
	}
	// create document fact
	fact := NewCreateDocumentsFact(token, sender, items)

	var fs []operation.FactSign
	// generate fact signature
	sig, err := operation.NewFactSignature(ownerPrvk, fact, nil)
	t.NoError(err)

	// make fact sign
	fs = append(fs, operation.NewBaseFactSign(ownerPrvk.Publickey(), sig))

	// create document with fact and fact sign
	cd, err := NewCreateDocuments(fact, fs, "")
	t.NoError(err)

	t.NoError(cd.IsValid(nil))

	t.Implements((*base.Fact)(nil), cd.Fact())
	t.Implements((*operation.Operation)(nil), cd)

	ufact := cd.Fact().(CreateDocumentsFact)
	// compare filedata from created document's fact with original filedata
	t.Equal(fh, ufact.Items()[0].FileHash())
	t.Equal(signer, ufact.Items()[0].Signers()[0])

}

func (t *testCreateDocumentsSingleFile) TestEmptyFileHash() {
	// owner private key
	ownerPrvk := key.MustNewBTCPrivatekey()

	// owner key
	ownerKey, err := currency.NewKey(ownerPrvk.Publickey(), 100)
	t.NoError(err)

	// owner keys
	ownerKeys, _ := currency.NewKeys([]currency.Key{ownerKey}, 100)
	// owner address calculated from owner keys
	owner, _ := currency.NewAddressFromKeys(ownerKeys)
	// sender address is same with owner address
	sender := owner

	// random token
	token := util.UUID().Bytes()

	// currency id
	cid := currency.CurrencyID("SHOWME")

	// Empty FileHash
	efh := FileHash("")
	documentid := currency.NewBig(0)
	signcode0 := "user0"
	title := "title01"
	size := currency.NewBig(555)

	// create document item
	items := []CreateDocumentsItem{
		NewCreateDocumentsItemSingleFile(
			efh,
			documentid,
			signcode0,
			title,
			size,
			[]base.Address{},
			[]string{},
			cid,
		),
	}
	// create document fact
	fact := NewCreateDocumentsFact(token, sender, items)

	// generate fact signature
	sig, err := operation.NewFactSignature(ownerPrvk, fact, nil)
	t.NoError(err)

	// make fact sign
	var fs []operation.FactSign
	fs = append(fs, operation.NewBaseFactSign(ownerPrvk.Publickey(), sig))

	// create document with fact and fact sign
	cd, err := NewCreateDocuments(fact, fs, "")
	t.NoError(err)

	err = cd.IsValid(nil)
	t.Contains(err.Error(), "empty fileHash")
}

func TestCreateDocumentsSingleFile(t *testing.T) {
	suite.Run(t, new(testCreateDocumentsSingleFile))
}

func testCreateDocumentsSingleFileEncode(enc encoder.Encoder) suite.TestingSuite {
	t := new(baseTestOperationEncode)

	t.enc = enc
	t.newObject = func() interface{} {
		// owner private key
		ownerPrvk := key.MustNewBTCPrivatekey()
		ownerKey, err := currency.NewKey(ownerPrvk.Publickey(), 100)
		t.NoError(err)
		ownerKeys, err := currency.NewKeys([]currency.Key{ownerKey}, 100)
		t.NoError(err)
		owner, _ := currency.NewAddressFromKeys(ownerKeys)

		signerPrvk0 := key.MustNewBTCPrivatekey()
		sigenrKey0, err := currency.NewKey(signerPrvk0.Publickey(), 100)
		t.NoError(err)
		signerKeys0, err := currency.NewKeys([]currency.Key{sigenrKey0}, 100)
		t.NoError(err)
		signer0, _ := currency.NewAddressFromKeys(signerKeys0)

		signerPrvk1 := key.MustNewBTCPrivatekey()
		sigenrKey1, err := currency.NewKey(signerPrvk1.Publickey(), 100)
		t.NoError(err)
		signerKeys1, err := currency.NewKeys([]currency.Key{sigenrKey1}, 100)
		t.NoError(err)
		signer1, _ := currency.NewAddressFromKeys(signerKeys1)

		sender := owner

		cid := currency.CurrencyID("SHOWME")

		filehash := FileHash("ABCD")
		documentid := currency.NewBig(0)
		signcode0 := "user0"
		title := "title01"
		size := currency.NewBig(555)
		signcode1 := "user1"
		signcode2 := "user2"
		// FileData for document
		item := NewCreateDocumentsItemSingleFile(filehash, documentid, signcode0, title, size, []base.Address{signer0, signer1}, []string{signcode1, signcode2}, cid)
		fact := NewCreateDocumentsFact(util.UUID().Bytes(), sender, []CreateDocumentsItem{item})

		var fs []operation.FactSign

		sig, err := operation.NewFactSignature(ownerPrvk, fact, nil)
		t.NoError(err)
		fs = append(fs, operation.NewBaseFactSign(ownerPrvk.Publickey(), sig))

		cd, err := NewCreateDocuments(fact, fs, util.UUID().String())
		t.NoError(err)

		return cd
	}

	t.compare = func(a, b interface{}) {
		da := a.(CreateDocuments)
		db := b.(CreateDocuments)

		t.Equal(da.Memo, db.Memo)

		fact := da.Fact().(CreateDocumentsFact)
		ufact := db.Fact().(CreateDocumentsFact)

		t.True(fact.Sender().Equal(ufact.Sender()))
		t.Equal(len(fact.Items()), len(ufact.Items()))

		for i := range fact.Items() {
			a := fact.Items()[i]
			b := ufact.Items()[i]

			t.True(a.FileHash().Equal(b.FileHash()))
			for i := range a.Signers() {
				t.Equal(a.Signers()[i].Bytes(), b.Signers()[i].Bytes())
			}

			t.Equal(a.Currency(), (b.Currency()))
		}
	}

	return t
}

func TestCreateDocumentsSingleFileEncodeJSON(t *testing.T) {
	suite.Run(t, testCreateDocumentsSingleFileEncode(jsonenc.NewEncoder()))
}

func TestCreateDocumentsSingleFileEncodeBSON(t *testing.T) {
	suite.Run(t, testCreateDocumentsSingleFileEncode(bsonenc.NewEncoder()))
}
