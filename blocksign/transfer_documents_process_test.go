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

type testTransferDocumentsOperations struct {
	baseTestOperationProcessor
	cid   currency.CurrencyID
	docid currency.Big
	fh    FileHash
	fee   currency.Big
}

func (t *testTransferDocumentsOperations) SetupSuite() {
	t.cid = currency.CurrencyID("SHOWME")
	t.docid = currency.NewBig(0)
	t.fh = FileHash("ABCD")
	t.fee = currency.NewBig(3)
}

func (t *testTransferDocumentsOperations) processor(cp *currency.CurrencyPool, pool *storage.Statepool) prprocessor.OperationProcessor {
	copr, err := NewOperationProcessor(cp).
		SetProcessor(TransferDocuments{}, NewTransferDocumentsProcessor(cp))
	t.NoError(err)

	if pool == nil {
		return copr
	}

	return copr.New(pool)
}

func (t *testTransferDocumentsOperations) newTransferDocumentsItem(docid currency.Big, owner base.Address, receiver base.Address, cid currency.CurrencyID) TransferDocumentsItem {

	return NewTransferDocumentsItemSingleFile(docid, owner, receiver, cid)
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

func (t *testTransferDocumentsOperations) newTestDocumentData(ca base.Address) DocumentData {
	doc := NewDocumentData(DocInfo{idx: t.docid, filehash: t.fh}, ca, ca, []DocSign{})
	return doc
}

func (t *testTransferDocumentsOperations) newTestBalance() []currency.Amount {
	return []currency.Amount{currency.NewAmount(currency.NewBig(33), t.cid)}
}

func (t *testTransferDocumentsOperations) newTestFixedFeeer(sa base.Address) currency.FixedFeeer {
	return currency.NewFixedFeeer(sa, t.fee)
}
func (t *testTransferDocumentsOperations) TestNormalCase() {
	balance := t.newTestBalance()
	sa, sta := t.newAccount(true, balance)
	ra, stb := t.newAccount(true, balance)
	dd := t.newTestDocumentData(sa.Address)

	sts := t.newStateDocument(sa.Address, dd)
	pool, _ := t.statepool(sta, stb, sts)

	feeer := t.newTestFixedFeeer(sa.Address)

	cp := currency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)

	items := []TransferDocumentsItem{t.newTransferDocumentsItem(t.docid, sa.Address, ra.Address, t.cid)}
	tfd := t.newTransferDocument(sa.Address, sa.Privs(), items)

	t.NoError(opr.Process(tfd))

	// check updated state
	// owner documents state
	var ons state.State
	// receiver documents state
	var rns state.State
	// document data state
	var dds state.State
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
		} else if (IsStateDocumentsKey(stu.Key())) && (stu.Key() == StateKeyDocuments(sa.Address)) {
			ons = stu.GetState()
		} else if (IsStateDocumentsKey(stu.Key())) && (stu.Key() == StateKeyDocuments(ra.Address)) {
			rns = stu.GetState()
		} else if (IsStateDocumentDataKey(stu.Key())) && (stu.Key() == StateKeyDocumentData(t.fh)) {
			dds = stu.GetState()
		}
	}

	t.NotNil(sb)

	sba, _ := currency.StateBalanceValue(sb)
	t.True(sba.Big().Equal(balance[0].Big().Sub(t.fee)))

	t.Equal(t.fee, sb.(currency.AmountState).Fee())

	usdoc := ons.Value().Interface().(DocumentInventory)
	t.True(!usdoc.Exists(t.docid))

	urdoc := rns.Value().Interface().(DocumentInventory)

	t.True(t.fh.Equal(urdoc.Documents()[0].FileHash()))
	t.True(t.docid.Equal(urdoc.Documents()[0].Index()))

	ndd, _ := StateDocumentDataValue(dds)
	t.True(ndd.FileHash().Equal(t.fh))
	t.True(ndd.Creator().Equal(sa.Address))
	t.True(ndd.Owner().Equal(ra.Address))
}

func (t *testTransferDocumentsOperations) TestSenderNotExist() {
	balance := t.newTestBalance()
	sa, _ := t.newAccount(false, nil)
	ra, sta := t.newAccount(true, balance)
	dd := t.newTestDocumentData(sa.Address)

	sts := t.newStateDocument(sa.Address, dd)
	pool, _ := t.statepool(sta, sts)

	feeer := t.newTestFixedFeeer(sa.Address)

	cp := currency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)

	items := []TransferDocumentsItem{t.newTransferDocumentsItem(t.docid, sa.Address, ra.Address, t.cid)}
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
	dd := t.newTestDocumentData(sa.Address)

	sts := t.newStateDocument(sa.Address, dd)
	pool, _ := t.statepool(sta, sts)

	feeer := t.newTestFixedFeeer(sa.Address)

	cp := currency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)

	items := []TransferDocumentsItem{t.newTransferDocumentsItem(t.docid, sa.Address, ra.Address, t.cid)}
	tfd := t.newTransferDocument(sa.Address, sa.Privs(), items)

	err := opr.Process(tfd)

	var oper operation.ReasonError
	t.True(xerrors.As(err, &oper))
	t.Contains(err.Error(), "receiver account not found")
}

func (t *testTransferDocumentsOperations) TestInsufficientBalanceForFee() {
	balance := []currency.Amount{currency.NewAmount(currency.NewBig(2), t.cid)}
	sa, sta := t.newAccount(true, balance)
	ra, stb := t.newAccount(true, balance)
	dd := t.newTestDocumentData(sa.Address)

	sts := t.newStateDocument(sa.Address, dd)
	pool, _ := t.statepool(sta, stb, sts)

	feeer := t.newTestFixedFeeer(sa.Address)

	cp := currency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)

	items := []TransferDocumentsItem{t.newTransferDocumentsItem(t.docid, sa.Address, ra.Address, t.cid)}
	tfd := t.newTransferDocument(sa.Address, sa.Privs(), items)

	err := opr.Process(tfd)

	var oper operation.ReasonError
	t.True(xerrors.As(err, &oper))
	t.Contains(err.Error(), "insufficient balance")
}

func (t *testTransferDocumentsOperations) TestMultipleItemsWithFee() {
	cid0 := currency.CurrencyID("SHOWME")
	cid1 := currency.CurrencyID("FINDME")
	balance0 := currency.NewAmount(currency.NewBig(33), cid0)
	balance1 := currency.NewAmount(currency.NewBig(33), cid1)
	sa, sta := t.newAccount(true, []currency.Amount{balance0, balance1})
	ra0, stb := t.newAccount(true, []currency.Amount{currency.NewAmount(currency.NewBig(0), cid0)})
	ra1, stc := t.newAccount(true, []currency.Amount{currency.NewAmount(currency.NewBig(0), cid0)})
	dd0 := t.newTestDocumentData(sa.Address)
	dd1 := NewDocumentData(DocInfo{idx: currency.NewBig(1), filehash: FileHash("EFGH")}, sa.Address, sa.Address, []DocSign{})
	sts0 := t.newStateDocument(sa.Address, dd0)
	dinv, _ := StateDocumentsValue(sts0[1])
	err := dinv.Append(DocInfo{idx: currency.NewBig(1), filehash: dd1.FileHash()})
	// sts0[1] = nst
	sts1 := t.newStateDocument(sa.Address, dd1)
	nst, _ := SetStateDocumentsValue(sts1[1], dinv)
	sts1[1] = nst
	sts := []state.State{sts0[0], sts0[2], sts1[0], sts1[1], sts1[2]}

	pool, _ := t.statepool(sta, stb, stc, sts)

	feeer := t.newTestFixedFeeer(sa.Address)

	cp := currency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(cid0, currency.NewBig(99), NewTestAddress(), feeer)))
	t.NoError(cp.Set(t.newCurrencyDesignState(cid1, currency.NewBig(99), NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)

	token := util.UUID().Bytes()
	items := []TransferDocumentsItem{
		t.newTransferDocumentsItem(t.docid, sa.Address, ra0.Address, cid0),
		t.newTransferDocumentsItem(currency.NewBig(1), sa.Address, ra1.Address, cid1),
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
	var nam0, nam1 currency.Amount
	for _, st := range pool.Updates() {
		if st.Key() == currency.StateKeyBalance(sa.Address, cid0) {
			nst0 = st.GetState()
			nam0, _ = currency.StateBalanceValue(nst0)
		} else if st.Key() == currency.StateKeyBalance(sa.Address, cid1) {
			nst1 = st.GetState()
			nam1, _ = currency.StateBalanceValue(nst1)
		}
	}

	t.Equal(balance0.Big().Sub(t.fee), nam0.Big())
	t.Equal(balance1.Big().Sub(t.fee), nam1.Big())
	t.Equal(t.fee, nst0.(currency.AmountState).Fee())
	t.Equal(t.fee, nst1.(currency.AmountState).Fee())
}

func (t *testTransferDocumentsOperations) TestInsufficientMultipleItemsWithFee() {
	cid0 := currency.CurrencyID("SHOWME")
	cid1 := currency.CurrencyID("FINDME")
	balance0 := currency.NewAmount(currency.NewBig(10), cid0)
	balance1 := currency.NewAmount(currency.NewBig(10), cid1)
	sa, sta := t.newAccount(true, []currency.Amount{balance0, balance1})
	ra0, stb := t.newAccount(true, []currency.Amount{currency.NewAmount(currency.NewBig(0), cid0)})
	ra1, stc := t.newAccount(true, []currency.Amount{currency.NewAmount(currency.NewBig(0), cid0)})
	dd0 := t.newTestDocumentData(sa.Address)
	dd1 := NewDocumentData(DocInfo{idx: currency.NewBig(1), filehash: FileHash("EFGH")}, sa.Address, sa.Address, []DocSign{})
	sts0 := t.newStateDocument(sa.Address, dd0)
	dinv, _ := StateDocumentsValue(sts0[1])
	err := dinv.Append(DocInfo{idx: currency.NewBig(1), filehash: dd1.FileHash()})
	// sts0[1] = nst
	sts1 := t.newStateDocument(sa.Address, dd1)
	nst, _ := SetStateDocumentsValue(sts1[1], dinv)
	sts1[1] = nst
	sts := []state.State{sts0[0], sts0[2], sts1[0], sts1[1], sts1[2]}

	pool, _ := t.statepool(sta, stb, stc, sts)

	fee0 := currency.NewBig(11)
	fee1 := currency.NewBig(3)
	feeer0 := currency.NewFixedFeeer(sa.Address, fee0)
	feeer1 := currency.NewFixedFeeer(sa.Address, fee1)

	cp := currency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(cid0, currency.NewBig(99), NewTestAddress(), feeer0)))
	t.NoError(cp.Set(t.newCurrencyDesignState(cid1, currency.NewBig(99), NewTestAddress(), feeer1)))

	opr := t.processor(cp, pool)

	token := util.UUID().Bytes()
	items := []TransferDocumentsItem{
		t.newTransferDocumentsItem(t.docid, sa.Address, ra0.Address, cid0),
		t.newTransferDocumentsItem(currency.NewBig(1), sa.Address, ra1.Address, cid1),
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
	cid0 := currency.CurrencyID("SHOWME")
	cid1 := currency.CurrencyID("FINDME")
	balance0 := currency.NewAmount(currency.NewBig(33), cid0)
	balance1 := currency.NewAmount(currency.NewBig(33), cid1)
	sa, sta := t.newAccount(true, []currency.Amount{balance0, balance1})
	ra0, stb := t.newAccount(true, []currency.Amount{currency.NewAmount(currency.NewBig(0), cid0)})
	ra1, stc := t.newAccount(true, []currency.Amount{currency.NewAmount(currency.NewBig(0), cid0)})
	dd0 := t.newTestDocumentData(sa.Address)
	dd1 := NewDocumentData(DocInfo{idx: currency.NewBig(1), filehash: FileHash("EFGH")}, sa.Address, sa.Address, []DocSign{})
	sts0 := t.newStateDocument(sa.Address, dd0)
	dinv, _ := StateDocumentsValue(sts0[1])
	err := dinv.Append(DocInfo{idx: currency.NewBig(1), filehash: dd1.FileHash()})
	// sts0[1] = nst
	sts1 := t.newStateDocument(sa.Address, dd1)
	nst, _ := SetStateDocumentsValue(sts1[1], dinv)
	sts1[1] = nst
	sts := []state.State{sts0[0], sts0[2], sts1[0], sts1[1], sts1[2]}

	pool, _ := t.statepool(sta, stb, stc, sts)

	feeer := t.newTestFixedFeeer(sa.Address)

	cp := currency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(cid0, currency.NewBig(99), NewTestAddress(), feeer)))
	t.NoError(cp.Set(t.newCurrencyDesignState(cid1, currency.NewBig(99), NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)

	items0 := []TransferDocumentsItem{
		t.newTransferDocumentsItem(t.docid, sa.Address, ra0.Address, cid0),
	}
	tfd0 := t.newTransferDocument(sa.Address, sa.Privs(), items0)

	t.NoError(opr.Process(tfd0))

	items1 := []TransferDocumentsItem{
		t.newTransferDocumentsItem(currency.NewBig(1), sa.Address, ra1.Address, cid1),
	}
	tfd1 := t.newTransferDocument(sa.Address, sa.Privs(), items1)

	err = opr.Process(tfd1)

	var oper operation.ReasonError
	t.True(xerrors.As(err, &oper))
	t.Contains(err.Error(), "violates only one sender")
}

func (t *testTransferDocumentsOperations) TestUnderThreshold() {
	spk := key.MustNewBTCPrivatekey()
	rpk := key.MustNewBTCPrivatekey()

	skey := t.newKey(spk.Publickey(), 50)
	rkey := t.newKey(rpk.Publickey(), 50)
	skeys, _ := currency.NewKeys([]currency.Key{skey, rkey}, 100)
	rkeys, _ := currency.NewKeys([]currency.Key{rkey}, 50)

	pks := []key.Privatekey{spk}
	sender, _ := currency.NewAddressFromKeys(skeys)
	receiver, _ := currency.NewAddressFromKeys(rkeys)

	// set sender state
	balance := currency.NewAmount(currency.NewBig(33), t.cid)
	dd := t.newTestDocumentData(sender)
	std := t.newStateDocument(sender, dd)

	var sts []state.State
	sts = append(sts, std...)
	sts = append(sts,
		t.newStateBalance(sender, balance.Big(), balance.Currency()),
		t.newStateBalance(receiver, balance.Big(), balance.Currency()),
		t.newStateKeys(sender, skeys),
		t.newStateKeys(receiver, rkeys),
	)

	pool, _ := t.statepool(sts)
	feeer := currency.NewFixedFeeer(sender, currency.ZeroBig)

	cp := currency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)

	items := []TransferDocumentsItem{t.newTransferDocumentsItem(t.docid, sender, receiver, t.cid)}
	tfd := t.newTransferDocument(sender, pks, items)

	err := opr.Process(tfd)

	var oper operation.ReasonError
	t.True(xerrors.As(err, &oper))
	t.Contains(err.Error(), "not passed threshold")
}

func TestTransferDocumentsOperations(t *testing.T) {
	suite.Run(t, new(testTransferDocumentsOperations))
}
