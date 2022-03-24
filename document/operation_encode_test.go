//go:build test
// +build test

package document

import (
	"github.com/stretchr/testify/suite"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/key"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/util/encoder"
	"github.com/spikeekips/mitum/util/localtime"
)

type baseTestEncode struct {
	suite.Suite

	enc       encoder.Encoder
	encs      *encoder.Encoders
	newObject func() interface{}
	encode    func(encoder.Encoder, interface{}) ([]byte, error)
	decode    func(encoder.Encoder, []byte) (interface{}, error)
	compare   func(interface{}, interface{})
}

func (t *baseTestEncode) SetupSuite() {
	t.encs = encoder.NewEncoders()
	t.encs.AddEncoder(t.enc)

	t.encs.TestAddHinter(key.BasePublickey{})
	t.encs.TestAddHinter(base.StringAddressHinter)
	t.encs.TestAddHinter(currency.AddressHinter)
	t.encs.TestAddHinter(base.BaseFactSignHinter)
	t.encs.TestAddHinter(currency.AccountKeyHinter)
	t.encs.TestAddHinter(currency.AccountKeysHinter)
	t.encs.TestAddHinter(CreateDocumentsFactHinter)
	t.encs.TestAddHinter(CreateDocumentsHinter)
	t.encs.TestAddHinter(UpdateDocumentsFactHinter)
	t.encs.TestAddHinter(UpdateDocumentsHinter)
	t.encs.TestAddHinter(SignDocumentsFactHinter)
	t.encs.TestAddHinter(SignDocumentsHinter)
	t.encs.TestAddHinter(currency.AccountHinter)
	t.encs.TestAddHinter(currency.AmountHinter)
	t.encs.TestAddHinter(CreateDocumentsItemImplHinter)
	t.encs.TestAddHinter(UpdateDocumentsItemImplHinter)
	t.encs.TestAddHinter(SignItemSingleDocumentHinter)
	t.encs.TestAddHinter(currency.CurrencyDesignHinter)
	t.encs.TestAddHinter(currency.NilFeeerHinter)
	t.encs.TestAddHinter(DocSignHinter)
	t.encs.TestAddHinter(BSDocDataHinter)
	t.encs.TestAddHinter(BCUserDataHinter)
	t.encs.TestAddHinter(BCLandDataHinter)
	t.encs.TestAddHinter(BCVotingDataHinter)
	t.encs.TestAddHinter(BCHistoryDataHinter)
	t.encs.TestAddHinter(UserStatisticsHinter)
	t.encs.TestAddHinter(DocInfoHinter)
	t.encs.TestAddHinter(VotingCandidateHinter)
	t.encs.TestAddHinter(BSDocIDHinter)
	t.encs.TestAddHinter(BCUserDocIDHinter)
	t.encs.TestAddHinter(BCLandDocIDHinter)
	t.encs.TestAddHinter(BCVotingDocIDHinter)
	t.encs.TestAddHinter(BCHistoryDocIDHinter)
	t.encs.TestAddHinter(DocumentInventoryHinter)
}

func (t *baseTestEncode) TestEncode() {
	i := t.newObject()

	var err error

	var b []byte
	if t.encode != nil {
		b, err = t.encode(t.enc, i)
		t.NoError(err)
	} else {
		b, err = t.enc.Marshal(i)
		t.NoError(err)
	}

	var v interface{}
	if t.decode != nil {
		v, err = t.decode(t.enc, b)
		t.NoError(err)
	} else {
		v, err = t.enc.Decode(b)
		t.NoError(err)
	}

	t.compare(i, v)
}

func (t *baseTestEncode) newAccount() *account {
	return generateAccount()
}

func (t *baseTestEncode) compareCurrencyDesign(a, b currency.CurrencyDesign) {
	t.True(a.Hint().Equal(b.Hint()))
	t.True(a.Amount.Equal(b.Amount))
	t.True(a.GenesisAccount().Equal(a.GenesisAccount()))
	t.Equal(a.Policy(), b.Policy())
	t.True(a.Aggregate().Equal(b.Aggregate()))
}

type baseTestOperationEncode struct {
	baseTestEncode
}

func (t *baseTestOperationEncode) TestEncode() {
	i := t.newObject()
	op, ok := i.(operation.Operation)
	t.True(ok)

	b, err := t.enc.Marshal(op)
	t.NoError(err)

	hinter, err := t.enc.Decode(b)
	t.NoError(err)

	uop, ok := hinter.(operation.Operation)
	t.True(ok)

	fact := op.Fact().(operation.OperationFact)
	ufact := uop.Fact().(operation.OperationFact)
	t.True(fact.Hash().Equal(ufact.Hash()))
	t.True(fact.Hint().Equal(ufact.Hint()))
	t.Equal(fact.Token(), ufact.Token())

	t.True(op.Hash().Equal(uop.Hash()))

	t.Equal(len(op.Signs()), len(uop.Signs()))
	for i := range op.Signs() {
		a := op.Signs()[i]
		b := uop.Signs()[i]
		t.True(a.Signer().Equal(b.Signer()))
		t.Equal(a.Signature(), b.Signature())
		t.True(localtime.Equal(a.SignedAt(), b.SignedAt()))
	}

	t.compare(op, uop)
}

type baseTestOperationItemEncode struct {
	baseTestEncode
}

func (t *baseTestOperationItemEncode) TestEncode() {
	i := t.newObject()
	var ok bool

	switch i.(type) {
	case CreateDocumentsItem:
		item, _ := i.(CreateDocumentsItem)
		ok = true
		b, err := t.enc.Marshal(item)
		t.NoError(err)

		hinter, err := t.enc.Decode(b)
		t.NoError(err)

		uitem, k := hinter.(CreateDocumentsItem)
		t.True(k)

		t.compare(item, uitem)
	case UpdateDocumentsItem:
		item, _ := i.(UpdateDocumentsItem)
		ok = true
		b, err := t.enc.Marshal(item)
		t.NoError(err)

		hinter, err := t.enc.Decode(b)
		t.NoError(err)

		uitem, k := hinter.(UpdateDocumentsItem)
		t.True(k)

		t.compare(item, uitem)
	case SignDocumentsItem:
		item, _ := i.(SignDocumentsItem)
		ok = true
		b, err := t.enc.Marshal(item)
		t.NoError(err)

		hinter, err := t.enc.Decode(b)
		t.NoError(err)

		uitem, k := hinter.(SignDocumentsItem)
		t.True(k)

		t.compare(item, uitem)
	}
	t.True(ok)
}
