package blocksign

import (
	"github.com/soonkuk/mitum-data/currency"
	"github.com/stretchr/testify/suite"

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

	t.encs.AddHinter(currency.Address(""))
	t.encs.AddHinter(operation.BaseFactSign{})
	t.encs.AddHinter(currency.Key{})
	t.encs.AddHinter(currency.Keys{})
	t.encs.AddHinter(currency.TransfersFact{})
	t.encs.AddHinter(currency.Transfers{})
	t.encs.AddHinter(currency.CreateAccountsFact{})
	t.encs.AddHinter(currency.CreateAccounts{})
	t.encs.AddHinter(currency.KeyUpdaterFact{})
	t.encs.AddHinter(currency.KeyUpdater{})
	t.encs.AddHinter(currency.FeeOperationFact{})
	t.encs.AddHinter(currency.FeeOperation{})
	t.encs.AddHinter(currency.Account{})
	t.encs.AddHinter(currency.GenesisCurrenciesFact{})
	t.encs.AddHinter(currency.GenesisCurrencies{})
	t.encs.AddHinter(currency.Amount{})
	t.encs.AddHinter(currency.CurrencyRegisterFact{})
	t.encs.AddHinter(currency.CurrencyRegister{})
	t.encs.AddHinter(currency.CurrencyDesign{})
	t.encs.AddHinter(currency.NilFeeer{})
	t.encs.AddHinter(currency.FixedFeeer{})
	t.encs.AddHinter(currency.RatioFeeer{})
	t.encs.AddHinter(currency.CurrencyPolicyUpdaterFact{})
	t.encs.AddHinter(currency.CurrencyPolicyUpdater{})
	t.encs.AddHinter(currency.CurrencyPolicy{})
	t.encs.AddHinter(CreateDocumentsFact{})
	t.encs.AddHinter(CreateDocuments{})
	t.encs.AddHinter(TransferDocumentsFact{})
	t.encs.AddHinter(TransferDocuments{})
	t.encs.AddHinter(FileData{})
	t.encs.AddHinter(FileData{})
	t.encs.AddHinter(FileID(""))
	t.encs.AddHinter(SignCode(""))
	t.encs.AddHinter(key.BTCPublickeyHinter)
	t.encs.AddHinter(CreateDocumentsItemSingleFile{})
	t.encs.AddHinter(CreateDocumentsItemSingleFileHinter)
	t.encs.AddHinter(TransfersItemSingleDocumentHinter)
	t.encs.AddHinter(currency.CreateAccountsItemMultiAmountsHinter)
	t.encs.AddHinter(currency.CreateAccountsItemSingleAmountHinter)
	t.encs.AddHinter(currency.TransfersItemMultiAmountsHinter)
	t.encs.AddHinter(currency.TransfersItemSingleAmountHinter)

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

func (t *baseTestEncode) compareCurrencyDesign(a, b currency.CurrencyDesign) {
	t.True(a.Amount.Equal(b.Amount))
	t.True(a.GenesisAccount().Equal(a.GenesisAccount()))
	t.Equal(a.Policy(), b.Policy())
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
