package blocksign

import (
	"testing"

	"github.com/soonkuk/mitum-data/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/key"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/base/prprocessor"
	"github.com/spikeekips/mitum/base/state"
	"github.com/spikeekips/mitum/storage"
	"github.com/spikeekips/mitum/util"
	"github.com/stretchr/testify/suite"
	"golang.org/x/xerrors"
)

type testCreateDocumentsOperation struct {
	baseTestOperationProcessor
}

func (t *testCreateDocumentsOperation) processor(cp *currency.CurrencyPool, pool *storage.Statepool) prprocessor.OperationProcessor {
	copr, err := NewOperationProcessor(nil).
		SetProcessor(CreateDocuments{}, NewCreateDocumentsProcessor(cp))
	t.NoError(err)

	if pool == nil {
		return copr
	}

	return copr.New(pool)
}

func (t *testCreateDocumentsOperation) newOperation(sender base.Address, items []CreateDocumentsItem, pks []key.Privatekey) CreateDocuments {
	token := util.UUID().Bytes()
	fact := NewCreateDocumentsFact(token, sender, items)

	var fs []operation.FactSign
	for _, pk := range pks {
		sig, err := operation.NewFactSignature(pk, fact, nil)
		if err != nil {
			panic(err)
		}

		fs = append(fs, operation.NewBaseFactSign(pk.Publickey(), sig))
	}

	cd, err := NewCreateDocuments(fact, fs, "")
	if err != nil {
		panic(err)
	}

	err = cd.IsValid(nil)
	if err != nil {
		panic(err)
	}

	return cd
}

func (t *testCreateDocumentsOperation) TestNormalCase() {
	cid := currency.CurrencyID("SHOWME")

	// sender initial balance
	balance := []currency.Amount{
		currency.NewAmount(currency.NewBig(33), cid),
	}

	// sender account
	sa, st := t.newAccount(true, balance)
	// document account
	na, _ := t.newAccount(false, nil)
	// owner account
	oa := sa

	pool, _ := t.statepool(st)

	fee := currency.NewBig(1)
	feeer := currency.NewFixedFeeer(sa.Address, fee)

	cp := currency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(cid, currency.NewBig(99), sa.Address, feeer)))

	opr := t.processor(cp, pool)

	// filedata
	sc := SignCode("ABCD")
	fd := NewFileData(sc, oa.Address)

	// create document of na(document account) with oa(owner) which is sent from sa(sender)
	items := []CreateDocumentsItem{NewCreateDocumentsItemSingleFile(na.Keys(), sc, oa.Address, cid)}
	cd := t.newOperation(sa.Address, items, sa.Privs())

	t.NoError(opr.Process(cd))

	// check updated state
	// new document account state
	var ns state.State
	// new document filedata state
	var fds state.State
	// sender balance state
	var sb state.State
	for _, stu := range pool.Updates() {
		if currency.IsStateBalanceKey(stu.Key()) {
			st := stu.GetState()

			i, err := currency.StateBalanceValue(st)
			t.NoError(err)

			if st.Key() == currency.StateKeyBalance(sa.Address, i.Currency()) {
				sb = st
			} else {
				continue
			}
		} else if (currency.IsStateDocumentKey(stu.Key())) && (stu.Key() == currency.StateKeyDocument(na.Address)) {
			ns = stu.GetState()
		} else if (IsStateFileDataKey(stu.Key())) && (stu.Key() == StateKeyFileData(na.Address)) {
			fds = stu.GetState()
		}
	}

	address, err := currency.NewAddressFromKeys(na.Keys())
	t.NoError(err)
	uac := ns.Value().Interface().(currency.Account)
	t.True(address.Equal(uac.Address()))

	ukeys := uac.Keys()

	t.Equal(len(oa.Keys().Keys()), len(ukeys.Keys()))
	t.Equal(oa.Keys().Threshold(), ukeys.Threshold())
	for i := range oa.Keys().Keys() {
		t.Equal(oa.Keys().Keys()[i].Weight(), ukeys.Keys()[i].Weight())
		t.True(oa.Keys().Keys()[i].Key().Equal(ukeys.Keys()[i].Key()))
	}

	t.NotNil(sb)

	sba, _ := currency.StateBalanceValue(sb)
	t.True(sba.Big().Equal(balance[0].Big().Sub(fee)))

	t.Equal(fee, sb.(currency.AmountState).Fee())

	nfd, _ := StateFileDataValue(fds)
	t.True(nfd.SignCode().Equal(fd.SignCode()))
	t.True(nfd.Owner().Equal(fd.Owner()))
}

func (t *testCreateDocumentsOperation) TestOwnerAccountsNotExist() {
	cid := currency.CurrencyID("FINDME")
	feeAmount := int64(1)

	balance := []currency.Amount{
		currency.NewAmount(currency.NewBig(33), cid),
	}
	sa, st := t.newAccount(true, balance)
	oa, _ := t.newAccount(false, nil)
	na, _ := t.newAccount(false, nil)

	pool, _ := t.statepool(st)
	fee := currency.NewBig(feeAmount)
	feeer := currency.NewFixedFeeer(sa.Address, fee)

	cp := currency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(cid, currency.NewBig(99), sa.Address, feeer)))

	opr := t.processor(cp, pool)

	// filedata
	sc := SignCode("ABCD")

	// create document of na(document account) with oa(owner) which is sent from sa(sender)
	items := []CreateDocumentsItem{NewCreateDocumentsItemSingleFile(na.Keys(), sc, oa.Address, cid)}
	cd := t.newOperation(sa.Address, items, sa.Privs())

	err := opr.Process(cd)

	var oper operation.ReasonError
	t.True(xerrors.As(err, &oper))
	t.Contains(err.Error(), "key of Owner does not exist")
}

func (t *testCreateDocumentsOperation) TestDocumentAlreadyExists() {
	cid := currency.CurrencyID("FINDME")
	feeAmount := int64(1)

	balance := []currency.Amount{
		currency.NewAmount(currency.NewBig(33), cid),
	}

	sa, sta := t.newAccount(true, balance)

	sc := SignCode("ABCD")
	oa := sa

	filedata := NewFileData(sc, oa.Address)

	na, stb := t.newDocument(true, filedata)

	pool, _ := t.statepool(sta, stb)

	fee := currency.NewBig(feeAmount)
	feeer := currency.NewFixedFeeer(sa.Address, fee)

	cp := currency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(cid, currency.NewBig(99), sa.Address, feeer)))

	opr := t.processor(cp, pool)

	// create document of na(document account) with oa(owner) which is sent from sa(sender)
	items := []CreateDocumentsItem{NewCreateDocumentsItemSingleFile(na.Keys(), sc, oa.Address, cid)}
	cd := t.newOperation(sa.Address, items, sa.Privs())

	err := opr.Process(cd)

	var oper operation.ReasonError
	t.True(xerrors.As(err, &oper))
	t.Contains(err.Error(), "key of Document already exists")
}

func (t *testCreateDocumentsOperation) TestSameSenders() {
	cid := currency.CurrencyID("FINDME")
	feeAmount := int64(1)

	balance := []currency.Amount{
		currency.NewAmount(currency.NewBig(33), cid),
	}

	sa, sta := t.newAccount(true, balance)
	oa := sa

	na0, _ := t.newAccount(false, nil)
	na1, _ := t.newAccount(false, nil)

	pool, _ := t.statepool(sta)

	fee := currency.NewBig(feeAmount)
	feeer := currency.NewFixedFeeer(sa.Address, fee)

	cp := currency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(cid, currency.NewBig(99), sa.Address, feeer)))

	opr := t.processor(cp, pool)

	sc := SignCode("ABCD")

	// create document of na(document account) with oa(owner) which is sent from sa(sender)
	items0 := []CreateDocumentsItem{NewCreateDocumentsItemSingleFile(na0.Keys(), sc, oa.Address, cid)}
	cd0 := t.newOperation(sa.Address, items0, sa.Privs())
	err := opr.Process(cd0)

	items1 := []CreateDocumentsItem{NewCreateDocumentsItemSingleFile(na1.Keys(), sc, oa.Address, cid)}
	cd1 := t.newOperation(sa.Address, items1, sa.Privs())

	raddresses, _ := cd1.Fact().(CreateDocumentsFact).Addresses()
	addresses := []base.Address{na1.Address, sa.Address}
	for i := range raddresses {
		t.True(addresses[i].Equal(raddresses[i]))
	}

	err = opr.Process(cd1)
	t.Contains(err.Error(), "violates only one sender")
}

func (t *testCreateDocumentsOperation) TestSameSendersWithInvalidOperation() {
	cid := currency.CurrencyID("FINDME")
	feeAmount := int64(1)

	balance := []currency.Amount{
		currency.NewAmount(currency.NewBig(33), cid),
	}

	sa, sta := t.newAccount(true, balance)
	oa := sa

	na0, _ := t.newAccount(false, nil)
	na1, _ := t.newAccount(false, nil)

	pool, _ := t.statepool(sta)

	fee := currency.NewBig(feeAmount)
	feeer := currency.NewFixedFeeer(sa.Address, fee)

	cp := currency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(cid, currency.NewBig(99), sa.Address, feeer)))

	opr := t.processor(cp, pool)

	sc := SignCode("ABCD")

	// insert invalid operation, under threshold signing. It can not be counted
	// to sender checking.
	{
		items := []CreateDocumentsItem{NewCreateDocumentsItemSingleFile(na0.Keys(), sc, sa.Address, cid)}
		cd := t.newOperation(sa.Address, items, []key.Privatekey{key.MustNewBTCPrivatekey()})
		err := opr.Process(cd)

		var oper operation.ReasonError
		t.True(xerrors.As(err, &oper))
	}

	// create document of na(document account) with oa(owner) which is sent from sa(sender)
	items0 := []CreateDocumentsItem{NewCreateDocumentsItemSingleFile(na0.Keys(), sc, oa.Address, cid)}
	cd0 := t.newOperation(sa.Address, items0, sa.Privs())
	err := opr.Process(cd0)

	items1 := []CreateDocumentsItem{NewCreateDocumentsItemSingleFile(na1.Keys(), sc, oa.Address, cid)}
	cd1 := t.newOperation(sa.Address, items1, sa.Privs())

	raddresses, _ := cd1.Fact().(CreateDocumentsFact).Addresses()
	addresses := []base.Address{na1.Address, sa.Address}
	for i := range raddresses {
		t.True(addresses[i].Equal(raddresses[i]))
	}

	err = opr.Process(cd1)
	t.Contains(err.Error(), "violates only one sender")
}

func (t *testCreateDocumentsOperation) TestSameDocumentAddress() {
	cid := currency.CurrencyID("FINDME")

	balance := []currency.Amount{
		currency.NewAmount(currency.NewBig(33), cid),
	}

	sa, _ := t.newAccount(true, balance)
	oa := sa

	na, _ := t.newAccount(false, nil)

	sc0 := SignCode("ABCD")
	sc1 := SignCode("DCBA")

	it0 := NewCreateDocumentsItemSingleFile(na.Keys(), sc0, oa.Address, cid)
	it1 := NewCreateDocumentsItemSingleFile(na.Keys(), sc1, oa.Address, cid)

	// create document of na(document account) with oa(owner) which is sent from sa(sender)
	items := []CreateDocumentsItem{it0, it1}
	t.Panicsf(func() { t.newOperation(sa.Address, items, sa.Privs()) }, "duplicated acocunt Keys found")

}

func (t *testCreateDocumentsOperation) TestSameDocumentAddressInMultipleOperations() {
	cid := currency.CurrencyID("FINDME")
	feeAmount := int64(1)

	balance := []currency.Amount{
		currency.NewAmount(currency.NewBig(33), cid),
	}

	sa, sta := t.newAccount(true, balance)
	sb, stb := t.newAccount(true, balance)
	oa := sa

	na, _ := t.newAccount(false, nil)

	pool, _ := t.statepool(sta, stb)

	fee := currency.NewBig(feeAmount)
	feeer := currency.NewFixedFeeer(sa.Address, fee)

	cp := currency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(cid, currency.NewBig(99), sa.Address, feeer)))

	opr := t.processor(cp, pool)

	sc := SignCode("ABCD")

	// create document of na(document account) with oa(owner) which is sent from sa(sender)
	items0 := []CreateDocumentsItem{NewCreateDocumentsItemSingleFile(na.Keys(), sc, oa.Address, cid)}
	cd0 := t.newOperation(sa.Address, items0, sa.Privs())
	err := opr.Process(cd0)

	items1 := []CreateDocumentsItem{NewCreateDocumentsItemSingleFile(na.Keys(), sc, oa.Address, cid)}
	cd1 := t.newOperation(sb.Address, items1, sb.Privs())

	err = opr.Process(cd1)
	t.Contains(err.Error(), "new address already processed")
}

func (t *testCreateDocumentsOperation) TestMultipleItemsWithFee() {
	cid0 := currency.CurrencyID("SHOWME")
	cid1 := currency.CurrencyID("FINDME")

	balance := []currency.Amount{
		currency.NewAmount(currency.NewBig(33), cid0),
		currency.NewAmount(currency.NewBig(33), cid1),
	}

	sa, st := t.newAccount(true, balance)
	na0, _ := t.newAccount(false, nil)
	na1, _ := t.newAccount(false, nil)
	oa := sa

	pool, _ := t.statepool(st)

	fee0 := currency.NewBig(1)
	fee1 := currency.NewBig(2)
	feeer0 := currency.NewFixedFeeer(sa.Address, fee0)
	feeer1 := currency.NewFixedFeeer(sa.Address, fee1)

	cp := currency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(cid0, currency.NewBig(99), sa.Address, feeer0)))
	t.NoError(cp.Set(t.newCurrencyDesignState(cid1, currency.NewBig(99), sa.Address, feeer1)))

	opr := t.processor(cp, pool)

	sc0 := SignCode("ABCD")
	sc1 := SignCode("EFGH")
	fd0 := NewFileData(sc0, oa.Address)
	fd1 := NewFileData(sc1, oa.Address)

	items := []CreateDocumentsItem{
		NewCreateDocumentsItemSingleFile(na0.Keys(), sc0, oa.Address, cid0),
		NewCreateDocumentsItemSingleFile(na1.Keys(), sc1, oa.Address, cid1),
	}
	cd := t.newOperation(sa.Address, items, sa.Privs())

	t.NoError(opr.Process(cd))

	// check updated state
	// new document account state
	var ns0, ns1 state.State
	// new document filedata state
	var fds0, fds1 state.State
	// sender balance state
	sb := map[currency.CurrencyID]state.State{}
	for _, stu := range pool.Updates() {
		if currency.IsStateBalanceKey(stu.Key()) {
			st := stu.GetState()

			i, err := currency.StateBalanceValue(st)
			t.NoError(err)

			if st.Key() == currency.StateKeyBalance(sa.Address, i.Currency()) {
				sb[i.Currency()] = st
			} else {
				continue
			}
		} else if currency.IsStateDocumentKey(stu.Key()) {
			if stu.Key() == currency.StateKeyDocument(na0.Address) {
				ns0 = stu.GetState()
			} else if stu.Key() == currency.StateKeyDocument(na1.Address) {
				ns1 = stu.GetState()
			}
		} else if IsStateFileDataKey(stu.Key()) {
			if stu.Key() == StateKeyFileData(na0.Address) {
				fds0 = stu.GetState()
			} else if stu.Key() == StateKeyFileData(na1.Address) {
				fds1 = stu.GetState()
			}
		}
	}

	address0, err := currency.NewAddressFromKeys(na0.Keys())
	t.NoError(err)
	address1, err := currency.NewAddressFromKeys(na1.Keys())
	t.NoError(err)

	uac0 := ns0.Value().Interface().(currency.Account)
	t.True(address0.Equal(uac0.Address()))

	uac1 := ns1.Value().Interface().(currency.Account)
	t.True(address1.Equal(uac1.Address()))

	ukeys0 := uac0.Keys()
	ukeys1 := uac1.Keys()

	t.Equal(len(oa.Keys().Keys()), len(ukeys0.Keys()))
	t.Equal(len(oa.Keys().Keys()), len(ukeys1.Keys()))
	t.Equal(oa.Keys().Threshold(), ukeys0.Threshold())
	t.Equal(oa.Keys().Threshold(), ukeys1.Threshold())
	for i := range na0.Keys().Keys() {
		t.Equal(oa.Keys().Keys()[i].Weight(), ukeys0.Keys()[i].Weight())
		t.True(oa.Keys().Keys()[i].Key().Equal(ukeys0.Keys()[i].Key()))
	}
	for i := range na1.Keys().Keys() {
		t.Equal(oa.Keys().Keys()[i].Weight(), ukeys1.Keys()[i].Weight())
		t.True(oa.Keys().Keys()[i].Key().Equal(ukeys1.Keys()[i].Key()))
	}

	t.Equal(len(balance), len(sb))

	sba0, _ := currency.StateBalanceValue(sb[cid0])
	t.True(sba0.Big().Equal(balance[0].Big().Sub(fee0)))

	sba1, _ := currency.StateBalanceValue(sb[cid1])
	t.True(sba1.Big().Equal(balance[1].Big().Sub(fee1)))

	t.Equal(fee0, sb[cid0].(currency.AmountState).Fee())
	t.Equal(fee1, sb[cid1].(currency.AmountState).Fee())

	nfd0, _ := StateFileDataValue(fds0)
	t.True(nfd0.SignCode().Equal(fd0.SignCode()))
	t.True(nfd0.Owner().Equal(fd0.Owner()))

	nfd1, _ := StateFileDataValue(fds1)
	t.True(nfd1.SignCode().Equal(fd1.SignCode()))
	t.True(nfd1.Owner().Equal(fd1.Owner()))
}

func (t *testCreateDocumentsOperation) TestInSufficientBalanceForFee() {
	// currency id
	cid := currency.CurrencyID("SHOWME")
	// sender balance
	senderBalance := int64(3)
	// fee amount
	feeAmount := int64(4)

	// sender initial balance
	balance := []currency.Amount{
		currency.NewAmount(currency.NewBig(senderBalance), cid),
	}

	// sender account
	sa, st := t.newAccount(true, balance)
	// document account
	na, _ := t.newAccount(false, nil)
	// owner account
	oa := sa

	pool, _ := t.statepool(st)

	fee := currency.NewBig(feeAmount)
	feeer := currency.NewFixedFeeer(sa.Address, fee)

	cp := currency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(cid, currency.NewBig(99), sa.Address, feeer)))

	opr := t.processor(cp, pool)

	// filedata
	sc := SignCode("ABCD")

	// create document of na(document account) with oa(owner) which is sent from sa(sender)
	items := []CreateDocumentsItem{NewCreateDocumentsItemSingleFile(na.Keys(), sc, oa.Address, cid)}
	cd := t.newOperation(sa.Address, items, sa.Privs())

	err := opr.Process(cd)

	var oper operation.ReasonError
	t.True(xerrors.As(err, &oper))
	t.Contains(err.Error(), "insufficient balance")
}

func (t *testCreateDocumentsOperation) TestUnknownCurrencyID() {
	// currency id of network
	cid0 := currency.CurrencyID("SHOWME")
	// currency id used in operation
	cid1 := currency.CurrencyID("FINDME")

	// sender initial balance
	balance := []currency.Amount{
		currency.NewAmount(currency.NewBig(33), cid0),
	}

	// sender account
	sa, st := t.newAccount(true, balance)
	// document account
	na, _ := t.newAccount(false, nil)
	// owner account
	oa := sa

	pool, _ := t.statepool(st)

	fee := currency.NewBig(2)
	feeer := currency.NewFixedFeeer(sa.Address, fee)

	cp := currency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(cid0, currency.NewBig(99), sa.Address, feeer)))

	opr := t.processor(cp, pool)

	// filedata
	sc := SignCode("ABCD")

	// create document of na(document account) with oa(owner) which is sent from sa(sender)
	items := []CreateDocumentsItem{NewCreateDocumentsItemSingleFile(na.Keys(), sc, oa.Address, cid1)}
	cd := t.newOperation(sa.Address, items, sa.Privs())

	err := opr.Process(cd)

	var oper operation.ReasonError
	t.True(xerrors.As(err, &oper))
	t.Contains(err.Error(), "unknown currency id found")
}

func (t *testCreateDocumentsOperation) TestEmptyCurrency() {
	cid0 := currency.CurrencyID("FINDME")
	cid1 := currency.CurrencyID("SHOWME")
	feeAmount := int64(1)

	balance := []currency.Amount{
		currency.NewAmount(currency.NewBig(33), cid0),
	}
	sa, st := t.newAccount(true, balance)
	oa, _ := t.newAccount(false, nil)
	na, _ := t.newAccount(false, nil)

	pool, _ := t.statepool(st)
	fee := currency.NewBig(feeAmount)
	feeer := currency.NewFixedFeeer(sa.Address, fee)

	cp := currency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(cid0, currency.NewBig(99), sa.Address, feeer)))
	t.NoError(cp.Set(t.newCurrencyDesignState(cid1, currency.NewBig(99), sa.Address, feeer)))

	opr := t.processor(cp, pool)

	// filedata
	sc := SignCode("ABCD")

	// create document of na(document account) with oa(owner) which is sent from sa(sender)
	items := []CreateDocumentsItem{NewCreateDocumentsItemSingleFile(na.Keys(), sc, oa.Address, cid1)}
	cd := t.newOperation(sa.Address, items, sa.Privs())

	err := opr.Process(cd)

	var oper operation.ReasonError
	t.True(xerrors.As(err, &oper))
	t.Contains(err.Error(), "currency of holder does not exist")
}

func (t *testCreateDocumentsOperation) TestSenderBalanceNotExist() {
	cid := currency.CurrencyID("FINDME")
	feeAmount := int64(1)

	sa, st := t.newAccount(true, nil)
	na, _ := t.newAccount(false, nil)
	oa := sa

	pool, _ := t.statepool(st)

	fee := currency.NewBig(feeAmount)
	feeer := currency.NewFixedFeeer(sa.Address, fee)

	cp := currency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(cid, currency.NewBig(99), sa.Address, feeer)))

	opr := t.processor(cp, pool)

	// filedata
	sc := SignCode("ABCD")

	// create document of na(document account) with oa(owner) which is sent from sa(sender)
	items := []CreateDocumentsItem{NewCreateDocumentsItemSingleFile(na.Keys(), sc, oa.Address, cid)}
	cd := t.newOperation(sa.Address, items, sa.Privs())

	err := opr.Process(cd)

	var oper operation.ReasonError
	t.True(xerrors.As(err, &oper))
	t.Contains(err.Error(), "currency of holder does not exist")
}

func TestCreateDocumentsOperation(t *testing.T) {
	suite.Run(t, new(testCreateDocumentsOperation))
}
