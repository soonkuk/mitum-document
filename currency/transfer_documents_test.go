package currency

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/key"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/util"
)

type testTransferDocuments struct {
	suite.Suite
	cid CurrencyID
}

func (t *testTransferDocuments) SetupSuite() {
	t.cid = CurrencyID("SHOWME")
	// t.sc = SignCode("ABCD")
	// t.fee = NewBig(1)
}

func (t *testTransferDocuments) TestNew() {
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

	tfd, err := NewTransferDocuments(fact, fs, "")
	t.NoError(err)

	t.NoError(tfd.IsValid(nil))

	t.Implements((*base.Fact)(nil), tfd.Fact())
	t.Implements((*operation.Operation)(nil), tfd)
}

func (t *testTransferDocuments) TestDuplicatedDocuments() {
	s := MustAddress(util.UUID().String())
	r := MustAddress(util.UUID().String())
	d := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()
	items := []TransferDocumentsItem{
		NewTransferDocumentsItemSingleFile(d, r, t.cid),
		NewTransferDocumentsItemSingleFile(d, r, t.cid),
	}
	fact := NewTransferDocumentsFact(token, s, items)

	pk := key.MustNewBTCPrivatekey()
	sig, err := operation.NewFactSignature(pk, fact, nil)
	t.NoError(err)

	fs := []operation.FactSign{operation.NewBaseFactSign(pk.Publickey(), sig)}

	tfd, err := NewTransferDocuments(fact, fs, "")
	t.NoError(err)

	err = tfd.IsValid(nil)
	t.Contains(err.Error(), "duplicated document found")
}

func (t *testTransferDocuments) TestReceiverSameWithSender() {
	s := MustAddress(util.UUID().String())
	d := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()
	items := []TransferDocumentsItem{
		NewTransferDocumentsItemSingleFile(d, s, t.cid),
	}
	fact := NewTransferDocumentsFact(token, s, items)

	pk := key.MustNewBTCPrivatekey()
	sig, err := operation.NewFactSignature(pk, fact, nil)
	t.NoError(err)

	fs := []operation.FactSign{operation.NewBaseFactSign(pk.Publickey(), sig)}

	tfd, err := NewTransferDocuments(fact, fs, "")
	t.NoError(err)

	err = tfd.IsValid(nil)
	t.Contains(err.Error(), "receiver is same with sender")
}

func (t *testTransferDocuments) TestOverSizeMemo() {
	s := MustAddress(util.UUID().String())
	r := MustAddress(util.UUID().String())
	d := MustAddress(util.UUID().String())

	token := util.UUID().Bytes()

	items := []TransferDocumentsItem{
		NewTransferDocumentsItemSingleFile(d, r, t.cid),
	}
	fact := NewTransferDocumentsFact(token, s, items)

	pk := key.MustNewBTCPrivatekey()
	sig, err := operation.NewFactSignature(pk, fact, nil)
	t.NoError(err)

	fs := []operation.FactSign{operation.NewBaseFactSign(pk.Publickey(), sig)}

	memo := strings.Repeat("a", MaxMemoSize) + "a"
	tf, err := NewTransferDocuments(fact, fs, memo)
	t.NoError(err)

	err = tf.IsValid(nil)
	t.Contains(err.Error(), "memo over max size")
}

func TestTransferDocuments(t *testing.T) {
	suite.Run(t, new(testTransferDocuments))
}
