package document

import (
	"bytes"
	"sort"

	"github.com/pkg/errors"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
	"github.com/spikeekips/mitum/util/valuehash"
)

var (
	DocumentType   = hint.Type("mitum-blockcity-document")
	DocumentHint   = hint.NewHint(DocumentType, "v0.0.1")
	DocumentHinter = Document{BaseHinter: hint.NewBaseHinter(DocumentHint)}
)

type Document struct {
	hint.BaseHinter
	data DocumentData
}

func NewDocument(doc DocumentData) Document {
	d := Document{
		BaseHinter: hint.NewBaseHinter(DocumentHint),
		data:       doc,
	}

	return d
}

func MustNewDocument(doc DocumentData) Document {
	d := NewDocument(doc)
	if err := d.data.IsValid(nil); err != nil {
		panic(err)
	}

	return d
}

func (doc Document) Owner() base.Address {
	return doc.data.Owner()
}

func (doc Document) Hint() hint.Hint {
	return doc.BaseHinter.Hint()
}

func (doc Document) DocumentData() DocumentData {
	return doc.data
}

func (doc Document) DocumentId() string {
	return doc.data.DocumentId()
}

func (doc Document) DocumentType() hint.Type {
	return doc.data.DocumentType()
}

func (doc Document) Bytes() []byte {
	return doc.DocumentData().Bytes()
}

func (doc Document) Hash() valuehash.Hash {
	return doc.GenerateHash()
}

func (doc Document) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(doc.Bytes())
}

func (doc Document) IsValid([]byte) error {

	if err := isvalid.Check(
		nil, false,
		doc.BaseHinter,
		doc.data,
	); err != nil {
		return isvalid.InvalidError.Errorf("invalid Document: %w", err)
	}
	return nil
}

type DocumentData interface {
	DocumentId() string
	DocumentType() hint.Type
	Hint() hint.Hint
	Bytes() []byte
	Hash() valuehash.Hash
	GenerateHash() valuehash.Hash
	Owner() base.Address
	Accounts() []base.Address
	Info() DocInfo
	IsValid([]byte) error
}

var (
	SignDocuDataType   = hint.Type("mitum-blocksign-document-data")
	SignDocuDataHint   = hint.NewHint(SignDocuDataType, "v0.0.1")
	SignDocuDataHinter = SignDocuData{BaseHinter: hint.NewBaseHinter(SignDocuDataHint)}
)

type SignDocuData struct {
	hint.BaseHinter
	info     DocInfo
	fileHash FileHash
	creator  DocSign
	title    string
	size     currency.Big
	signers  []DocSign
}

var (
	BCUserDataType   = hint.Type("mitum-blockcity-document-user-data")
	BCUserDataHint   = hint.NewHint(BCUserDataType, "v0.0.1")
	BCUserDataHinter = BCUserData{BaseHinter: hint.NewBaseHinter(BCUserDataHint)}
)

type BCUserData struct {
	hint.BaseHinter
	info       DocInfo
	owner      base.Address
	gold       uint
	bankgold   uint
	statistics UserStatistics
}

func NewBCUserData(info DocInfo,
	owner base.Address,
	gold,
	bankgold uint,
	statistics UserStatistics,
) BCUserData {
	doc := BCUserData{
		BaseHinter: hint.NewBaseHinter(BCUserDataHint),
		info:       info,
		owner:      owner,
		gold:       gold,
		bankgold:   bankgold,
		statistics: statistics,
	}

	return doc
}

func MustNewBCUserData(info DocInfo, owner base.Address, gold, bankgold uint, statistics UserStatistics) BCUserData {
	doc := NewBCUserData(info, owner, gold, bankgold, statistics)
	if err := doc.IsValid(nil); err != nil {
		panic(err)
	}

	return doc
}

func (doc BCUserData) DocumentId() string {
	return doc.info.DocumentId()
}

func (doc BCUserData) DocumentType() hint.Type {
	return doc.info.docType
}

func (doc BCUserData) Bytes() []byte {
	bs := make([][]byte, 5)

	bs[0] = doc.info.Bytes()
	bs[1] = doc.owner.Bytes()
	bs[2] = util.UintToBytes(doc.gold)
	bs[3] = util.UintToBytes(doc.bankgold)
	bs[4] = doc.statistics.Bytes()

	return util.ConcatBytesSlice(bs...)
}

func (doc BCUserData) Hash() valuehash.Hash {
	return doc.GenerateHash()
}

func (doc BCUserData) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(doc.Bytes())
}

func (doc BCUserData) IsEmpty() bool {
	return len(doc.info.DocType()) < 1
}

func (doc BCUserData) IsValid([]byte) error {
	if doc.info.docType != doc.Hint().Type() {
		return errors.Errorf("DocInfo not matched with DocumentData Type : DocInfo type %v, DocumentData type %v", doc.info.docType, doc.Hint().Type())
	}

	if err := isvalid.Check(
		nil, false,
		doc.BaseHinter,
		doc.info,
		doc.owner,
		doc.statistics,
	); err != nil {
		return isvalid.InvalidError.Errorf("invalid User Document Data: %w", err)
	}

	return nil
}

func (doc BCUserData) Owner() base.Address {
	return doc.owner
}

func (doc BCUserData) Accounts() []base.Address {
	return nil
}

func (doc BCUserData) Info() DocInfo {
	return doc.info
}

func (doc BCUserData) Equal(b BCUserData) bool {

	if doc.info.DocType() != b.info.DocType() {
		return false
	}

	if !doc.owner.Equal(b.owner) {
		return false
	}

	if doc.gold != b.gold {
		return false
	}

	if doc.bankgold != b.bankgold {
		return false
	}

	if !doc.statistics.Equal(b.statistics) {
		return false
	}

	return true
}

var (
	BCLandDataType   = hint.Type("mitum-blockcity-document-land-data")
	BCLandDataHint   = hint.NewHint(BCLandDataType, "v0.0.1")
	BCLandDataHinter = BCLandData{BaseHinter: hint.NewBaseHinter(BCLandDataHint)}
)

type BCLandData struct {
	hint.BaseHinter
	info      DocInfo
	owner     base.Address
	address   string
	area      string
	renter    string
	account   base.Address
	rentdate  string
	periodday uint
}

func NewBCLandData(info DocInfo,
	owner base.Address,
	address, area, renter string,
	account base.Address,
	rentdate string,
	periodday uint,
) BCLandData {
	doc := BCLandData{
		BaseHinter: hint.NewBaseHinter(BCLandDataHint),
		info:       info,
		owner:      owner,
		address:    address,
		area:       area,
		renter:     renter,
		account:    account,
		rentdate:   rentdate,
		periodday:  periodday,
	}
	return doc
}

func MustNewBCLandData(info DocInfo,
	owner base.Address,
	address, area, renter string,
	account base.Address,
	rentdate string,
	periodday uint,
) BCLandData {
	doc := NewBCLandData(info, owner, address, area, renter, account, rentdate, periodday)
	if err := doc.IsValid(nil); err != nil {
		panic(err)
	}
	return doc
}

func (doc BCLandData) DocumentId() string {
	return doc.info.DocumentId()
}

func (doc BCLandData) DocumentType() hint.Type {
	return doc.info.docType
}

func (doc BCLandData) Bytes() []byte {
	bs := make([][]byte, 8)

	bs[0] = doc.info.Bytes()
	bs[1] = doc.owner.Bytes()
	bs[2] = []byte(doc.address)
	bs[3] = []byte(doc.area)
	bs[4] = []byte(doc.renter)
	bs[5] = doc.account.Bytes()
	bs[6] = []byte(doc.rentdate)
	bs[7] = util.UintToBytes(doc.periodday)
	return util.ConcatBytesSlice(bs...)
}

func (doc BCLandData) Hash() valuehash.Hash {
	return doc.GenerateHash()
}

func (doc BCLandData) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(doc.Bytes())
}

func (doc BCLandData) IsValid([]byte) error {
	if doc.info.docType != doc.Hint().Type() {
		return errors.Errorf("DocInfo not matched with DocumentData Type : DocInfo type %v, DocumentData type %v", doc.info.docType, doc.Hint().Type())
	}

	if err := isvalid.Check(
		nil, false,
		doc.BaseHinter,
		doc.info,
		doc.owner,
	); err != nil {
		return errors.Wrap(err, "Invalid Land document data")
	}
	return nil
}

func (doc BCLandData) Info() DocInfo {
	return doc.info
}

func (doc BCLandData) Accounts() []base.Address {
	return []base.Address{}
}

func (doc BCLandData) Owner() base.Address {
	return doc.owner
}

func (doc BCLandData) Equal(b BCLandData) bool {

	if !doc.info.Equal(b.info) {
		return false
	}

	if doc.address != b.address {
		return false
	}

	if doc.area != b.area {
		return false
	}

	if doc.renter != b.renter {
		return false
	}

	if !doc.account.Equal(b.account) {
		return false
	}

	if doc.rentdate != b.rentdate {
		return false
	}

	if doc.periodday != b.periodday {
		return false
	}

	return true
}

var (
	BCVotingDataType   = hint.Type("mitum-blockcity-document-voting-data")
	BCVotingDataHint   = hint.NewHint(BCVotingDataType, "v0.0.1")
	BCVotingDataHinter = BCVotingData{BaseHinter: hint.NewBaseHinter(BCVotingDataHint)}
)

type BCVotingData struct {
	hint.BaseHinter
	info         DocInfo
	owner        base.Address
	round        uint
	endVoteTime  string
	candidates   []VotingCandidate
	bossname     string
	account      base.Address
	termofoffice string
}

func NewBCVotingData(info DocInfo,
	owner base.Address,
	round uint,
	endVoteTime string,
	candidates []VotingCandidate,
	bossname string,
	account base.Address,
	termofoffice string,
) BCVotingData {
	doc := BCVotingData{
		BaseHinter:   hint.NewBaseHinter(BCVotingDataHint),
		info:         info,
		owner:        owner,
		round:        round,
		endVoteTime:  endVoteTime,
		candidates:   candidates,
		bossname:     bossname,
		account:      account,
		termofoffice: termofoffice,
	}
	return doc
}

func MustNewBCVotingData(
	info DocInfo,
	owner base.Address,
	round uint,
	endVoteTime string,
	candidates []VotingCandidate,
	bossname string,
	account base.Address,
	termofoffice string,
) BCVotingData {
	doc := NewBCVotingData(info, owner, round, endVoteTime, candidates, bossname, account, termofoffice)
	if err := doc.IsValid(nil); err != nil {
		panic(err)
	}
	return doc
}

func (doc BCVotingData) DocumentId() string {
	return doc.info.DocumentId()
}

func (doc BCVotingData) DocumentType() hint.Type {
	return doc.info.docType
}

func (doc BCVotingData) Bytes() []byte {
	bs := make([][]byte, len(doc.candidates)+7)

	sort.Slice(doc.candidates, func(i, j int) bool {
		return bytes.Compare(doc.candidates[i].Bytes(), doc.candidates[j].Bytes()) < 0
	})

	bs[0] = doc.info.Bytes()
	bs[1] = doc.owner.Bytes()
	bs[2] = util.UintToBytes(doc.round)
	bs[3] = []byte(doc.endVoteTime)
	bs[4] = []byte(doc.bossname)
	bs[5] = doc.account.Bytes()
	bs[6] = []byte(doc.termofoffice)
	for i := range doc.candidates {
		bs[i+7] = doc.candidates[i].Bytes()
	}
	return util.ConcatBytesSlice(bs...)
}

func (doc BCVotingData) Hash() valuehash.Hash {
	return doc.GenerateHash()
}

func (doc BCVotingData) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(doc.Bytes())
}

func (doc BCVotingData) IsValid([]byte) error {
	if doc.info.docType != doc.Hint().Type() {
		return errors.Errorf("DocInfo not matched with DocumentData Type : DocInfo type %v, DocumentData type %v", doc.info.docType, doc.Hint().Type())
	}

	if err := isvalid.Check(
		nil, false,
		doc.BaseHinter,
		doc.info,
		doc.owner,
	); err != nil {
		return errors.Wrap(err, "Invalid Voting document data")
	}

	for i := range doc.candidates {
		if err := doc.candidates[i].IsValid(nil); err != nil {
			return err
		}
	}

	return nil
}

func (doc BCVotingData) Info() DocInfo {
	return doc.info
}

func (doc BCVotingData) Accounts() []base.Address {
	var accounts []base.Address
	accounts = append(accounts, doc.account)
	for i := range doc.candidates {
		accounts = append(accounts, doc.candidates[i].address)
	}
	return accounts
}

func (doc BCVotingData) Owner() base.Address {
	return doc.owner
}

func (doc BCVotingData) Candidates() []VotingCandidate {
	sort.Slice(doc.candidates, func(i, j int) bool {
		return bytes.Compare(doc.candidates[i].Bytes(), doc.candidates[j].Bytes()) < 0
	})
	return doc.candidates
}

func (doc BCVotingData) Equal(b BCVotingData) bool {

	if !doc.info.Equal(b.info) {
		return false
	}

	if !doc.owner.Equal(b.owner) {
		return false
	}

	sort.Slice(doc.candidates, func(i, j int) bool {
		return bytes.Compare(doc.candidates[i].Bytes(), doc.candidates[j].Bytes()) < 0
	})

	sort.Slice(b.candidates, func(i, j int) bool {
		return bytes.Compare(b.candidates[i].Bytes(), b.candidates[j].Bytes()) < 0
	})

	for i := range doc.candidates {
		if !doc.candidates[i].Equal(b.candidates[i]) {
			return false
		}
	}
	return true
}

var (
	BCHistoryDataType   = hint.Type("mitum-blockcity-document-history-data")
	BCHistoryDataHint   = hint.NewHint(BCHistoryDataType, "v0.0.1")
	BCHistoryDataHinter = BCHistoryData{BaseHinter: hint.NewBaseHinter(BCHistoryDataHint)}
)

type BCHistoryData struct {
	hint.BaseHinter
	info        DocInfo
	owner       base.Address
	name        string
	account     base.Address
	date        string
	usage       string
	application string
}

func NewBCHistoryData(info DocInfo,
	owner base.Address,
	name string,
	account base.Address,
	date, usage, application string,
) BCHistoryData {
	doc := BCHistoryData{
		BaseHinter:  hint.NewBaseHinter(BCHistoryDataHint),
		info:        info,
		owner:       owner,
		name:        name,
		account:     account,
		date:        date,
		usage:       usage,
		application: application,
	}
	return doc
}

func MustNewBCHistoryData(info DocInfo,
	owner base.Address,
	name string,
	account base.Address,
	date, usage, application string,
) BCHistoryData {
	doc := NewBCHistoryData(info, owner, name, account, date, usage, application)
	if err := doc.IsValid(nil); err != nil {
		panic(err)
	}
	return doc
}

func (doc BCHistoryData) DocumentId() string {
	return doc.info.DocumentId()
}

func (doc BCHistoryData) DocumentType() hint.Type {
	return doc.info.docType
}

func (doc BCHistoryData) Bytes() []byte {
	bs := make([][]byte, 7)

	bs[0] = doc.info.Bytes()
	bs[1] = doc.owner.Bytes()
	bs[2] = []byte(doc.name)
	bs[3] = doc.account.Bytes()
	bs[4] = []byte(doc.date)
	bs[5] = []byte(doc.usage)
	bs[6] = []byte(doc.application)

	return util.ConcatBytesSlice(bs...)
}

func (doc BCHistoryData) Hash() valuehash.Hash {
	return doc.GenerateHash()
}

func (doc BCHistoryData) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(doc.Bytes())
}

func (doc BCHistoryData) IsValid([]byte) error {
	if doc.info.docType != doc.Hint().Type() {
		return errors.Errorf("DocInfo not matched with DocumentData Type : DocInfo type %v, DocumentData type %v", doc.info.docType, doc.Hint().Type())
	}

	if err := isvalid.Check(
		nil, false,
		doc.BaseHinter,
		doc.info,
		doc.owner,
		doc.account,
	); err != nil {
		return errors.Wrap(err, "Invalid history document data")
	}
	return nil
}

func (doc BCHistoryData) Info() DocInfo {
	return doc.info
}

func (doc BCHistoryData) Accounts() []base.Address {
	return []base.Address{doc.account}
}

func (doc BCHistoryData) Owner() base.Address {
	return doc.owner
}

func (doc BCHistoryData) Equal(b BCHistoryData) bool {

	if !doc.info.Equal(b.info) {
		return false
	}

	if doc.name != b.name {
		return false
	}

	if !doc.account.Equal(b.account) {
		return false
	}

	if doc.date != b.date {
		return false
	}

	if doc.usage != b.usage {
		return false
	}

	if doc.application != b.application {
		return false
	}

	return true
}
