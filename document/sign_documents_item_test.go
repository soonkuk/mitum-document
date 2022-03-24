//go:build test
// +build test

package document

import (
	"testing"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/util/encoder"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
	"github.com/stretchr/testify/suite"
)

type testSignDocumentsItemImpl struct {
	baseTest
}

func (t *testSignDocumentsItemImpl) TestNewSignDocumentsItem() {
	bsDocID := "1sdi"
	ownerAccount := generateAccount()
	cid := currency.CurrencyID("SHOWME")
	// create document item
	sd := NewSignDocumentsItemSingleFile(bsDocID, ownerAccount.Address, cid)

	// compare filedata from created item's BSDocData with original filedata
	t.Equal(bsDocID, sd.id)
	t.Equal(ownerAccount.Address, sd.owner)
	t.Equal(cid, sd.cid)
}

func (t *testSignDocumentsItemImpl) TestInvaliDocumentType() {
	// invalid type docID
	bsDocID := "1cui"
	ownerAccount := generateAccount()
	cid := currency.CurrencyID("SHOWME")
	// create document item
	sd := NewSignDocumentsItemSingleFile(bsDocID, ownerAccount.Address, cid)
	err := sd.IsValid(nil)
	t.Contains(err.Error(), "invalid docID type")
}

func TestSignDocumentsItemSingleFile(t *testing.T) {
	suite.Run(t, new(testSignDocumentsItemImpl))
}

func testSignDocumentsItemSingleFileEncode(enc encoder.Encoder) suite.TestingSuite {
	t := new(baseTestOperationItemEncode)

	t.enc = enc
	t.newObject = func() interface{} {
		bsDocID := "1sdi"
		ownerAccount := generateAccount()
		cid := currency.CurrencyID("SHOWME")
		sd := NewSignDocumentsItemSingleFile(bsDocID, ownerAccount.Address, cid)
		err := sd.IsValid(nil)
		t.NoError(err)

		return sd
	}

	t.compare = func(a, b interface{}) {
		da := a.(SignDocumentsItemSingleDocument)
		db := b.(SignDocumentsItemSingleDocument)

		t.Equal(da.Hint(), db.Hint())
		t.Equal(da.DocumentID(), db.DocumentID())
		t.Equal(da.Owner(), db.Owner())
		t.Equal(da.Currency(), db.Currency())
	}

	return t
}

func TestSignDocumentsItemImplEncodeJSON(t *testing.T) {
	suite.Run(t, testSignDocumentsItemSingleFileEncode(jsonenc.NewEncoder()))
}

func TestSignDocumentsItemImplEncodeBSON(t *testing.T) {
	suite.Run(t, testSignDocumentsItemSingleFileEncode(bsonenc.NewEncoder()))
}
