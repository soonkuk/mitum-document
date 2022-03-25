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

type testUpdateDocumentsOperation struct {
	baseTestOperationProcessor
}

func (t *testUpdateDocumentsOperation) processor(cp *currency.CurrencyPool, pool *storage.Statepool) prprocessor.OperationProcessor {
	copr, err := NewOperationProcessor(nil).
		SetProcessor(UpdateDocumentsHinter, NewUpdateDocumentsProcessor(cp))
	t.NoError(err)

	if pool == nil {
		return copr
	}

	return copr.New(pool)
}

func (t *testUpdateDocumentsOperation) newOperation(sender base.Address, items []UpdateDocumentsItem, pks []key.Privatekey) UpdateDocuments {
	token := util.UUID().Bytes()
	fact := NewUpdateDocumentsFact(token, sender, items)

	var fs []base.FactSign
	for _, pk := range pks {
		sig, err := base.NewFactSignature(pk, fact, nil)
		if err != nil {
			panic(err)
		}

		fs = append(fs, base.NewBaseFactSign(pk.Publickey(), sig))
	}

	cd, err := NewUpdateDocuments(fact, fs, "")
	if err != nil {
		panic(err)
	}

	err = cd.IsValid(nil)
	if err != nil {
		panic(err)
	}

	return cd
}

func (t *testUpdateDocumentsOperation) TestNormalCase() {
	bsDocIDStr := "1sdi"
	bcUserDocIDStr := "1cui"
	bcLandDocIDStr := "1cli"
	bcVotingDocIDStr := "1cvi"
	bcHistoryDocIDStr := "1chi"
	// create BSDocData
	bsDocData, ownerAccount, signerAccount := newBSDocData("filehash", bsDocIDStr, account{})
	// create BCUserData
	bcUserData, _, _ := newBCUserData(bcUserDocIDStr, *ownerAccount)
	// create BCLandData
	bcLandData, _, renterAccount := newBCLandData(bcLandDocIDStr, *ownerAccount)
	// create BCVotingData
	bcVotingData, _, bossAccount := newBCVotingData(bcVotingDocIDStr, *ownerAccount)
	// create BCHistoryData
	bcHistoryData, _, cityAdminAccount := newBCHistoryData(bcHistoryDocIDStr, *ownerAccount)
	// currency id
	cid := currency.CurrencyID("SHOWME")
	// sender initial balance
	balance := []currency.Amount{
		currency.NewAmount(currency.NewBig(33), cid),
	}
	// prepare sender account state
	sender, senderState := t.newStateAccount(true, balance, ownerAccount)
	// prepare signer account state
	_, signerState := t.newStateAccount(true, balance, signerAccount)
	// prepare renter account state
	_, renterState := t.newStateAccount(true, balance, renterAccount)
	// prepare boss account state
	_, bossState := t.newStateAccount(true, balance, bossAccount)
	// prepare cityAdmin account state
	_, cityAdminState := t.newStateAccount(true, balance, cityAdminAccount)

	// prepare BSDocData state
	bsDocDataState := t.newStateDocument(ownerAccount.Address, bsDocData, DocumentInventory{})
	d, err := StateDocumentsValue(bsDocDataState[1])
	t.NoError(err)
	// prepare BCUserData state
	bcUserDataState := t.newStateDocument(ownerAccount.Address, bcUserData, d)
	d, err = StateDocumentsValue(bcUserDataState[1])
	// prepare BCLandData state
	bcLandDataState := t.newStateDocument(ownerAccount.Address, bcLandData, d)
	d, err = StateDocumentsValue(bcLandDataState[1])
	// prepare BCVotingData state
	bcVotingDataState := t.newStateDocument(ownerAccount.Address, bcVotingData, d)
	d, err = StateDocumentsValue(bcVotingDataState[1])
	// prepare BCHistoryData state
	bcHistoryDataState := t.newStateDocument(ownerAccount.Address, bcHistoryData, d)

	fee := currency.NewBig(1)
	feeer := currency.NewFixedFeeer(sender.Address, fee)
	cp := currency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(cid, currency.NewBig(99), sender.Address, feeer)))

	newSigner, newSignerState := t.newStateAccount(true, balance, nil)
	newbsDocData := *bsDocData
	(&newbsDocData).title = "update title"
	(&newbsDocData).size = currency.NewBig(99)
	(&newbsDocData).signers = []DocSign{MustNewDocSign(newSigner.Address, "signcode2", false)}

	newbcUserData := *bcUserData
	newStatistics := MustNewUserStatistics(
		uint(99), uint(99), uint(99), uint(99), uint(99), uint(99), uint(99),
	)
	(&newbcUserData).gold = uint(99)
	(&newbcUserData).bankgold = uint(99)
	(&newbcUserData).statistics = newStatistics

	newRenter, newRenterState := t.newStateAccount(true, balance, nil)
	newbcLandData := *bcLandData
	(&newbcLandData).account = newRenter.Address
	(&newbcLandData).address = "update address"
	(&newbcLandData).area = "update area"
	(&newbcLandData).periodday = uint(99)
	(&newbcLandData).rentdate = "update rentdata"
	(&newbcLandData).renter = "update renter"

	newBoss, newBossState := t.newStateAccount(true, balance, nil)
	newbcVotingData := *bcVotingData
	(&newbcVotingData).account = newBoss.Address
	(&newbcVotingData).bossname = "update bossname"
	(&newbcVotingData).candidates = []VotingCandidate{MustNewVotingCandidate(newBoss.Address, "new boss", "new manifest", uint(99))}
	(&newbcVotingData).endVoteTime = "update endVoteTime"
	(&newbcVotingData).round = uint(99)
	(&newbcVotingData).termofoffice = "update termofoffice"

	newCityAdmin, newCityAdminState := t.newStateAccount(true, balance, nil)
	newbcHistoryData := *bcHistoryData
	(&newbcHistoryData).account = newCityAdmin.Address
	(&newbcHistoryData).application = "update application"
	(&newbcHistoryData).date = "update date"
	(&newbcHistoryData).name = "update name"
	(&newbcHistoryData).usage = "update usage"

	// create UpdateDocumentsItems
	items := []UpdateDocumentsItem{
		NewUpdateDocumentsItemImpl(newbsDocData, cid),
		NewUpdateDocumentsItemImpl(newbcUserData, cid),
		NewUpdateDocumentsItemImpl(newbcLandData, cid),
		NewUpdateDocumentsItemImpl(newbcVotingData, cid),
		NewUpdateDocumentsItemImpl(newbcHistoryData, cid),
	}

	// state pool
	pool, _ := t.statepool(
		senderState,
		signerState,
		newSignerState,
		renterState,
		newRenterState,
		bossState,
		newBossState,
		cityAdminState,
		newCityAdminState,
		[]state.State{bsDocDataState[0]},
		[]state.State{bcUserDataState[0]},
		[]state.State{bcLandDataState[0]},
		[]state.State{bcVotingDataState[0]},
		bcHistoryDataState,
	)

	opr := t.processor(cp, pool)
	cd := t.newOperation(sender.Address, items, sender.Privs())

	t.NoError(opr.Process(cd))

	// new document state
	var newBSDocDataDocumentState, newBCUserDataDocumentState, newBCLandDataDocumentState state.State
	var newBCVotingDataDocumentState, newBCHistoryDataDocumentState state.State
	// sender balance state
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

	// sender balance state should not be nil
	t.NotNil(senderBalanceState)

	senderBalanceAmount, _ := currency.StateBalanceValue(senderBalanceState)
	// amount of sender balance state should be deductes by fee
	t.True(senderBalanceAmount.Big().Equal(balance[0].Big().Sub(currency.NewBig(5).Mul(fee))))
	// processed fee amount should be same with currencydesign fee
	t.Equal(currency.NewBig(5).Mul(fee), senderBalanceState.(currency.AmountState).Fee())
	// check new documentData state
	newBSDocData, _ := StateDocumentDataValue(newBSDocDataDocumentState)
	nd0, ok := newBSDocData.(BSDocData)
	t.True(ok)
	t.True(!nd0.Equal(*bsDocData))
	t.True(nd0.Equal(newbsDocData))

	newBCUserData, _ := StateDocumentDataValue(newBCUserDataDocumentState)
	nd1, ok := newBCUserData.(BCUserData)
	t.True(ok)
	t.True(!nd1.Equal(*bcUserData))
	t.True(nd1.Equal(newbcUserData))

	newBCLandData, _ := StateDocumentDataValue(newBCLandDataDocumentState)
	nd2, ok := newBCLandData.(BCLandData)
	t.True(ok)
	t.True(!nd2.Equal(*bcLandData))
	t.True(nd2.Equal(newbcLandData))

	newBCVotingData, _ := StateDocumentDataValue(newBCVotingDataDocumentState)
	nd3, ok := newBCVotingData.(BCVotingData)
	t.True(ok)
	t.True(!nd3.Equal(*bcVotingData))
	t.True(nd3.Equal(newbcVotingData))

	newBCHistoryData, _ := StateDocumentDataValue(newBCHistoryDataDocumentState)
	nd4, ok := newBCHistoryData.(BCHistoryData)
	t.True(ok)
	t.True(!nd4.Equal(*bcHistoryData))
	t.True(nd4.Equal(newbcHistoryData))
}

func (t *testUpdateDocumentsOperation) TestAccountsNotExist() {
	// currency id
	cid := currency.CurrencyID("SHOWME")
	// sender initial balance
	balance := []currency.Amount{
		currency.NewAmount(currency.NewBig(33), cid),
	}
	// create BSDocID
	bsDocIDStr := "1sdi"
	// create BSDocData
	bsDocData, ownerAccount, _ := newBSDocData("filehash", bsDocIDStr, account{})
	// sender account
	sender, senderState := t.newStateAccount(true, balance, ownerAccount)
	// prepare BSDocData state
	bsDocDataState := t.newStateDocument(ownerAccount.Address, bsDocData, DocumentInventory{})
	// feeer
	fee := currency.NewBig(1)
	feeer := currency.NewFixedFeeer(sender.Address, fee)
	// currencyPool
	cp := currency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(cid, currency.NewBig(99), sender.Address, feeer)))
	// create UpdateDocumentsItem with BSDocData
	items := []UpdateDocumentsItem{
		NewUpdateDocumentsItemImpl(*bsDocData, cid),
	}
	pool, _ := t.statepool(senderState, bsDocDataState)
	// state pool
	opr := t.processor(cp, pool)
	cd := t.newOperation(sender.Address, items, sender.Privs())
	err := opr.Process(cd)

	t.Contains(err.Error(), "documentData related accounts not found")

	// create BCLandDocID
	bcLandDocIDStr := "1cli"
	// create BCLandData
	bcLandData, _, _ := newBCLandData(bcLandDocIDStr, *ownerAccount)
	// prepare BCLandData state
	bcLandDataState := t.newStateDocument(ownerAccount.Address, bcLandData, DocumentInventory{})
	// create document item with BCLandData
	items = []UpdateDocumentsItem{
		NewUpdateDocumentsItemImpl(*bcLandData, cid),
	}
	pool, _ = t.statepool(senderState, bcLandDataState)
	opr = t.processor(cp, pool)
	cd = t.newOperation(sender.Address, items, sender.Privs())
	err = opr.Process(cd)
	t.Contains(err.Error(), "documentData related accounts not found")

	// create BCVotingDocID
	bcVotingDocIDStr := "1cvi"
	// create BCVotingData
	bcVotingData, _, _ := newBCVotingData(bcVotingDocIDStr, *ownerAccount)
	// prepare BCVotingData state
	bcVotingDataState := t.newStateDocument(ownerAccount.Address, bcVotingData, DocumentInventory{})
	// create document item with BCUserData
	items = []UpdateDocumentsItem{
		NewUpdateDocumentsItemImpl(*bcVotingData, cid),
	}
	pool, _ = t.statepool(senderState, bcVotingDataState)
	opr = t.processor(cp, pool)
	cd = t.newOperation(sender.Address, items, sender.Privs())
	err = opr.Process(cd)
	t.Contains(err.Error(), "documentData related accounts not found")

	// create BCHistoryDocID
	bcHistoryDocIDStr := "1chi"
	// create BCHistoryData
	bcHistoryData, _, _ := newBCHistoryData(bcHistoryDocIDStr, *ownerAccount)
	// prepare BCHistoryData state
	bcHistoryDataState := t.newStateDocument(ownerAccount.Address, bcHistoryData, DocumentInventory{})
	// create document item with BCUserData
	items = []UpdateDocumentsItem{
		NewUpdateDocumentsItemImpl(*bcHistoryData, cid),
	}
	pool, _ = t.statepool(senderState, bcHistoryDataState)
	opr = t.processor(cp, pool)
	cd = t.newOperation(sender.Address, items, sender.Privs())
	err = opr.Process(cd)

	t.Contains(err.Error(), "documentData related accounts not found")
}

func (t *testUpdateDocumentsOperation) TestDocumentNotExists() {
	// currency id
	cid := currency.CurrencyID("SHOWME")
	// sender initial balance
	balance := []currency.Amount{
		currency.NewAmount(currency.NewBig(33), cid),
	}
	// create BSDocID
	bsDocIDStr := "1sdi"
	// create BSDocData
	bsDocData, ownerAccount, _ := newBSDocData("filehash", bsDocIDStr, account{})
	// sender account
	sender, senderState := t.newStateAccount(true, balance, ownerAccount)
	// feeer
	fee := currency.NewBig(1)
	feeer := currency.NewFixedFeeer(sender.Address, fee)
	// currencyPool
	cp := currency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(cid, currency.NewBig(99), sender.Address, feeer)))
	// create UpdateDocumentsItem with BSDocData
	items := []UpdateDocumentsItem{
		NewUpdateDocumentsItemImpl(*bsDocData, cid),
	}
	// no document state and documents state
	pool, _ := t.statepool(senderState)
	// state pool
	opr := t.processor(cp, pool)
	cd := t.newOperation(sender.Address, items, sender.Privs())
	err := opr.Process(cd)

	t.Contains(err.Error(), "owner has no document inventory")

	// prepare BSDocData state
	bsDocDataState := t.newStateDocument(ownerAccount.Address, bsDocData, DocumentInventory{})
	documentsState := []state.State{bsDocDataState[1]}
	// put documents state only
	pool, _ = t.statepool(senderState, documentsState)
	// state pool
	opr = t.processor(cp, pool)
	cd = t.newOperation(sender.Address, items, sender.Privs())
	err = opr.Process(cd)

	t.Contains(err.Error(), "document not registered with documentid")
}

func (t *testUpdateDocumentsOperation) TestSameSenders() {
	// currency id
	cid := currency.CurrencyID("SHOWME")
	// sender initial balance
	balance := []currency.Amount{
		currency.NewAmount(currency.NewBig(33), cid),
	}
	// create BSDocID
	bsDocIDStr := "1sdi"
	// create BCUserDocID
	bcUserDocIDStr := "1cui"
	// create BSDocData
	bsDocData, ownerAccount, signerAccount := newBSDocData("filehash", bsDocIDStr, account{})
	// create BCLandData
	bcUserData, _, _ := newBCUserData(bcUserDocIDStr, *ownerAccount)
	// sender account state
	sender, senderState := t.newStateAccount(true, balance, ownerAccount)
	// signer account state
	_, signerState := t.newStateAccount(true, balance, signerAccount)
	// prepare document and documents state
	bsDocDataState := t.newStateDocument(ownerAccount.Address, bsDocData, DocumentInventory{})
	d, err := StateDocumentsValue(bsDocDataState[1])
	bcUserDataState := t.newStateDocument(ownerAccount.Address, bcUserData, d)
	// feeer
	fee := currency.NewBig(1)
	feeer := currency.NewFixedFeeer(sender.Address, fee)
	// currencyPool
	cp := currency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(cid, currency.NewBig(99), sender.Address, feeer)))
	// create UpdateDocumentsItem with BSDocData
	items0 := []UpdateDocumentsItem{
		NewUpdateDocumentsItemImpl(*bsDocData, cid),
	}
	items1 := []UpdateDocumentsItem{
		NewUpdateDocumentsItemImpl(*bcUserData, cid),
	}
	// no document state and documents state
	pool, _ := t.statepool(senderState, signerState, []state.State{bsDocDataState[0]}, bcUserDataState)
	// state pool
	opr := t.processor(cp, pool)
	cd0 := t.newOperation(sender.Address, items0, sender.Privs())
	cd1 := t.newOperation(sender.Address, items1, sender.Privs())
	err = opr.Process(cd0)
	t.NoError(err)
	err = opr.Process(cd1)
	t.Contains(err.Error(), "violates only one sender")
}

func (t *testUpdateDocumentsOperation) TestInSufficientBalanceForFee() {
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
	// sender account
	sender, senderState := t.newStateAccount(true, balance, ownerAccount)
	// sigenr account
	_, signerState := t.newStateAccount(true, balance, signerAccount)
	// prepare BSDocData state
	bsDocDataState := t.newStateDocument(ownerAccount.Address, bsDocData, DocumentInventory{})
	// feeer
	fee := currency.NewBig(11)
	feeer := currency.NewFixedFeeer(sender.Address, fee)
	// currencyPool
	cp := currency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(cid, currency.NewBig(99), sender.Address, feeer)))
	// create UpdateDocumentsItem with BSDocData
	items := []UpdateDocumentsItem{
		NewUpdateDocumentsItemImpl(*bsDocData, cid),
	}
	pool, _ := t.statepool(senderState, signerState, bsDocDataState)
	// state pool
	opr := t.processor(cp, pool)
	cd := t.newOperation(sender.Address, items, sender.Privs())
	err := opr.Process(cd)

	t.Contains(err.Error(), "insufficient balance")
}

func (t *testUpdateDocumentsOperation) TestUnknownCurrencyID() {
	// currency id
	cid0 := currency.CurrencyID("SHOWME")
	cid1 := currency.CurrencyID("FINDME")
	// sender initial balance
	balance := []currency.Amount{
		currency.NewAmount(currency.NewBig(33), cid1),
	}
	// create BSDocID
	bsDocIDStr := "1sdi"
	// create BSDocData
	bsDocData, ownerAccount, signerAccount := newBSDocData("filehash", bsDocIDStr, account{})
	// sender account
	sender, senderState := t.newStateAccount(true, balance, ownerAccount)
	// sigenr account
	_, signerState := t.newStateAccount(true, balance, signerAccount)
	// prepare BSDocData state
	bsDocDataState := t.newStateDocument(ownerAccount.Address, bsDocData, DocumentInventory{})
	// feeer
	fee := currency.NewBig(1)
	feeer := currency.NewFixedFeeer(sender.Address, fee)
	// currencyPool
	cp := currency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(cid1, currency.NewBig(99), sender.Address, feeer)))
	// create UpdateDocumentsItem with BSDocData
	items := []UpdateDocumentsItem{
		NewUpdateDocumentsItemImpl(*bsDocData, cid0),
	}
	pool, _ := t.statepool(senderState, signerState, bsDocDataState)
	// state pool
	opr := t.processor(cp, pool)
	cd := t.newOperation(sender.Address, items, sender.Privs())
	err := opr.Process(cd)

	t.Contains(err.Error(), "unknown currency id found")
}

func (t *testUpdateDocumentsOperation) TestSenderNotHaveCurrency() {
	// currency id
	cid0 := currency.CurrencyID("SHOWME")
	cid1 := currency.CurrencyID("FINDME")
	// sender initial balance
	balance := []currency.Amount{
		currency.NewAmount(currency.NewBig(33), cid1),
	}
	// create BSDocID
	bsDocIDStr := "1sdi"
	// create BSDocData
	bsDocData, ownerAccount, signerAccount := newBSDocData("filehash", bsDocIDStr, account{})
	// sender account
	sender, senderState := t.newStateAccount(true, balance, ownerAccount)
	// sigenr account
	_, signerState := t.newStateAccount(true, balance, signerAccount)
	// prepare BSDocData state
	bsDocDataState := t.newStateDocument(ownerAccount.Address, bsDocData, DocumentInventory{})
	// feeer
	fee := currency.NewBig(1)
	feeer := currency.NewFixedFeeer(sender.Address, fee)
	// currencyPool
	cp := currency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(cid0, currency.NewBig(99), sender.Address, feeer)))
	t.NoError(cp.Set(t.newCurrencyDesignState(cid1, currency.NewBig(99), sender.Address, feeer)))
	// create UpdateDocumentsItem with BSDocData
	items := []UpdateDocumentsItem{
		NewUpdateDocumentsItemImpl(*bsDocData, cid0),
	}
	pool, _ := t.statepool(senderState, signerState, bsDocDataState)
	// state pool
	opr := t.processor(cp, pool)
	cd := t.newOperation(sender.Address, items, sender.Privs())
	err := opr.Process(cd)

	t.Contains(err.Error(), "currency of holder does not exist")
}

func (t *testUpdateDocumentsOperation) TestSenderBalanceNotExist() {
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
	// sender account
	sender, senderState := t.newStateAccount(true, nil, ownerAccount)
	// sigenr account
	_, signerState := t.newStateAccount(true, balance, signerAccount)
	// prepare BSDocData state
	bsDocDataState := t.newStateDocument(ownerAccount.Address, bsDocData, DocumentInventory{})
	// feeer
	fee := currency.NewBig(1)
	feeer := currency.NewFixedFeeer(sender.Address, fee)
	// currencyPool
	cp := currency.NewCurrencyPool()
	t.NoError(cp.Set(t.newCurrencyDesignState(cid, currency.NewBig(99), sender.Address, feeer)))
	// create UpdateDocumentsItem with BSDocData
	items := []UpdateDocumentsItem{
		NewUpdateDocumentsItemImpl(*bsDocData, cid),
	}
	pool, _ := t.statepool(senderState, signerState, bsDocDataState)
	// state pool
	opr := t.processor(cp, pool)
	cd := t.newOperation(sender.Address, items, sender.Privs())
	err := opr.Process(cd)

	t.Contains(err.Error(), "currency of holder does not exist")
}

func TestUpdateDocumentsOperation(t *testing.T) {
	suite.Run(t, new(testUpdateDocumentsOperation))
}
