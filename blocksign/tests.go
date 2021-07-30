// +build test

package blocksign

import (
	"github.com/soonkuk/mitum-data/currency"
	"github.com/stretchr/testify/suite"

	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/key"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/base/prprocessor"
	"github.com/spikeekips/mitum/base/state"
	"github.com/spikeekips/mitum/isaac"
	"github.com/spikeekips/mitum/storage"
)

type account struct { // nolint: unused
	Address base.Address
	Priv    key.Privatekey
	Key     currency.Key
}

func (ac *account) Privs() []key.Privatekey {
	return []key.Privatekey{ac.Priv}
}

func (ac *account) Keys() currency.Keys {
	keys, _ := currency.NewKeys([]currency.Key{ac.Key}, 100)

	return keys
}

func generateAccount() *account { // nolint: unused
	priv := key.MustNewBTCPrivatekey()

	key, err := currency.NewKey(priv.Publickey(), 100)
	if err != nil {
		panic(err)
	}

	keys, err := currency.NewKeys([]currency.Key{key}, 100)
	if err != nil {
		panic(err)
	}

	address, _ := currency.NewAddressFromKeys(keys)

	return &account{Address: address, Priv: priv, Key: key}
}

type baseTest struct { // nolint: unused
	suite.Suite
	isaac.StorageSupportTest
	cid currency.CurrencyID
}

func (t *baseTest) SetupSuite() {
	t.StorageSupportTest.SetupSuite()

	_ = t.Encs.TestAddHinter(key.BTCPublickey{})
	_ = t.Encs.TestAddHinter(operation.BaseFactSign{})
	_ = t.Encs.TestAddHinter(currency.Key{})
	_ = t.Encs.TestAddHinter(currency.Keys{})
	_ = t.Encs.TestAddHinter(currency.Address(""))
	_ = t.Encs.TestAddHinter(currency.CreateAccounts{})
	_ = t.Encs.TestAddHinter(currency.Transfers{})
	_ = t.Encs.TestAddHinter(CreateDocuments{})
	_ = t.Encs.TestAddHinter(CreateDocumentsFact{})
	_ = t.Encs.TestAddHinter(SignDocuments{})
	_ = t.Encs.TestAddHinter(SignDocumentsFact{})
	_ = t.Encs.TestAddHinter(TransferDocuments{})
	_ = t.Encs.TestAddHinter(TransferDocumentsFact{})
	_ = t.Encs.TestAddHinter(currency.KeyUpdaterFact{})
	_ = t.Encs.TestAddHinter(currency.KeyUpdater{})
	_ = t.Encs.TestAddHinter(currency.FeeOperationFact{})
	_ = t.Encs.TestAddHinter(currency.FeeOperation{})
	_ = t.Encs.TestAddHinter(currency.Account{})
	_ = t.Encs.TestAddHinter(currency.CurrencyDesign{})
	_ = t.Encs.TestAddHinter(currency.CurrencyPolicyUpdaterFact{})
	_ = t.Encs.TestAddHinter(currency.CurrencyPolicyUpdater{})
	_ = t.Encs.TestAddHinter(currency.CurrencyPolicy{})
	_ = t.Encs.TestAddHinter(DocSign{})
	_ = t.Encs.TestAddHinter(DocInfo{})
	_ = t.Encs.TestAddHinter(DocId{})
	_ = t.Encs.TestAddHinter(DocumentData{})
	_ = t.Encs.TestAddHinter(DocumentInventory{})

	t.cid = currency.CurrencyID("SEEME")
}

func (t *baseTest) newAccount() *account {
	return generateAccount()
}

func (t *baseTest) currencyDesign(big currency.Big, cid currency.CurrencyID) currency.CurrencyDesign {
	return currency.NewCurrencyDesign(currency.NewAmount(big, cid), NewTestAddress(), currency.NewCurrencyPolicy(currency.ZeroBig, currency.NewNilFeeer()))
}

func (t *baseTest) compareCurrencyDesign(a, b currency.CurrencyDesign) {
	t.True(a.Amount.Equal(b.Amount))
	if a.GenesisAccount() != nil {
		t.True(a.GenesisAccount().Equal(a.GenesisAccount()))
	}
	t.Equal(a.Policy(), b.Policy())
}

type baseTestOperationProcessor struct { // nolint: unused
	baseTest
}

func (t *baseTestOperationProcessor) statepool(s ...[]state.State) (*storage.Statepool, prprocessor.OperationProcessor) {
	base := map[string]state.State{}
	for _, l := range s {
		for _, st := range l {
			base[st.Key()] = st
		}
	}

	pool, err := storage.NewStatepoolWithBase(t.Database(nil, nil), base)
	t.NoError(err)

	opr := (NewOperationProcessor(nil)).New(pool)

	return pool, opr
}

func (t *baseTestOperationProcessor) newStateKeys(a base.Address, keys currency.Keys) state.State {
	key := currency.StateKeyAccount(a)

	ac, err := currency.NewAccount(a, keys)
	t.NoError(err)

	value, _ := state.NewHintedValue(ac)
	su, err := state.NewStateV0(key, value, base.NilHeight)
	t.NoError(err)

	return su
}

func (t *baseTestOperationProcessor) newKey(pub key.Publickey, w uint) currency.Key {
	k, err := currency.NewKey(pub, w)
	if err != nil {
		panic(err)
	}

	return k
}

func (t *baseTestOperationProcessor) newAccount(exists bool, amounts []currency.Amount) (*account, []state.State) {
	ac := t.baseTest.newAccount()

	if !exists {
		return ac, nil
	}

	var sts []state.State
	sts = append(sts, t.newStateKeys(ac.Address, ac.Keys()))

	for _, am := range amounts {
		sts = append(sts, t.newStateAmount(ac.Address, am))
	}

	return ac, sts
}

func (t *baseTestOperationProcessor) newStateAmount(a base.Address, amount currency.Amount) state.State {
	key := currency.StateKeyBalance(a, amount.Currency())
	value, _ := state.NewHintedValue(amount)
	su, err := state.NewStateV0(key, value, base.NilHeight)
	t.NoError(err)

	return su
}

func (t *baseTestOperationProcessor) newStateBalance(a base.Address, big currency.Big, cid currency.CurrencyID) state.State {
	key := currency.StateKeyBalance(a, cid)
	value, _ := state.NewHintedValue(currency.NewAmount(big, cid))
	su, err := state.NewStateV0(key, value, base.NilHeight)
	t.NoError(err)

	return su
}

func (t *baseTestOperationProcessor) newStateDocumentId(doc DocId) state.State {
	value, _ := state.NewHintedValue(doc)
	su, err := state.NewStateV0(StateKeyLastDocumentId, value, base.NilHeight)
	t.NoError(err)

	return su
}

func (t *baseTestOperationProcessor) newStateDocuments(a base.Address, doc DocInfo) state.State {
	key := StateKeyDocuments(a)

	docinv := NewDocumentInventory([]DocInfo{doc})

	value, _ := state.NewHintedValue(docinv)
	su, err := state.NewStateV0(key, value, base.NilHeight)
	t.NoError(err)

	return su
}

func (t *baseTestOperationProcessor) newStateDocument(a base.Address, docData DocumentData) []state.State {

	var sts []state.State

	sts = append(sts, t.newStateDocumentId(DocId{idx: docData.Info().Index()}))

	sts = append(sts, t.newStateDocuments(a, docData.Info()))

	sts = append(sts, t.newStateDocumentData(docData))

	return sts
}

func (t *baseTestOperationProcessor) newStateDocumentData(docData DocumentData) state.State {
	key := StateKeyDocumentData(docData.FileHash())
	value, _ := state.NewHintedValue(docData)
	su, err := state.NewStateV0(key, value, base.NilHeight)
	t.NoError(err)

	return su
}

func (t *baseTestOperationProcessor) newCurrencyDesignState(cid currency.CurrencyID, big currency.Big, genesisAccount base.Address, feeer currency.Feeer) state.State {
	de := currency.NewCurrencyDesign(currency.NewAmount(big, cid), genesisAccount, currency.NewCurrencyPolicy(currency.ZeroBig, feeer))

	st, err := state.NewStateV0(currency.StateKeyCurrencyDesign(cid), nil, base.NilHeight)
	t.NoError(err)

	nst, err := currency.SetStateCurrencyDesignValue(st, de)
	t.NoError(err)

	return nst
}

func NewTestAddress() base.Address {
	k, err := currency.NewKey(key.MustNewBTCPrivatekey().Publickey(), 100)
	if err != nil {
		panic(err)
	}

	keys, err := currency.NewKeys([]currency.Key{k}, 100)
	if err != nil {
		panic(err)
	}

	a, err := currency.NewAddressFromKeys(keys)
	if err != nil {
		panic(err)
	}

	return a
}
