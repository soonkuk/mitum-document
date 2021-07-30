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

type testSignDocumentsItemSingleFile struct {
	suite.Suite
	cid   currency.CurrencyID
	docId currency.Big
}

func (t *testSignDocumentsItemSingleFile) SetupSuite() {
	t.cid = currency.CurrencyID("SHOWME")
	t.docId = currency.NewBig(0)
}

func (t *testSignDocumentsItemSingleFile) TestNew() {
	s := MustAddress(util.UUID().String())
	g := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()
	items := []SignDocumentItem{NewSignDocumentsItemSingleFile(t.docId, s, t.cid)}
	fact := NewSignDocumentsFact(token, g, items)

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

	tf, err := NewSignDocuments(fact, fs, "")
	t.NoError(err)

	t.NoError(tf.IsValid(nil))

	t.Implements((*base.Fact)(nil), tf.Fact())
	t.Implements((*operation.Operation)(nil), tf)
}

func (t *testSignDocumentsItemSingleFile) TestZeroBig() {
	s := MustAddress(util.UUID().String())
	g := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()
	cid := currency.CurrencyID("")
	items := []SignDocumentItem{NewSignDocumentsItemSingleFile(t.docId, s, cid)}

	err := items[0].IsValid(nil)
	t.Contains(err.Error(), "invalid length of currency id")

	fact := NewSignDocumentsFact(token, g, items)

	pk := key.MustNewBTCPrivatekey()
	sig, err := operation.NewFactSignature(pk, fact, nil)
	t.NoError(err)

	fs := []operation.FactSign{operation.NewBaseFactSign(pk.Publickey(), sig)}

	tfd, err := NewSignDocuments(fact, fs, "")
	t.NoError(err)

	err = tfd.IsValid(nil)
	t.Contains(err.Error(), "invalid length of currency id")
}

func TestSignDocumentsItemSingleFile(t *testing.T) {
	suite.Run(t, new(testSignDocumentsItemSingleFile))
}

func testSignDocumentsItemSingleFileEncode(enc encoder.Encoder) suite.TestingSuite {
	t := new(baseTestOperationEncode)

	docId0 := currency.NewBig(0)
	docId1 := currency.NewBig(1)
	t.enc = enc
	t.newObject = func() interface{} {
		s := MustAddress(util.UUID().String())
		g := MustAddress(util.UUID().String())

		token := util.UUID().Bytes()
		items := []SignDocumentItem{
			NewSignDocumentsItemSingleFile(docId0, s, currency.CurrencyID("SHOWME")),
			NewSignDocumentsItemSingleFile(docId1, s, currency.CurrencyID("FINDME")),
		}
		fact := NewSignDocumentsFact(token, g, items)

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

		tfd, err := NewSignDocuments(fact, fs, util.UUID().String())
		t.NoError(err)

		return tfd
	}

	t.compare = func(a, b interface{}) {
		ta := a.(SignDocuments)
		tb := b.(SignDocuments)

		t.Equal(ta.Memo, tb.Memo)

		fact := ta.Fact().(SignDocumentsFact)
		ufact := tb.Fact().(SignDocumentsFact)

		t.True(fact.sender.Equal(ufact.sender))
		t.Equal(len(fact.Items()), len(ufact.Items()))

		for i := range fact.Items() {
			a := fact.Items()[i]
			b := ufact.Items()[i]
			t.True(a.DocumentId().Equal(b.DocumentId()))
			t.True(a.Owner().Equal(b.Owner()))
			t.True(a.Currency().Equal(b.Currency()))
		}

	}

	return t
}

func TestSignDocumentsItemSingleleFiEncodeJSON(t *testing.T) {
	suite.Run(t, testSignDocumentsItemSingleFileEncode(jsonenc.NewEncoder()))
}

func TestSignDocumentssItemSingleFileEncodeBSON(t *testing.T) {
	suite.Run(t, testSignDocumentsItemSingleFileEncode(bsonenc.NewEncoder()))
}
