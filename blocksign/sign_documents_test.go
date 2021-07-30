package blocksign

import (
	"strings"
	"testing"

	"github.com/soonkuk/mitum-data/currency"
	"github.com/stretchr/testify/suite"

	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/key"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/util"
)

type testSignDocuments struct {
	suite.Suite
	cid   currency.CurrencyID
	docId currency.Big
}

func (t *testSignDocuments) SetupSuite() {
	t.cid = currency.CurrencyID("SHOWME")
	t.docId = currency.NewBig(0)
}

func (t *testSignDocuments) TestNew() {
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

	tfd, err := NewSignDocuments(fact, fs, "")
	t.NoError(err)

	t.NoError(tfd.IsValid(nil))

	t.Implements((*base.Fact)(nil), tfd.Fact())
	t.Implements((*operation.Operation)(nil), tfd)
}

func (t *testSignDocuments) TestDuplicatedDocuments() {
	s := MustAddress(util.UUID().String())
	g := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()
	items := []SignDocumentItem{
		NewSignDocumentsItemSingleFile(t.docId, s, t.cid),
		NewSignDocumentsItemSingleFile(t.docId, s, t.cid),
	}
	fact := NewSignDocumentsFact(token, g, items)

	pk := key.MustNewBTCPrivatekey()
	sig, err := operation.NewFactSignature(pk, fact, nil)
	t.NoError(err)

	fs := []operation.FactSign{operation.NewBaseFactSign(pk.Publickey(), sig)}

	tfd, err := NewSignDocuments(fact, fs, "")
	t.NoError(err)

	err = tfd.IsValid(nil)
	t.Contains(err.Error(), "duplicated document found")
}

func (t *testSignDocuments) TestOverSizeMemo() {
	s := MustAddress(util.UUID().String())
	g := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()

	items := []SignDocumentItem{
		NewSignDocumentsItemSingleFile(t.docId, s, t.cid),
	}
	fact := NewSignDocumentsFact(token, g, items)

	pk := key.MustNewBTCPrivatekey()
	sig, err := operation.NewFactSignature(pk, fact, nil)
	t.NoError(err)

	fs := []operation.FactSign{operation.NewBaseFactSign(pk.Publickey(), sig)}

	memo := strings.Repeat("a", currency.MaxMemoSize) + "a"
	tf, err := NewSignDocuments(fact, fs, memo)
	t.NoError(err)

	err = tf.IsValid(nil)
	t.Contains(err.Error(), "memo over max size")
}

func TestSignDocuments(t *testing.T) {
	suite.Run(t, new(testSignDocuments))
}
