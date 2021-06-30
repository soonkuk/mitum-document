package blocksign

import (
	"testing"

	"github.com/soonkuk/mitum-data/currency"
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
	// document private key
	docPrvk := key.MustNewBTCPrivatekey()
	// owner private key
	ownerPrvk := key.MustNewBTCPrivatekey()

	// document key(pubkey, weight)
	docKey, err := currency.NewKey(docPrvk.Publickey(), 100)
	t.NoError(err)
	// owner key
	ownerKey, err := currency.NewKey(ownerPrvk.Publickey(), 100)
	t.NoError(err)

	// document keys(keys, threshold)
	docKeys, _ := currency.NewKeys([]currency.Key{docKey}, 100)
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

	// uploaderSignCode for document
	sc := SignCode("ABCD")

	// create document item
	item := NewCreateDocumentsItemSingleFile(docKeys, sc, owner, cid)
	// create document fact
	fact := NewCreateDocumentsFact(token, sender, []CreateDocumentsItem{item})

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
	t.Equal(sc, ufact.Items()[0].SignCode())
	t.Equal(owner, ufact.Items()[0].Owner())

}

func (t *testCreateDocumentsSingleFile) TestEmptyFileData() {
	// document private key
	docPrvk := key.MustNewBTCPrivatekey()
	// owner private key
	ownerPrvk := key.MustNewBTCPrivatekey()

	// document key(pubkey, weight)
	docKey, err := currency.NewKey(docPrvk.Publickey(), 100)
	t.NoError(err)
	// owner key
	ownerKey, err := currency.NewKey(ownerPrvk.Publickey(), 100)
	t.NoError(err)

	// document keys(keys, threshold)
	docKeys, _ := currency.NewKeys([]currency.Key{docKey}, 100)
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

	// Empty FileData
	esc := SignCode("")
	eow := currency.EmptyAddress

	// create document item
	item := NewCreateDocumentsItemSingleFile(docKeys, esc, eow, cid)
	// create document fact
	fact := NewCreateDocumentsFact(token, sender, []CreateDocumentsItem{item})

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
	t.Contains(err.Error(), "empty filedata")
}

func TestCreateDocumentsSingleFile(t *testing.T) {
	suite.Run(t, new(testCreateDocumentsSingleFile))
}

func testCreateDocumentsSingleFileEncode(enc encoder.Encoder) suite.TestingSuite {
	t := new(baseTestOperationEncode)

	t.enc = enc
	t.newObject = func() interface{} {
		// document private key
		docPrvk := key.MustNewBTCPrivatekey()
		// owner private key
		ownerPrvk := key.MustNewBTCPrivatekey()

		docKey, err := currency.NewKey(docPrvk.Publickey(), 100)
		t.NoError(err)
		ownerKey, err := currency.NewKey(ownerPrvk.Publickey(), 100)
		t.NoError(err)
		docKeys, err := currency.NewKeys([]currency.Key{docKey}, 100)
		t.NoError(err)
		ownerKeys, err := currency.NewKeys([]currency.Key{ownerKey}, 100)
		t.NoError(err)

		owner, _ := currency.NewAddressFromKeys(ownerKeys)
		sender := owner

		cid := currency.CurrencyID("SHOWME")

		// uploaderSignCode for document
		sc := SignCode("signcode")
		// FileData for document
		item := NewCreateDocumentsItemSingleFile(docKeys, sc, owner, cid)
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

			t.True(a.Keys().Hash().Equal(b.Keys().Hash()))
			for i := range a.Keys().Keys() {
				t.Equal(a.Keys().Keys()[i].Bytes(), b.Keys().Keys()[i].Bytes())
			}

			t.Equal(a.Keys().Threshold(), b.Keys().Threshold())
			asc := a.SignCode()
			bsc := b.SignCode()
			aowner := a.Owner()
			bowner := b.Owner()
			t.True(asc.Equal(bsc))
			t.True(aowner.Equal(bowner))
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
