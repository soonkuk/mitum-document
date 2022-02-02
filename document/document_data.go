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
	CityUserDataType   = hint.Type("mitum-blockcity-document-user-data")
	CityUserDataHint   = hint.NewHint(CityUserDataType, "v0.0.1")
	CityUserDataHinter = CityUserData{BaseHinter: hint.NewBaseHinter(CityUserDataHint)}
)

type CityUserData struct {
	hint.BaseHinter
	info       DocInfo
	owner      base.Address
	gold       currency.Big
	bankgold   currency.Big
	statistics UserStatistics
}

func NewCityUserData(info DocInfo,
	owner base.Address,
	gold,
	bankgold currency.Big,
	statistics UserStatistics,
) CityUserData {
	doc := CityUserData{
		BaseHinter: hint.NewBaseHinter(CityUserDataHint),
		info:       info,
		owner:      owner,
		gold:       gold,
		bankgold:   bankgold,
		statistics: statistics,
	}

	return doc
}

func MustNewCityUserData(info DocInfo, owner base.Address, gold, bankgold currency.Big, statistics UserStatistics) CityUserData {
	doc := NewCityUserData(info, owner, gold, bankgold, statistics)
	if err := doc.IsValid(nil); err != nil {
		panic(err)
	}

	return doc
}

func (doc CityUserData) DocumentId() string {
	return doc.info.DocumentId()
}

func (doc CityUserData) DocumentType() hint.Type {
	return doc.info.docType
}

func (doc CityUserData) Bytes() []byte {
	bs := make([][]byte, 5)

	bs[0] = doc.info.Bytes()
	bs[1] = doc.owner.Bytes()
	bs[2] = doc.gold.Bytes()
	bs[3] = doc.bankgold.Bytes()
	bs[4] = doc.statistics.Bytes()

	return util.ConcatBytesSlice(bs...)
}

func (doc CityUserData) Hash() valuehash.Hash {
	return doc.GenerateHash()
}

func (doc CityUserData) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(doc.Bytes())
}

func (doc CityUserData) IsEmpty() bool {
	return len(doc.info.DocType()) < 1
}

func (doc CityUserData) IsValid([]byte) error {
	if err := isvalid.Check(
		nil, false,
		doc.BaseHinter,
		doc.info,
		doc.owner,
		doc.gold,
		doc.bankgold,
		doc.statistics,
	); err != nil {
		return isvalid.InvalidError.Errorf("invalid Document User Data: %w", err)
	}

	return nil
}

func (doc CityUserData) Owner() base.Address {
	return doc.owner
}

func (doc CityUserData) Gold() currency.Big {
	return doc.gold
}

func (doc CityUserData) Bankgold() currency.Big {
	return doc.bankgold
}

func (doc CityUserData) Info() DocInfo {
	return doc.info
}

func (doc CityUserData) Equal(b CityUserData) bool {

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

func (doc CityUserData) WithData(info DocInfo, owner base.Address, gold, bankgold currency.Big, statistics UserStatistics) CityUserData {
	doc.info = info
	doc.owner = owner
	doc.gold = gold
	doc.bankgold = bankgold
	doc.statistics = statistics
	return doc
}

var (
	CityLandDataType   = hint.Type("mitum-blockcity-document-land-data")
	CityLandDataHint   = hint.NewHint(CityLandDataType, "v0.0.1")
	CityLandDataHinter = CityLandData{BaseHinter: hint.NewBaseHinter(CityLandDataHint)}
)

type CityLandData struct {
	hint.BaseHinter
	info      DocInfo
	owner     base.Address
	lender    base.Address
	starttime string
	periodday uint
}

func NewCityLandData(info DocInfo, owner, lender base.Address, starttime string, periodday uint) CityLandData {
	doc := CityLandData{
		BaseHinter: hint.NewBaseHinter(CityLandDataHint),
		info:       info,
		owner:      owner,
		lender:     lender,
		starttime:  starttime,
		periodday:  periodday,
	}
	return doc
}

func MustNewCityLandData(info DocInfo, owner, lender base.Address, starttime string, periodday uint) CityLandData {
	doc := NewCityLandData(info, owner, lender, starttime, periodday)
	if err := doc.IsValid(nil); err != nil {
		panic(err)
	}
	return doc
}

func (doc CityLandData) DocumentId() string {
	return doc.info.DocumentId()
}

func (doc CityLandData) DocumentType() hint.Type {
	return doc.info.docType
}

func (doc CityLandData) Bytes() []byte {
	bs := make([][]byte, 5)

	bs[0] = doc.info.Bytes()
	bs[1] = doc.owner.Bytes()
	bs[2] = doc.lender.Bytes()
	bs[3] = []byte(doc.starttime)
	bs[4] = util.UintToBytes(doc.periodday)
	return util.ConcatBytesSlice(bs...)
}

func (doc CityLandData) Hash() valuehash.Hash {
	return doc.GenerateHash()
}

func (doc CityLandData) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(doc.Bytes())
}

func (doc CityLandData) IsValid([]byte) error {
	return nil
}

func (doc CityLandData) Info() DocInfo {
	return doc.info
}

func (doc CityLandData) Owner() base.Address {
	return doc.owner
}

func (doc CityLandData) Lender() base.Address {
	return doc.lender
}

func (doc CityLandData) Starttime() string {
	return doc.starttime
}

func (doc CityLandData) Periodday() uint {
	return doc.periodday
}

func (doc CityLandData) Equal(b CityLandData) bool {

	if !doc.info.Equal(b.info) {
		return false
	}

	if doc.starttime != b.starttime {
		return false
	}

	if doc.periodday != b.periodday {
		return false
	}

	return true
}

var (
	CityVotingDataType   = hint.Type("mitum-blockcity-document-voting-data")
	CityVotingDataHint   = hint.NewHint(CityVotingDataType, "v0.0.1")
	CityVotingDataHinter = CityVotingData{BaseHinter: hint.NewBaseHinter(CityVotingDataHint)}
)

type CityVotingData struct {
	hint.BaseHinter
	info       DocInfo
	owner      base.Address
	round      uint
	candidates []VotingCandidate
}

func NewCityVotingData(info DocInfo, owner base.Address, round uint, candidate []VotingCandidate) CityVotingData {
	doc := CityVotingData{
		BaseHinter: hint.NewBaseHinter(CityVotingDataHint),
		info:       info,
		owner:      owner,
		round:      round,
		candidates: candidate,
	}
	return doc
}

func MustNewCityVotingData(info DocInfo, owner base.Address, round uint, candidate []VotingCandidate) CityVotingData {
	doc := NewCityVotingData(info, owner, round, candidate)
	if err := doc.IsValid(nil); err != nil {
		panic(err)
	}
	return doc
}

func (doc CityVotingData) DocumentId() string {
	return doc.info.DocumentId()
}

func (doc CityVotingData) DocumentType() hint.Type {
	return doc.info.docType
}

func (doc CityVotingData) Bytes() []byte {
	bs := make([][]byte, len(doc.candidates)+3)

	sort.Slice(doc.candidates, func(i, j int) bool {
		return bytes.Compare(doc.candidates[i].Bytes(), doc.candidates[j].Bytes()) < 0
	})

	bs[0] = doc.info.Bytes()
	bs[1] = doc.owner.Bytes()
	bs[2] = util.UintToBytes(doc.round)
	for i := range doc.candidates {
		bs[i+3] = doc.candidates[i].Bytes()
	}
	return util.ConcatBytesSlice(bs...)
}

func (doc CityVotingData) Hash() valuehash.Hash {
	return doc.GenerateHash()
}

func (doc CityVotingData) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(doc.Bytes())
}

func (doc CityVotingData) IsValid([]byte) error {
	if err := isvalid.Check(
		nil, false,
		doc.BaseHinter,
		doc.info,
		doc.owner); err != nil {
		return errors.Wrap(err, "invalid document data")
	}

	for i := range doc.candidates {
		c := doc.candidates[i]
		if err := c.IsValid(nil); err != nil {
			return err
		}
	}

	return nil
}

func (doc CityVotingData) Info() DocInfo {
	return doc.info
}

func (doc CityVotingData) Owner() base.Address {
	return doc.owner
}

func (doc CityVotingData) Round() uint {
	return doc.round
}

func (doc CityVotingData) Candidates() []VotingCandidate {
	sort.Slice(doc.candidates, func(i, j int) bool {
		return bytes.Compare(doc.candidates[i].Bytes(), doc.candidates[j].Bytes()) < 0
	})
	return doc.candidates
}

func (doc CityVotingData) Equal(b CityVotingData) bool {

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
