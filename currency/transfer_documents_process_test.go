package currency

import (
	"testing"

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

type testTransferDocumentsOperations struct {
	baseTestOperationProcessor
	cid CurrencyID
	sc  SignCode
	fee Big
}

func (t *testTransferDocumentsOperations) SetupSuite() {
	t.cid = CurrencyID("SHOWME")
	t.sc = SignCode("ABCD")
	t.fee = NewBig(1)
}

func (t *testTransferDocumentsOperations) processor(cp *CurrencyPool, pool *storage.Statepool) prprocessor.OperationProcessor {
	copr, err := NewOperationProcessor(cp).
		SetProcessor(TransferDocuments{}, NewTransferDocumentsProcessor(cp))
	t.NoError(err)

	if pool == nil {
		return copr
	}

	return copr.New(pool)
}

func (t *testTransferDocumentsOperations) newTransferDocumentsItem(document base.Address, receiver base.Address, cid CurrencyID) TransferDocumentsItem {

	return NewTransferDocumentsItemSingleFile(document, receiver, cid)
}

func (t *testTransferDocumentsOperations) newTransferDocument(sender base.Address, keys []key.Privatekey, items []TransferDocumentsItem) TransferDocuments {
	token := util.UUID().Bytes()
	fact := NewTransferDocumentsFact(token, sender, items)

	var fs []operation.FactSign
	for _, pk := range keys {
		sig, err := operation.NewFactSignature(pk, fact, nil)
		t.NoError(err)

		fs = append(fs, operation.NewBaseFactSign(pk.Publickey(), sig))
	}

	tfd, err := NewTransferDocuments(fact, fs, "")
	t.NoError(err)

	t.NoError(tfd.IsValid(nil))

	return tfd
}

func (t *testTransferDocumentsOperations) newTestFileData(oa base.Address) FileData {
	return NewFileData(t.sc, oa)
}

func (t *testTransferDocumentsOperations) newTestBalance() []Amount {
	return []Amount{NewAmount(NewBig(33), t.cid)}
}

func (t *testTransferDocumentsOperations) newTestFixedFeeer(sa base.Address) FixedFeeer {
	return NewFixedFeeer(sa, t.fee)
}
func (t *testTransferDocumentsOperations) TestNormalCase() {
	balance := t.newTestBalance()
	sa, sta := t.newAccount(true, balance)
	ra, stb := t.newAccount(true, balance)
	oa := sa
	fd := t.newTestFileData(oa.Address)

	da, std := t.newDocument(true, fd)
	pool, _ := t.statepool(sta, stb, std)

	feeer := t.newTestFixedFeeer(sa.Address)

	cp := NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, NewBig(99), NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)

	items := []TransferDocumentsItem{t.newTransferDocumentsItem(da.Address, ra.Address, t.cid)}
	tfd := t.newTransferDocument(sa.Address, sa.Privs(), items)

	err := opr.Process(tfd)

	// check updated state
	// document account state
	var ns state.State
	// document filedata state
	var fds state.State
	// sender balance state
	var sb state.State
	for _, stu := range pool.Updates() {
		if IsStateBalanceKey(stu.Key()) {
			st := stu.GetState()

			i, err := StateBalanceValue(st)
			t.NoError(err)

			if st.Key() == StateKeyBalance(sa.Address, i.Currency()) {
				sb = st
			} else {
				continue
			}
		} else if (IsStateDocumentKey(stu.Key())) && (stu.Key() == StateKeyDocument(da.Address)) {
			ns = stu.GetState()
		} else if (IsStateFileDataKey(stu.Key())) && (stu.Key() == StateKeyFileData(da.Address)) {
			fds = stu.GetState()
		}
	}

	address := da.Address
	t.NoError(err)
	uac := ns.Value().Interface().(Account)
	t.True(address.Equal(uac.Address()))

	ukeys := uac.Keys()

	t.Equal(len(ra.Keys().Keys()), len(ukeys.Keys()))
	t.Equal(ra.Keys().Threshold(), ukeys.Threshold())
	for i := range ra.Keys().Keys() {
		t.Equal(ra.Keys().Keys()[i].Weight(), ukeys.Keys()[i].Weight())
		t.True(ra.Keys().Keys()[i].Key().Equal(ukeys.Keys()[i].Key()))
	}

	t.NotNil(sb)

	sba, _ := StateBalanceValue(sb)
	t.True(sba.Big().Equal(balance[0].Big().Sub(t.fee)))

	t.Equal(t.fee, sb.(AmountState).Fee())

	nfd, _ := StateFileDataValue(fds)
	t.True(nfd.SignCode().Equal(fd.SignCode()))
	t.True(nfd.Owner().Equal(ra.Address))

}

func (t *testTransferDocumentsOperations) TestSenderNotExist() {
	balance := t.newTestBalance()
	sa, _ := t.newAccount(false, nil)
	ra, sta := t.newAccount(true, balance)
	oa, stb := t.newAccount(true, balance)
	fd := t.newTestFileData(oa.Address)

	da, std := t.newDocument(true, fd)
	pool, _ := t.statepool(sta, stb, std)

	feeer := t.newTestFixedFeeer(sa.Address)

	cp := NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, NewBig(99), NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)

	items := []TransferDocumentsItem{t.newTransferDocumentsItem(da.Address, ra.Address, t.cid)}
	tfd := t.newTransferDocument(sa.Address, sa.Privs(), items)

	err := opr.Process(tfd)

	var oper operation.ReasonError
	t.True(xerrors.As(err, &oper))
	t.Contains(err.Error(), "does not exist")
}

func (t *testTransferDocumentsOperations) TestReceiverNotExist() {
	balance := t.newTestBalance()
	sa, sta := t.newAccount(true, balance)
	ra, _ := t.newAccount(false, nil)
	oa := sa
	fd := t.newTestFileData(oa.Address)

	da, std := t.newDocument(true, fd)
	pool, _ := t.statepool(sta, std)

	feeer := t.newTestFixedFeeer(sa.Address)

	cp := NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, NewBig(99), NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)

	items := []TransferDocumentsItem{t.newTransferDocumentsItem(da.Address, ra.Address, t.cid)}
	tfd := t.newTransferDocument(sa.Address, sa.Privs(), items)

	err := opr.Process(tfd)

	var oper operation.ReasonError
	t.True(xerrors.As(err, &oper))
	t.Contains(err.Error(), "receiver does not exist")
}

func (t *testTransferDocumentsOperations) TestInsufficientBalanceForFee() {
	balance := []Amount{NewAmount(NewBig(2), t.cid)}
	sa, sta := t.newAccount(true, balance)
	ra, stb := t.newAccount(true, balance)
	oa := sa
	fd := t.newTestFileData(oa.Address)

	da, std := t.newDocument(true, fd)
	pool, _ := t.statepool(sta, stb, std)

	fee := NewBig(3)
	feeer := NewFixedFeeer(sa.Address, fee)

	cp := NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, NewBig(99), NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)

	items := []TransferDocumentsItem{t.newTransferDocumentsItem(da.Address, ra.Address, t.cid)}
	tfd := t.newTransferDocument(sa.Address, sa.Privs(), items)

	err := opr.Process(tfd)

	var oper operation.ReasonError
	t.True(xerrors.As(err, &oper))
	t.Contains(err.Error(), "insufficient balance")
}

func (t *testTransferDocumentsOperations) TestMultipleItemsWithFee() {
	cid0 := CurrencyID("SHOWME")
	cid1 := CurrencyID("FINDME")
	balance0 := NewAmount(NewBig(33), cid0)
	balance1 := NewAmount(NewBig(33), cid1)
	sa, sta := t.newAccount(true, []Amount{balance0, balance1})
	ra0, stb := t.newAccount(true, []Amount{NewAmount(NewBig(0), cid0)})
	ra1, stc := t.newAccount(true, []Amount{NewAmount(NewBig(0), cid0)})
	fd := t.newTestFileData(sa.Address)
	da0, std0 := t.newDocument(true, fd)
	da1, std1 := t.newDocument(true, fd)

	pool, _ := t.statepool(sta, stb, stc, std0, std1)

	feeer := t.newTestFixedFeeer(sa.Address)

	cp := NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(cid0, NewBig(99), NewTestAddress(), feeer)))
	t.NoError(cp.Set(t.newCurrencyDesignState(cid1, NewBig(99), NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)

	token := util.UUID().Bytes()
	items := []TransferDocumentsItem{
		t.newTransferDocumentsItem(da0.Address, ra0.Address, cid0),
		t.newTransferDocumentsItem(da1.Address, ra1.Address, cid1),
	}
	fact := NewTransferDocumentsFact(token, sa.Address, items)
	sig, err := operation.NewFactSignature(sa.Privs()[0], fact, nil)
	t.NoError(err)
	fs := []operation.FactSign{operation.NewBaseFactSign(sa.Privs()[0].Publickey(), sig)}
	tfd, err := NewTransferDocuments(fact, fs, "")
	t.NoError(err)

	err = opr.Process(tfd)
	t.NoError(err)

	var nst0, nst1 state.State
	var nam0, nam1 Amount
	for _, st := range pool.Updates() {
		if st.Key() == StateKeyBalance(sa.Address, cid0) {
			nst0 = st.GetState()
			nam0, _ = StateBalanceValue(nst0)
		} else if st.Key() == StateKeyBalance(sa.Address, cid1) {
			nst1 = st.GetState()
			nam1, _ = StateBalanceValue(nst1)
		}
	}

	t.Equal(balance0.Big().Sub(t.fee), nam0.Big())
	t.Equal(balance1.Big().Sub(t.fee), nam1.Big())
	t.Equal(t.fee, nst0.(AmountState).Fee())
	t.Equal(t.fee, nst1.(AmountState).Fee())
}

func (t *testTransferDocumentsOperations) TestInsufficientMultipleItemsWithFee() {
	cid0 := CurrencyID("SHOWME")
	cid1 := CurrencyID("FINDME")
	balance0 := NewAmount(NewBig(10), cid0)
	balance1 := NewAmount(NewBig(10), cid1)
	sa, sta := t.newAccount(true, []Amount{balance0, balance1})
	ra0, stb := t.newAccount(true, []Amount{NewAmount(NewBig(0), cid0)})
	ra1, stc := t.newAccount(true, []Amount{NewAmount(NewBig(0), cid0)})
	fd := t.newTestFileData(sa.Address)
	da0, std0 := t.newDocument(true, fd)
	da1, std1 := t.newDocument(true, fd)

	pool, _ := t.statepool(sta, stb, stc, std0, std1)

	fee0 := NewBig(11)
	fee1 := NewBig(3)
	feeer0 := NewFixedFeeer(sa.Address, fee0)
	feeer1 := NewFixedFeeer(sa.Address, fee1)

	cp := NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(cid0, NewBig(99), NewTestAddress(), feeer0)))
	t.NoError(cp.Set(t.newCurrencyDesignState(cid1, NewBig(99), NewTestAddress(), feeer1)))

	opr := t.processor(cp, pool)

	token := util.UUID().Bytes()
	items := []TransferDocumentsItem{
		t.newTransferDocumentsItem(da0.Address, ra0.Address, cid0),
		t.newTransferDocumentsItem(da1.Address, ra1.Address, cid1),
	}
	fact := NewTransferDocumentsFact(token, sa.Address, items)
	sig, err := operation.NewFactSignature(sa.Privs()[0], fact, nil)
	t.NoError(err)
	fs := []operation.FactSign{operation.NewBaseFactSign(sa.Privs()[0].Publickey(), sig)}
	tfd, err := NewTransferDocuments(fact, fs, "")
	t.NoError(err)

	err = opr.Process(tfd)

	var oper operation.ReasonError
	t.True(xerrors.As(err, &oper))
	t.Contains(err.Error(), "insufficient balance")
}

func (t *testTransferDocumentsOperations) TestSameSenders() {
	cid0 := CurrencyID("SHOWME")
	cid1 := CurrencyID("FINDME")
	balance0 := NewAmount(NewBig(33), cid0)
	balance1 := NewAmount(NewBig(33), cid1)
	sa, sta := t.newAccount(true, []Amount{balance0, balance1})
	ra0, stb := t.newAccount(true, []Amount{NewAmount(NewBig(0), cid0)})
	ra1, stc := t.newAccount(true, []Amount{NewAmount(NewBig(0), cid0)})
	fd := t.newTestFileData(sa.Address)
	da0, std0 := t.newDocument(true, fd)
	da1, std1 := t.newDocument(true, fd)

	pool, _ := t.statepool(sta, stb, stc, std0, std1)

	feeer := t.newTestFixedFeeer(sa.Address)

	cp := NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(cid0, NewBig(99), NewTestAddress(), feeer)))
	t.NoError(cp.Set(t.newCurrencyDesignState(cid1, NewBig(99), NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)

	items0 := []TransferDocumentsItem{t.newTransferDocumentsItem(da0.Address, ra0.Address, cid0)}
	tfd0 := t.newTransferDocument(sa.Address, sa.Privs(), items0)

	t.NoError(opr.Process(tfd0))

	items1 := []TransferDocumentsItem{t.newTransferDocumentsItem(da1.Address, ra1.Address, cid1)}
	tfd1 := t.newTransferDocument(sa.Address, sa.Privs(), items1)

	err := opr.Process(tfd1)

	var oper operation.ReasonError
	t.True(xerrors.As(err, &oper))
	t.Contains(err.Error(), "violates only one sender")
}

func (t *testTransferDocumentsOperations) TestUnderThreshold() {
	spk := key.MustNewBTCPrivatekey()
	rpk := key.MustNewBTCPrivatekey()
	dpk := key.MustNewBTCPrivatekey()

	skey := t.newKey(spk.Publickey(), 50)
	rkey := t.newKey(rpk.Publickey(), 50)
	dkey := t.newKey(dpk.Publickey(), 100)
	skeys, _ := NewKeys([]Key{skey, rkey}, 100)
	rkeys, _ := NewKeys([]Key{rkey}, 50)
	dkeys, _ := NewKeys([]Key{dkey}, 100)

	pks := []key.Privatekey{spk}
	sender, _ := NewAddressFromKeys(skeys)
	receiver, _ := NewAddressFromKeys(rkeys)
	document, _ := NewAddressFromKeys(dkeys)

	// set sender state
	senderBalance := NewAmount(NewBig(33), t.cid)
	filedata := t.newTestFileData(sender)

	var sts []state.State
	sts = append(sts,
		t.newStateBalance(sender, senderBalance.Big(), senderBalance.Currency()),
		t.newStateKeys(sender, skeys),
		t.newStateFileData(document, filedata),
		t.newStateDocumentKeys(document, dkeys),
		t.newStateKeys(receiver, skeys),
	)

	pool, _ := t.statepool(sts)
	feeer := NewFixedFeeer(sender, ZeroBig)

	cp := NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, NewBig(99), NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)

	items := []TransferDocumentsItem{t.newTransferDocumentsItem(document, receiver, t.cid)}
	tfd := t.newTransferDocument(sender, pks, items)

	err := opr.Process(tfd)

	var oper operation.ReasonError
	t.True(xerrors.As(err, &oper))
	t.Contains(err.Error(), "not passed threshold")
}

func (t *testTransferDocumentsOperations) TestUnknownKey() {
	balance := []Amount{NewAmount(NewBig(2), t.cid)}
	sa, sta := t.newAccount(true, balance)
	ra, stb := t.newAccount(true, balance)
	oa := sa
	fd := t.newTestFileData(oa.Address)

	da, std := t.newDocument(true, fd)
	pool, _ := t.statepool(sta, stb, std)

	feeer := t.newTestFixedFeeer(sa.Address)

	cp := NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, NewBig(99), NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)

	items := []TransferDocumentsItem{t.newTransferDocumentsItem(da.Address, ra.Address, t.cid)}
	tfd := t.newTransferDocument(sa.Address, []key.Privatekey{sa.Priv, key.MustNewBTCPrivatekey()}, items)

	err := opr.Process(tfd)

	var oper operation.ReasonError
	t.True(xerrors.As(err, &oper))
	t.Contains(err.Error(), "unknown key found")
}

func TestTransferDocumentsOperations(t *testing.T) {
	suite.Run(t, new(testTransferDocumentsOperations))
}
