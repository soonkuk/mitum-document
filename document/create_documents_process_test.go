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
	"github.com/spikeekips/mitum/util/hint"
	"github.com/stretchr/testify/suite"
)

type testCreateDocumentsOperation struct {
	baseTestOperationProcessor
}

func (t *testCreateDocumentsOperation) processor(cp *currency.CurrencyPool, pool *storage.Statepool) prprocessor.OperationProcessor {
	copr, err := NewOperationProcessor(nil).
		SetProcessor(CreateDocumentsHinter, NewCreateDocumentsProcessor(cp))
	t.NoError(err)

	if pool == nil {
		return copr
	}

	return copr.New(pool)
}

func (t *testCreateDocumentsOperation) newOperation(sender base.Address, items []CreateDocumentsItem, pks []key.Privatekey) CreateDocuments {
	token := util.UUID().Bytes()
	fact := NewCreateDocumentsFact(token, sender, items)

	var fs []base.FactSign
	for _, pk := range pks {
		sig, err := base.NewFactSignature(pk, fact, nil)
		if err != nil {
			panic(err)
		}

		fs = append(fs, base.NewBaseFactSign(pk.Publickey(), sig))
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
	bsDocID := "1sdi"
	bcUserDocID := "1cui"
	bcLandDocID := "1cli"
	bcVotingDocID := "1cvi"
	bcHistoryDocID := "1chi"
	// create BSDocData
	bsDocData, ownerAccount, signerAccount := newBSDocData("filehash", bsDocID, account{})
	// create BCUserData
	bcUserData, _, _ := newBCUserData(bcUserDocID, *ownerAccount)
	// create BCLandData
	bcLandData, _, renterAccount := newBCLandData(bcLandDocID, *ownerAccount)
	// create BCVotingData
	bcVotingData, _, bossAccount := newBCVotingData(bcVotingDocID, *ownerAccount)
	// create BCHistoryData
	bcHistoryData, _, cityAdminAccount := newBCHistoryData(bcHistoryDocID, *ownerAccount)
	// currency id
	cid := currency.CurrencyID("SHOWME")
	// create CreateDocumentsItems
	items := []CreateDocumentsItem{
		NewCreateDocumentsItemImpl(*bsDocData, cid),
		NewCreateDocumentsItemImpl(*bcUserData, cid),
		NewCreateDocumentsItemImpl(*bcLandData, cid),
		NewCreateDocumentsItemImpl(*bcVotingData, cid),
		NewCreateDocumentsItemImpl(*bcHistoryData, cid),
	}
	// sender initial balance
	balance := []currency.Amount{
		currency.NewAmount(currency.NewBig(33), cid),
	}
	// sender account
	sender, senderState := t.newStateAccount(true, balance, ownerAccount)
	// signer account
	_, signerState := t.newStateAccount(true, balance, signerAccount)
	// renter account
	_, renterState := t.newStateAccount(true, balance, renterAccount)
	// boss account
	_, bossState := t.newStateAccount(true, balance, bossAccount)
	// cityAdmin account
	_, cityAdminState := t.newStateAccount(true, balance, cityAdminAccount)

	// state pool
	pool, _ := t.statepool(
		senderState,
		signerState,
		renterState,
		bossState,
		cityAdminState,
	)

	fee := currency.NewBig(1)
	feeer := currency.NewFixedFeeer(sender.Address, fee)
	cp := currency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(cid, currency.NewBig(99), sender.Address, feeer)))

	opr := t.processor(cp, pool)
	cd := t.newOperation(sender.Address, items, sender.Privs())
	// t.NoError(opr.Process(cd))

	t.NoError(opr.Process(cd))

	// new document state
	var newBSDocDataDocumentState, newBCUserDataDocumentState, newBCLandDataDocumentState state.State
	var newBCVotingDataDocumentState, newBCHistoryDataDocumentState state.State
	// new documents data state
	// sender balance state
	var newDocumentsState state.State
	var senderBalanceState state.State
	for _, stu := range pool.Updates() {
		if currency.IsStateBalanceKey(stu.Key()) {
			st := stu.GetState()

			i, err := currency.StateBalanceValue(st)
			t.NoError(err)

			if st.Key() == currency.StateKeyBalance(sender.Address, i.Currency()) {
				senderBalanceState = st
			} else {
				continue
			}
		} else if (IsStateDocumentsKey(stu.Key())) && (stu.Key() == StateKeyDocuments(sender.Address)) {
			newDocumentsState = stu.GetState()
		} else if IsStateDocumentDataKey(stu.Key()) {
			if stu.Key() == StateKeyDocumentData(bsDocData.DocumentID()) {
				newBSDocDataDocumentState = stu.GetState()
			} else if stu.Key() == StateKeyDocumentData(bcUserData.DocumentID()) {
				newBCUserDataDocumentState = stu.GetState()
			} else if stu.Key() == StateKeyDocumentData(bcLandData.DocumentID()) {
				newBCLandDataDocumentState = stu.GetState()
			} else if stu.Key() == StateKeyDocumentData(bcVotingData.DocumentID()) {
				newBCVotingDataDocumentState = stu.GetState()
			} else if stu.Key() == StateKeyDocumentData(bcHistoryData.DocumentID()) {
				newBCHistoryDataDocumentState = stu.GetState()
			}
		}
	}

	// sender balance state should be processed
	t.NotNil(senderBalanceState)

	senderBalanceAmount, _ := currency.StateBalanceValue(senderBalanceState)
	// value of sender balance state should be deductes by fee
	t.True(senderBalanceAmount.Big().Equal(balance[0].Big().Sub(currency.NewBig(5).Mul(fee))))
	// processed fee amount should be same with currencydesign fee
	t.Equal(currency.NewBig(5).Mul(fee), senderBalanceState.(currency.AmountState).Fee())
	// check new documentData state
	newBSDocData, _ := StateDocumentDataValue(newBSDocDataDocumentState)
	t.Equal(newBSDocData.DocumentID(), bsDocData.DocumentID())
	t.True(newBSDocData.Owner().Equal(sender.Address))

	newBCUserData, _ := StateDocumentDataValue(newBCUserDataDocumentState)
	t.Equal(newBCUserData.DocumentID(), bcUserData.DocumentID())
	t.True(newBCUserData.Owner().Equal(sender.Address))

	newBCLandData, _ := StateDocumentDataValue(newBCLandDataDocumentState)
	t.Equal(newBCLandData.DocumentID(), bcLandData.DocumentID())
	t.True(newBCLandData.Owner().Equal(sender.Address))

	newBCVotingData, _ := StateDocumentDataValue(newBCVotingDataDocumentState)
	t.Equal(newBCVotingData.DocumentID(), bcVotingData.DocumentID())
	t.True(newBCVotingData.Owner().Equal(sender.Address))

	newBCHistoryData, _ := StateDocumentDataValue(newBCHistoryDataDocumentState)
	t.Equal(newBCHistoryData.DocumentID(), bcHistoryData.DocumentID())
	t.True(newBCHistoryData.Owner().Equal(sender.Address))

	newDocumentInventory, _ := StateDocumentsValue(newDocumentsState)
	docIDs := make([]string, len((newDocumentInventory.Documents())))
	docTypes := make([]hint.Type, len((newDocumentInventory.Documents())))
	for i := range newDocumentInventory.Documents() {
		docIDs[i] = newDocumentInventory.Documents()[i].DocumentID()
		docTypes[i] = newDocumentInventory.Documents()[i].DocType()
	}
	var bsDocDataExist bool
	var bcUserDataExist bool
	var bcLandDataExist bool
	var bcVotingDataExist bool
	var bcHistoryDataExist bool

	for i := range docIDs {
		if docIDs[i] == bsDocData.DocumentID() && docTypes[i] == bsDocData.DocumentType() {
			bsDocDataExist = true
		}
		if docIDs[i] == bcUserData.DocumentID() && docTypes[i] == bcUserData.DocumentType() {
			bcUserDataExist = true
		}
		if docIDs[i] == bcLandData.DocumentID() && docTypes[i] == bcLandData.DocumentType() {
			bcLandDataExist = true
		}
		if docIDs[i] == bcVotingData.DocumentID() && docTypes[i] == bcVotingData.DocumentType() {
			bcVotingDataExist = true
		}
		if docIDs[i] == bcHistoryData.DocumentID() && docTypes[i] == bcHistoryData.DocumentType() {
			bcHistoryDataExist = true
		}
	}
	t.True(bsDocDataExist)
	t.True(bcUserDataExist)
	t.True(bcLandDataExist)
	t.True(bcVotingDataExist)
	t.True(bcHistoryDataExist)
}

func (t *testCreateDocumentsOperation) TestAccountsNotExist() {
	bsDocID := "1sdi"
	bcLandDocID := "1cli"
	bcVotingDocID := "1cvi"
	bcHistoryDocID := "1chi"
	// create BSDocData
	bsDocData, ownerAccount, _ := newBSDocData("filehash", bsDocID, account{})
	// create BCLandData
	bcLandData, _, _ := newBCLandData(bcLandDocID, *ownerAccount)
	// create BCVotingData
	bcVotingData, _, _ := newBCVotingData(bcVotingDocID, *ownerAccount)
	// create BCHistoryData
	bcHistoryData, _, _ := newBCHistoryData(bcHistoryDocID, *ownerAccount)

	// currency id
	cid := currency.CurrencyID("SHOWME")
	// sender initial balance
	balance := []currency.Amount{
		currency.NewAmount(currency.NewBig(33), cid),
	}
	// sender account
	sender, senderState := t.newStateAccount(true, balance, ownerAccount)
	// create document item with BSDocData
	fee := currency.NewBig(1)
	feeer := currency.NewFixedFeeer(sender.Address, fee)
	cp := currency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(cid, currency.NewBig(99), sender.Address, feeer)))
	items := []CreateDocumentsItem{
		NewCreateDocumentsItemImpl(*bsDocData, cid),
	}
	pool, _ := t.statepool(senderState)
	// state pool
	opr := t.processor(cp, pool)
	cd := t.newOperation(sender.Address, items, sender.Privs())
	err := opr.Process(cd)

	t.Contains(err.Error(), "documentData related accounts not found")

	// create CreateDocumentsItem with BCUserData
	items = []CreateDocumentsItem{
		NewCreateDocumentsItemImpl(*bcLandData, cid),
	}
	pool, _ = t.statepool(senderState)
	opr = t.processor(cp, pool)
	cd = t.newOperation(sender.Address, items, sender.Privs())
	err = opr.Process(cd)

	t.Contains(err.Error(), "documentData related accounts not found")

	// create CreateDocumentsItem with BCVotingData
	items = []CreateDocumentsItem{
		NewCreateDocumentsItemImpl(*bcVotingData, cid),
	}
	pool, _ = t.statepool(senderState)
	opr = t.processor(cp, pool)
	cd = t.newOperation(sender.Address, items, sender.Privs())
	err = opr.Process(cd)

	t.Contains(err.Error(), "documentData related accounts not found")

	// create CreateDocumentsItem with BCUserData
	items = []CreateDocumentsItem{
		NewCreateDocumentsItemImpl(*bcHistoryData, cid),
	}
	pool, _ = t.statepool(senderState)
	opr = t.processor(cp, pool)
	cd = t.newOperation(sender.Address, items, sender.Privs())
	err = opr.Process(cd)

	t.Contains(err.Error(), "documentData related accounts not found")
}

func (t *testCreateDocumentsOperation) TestDocumentAlreadyExists() {
	bsDocID := "1sdi"
	// create BSDocData
	bsDocData, ownerAccount, _ := newBSDocData("filehash", bsDocID, account{})
	// currency id
	cid := currency.CurrencyID("SHOWME")
	balance := []currency.Amount{
		currency.NewAmount(currency.NewBig(33), cid),
	}
	// sender account
	sender, senderState := t.newStateAccount(true, balance, ownerAccount)
	// document state
	documentState := t.newStateDocument(sender.Address, bsDocData, DocumentInventory{})
	fee := currency.NewBig(1)
	feeer := currency.NewFixedFeeer(sender.Address, fee)
	cp := currency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(cid, currency.NewBig(99), sender.Address, feeer)))
	pool, _ := t.statepool(senderState, documentState)
	opr := t.processor(cp, pool)

	items := []CreateDocumentsItem{
		NewCreateDocumentsItemImpl(*bsDocData, cid),
	}

	cd := t.newOperation(sender.Address, items, sender.Privs())

	err := opr.Process(cd)

	t.Contains(err.Error(), "documentid already registered")
}

func (t *testCreateDocumentsOperation) TestSameSenders() {
	bsDocID := "1sdi"
	bcLandDocID := "1cli"
	// create BSDocData
	bsDocData, ownerAccount, signerAccount := newBSDocData("filehash", bsDocID, account{})
	// create BCLandData
	bcLandData, _, renterAccount := newBCLandData(bcLandDocID, *ownerAccount)
	// currency id
	cid := currency.CurrencyID("SHOWME")
	balance := []currency.Amount{
		currency.NewAmount(currency.NewBig(33), cid),
	}
	// sender account state
	sender, senderState := t.newStateAccount(true, balance, ownerAccount)
	// renter account state
	_, renterState := t.newStateAccount(true, balance, renterAccount)
	// signer account state
	_, signerState := t.newStateAccount(true, balance, signerAccount)
	fee := currency.NewBig(1)
	feeer := currency.NewFixedFeeer(sender.Address, fee)
	cp := currency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(cid, currency.NewBig(99), sender.Address, feeer)))
	pool, _ := t.statepool(senderState, renterState, signerState)
	opr := t.processor(cp, pool)

	items0 := []CreateDocumentsItem{
		NewCreateDocumentsItemImpl(*bsDocData, cid),
	}
	items1 := []CreateDocumentsItem{
		NewCreateDocumentsItemImpl(*bcLandData, cid),
	}

	cd0 := t.newOperation(sender.Address, items0, sender.Privs())
	cd1 := t.newOperation(sender.Address, items1, sender.Privs())
	opr.Process(cd0)
	err := opr.Process(cd1)
	t.Contains(err.Error(), "violates only one sender")
}

func (t *testCreateDocumentsOperation) TestDuplicatedSigner() {
	bsDocID := "1sdi"
	// create BSDocData
	bsDocData, ownerAccount, signerAccount := newBSDocData("filehash", bsDocID, account{})
	docSign := []DocSign{
		MustNewDocSign(signerAccount.Address, "signcode0", false),
		MustNewDocSign(signerAccount.Address, "signcode1", false),
	}
	bsDocData.signers = docSign
	// currency id
	cid := currency.CurrencyID("SHOWME")

	items := []CreateDocumentsItem{
		NewCreateDocumentsItemImpl(*bsDocData, cid),
	}

	panicfunc := func() {
		t.newOperation(ownerAccount.Address, items, ownerAccount.Privs())
	}
	recovered := false
	assertPanic(t.baseTest, panicfunc, "duplicated signer", &recovered)
	t.True(recovered)
}

func (t *testCreateDocumentsOperation) TestInSufficientBalanceForFee() {
	bsDocID := "1sdi"
	// create BSDocData
	bsDocData, ownerAccount, _ := newBSDocData("filehash", bsDocID, account{})
	// currency id
	cid := currency.CurrencyID("SHOWME")
	balance := []currency.Amount{
		currency.NewAmount(currency.NewBig(9), cid),
	}
	// sender account
	sender, senderState := t.newStateAccount(true, balance, ownerAccount)
	fee := currency.NewBig(10)
	feeer := currency.NewFixedFeeer(sender.Address, fee)
	cp := currency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(cid, currency.NewBig(99), sender.Address, feeer)))
	pool, _ := t.statepool(senderState)
	opr := t.processor(cp, pool)

	items := []CreateDocumentsItem{
		NewCreateDocumentsItemImpl(*bsDocData, cid),
	}

	cd := t.newOperation(sender.Address, items, sender.Privs())

	err := opr.Process(cd)

	t.Contains(err.Error(), "insufficient balance")
}

func (t *testCreateDocumentsOperation) TestUnknownCurrencyID() {
	bsDocID := "1sdi"
	// create BSDocData
	bsDocData, ownerAccount, _ := newBSDocData("filehash", bsDocID, account{})
	// currency id
	cid0 := currency.CurrencyID("SHOWME")
	cid1 := currency.CurrencyID("FINDME")
	balance := []currency.Amount{
		currency.NewAmount(currency.NewBig(9), cid1),
	}
	// sender account
	sender, senderState := t.newStateAccount(true, balance, ownerAccount)
	fee := currency.NewBig(10)
	feeer := currency.NewFixedFeeer(sender.Address, fee)
	cp := currency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(cid1, currency.NewBig(99), sender.Address, feeer)))
	pool, _ := t.statepool(senderState)
	opr := t.processor(cp, pool)

	items := []CreateDocumentsItem{
		NewCreateDocumentsItemImpl(*bsDocData, cid0),
	}

	cd := t.newOperation(sender.Address, items, sender.Privs())

	err := opr.Process(cd)

	t.Contains(err.Error(), "unknown currency id found")
}

func (t *testCreateDocumentsOperation) TestSenderNotHaveCurrency() {
	bsDocID := "1sdi"
	// create BSDocData
	bsDocData, ownerAccount, signerAccount := newBSDocData("filehash", bsDocID, account{})
	// currency id
	cid0 := currency.CurrencyID("SHOWME")
	cid1 := currency.CurrencyID("FINDME")
	balance := []currency.Amount{
		currency.NewAmount(currency.NewBig(33), cid1),
	}
	// sender account
	sender, senderState := t.newStateAccount(true, balance, ownerAccount)
	// signer account
	_, signerState := t.newStateAccount(true, balance, signerAccount)
	fee := currency.NewBig(1)
	feeer := currency.NewFixedFeeer(sender.Address, fee)
	cp := currency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(cid0, currency.NewBig(99), sender.Address, feeer)))
	t.NoError(cp.Set(t.newCurrencyDesignState(cid1, currency.NewBig(99), sender.Address, feeer)))
	pool, _ := t.statepool(senderState, signerState)
	opr := t.processor(cp, pool)

	items := []CreateDocumentsItem{
		NewCreateDocumentsItemImpl(*bsDocData, cid0),
	}

	cd := t.newOperation(sender.Address, items, sender.Privs())

	err := opr.Process(cd)
	t.Contains(err.Error(), "currency of holder does not exist")
}

func (t *testCreateDocumentsOperation) TestSenderBalanceNotExist() {
	bsDocID := "1sdi"
	// create BSDocData
	bsDocData, ownerAccount, signerAccount := newBSDocData("filehash", bsDocID, account{})
	// currency id
	cid := currency.CurrencyID("SHOWME")
	balance := []currency.Amount{
		currency.NewAmount(currency.NewBig(33), cid),
	}
	// sender account
	sender, senderState := t.newStateAccount(true, nil, ownerAccount)
	// signer account
	_, signerState := t.newStateAccount(true, balance, signerAccount)
	fee := currency.NewBig(1)
	feeer := currency.NewFixedFeeer(sender.Address, fee)
	cp := currency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(cid, currency.NewBig(99), sender.Address, feeer)))
	pool, _ := t.statepool(senderState, signerState)
	opr := t.processor(cp, pool)

	items := []CreateDocumentsItem{
		NewCreateDocumentsItemImpl(*bsDocData, cid),
	}

	cd := t.newOperation(sender.Address, items, sender.Privs())

	err := opr.Process(cd)
	t.Contains(err.Error(), "currency of holder does not exist")
}

func TestCreateDocumentsOperation(t *testing.T) {
	suite.Run(t, new(testCreateDocumentsOperation))
}
