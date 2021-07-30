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

type testSignDocumentsOperations struct {
	baseTestOperationProcessor
	cid   currency.CurrencyID
	docid currency.Big
	fh    FileHash
	fee   currency.Big
}

func (t *testSignDocumentsOperations) SetupSuite() {
	t.cid = currency.CurrencyID("SHOWME")
	t.docid = currency.NewBig(0)
	t.fh = FileHash("ABCD")
	t.fee = currency.NewBig(3)
}

func (t *testSignDocumentsOperations) processor(cp *currency.CurrencyPool, pool *storage.Statepool) prprocessor.OperationProcessor {
	copr, err := NewOperationProcessor(cp).
		SetProcessor(SignDocuments{}, NewSignDocumentsProcessor(cp))
	t.NoError(err)

	if pool == nil {
		return copr
	}

	return copr.New(pool)
}

func (t *testSignDocumentsOperations) newSignDocumentsItem(
	docid currency.Big,
	owner base.Address,
	cid currency.CurrencyID,
) SignDocumentItem {

	return NewSignDocumentsItemSingleFile(docid, owner, cid)
}

func (t *testSignDocumentsOperations) newSignDocument(
	sender base.Address,
	keys []key.Privatekey,
	items []SignDocumentItem,
) SignDocuments {
	token := util.UUID().Bytes()
	fact := NewSignDocumentsFact(token, sender, items)

	var fs []operation.FactSign
	for _, pk := range keys {
		sig, err := operation.NewFactSignature(pk, fact, nil)
		t.NoError(err)

		fs = append(fs, operation.NewBaseFactSign(pk.Publickey(), sig))
	}

	tfd, err := NewSignDocuments(fact, fs, "")
	t.NoError(err)

	t.NoError(tfd.IsValid(nil))

	return tfd
}

func (t *testSignDocumentsOperations) newTestDocumentData(ca base.Address, ga base.Address) DocumentData {
	var doc DocumentData
	if ga == nil {
		doc = NewDocumentData(t.fh, ca, ca, []DocSign{})
	} else {
		doc = NewDocumentData(t.fh, ca, ca, []DocSign{{address: ga, signed: false}})
	}
	doc = doc.WithData(doc.FileHash(), DocInfo{idx: t.docid, filehash: t.fh}, doc.Creator(), doc.Owner(), doc.Signers())
	return doc
}

func (t *testSignDocumentsOperations) newTestBalance() []currency.Amount {
	return []currency.Amount{currency.NewAmount(currency.NewBig(33), t.cid)}
}

func (t *testSignDocumentsOperations) newTestFixedFeeer(sa base.Address) currency.FixedFeeer {
	return currency.NewFixedFeeer(sa, t.fee)
}
func (t *testSignDocumentsOperations) TestNormalCase() {
	balance := t.newTestBalance()
	sa, sta := t.newAccount(true, balance) // sender, signer
	ca, stb := t.newAccount(true, balance) // creator, owner
	dd := t.newTestDocumentData(ca.Address, sa.Address)

	sts := t.newStateDocument(ca.Address, dd)
	pool, _ := t.statepool(sta, stb, sts)

	feeer := t.newTestFixedFeeer(ca.Address)

	cp := currency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)

	items := []SignDocumentItem{t.newSignDocumentsItem(t.docid, ca.Address, t.cid)}
	tfd := t.newSignDocument(sa.Address, sa.Privs(), items)

	t.NoError(opr.Process(tfd))

	// check updated state
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
		} else if (IsStateDocumentDataKey(stu.Key())) && (stu.Key() == StateKeyDocumentData(t.fh)) {
			dds = stu.GetState()
		}
	}

	t.NotNil(sb)

	sba, _ := currency.StateBalanceValue(sb)
	t.True(sba.Big().Equal(balance[0].Big().Sub(t.fee)))

	t.Equal(t.fee, sb.(currency.AmountState).Fee())

	ndd, _ := StateDocumentDataValue(dds)
	t.True(ndd.FileHash().Equal(t.fh))
	t.True(ndd.Creator().Equal(ca.Address))
	t.True(ndd.Signers()[0].Address().Equal(sa.Address))
	t.True(ndd.Signers()[0].Signed() == true)
}

func (t *testSignDocumentsOperations) TestSenderNotExist() {
	balance := t.newTestBalance()
	sa, _ := t.newAccount(false, nil)     // sender, signer
	ca, st := t.newAccount(true, balance) // creator, owner
	dd := t.newTestDocumentData(ca.Address, sa.Address)

	sts := t.newStateDocument(ca.Address, dd)
	pool, _ := t.statepool(st, sts)

	feeer := t.newTestFixedFeeer(ca.Address)

	cp := currency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)

	items := []SignDocumentItem{t.newSignDocumentsItem(t.docid, ca.Address, t.cid)}
	tfd := t.newSignDocument(sa.Address, sa.Privs(), items)

	err := opr.Process(tfd)

	var oper operation.ReasonError
	t.True(xerrors.As(err, &oper))
	t.Contains(err.Error(), "does not exist")
}

func (t *testSignDocumentsOperations) TestOwnerNotExist() {
	balance := t.newTestBalance()
	sa, st := t.newAccount(true, balance) // sender, signer
	ca, _ := t.newAccount(false, nil)     // creator, owner
	dd := t.newTestDocumentData(ca.Address, sa.Address)

	sts := t.newStateDocument(ca.Address, dd)
	pool, _ := t.statepool(st, sts)

	feeer := t.newTestFixedFeeer(ca.Address)

	cp := currency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)

	items := []SignDocumentItem{t.newSignDocumentsItem(t.docid, ca.Address, t.cid)}
	tfd := t.newSignDocument(sa.Address, sa.Privs(), items)

	err := opr.Process(tfd)

	var oper operation.ReasonError
	t.True(xerrors.As(err, &oper))
	t.Contains(err.Error(), "does not exist")
}

func (t *testSignDocumentsOperations) TestSenderNotExistInSignersList() {
	balance := t.newTestBalance()
	sa, sta := t.newAccount(true, balance) // sender, signer
	ca, stb := t.newAccount(true, balance) // creator, owner
	dd := t.newTestDocumentData(ca.Address, nil)

	sts := t.newStateDocument(ca.Address, dd)
	pool, _ := t.statepool(sta, stb, sts)

	feeer := t.newTestFixedFeeer(ca.Address)

	cp := currency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)

	items := []SignDocumentItem{t.newSignDocumentsItem(t.docid, ca.Address, t.cid)}
	tfd := t.newSignDocument(sa.Address, sa.Privs(), items)

	err := opr.Process(tfd)

	var oper operation.ReasonError
	t.True(xerrors.As(err, &oper))
	t.Contains(err.Error(), "sender not found in document Signers")
}

func (t *testSignDocumentsOperations) TestInsufficientBalanceForFee() {
	balance := []currency.Amount{currency.NewAmount(currency.NewBig(2), t.cid)}
	sa, st := t.newAccount(true, balance) // sender, signer
	ca, _ := t.newAccount(true, balance)  // creator, owner
	dd := t.newTestDocumentData(ca.Address, sa.Address)

	sts := t.newStateDocument(ca.Address, dd)
	pool, _ := t.statepool(st, sts)

	feeer := t.newTestFixedFeeer(sa.Address)

	cp := currency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(t.cid, currency.NewBig(99), NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)

	items := []SignDocumentItem{t.newSignDocumentsItem(t.docid, ca.Address, t.cid)}
	tfd := t.newSignDocument(sa.Address, sa.Privs(), items)

	err := opr.Process(tfd)

	var oper operation.ReasonError
	t.True(xerrors.As(err, &oper))
	t.Contains(err.Error(), "insufficient balance")
}

func (t *testSignDocumentsOperations) TestMultipleItemsWithFee() {
	cid0 := currency.CurrencyID("SHOWME")
	cid1 := currency.CurrencyID("FINDME")
	balance0 := currency.NewAmount(currency.NewBig(33), cid0)
	balance1 := currency.NewAmount(currency.NewBig(33), cid1)
	sa, sta := t.newAccount(true, []currency.Amount{balance0, balance1})
	ca, stb := t.newAccount(true, []currency.Amount{currency.NewAmount(currency.NewBig(0), cid0)})

	dd0 := t.newTestDocumentData(ca.Address, sa.Address)
	dd1 := NewDocumentData(FileHash("EFGH"), ca.Address, ca.Address, []DocSign{{address: sa.Address, signed: false}})
	dd1 = dd1.WithData(dd1.FileHash(), DocInfo{idx: currency.NewBig(1), filehash: dd1.FileHash()}, dd1.Creator(), dd1.Owner(), dd1.Signers())
	sts0 := t.newStateDocument(ca.Address, dd0)
	dinv, _ := StateDocumentsValue(sts0[1])
	err := dinv.Append(DocInfo{idx: currency.NewBig(1), filehash: dd1.FileHash()})
	// sts0[1] = nst
	sts1 := t.newStateDocument(ca.Address, dd1)
	nst, _ := SetStateDocumentsValue(sts1[1], dinv)
	sts1[1] = nst
	sts := []state.State{sts0[0], sts0[2], sts1[0], sts1[1], sts1[2]}

	pool, _ := t.statepool(sta, stb, sts)

	feeer := t.newTestFixedFeeer(sa.Address)

	cp := currency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(cid0, currency.NewBig(99), NewTestAddress(), feeer)))
	t.NoError(cp.Set(t.newCurrencyDesignState(cid1, currency.NewBig(99), NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)

	token := util.UUID().Bytes()
	items := []SignDocumentItem{
		t.newSignDocumentsItem(t.docid, ca.Address, cid0),
		t.newSignDocumentsItem(currency.NewBig(1), ca.Address, cid1),
	}
	fact := NewSignDocumentsFact(token, sa.Address, items)
	sig, err := operation.NewFactSignature(sa.Privs()[0], fact, nil)
	t.NoError(err)
	fs := []operation.FactSign{operation.NewBaseFactSign(sa.Privs()[0].Publickey(), sig)}
	tfd, err := NewSignDocuments(fact, fs, "")
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

func (t *testSignDocumentsOperations) TestInsufficientMultipleItemsWithFee() {
	cid0 := currency.CurrencyID("SHOWME")
	cid1 := currency.CurrencyID("FINDME")
	balance0 := currency.NewAmount(currency.NewBig(10), cid0)
	balance1 := currency.NewAmount(currency.NewBig(10), cid1)
	sa, sta := t.newAccount(true, []currency.Amount{balance0, balance1})
	ca, stb := t.newAccount(true, []currency.Amount{currency.NewAmount(currency.NewBig(0), cid0)})

	dd0 := t.newTestDocumentData(ca.Address, sa.Address)
	dd1 := NewDocumentData(FileHash("EFGH"), ca.Address, ca.Address, []DocSign{{address: sa.Address, signed: false}})
	dd1 = dd1.WithData(dd1.FileHash(), DocInfo{idx: currency.NewBig(1), filehash: dd1.FileHash()}, dd1.Creator(), dd1.Owner(), dd1.Signers())
	sts0 := t.newStateDocument(ca.Address, dd0)
	dinv, _ := StateDocumentsValue(sts0[1])
	err := dinv.Append(DocInfo{idx: currency.NewBig(1), filehash: dd1.FileHash()})
	// sts0[1] = nst
	sts1 := t.newStateDocument(ca.Address, dd1)
	nst, _ := SetStateDocumentsValue(sts1[1], dinv)
	sts1[1] = nst
	sts := []state.State{sts0[0], sts0[2], sts1[0], sts1[1], sts1[2]}

	pool, _ := t.statepool(sta, stb, sts)

	fee0 := currency.NewBig(11)
	fee1 := currency.NewBig(3)
	feeer0 := currency.NewFixedFeeer(sa.Address, fee0)
	feeer1 := currency.NewFixedFeeer(sa.Address, fee1)

	cp := currency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(cid0, currency.NewBig(99), NewTestAddress(), feeer0)))
	t.NoError(cp.Set(t.newCurrencyDesignState(cid1, currency.NewBig(99), NewTestAddress(), feeer1)))

	opr := t.processor(cp, pool)

	token := util.UUID().Bytes()
	items := []SignDocumentItem{
		t.newSignDocumentsItem(t.docid, ca.Address, cid0),
		t.newSignDocumentsItem(currency.NewBig(1), ca.Address, cid1),
	}
	fact := NewSignDocumentsFact(token, sa.Address, items)
	sig, err := operation.NewFactSignature(sa.Privs()[0], fact, nil)
	t.NoError(err)
	fs := []operation.FactSign{operation.NewBaseFactSign(sa.Privs()[0].Publickey(), sig)}
	tfd, err := NewSignDocuments(fact, fs, "")
	t.NoError(err)

	err = opr.Process(tfd)

	var oper operation.ReasonError
	t.True(xerrors.As(err, &oper))
	t.Contains(err.Error(), "insufficient balance")
}

func (t *testSignDocumentsOperations) TestSameSenders() {
	cid0 := currency.CurrencyID("SHOWME")
	cid1 := currency.CurrencyID("FINDME")
	balance0 := currency.NewAmount(currency.NewBig(33), cid0)
	balance1 := currency.NewAmount(currency.NewBig(33), cid1)
	sa, sta := t.newAccount(true, []currency.Amount{balance0, balance1})
	ca, stb := t.newAccount(true, []currency.Amount{currency.NewAmount(currency.NewBig(0), cid0)})

	dd0 := t.newTestDocumentData(ca.Address, sa.Address)
	dd1 := NewDocumentData(FileHash("EFGH"), ca.Address, ca.Address, []DocSign{{address: sa.Address, signed: false}})
	dd1 = dd1.WithData(dd1.FileHash(), DocInfo{idx: currency.NewBig(1), filehash: dd1.FileHash()}, dd1.Creator(), dd1.Owner(), dd1.Signers())
	sts0 := t.newStateDocument(ca.Address, dd0)
	dinv, _ := StateDocumentsValue(sts0[1])
	err := dinv.Append(DocInfo{idx: currency.NewBig(1), filehash: dd1.FileHash()})
	// sts0[1] = nst
	sts1 := t.newStateDocument(ca.Address, dd1)
	nst, _ := SetStateDocumentsValue(sts1[1], dinv)
	sts1[1] = nst
	sts := []state.State{sts0[0], sts0[2], sts1[0], sts1[1], sts1[2]}

	pool, _ := t.statepool(sta, stb, sts)

	feeer := t.newTestFixedFeeer(sa.Address)

	cp := currency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(cid0, currency.NewBig(99), NewTestAddress(), feeer)))
	t.NoError(cp.Set(t.newCurrencyDesignState(cid1, currency.NewBig(99), NewTestAddress(), feeer)))

	opr := t.processor(cp, pool)

	items0 := []SignDocumentItem{
		t.newSignDocumentsItem(t.docid, ca.Address, cid0),
	}
	tfd0 := t.newSignDocument(sa.Address, sa.Privs(), items0)

	t.NoError(opr.Process(tfd0))

	items1 := []SignDocumentItem{
		t.newSignDocumentsItem(currency.NewBig(1), ca.Address, cid1),
	}
	tfd1 := t.newSignDocument(sa.Address, sa.Privs(), items1)

	err = opr.Process(tfd1)

	var oper operation.ReasonError
	t.True(xerrors.As(err, &oper))
	t.Contains(err.Error(), "violates only one sender")
}

func TestSignDocumentsOperations(t *testing.T) {
	suite.Run(t, new(testSignDocumentsOperations))
}
