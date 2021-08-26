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

type testTransferDocumentsItemSingleFile struct {
	suite.Suite
	cid   currency.CurrencyID
	docId currency.Big
	fh    FileHash
}

func (t *testTransferDocumentsItemSingleFile) SetupSuite() {
	t.cid = currency.CurrencyID("SHOWME")
	t.fh = FileHash("ABCD")
	t.docId = currency.NewBig(0)
}

func (t *testTransferDocumentsItemSingleFile) TestNew() {
	s := MustAddress(util.UUID().String())
	r := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()
	items := []TransferDocumentsItem{NewTransferDocumentsItemSingleFile(t.docId, s, r, t.cid)}
	fact := NewTransferDocumentsFact(token, s, items)

	var fs []operation.FactSign

	for _, pk := range []key.Privatekey{
		key.MustNewBTCPrivatekey(),
		key.MustNewBTCPrivatekey(),
		key.MustNewBTCPrivatekey(),
	} {
		sig, err := operation.NewFactSignature(pk, fact, nil)
		t.NoError(err)

		fs = append(fs, operation.NewBaseFactSign(pk.Publickey(), sig))
	}

	tf, err := NewTransferDocuments(fact, fs, "")
	t.NoError(err)

	t.NoError(tf.IsValid(nil))

	t.Implements((*base.Fact)(nil), tf.Fact())
	t.Implements((*operation.Operation)(nil), tf)
}

func (t *testTransferDocumentsItemSingleFile) TestZeroBig() {
	s := MustAddress(util.UUID().String())
	r := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()
	cid := currency.CurrencyID("")
	items := []TransferDocumentsItem{NewTransferDocumentsItemSingleFile(t.docId, s, r, cid)}

	err := items[0].IsValid(nil)
	t.Contains(err.Error(), "invalid length of currency id")

	fact := NewTransferDocumentsFact(token, s, items)

	pk := key.MustNewBTCPrivatekey()
	sig, err := operation.NewFactSignature(pk, fact, nil)
	t.NoError(err)

	fs := []operation.FactSign{operation.NewBaseFactSign(pk.Publickey(), sig)}

	tfd, err := NewTransferDocuments(fact, fs, "")
	t.NoError(err)

	err = tfd.IsValid(nil)
	t.Contains(err.Error(), "invalid length of currency id")
}

func TestTransferDocumentsItemSingleFile(t *testing.T) {
	suite.Run(t, new(testTransferDocumentsItemSingleFile))
}

func testTransferDocumentsItemSingleFileEncode(enc encoder.Encoder) suite.TestingSuite {
	t := new(baseTestOperationEncode)

	docId0 := currency.NewBig(0)
	docId1 := currency.NewBig(1)
	t.enc = enc
	t.newObject = func() interface{} {
		s := MustAddress(util.UUID().String())
		r := MustAddress(util.UUID().String())

		token := util.UUID().Bytes()
		items := []TransferDocumentsItem{
			NewTransferDocumentsItemSingleFile(docId0, s, r, currency.CurrencyID("SHOWME")),
			NewTransferDocumentsItemSingleFile(docId1, s, r, currency.CurrencyID("FINDME")),
		}
		fact := NewTransferDocumentsFact(token, s, items)

		var fs []operation.FactSign

		for _, pk := range []key.Privatekey{
			key.MustNewBTCPrivatekey(),
			key.MustNewBTCPrivatekey(),
			key.MustNewBTCPrivatekey(),
		} {
			sig, err := operation.NewFactSignature(pk, fact, nil)
			t.NoError(err)

			fs = append(fs, operation.NewBaseFactSign(pk.Publickey(), sig))
		}

		tfd, err := NewTransferDocuments(fact, fs, util.UUID().String())
		t.NoError(err)

		return tfd
	}

	t.compare = func(a, b interface{}) {
		ta := a.(TransferDocuments)
		tb := b.(TransferDocuments)

		t.Equal(ta.Memo, tb.Memo)

		fact := ta.Fact().(TransferDocumentsFact)
		ufact := tb.Fact().(TransferDocumentsFact)

		t.True(fact.sender.Equal(ufact.sender))
		t.Equal(len(fact.Items()), len(ufact.Items()))

		for i := range fact.Items() {
			a := fact.Items()[i]
			b := ufact.Items()[i]
			t.True(a.DocumentId().Equal(b.DocumentId()))
			t.True(a.Owner().Equal(b.Owner()))
			t.True(a.Receiver().Equal(b.Receiver()))
			t.Equal(a.Currency(), (b.Currency()))
		}

	}

	return t
}

func TestTransferDocumentsItemSingleleFiEncodeJSON(t *testing.T) {
	suite.Run(t, testTransferDocumentsItemSingleFileEncode(jsonenc.NewEncoder()))
}

func TestTransferDocumentssItemSingleFileEncodeBSON(t *testing.T) {
	suite.Run(t, testTransferDocumentsItemSingleFileEncode(bsonenc.NewEncoder()))
}
