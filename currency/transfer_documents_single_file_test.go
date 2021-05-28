package currency

import (
	"testing"

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
	cid CurrencyID
	sc  SignCode
}

func (t *testTransferDocumentsItemSingleFile) SetupSuite() {
	t.cid = CurrencyID("SHOWME")
	t.sc = SignCode("ABCD")
}

func (t *testTransferDocumentsItemSingleFile) newTestFileData(oa base.Address) FileData {
	return NewFileData(t.sc, oa)
}

func (t *testTransferDocumentsItemSingleFile) TestNew() {
	s := MustAddress(util.UUID().String())
	r := MustAddress(util.UUID().String())
	d := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()
	items := []TransferDocumentsItem{NewTransferDocumentsItemSingleFile(d, r, t.cid)}
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
	d := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()
	cid := CurrencyID("")
	items := []TransferDocumentsItem{NewTransferDocumentsItemSingleFile(d, r, cid)}

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

// TODO :empty document address check

func TestTransferDocumentsItemSingleFile(t *testing.T) {
	suite.Run(t, new(testTransferDocumentsItemSingleFile))
}

func testTransferDocumentsItemSingleFileEncode(enc encoder.Encoder) suite.TestingSuite {
	t := new(baseTestOperationEncode)

	t.enc = enc
	t.newObject = func() interface{} {
		s := MustAddress(util.UUID().String())
		r := MustAddress(util.UUID().String())
		d0 := MustAddress(util.UUID().String())
		d1 := MustAddress(util.UUID().String())

		token := util.UUID().Bytes()
		items := []TransferDocumentsItem{
			NewTransferDocumentsItemSingleFile(d0, r, CurrencyID("SHOWME")),
			NewTransferDocumentsItemSingleFile(d1, r, CurrencyID("FINDME")),
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
			t.True(a.Receiver().Equal(b.Receiver()))
			t.True(a.Document().Equal(b.Document()))
			t.True(a.Currency().Equal(b.Currency()))
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
