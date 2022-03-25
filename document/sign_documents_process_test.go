//go:build test
// +build test

package document

import (
	"testing"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/key"
	"github.com/spikeekips/mitum/base/prprocessor"
	"github.com/spikeekips/mitum/base/state"
	"github.com/spikeekips/mitum/storage"
	"github.com/spikeekips/mitum/util"
	"github.com/stretchr/testify/suite"
)

type testSignDocumentsOperation struct {
	baseTestOperationProcessor
}

func (t *testSignDocumentsOperation) processor(cp *currency.CurrencyPool, pool *storage.Statepool) prprocessor.OperationProcessor {
	copr, err := NewOperationProcessor(cp).
		SetProcessor(SignDocumentsHinter, NewSignDocumentsProcessor(cp))
	t.NoError(err)

	if pool == nil {
		return copr
	}

	return copr.New(pool)
}

func (t *testSignDocumentsOperation) newOperation(sender base.Address, items []SignDocumentsItem, pks []key.Privatekey) SignDocuments {
	token := util.UUID().Bytes()
	fact := NewSignDocumentsFact(token, sender, items)

	var fs []base.FactSign
	for _, pk := range pks {
		sig, err := base.NewFactSignature(pk, fact, nil)
		if err != nil {
			panic(err)
		}

		fs = append(fs, base.NewBaseFactSign(pk.Publickey(), sig))
	}

	cd, err := NewSignDocuments(fact, fs, "")
	if err != nil {
		panic(err)
	}

	err = cd.IsValid(nil)
	if err != nil {
		panic(err)
	}

	return cd
}

func (t *testSignDocumentsOperation) TestNormalCase() {
	// currency id
	cid := currency.CurrencyID("SHOWME")
	// sender initial balance
	balance := []currency.Amount{
		currency.NewAmount(currency.NewBig(33), cid),
	}
	// create BSDocID
	bsDocIDStr := "1sdi"
	// create BSDocData
	bsDocData, ownerAccount, signerAccount := newBSDocData("filehash", bsDocIDStr, account{})
	// sender account state
	sender, senderState := t.newStateAccount(true, balance, ownerAccount)
	// signer account state
	signer, signerState := t.newStateAccount(true, balance, signerAccount)
	// prepare BSDocData state
	bsDocDataState := t.newStateDocument(ownerAccount.Address, *bsDocData, DocumentInventory{})
	// feeer
	fee := currency.NewBig(1)
	feeer := currency.NewFixedFeeer(sender.Address, fee)
	// currencyPool
	cp := currency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(cid, currency.NewBig(99), sender.Address, feeer)))
	// create SignDocumentsItem with BSDocData
	items := []SignDocumentsItem{
		NewSignDocumentsItemSingleFile((*bsDocData).DocumentID(), sender.Address, cid),
	}
	pool, _ := t.statepool(senderState, signerState, bsDocDataState)
	opr := t.processor(cp, pool)
	sd := t.newOperation(signer.Address, items, signer.Privs())
	t.NoError(opr.Process(sd))

	// document data state
	var newDocumentState state.State
	// signer balance state
	var signerBalanceState state.State
	for _, stu := range pool.Updates() {
		if currency.IsStateBalanceKey(stu.Key()) {
			st := stu.GetState()

			i, err := currency.StateBalanceValue(st)
			t.NoError(err)

			if st.Key() == currency.StateKeyBalance(signer.Address, i.Currency()) {
				signerBalanceState = st
			} else {
				continue
			}
		} else if (IsStateDocumentDataKey(stu.Key())) && (stu.Key() == StateKeyDocumentData(bsDocData.DocumentID())) {
			newDocumentState = stu.GetState()
		}
	}

	// signer balance state should not be nil
	t.NotNil(signerBalanceState)

	signerBalanceAmount, _ := currency.StateBalanceValue(signerBalanceState)
	// amount of signer balance state should be deductes by fee
	t.True(signerBalanceAmount.Big().Equal(balance[0].Big().Sub(fee)))
	// processed fee amount should be same with currencydesign fee
	t.Equal(fee, signerBalanceState.(currency.AmountState).Fee())

	d, _ := StateDocumentDataValue(newDocumentState)
	newDocumentData, ok := d.(BSDocData)
	t.True(ok)
	t.Equal(newDocumentData.fileHash, FileHash("filehash"))
	t.True(newDocumentData.Creator().Equal(MustNewDocSign(sender.Address, "signcode0", true)))
	t.True(newDocumentData.Signers()[0].Address().Equal(signer.Address))
	t.True(newDocumentData.Signers()[0].Signed() == true)
}

func (t *testSignDocumentsOperation) TestSenderNotExist() {
	// currency id
	cid := currency.CurrencyID("SHOWME")
	// sender initial balance
	balance := []currency.Amount{
		currency.NewAmount(currency.NewBig(33), cid),
	}
	// create BSDocID
	bsDocIDStr := "1sdi"
	// create BSDocData
	bsDocData, ownerAccount, signerAccount := newBSDocData("filehash", bsDocIDStr, account{})
	// signer account state
	_, ownerState := t.newStateAccount(true, balance, ownerAccount)
	// prepare BSDocData state
	bsDocDataState := t.newStateDocument(ownerAccount.Address, *bsDocData, DocumentInventory{})
	// feeer
	fee := currency.NewBig(1)
	feeer := currency.NewFixedFeeer(ownerAccount.Address, fee)
	// currencyPool
	cp := currency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(cid, currency.NewBig(99), ownerAccount.Address, feeer)))
	// create SignDocumentsItem with BSDocData
	items := []SignDocumentsItem{
		NewSignDocumentsItemSingleFile((*bsDocData).DocumentID(), ownerAccount.Address, cid),
	}
	pool, _ := t.statepool(ownerState, bsDocDataState)
	opr := t.processor(cp, pool)
	sd := t.newOperation(signerAccount.Address, items, signerAccount.Privs())
	err := opr.Process(sd)

	t.Contains(err.Error(), "does not exist")
}

func (t *testSignDocumentsOperation) TestOwnerNotExist() {
	// currency id
	cid := currency.CurrencyID("SHOWME")
	// sender initial balance
	balance := []currency.Amount{
		currency.NewAmount(currency.NewBig(33), cid),
	}
	// create BSDocID
	bsDocIDStr := "1sdi"
	// create BSDocData
	bsDocData, ownerAccount, signerAccount := newBSDocData("filehash", bsDocIDStr, account{})
	// signer account state
	_, signerState := t.newStateAccount(true, balance, signerAccount)
	// prepare BSDocData state
	bsDocDataState := t.newStateDocument(ownerAccount.Address, *bsDocData, DocumentInventory{})
	// feeer
	fee := currency.NewBig(1)
	feeer := currency.NewFixedFeeer(ownerAccount.Address, fee)
	// currencyPool
	cp := currency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(cid, currency.NewBig(99), ownerAccount.Address, feeer)))
	// create SignDocumentsItem with BSDocData
	items := []SignDocumentsItem{
		NewSignDocumentsItemSingleFile((*bsDocData).DocumentID(), ownerAccount.Address, cid),
	}
	pool, _ := t.statepool(signerState, bsDocDataState)
	opr := t.processor(cp, pool)
	sd := t.newOperation(signerAccount.Address, items, signerAccount.Privs())
	err := opr.Process(sd)

	t.Contains(err.Error(), "does not exist")
}

func (t *testSignDocumentsOperation) TestSenderNotExistInSignersList() {
	// currency id
	cid := currency.CurrencyID("SHOWME")
	// sender initial balance
	balance := []currency.Amount{
		currency.NewAmount(currency.NewBig(33), cid),
	}
	// create BSDocID
	bsDocIDStr := "1sdi"
	// create BSDocData
	bsDocData, ownerAccount, signerAccount := newBSDocData("filehash", bsDocIDStr, account{})
	// owner account state
	_, ownerState := t.newStateAccount(true, balance, ownerAccount)
	// signer account state
	_, signerState := t.newStateAccount(true, balance, signerAccount)
	// new signer account state
	newSigner, newSignerState := t.newStateAccount(true, balance, nil)

	// prepare BSDocData state
	bsDocData.signers = []DocSign{MustNewDocSign(newSigner.Address, "signcode2", false)}
	bsDocDataState := t.newStateDocument(ownerAccount.Address, *bsDocData, DocumentInventory{})
	// feeer
	fee := currency.NewBig(1)
	feeer := currency.NewFixedFeeer(ownerAccount.Address, fee)
	// currencyPool
	cp := currency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(cid, currency.NewBig(99), ownerAccount.Address, feeer)))
	// create SignDocumentsItem with BSDocData
	items := []SignDocumentsItem{
		NewSignDocumentsItemSingleFile((*bsDocData).DocumentID(), ownerAccount.Address, cid),
	}
	pool, _ := t.statepool(ownerState, signerState, newSignerState, bsDocDataState)
	opr := t.processor(cp, pool)
	sd := t.newOperation(signerAccount.Address, items, signerAccount.Privs())
	err := opr.Process(sd)

	t.Contains(err.Error(), "sender not found in document Signers")
}

func (t *testSignDocumentsOperation) TestInsufficientBalanceForFee() {
	// currency id
	cid := currency.CurrencyID("SHOWME")
	// sender initial balance
	balance := []currency.Amount{
		currency.NewAmount(currency.NewBig(10), cid),
	}
	// create BSDocID
	bsDocIDStr := "1sdi"
	// create BSDocData
	bsDocData, ownerAccount, signerAccount := newBSDocData("filehash", bsDocIDStr, account{})
	// sender account state
	sender, senderState := t.newStateAccount(true, balance, ownerAccount)
	// signer account state
	signer, signerState := t.newStateAccount(true, balance, signerAccount)
	// prepare BSDocData state
	bsDocDataState := t.newStateDocument(ownerAccount.Address, *bsDocData, DocumentInventory{})
	// feeer
	fee := currency.NewBig(11)
	feeer := currency.NewFixedFeeer(sender.Address, fee)
	// currencyPool
	cp := currency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(cid, currency.NewBig(99), sender.Address, feeer)))
	// create SignDocumentsItem with BSDocData
	items := []SignDocumentsItem{
		NewSignDocumentsItemSingleFile((*bsDocData).DocumentID(), sender.Address, cid),
	}
	pool, _ := t.statepool(senderState, signerState, bsDocDataState)
	opr := t.processor(cp, pool)
	sd := t.newOperation(signer.Address, items, signer.Privs())
	err := opr.Process(sd)

	t.Contains(err.Error(), "insufficient balance")
}

func (t *testSignDocumentsOperation) TestMultipleItemsWithFee() {
	// currency id
	cid0 := currency.CurrencyID("SHOWME")
	cid1 := currency.CurrencyID("FINDME")
	// sender initial balance
	balance := []currency.Amount{
		currency.NewAmount(currency.NewBig(33), cid0),
		currency.NewAmount(currency.NewBig(33), cid1),
	}
	// create BSDocID
	bsDocIDStr0 := "1sdi"
	bsDocIDStr1 := "2sdi"
	// create BSDocData
	bsDocData0, ownerAccount, signerAccount := newBSDocData("filehash0", bsDocIDStr0, account{})
	bsDocData1, _, _ := newBSDocData("filehash1", bsDocIDStr1, *ownerAccount)
	bsDocData1.signers = []DocSign{MustNewDocSign(signerAccount.Address, "signcode1", false)}
	// owner account state
	_, ownerState := t.newStateAccount(true, balance, ownerAccount)
	// signer account state
	signer, signerState := t.newStateAccount(true, balance, signerAccount)
	// prepare BSDocData state0
	bsDocDataState0 := t.newStateDocument(ownerAccount.Address, *bsDocData0, DocumentInventory{})
	d, err := StateDocumentsValue(bsDocDataState0[1])
	t.NoError(err)
	bsDocDataState1 := t.newStateDocument(ownerAccount.Address, *bsDocData1, d)
	// feeer
	fee := currency.NewBig(1)
	feeer := currency.NewFixedFeeer(signer.Address, fee)
	// currencyPool
	cp := currency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(cid0, currency.NewBig(99), signer.Address, feeer)))
	t.NoError(cp.Set(t.newCurrencyDesignState(cid1, currency.NewBig(99), signer.Address, feeer)))
	// create SignDocumentsItem with BSDocData
	items := []SignDocumentsItem{
		NewSignDocumentsItemSingleFile(bsDocData0.DocumentID(), ownerAccount.Address, cid0),
		NewSignDocumentsItemSingleFile(bsDocData1.DocumentID(), ownerAccount.Address, cid1),
	}
	pool, _ := t.statepool(ownerState, signerState, []state.State{bsDocDataState0[0]}, bsDocDataState1)
	opr := t.processor(cp, pool)
	sd := t.newOperation(signer.Address, items, signer.Privs())
	t.NoError(opr.Process(sd))

	// document data state
	var newDocumentState0 state.State
	var newDocumentState1 state.State
	// signer balance state
	var signerBalanceState state.State
	for _, stu := range pool.Updates() {
		if currency.IsStateBalanceKey(stu.Key()) {
			st := stu.GetState()

			i, err := currency.StateBalanceValue(st)
			t.NoError(err)

			if st.Key() == currency.StateKeyBalance(signer.Address, i.Currency()) {
				signerBalanceState = st
			} else {
				continue
			}
		} else if IsStateDocumentDataKey(stu.Key()) {
			if stu.Key() == StateKeyDocumentData(bsDocData0.DocumentID()) {
				newDocumentState0 = stu.GetState()
			} else if stu.Key() == StateKeyDocumentData(bsDocData1.DocumentID()) {
				newDocumentState1 = stu.GetState()
			}
		}
	}

	// signer balance state should not be nil
	t.NotNil(signerBalanceState)

	signerBalanceAmount, _ := currency.StateBalanceValue(signerBalanceState)
	// amount of signer balance state should be deductes by fee
	t.True(signerBalanceAmount.Big().Equal(balance[0].Big().Sub(fee)))
	t.True(signerBalanceAmount.Big().Equal(balance[1].Big().Sub(fee)))
	// processed fee amount should be same with currencydesign fee
	t.Equal(fee, signerBalanceState.(currency.AmountState).Fee())

	doc, _ := StateDocumentDataValue(newDocumentState0)
	newDocumentData, ok := doc.(BSDocData)
	t.True(ok)
	t.True(newDocumentData.DocumentID() == "1sdi")
	t.Equal(newDocumentData.fileHash, FileHash("filehash0"))
	t.True(newDocumentData.Creator().Equal(MustNewDocSign(ownerAccount.Address, "signcode0", true)))
	t.True(newDocumentData.Signers()[0].Address().Equal(signer.Address))
	t.True(newDocumentData.Signers()[0].Signed() == true)

	doc, _ = StateDocumentDataValue(newDocumentState1)
	newDocumentData, ok = doc.(BSDocData)
	t.True(ok)
	t.True(newDocumentData.DocumentID() == "2sdi")
	t.Equal(newDocumentData.fileHash, FileHash("filehash1"))
	t.True(newDocumentData.Creator().Equal(MustNewDocSign(ownerAccount.Address, "signcode0", true)))
	t.True(newDocumentData.Signers()[0].Address().Equal(signer.Address))
	t.True(newDocumentData.Signers()[0].Signed() == true)
}

func (t *testSignDocumentsOperation) TestInsufficientMultipleItemsWithFee() {
	// currency id
	cid0 := currency.CurrencyID("SHOWME")
	cid1 := currency.CurrencyID("FINDME")
	// sender initial balance
	balance := []currency.Amount{
		currency.NewAmount(currency.NewBig(10), cid0),
		currency.NewAmount(currency.NewBig(10), cid1),
	}
	// create BSDocID
	bsDocIDStr0 := "1sdi"
	bsDocIDStr1 := "2sdi"
	// create BSDocData
	bsDocData0, ownerAccount, signerAccount := newBSDocData("filehash0", bsDocIDStr0, account{})
	bsDocData1, _, _ := newBSDocData("filehash1", bsDocIDStr1, *ownerAccount)
	bsDocData1.signers = []DocSign{MustNewDocSign(signerAccount.Address, "signcode1", false)}
	// owner account state
	_, ownerState := t.newStateAccount(true, balance, ownerAccount)
	// signer account state
	signer, signerState := t.newStateAccount(true, balance, signerAccount)
	// prepare BSDocData state0
	bsDocDataState0 := t.newStateDocument(ownerAccount.Address, *bsDocData0, DocumentInventory{})
	d, err := StateDocumentsValue(bsDocDataState0[1])
	t.NoError(err)
	bsDocDataState1 := t.newStateDocument(ownerAccount.Address, *bsDocData1, d)
	// feeer
	fee0 := currency.NewBig(1)
	feeer0 := currency.NewFixedFeeer(signer.Address, fee0)
	fee1 := currency.NewBig(11)
	feeer1 := currency.NewFixedFeeer(signer.Address, fee1)
	// currencyPool
	cp := currency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(cid0, currency.NewBig(99), signer.Address, feeer0)))
	t.NoError(cp.Set(t.newCurrencyDesignState(cid1, currency.NewBig(99), signer.Address, feeer1)))
	// create SignDocumentsItem with BSDocData
	items := []SignDocumentsItem{
		NewSignDocumentsItemSingleFile(bsDocData0.DocumentID(), ownerAccount.Address, cid0),
		NewSignDocumentsItemSingleFile(bsDocData1.DocumentID(), ownerAccount.Address, cid1),
	}
	pool, _ := t.statepool(ownerState, signerState, []state.State{bsDocDataState0[0]}, bsDocDataState1)
	opr := t.processor(cp, pool)
	sd := t.newOperation(signer.Address, items, signer.Privs())
	err = opr.Process(sd)
	t.Contains(err.Error(), "insufficient balance")
}

func (t *testSignDocumentsOperation) TestSameSenders() {
	// currency id
	cid0 := currency.CurrencyID("SHOWME")
	cid1 := currency.CurrencyID("FINDME")
	// sender initial balance
	balance := []currency.Amount{
		currency.NewAmount(currency.NewBig(33), cid0),
		currency.NewAmount(currency.NewBig(33), cid1),
	}
	// create BSDocID
	bsDocIDStr0 := "1sdi"
	bsDocIDStr1 := "2sdi"
	// create BSDocData
	bsDocData0, ownerAccount, signerAccount := newBSDocData("filehash0", bsDocIDStr0, account{})
	bsDocData1, _, _ := newBSDocData("filehash1", bsDocIDStr1, *ownerAccount)
	bsDocData1.signers = []DocSign{MustNewDocSign(signerAccount.Address, "signcode1", false)}
	// owner account state
	_, ownerState := t.newStateAccount(true, balance, ownerAccount)
	// signer account state
	signer, signerState := t.newStateAccount(true, balance, signerAccount)
	// prepare BSDocData state0
	bsDocDataState0 := t.newStateDocument(ownerAccount.Address, *bsDocData0, DocumentInventory{})
	d, err := StateDocumentsValue(bsDocDataState0[1])
	t.NoError(err)
	bsDocDataState1 := t.newStateDocument(ownerAccount.Address, *bsDocData1, d)
	// feeer
	fee := currency.NewBig(1)
	feeer := currency.NewFixedFeeer(signer.Address, fee)
	// currencyPool
	cp := currency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(cid0, currency.NewBig(99), signer.Address, feeer)))
	t.NoError(cp.Set(t.newCurrencyDesignState(cid1, currency.NewBig(99), signer.Address, feeer)))
	// create SignDocumentsItem with BSDocData
	items0 := []SignDocumentsItem{
		NewSignDocumentsItemSingleFile(bsDocData0.DocumentID(), ownerAccount.Address, cid0),
	}
	items1 := []SignDocumentsItem{
		NewSignDocumentsItemSingleFile(bsDocData1.DocumentID(), ownerAccount.Address, cid1),
	}

	pool, _ := t.statepool(ownerState, signerState, []state.State{bsDocDataState0[0]}, bsDocDataState1)
	opr := t.processor(cp, pool)
	sd0 := t.newOperation(signer.Address, items0, signer.Privs())
	t.NoError(opr.Process(sd0))
	sd1 := t.newOperation(signer.Address, items1, signer.Privs())
	err = opr.Process(sd1)
	t.Contains(err.Error(), "violates only one sender")
}

func TestSignDocumentsOperations(t *testing.T) {
	suite.Run(t, new(testSignDocumentsOperation))
}
