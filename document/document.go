package document // nolint: dupl, revive

import (
	"bytes"
	"fmt"

	"github.com/pkg/errors"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
	"github.com/spikeekips/mitum/util/valuehash"
)

var (
	UserStatisticsType   = hint.Type("mitum-blockcity-user-statistics")
	UserStatisticsHint   = hint.NewHint(UserStatisticsType, "v0.0.1")
	UserStatisticsHinter = UserStatistics{BaseHinter: hint.NewBaseHinter(UserStatisticsHint)}
)

type UserStatistics struct {
	hint.BaseHinter
	hp           uint
	strength     uint
	agility      uint
	dexterity    uint
	charisma     uint
	intelligence uint
	vital        uint
}

func NewUserStatistics(hp, strength, agility, dexterity, charisma, intelligence, vital uint) UserStatistics {
	us := UserStatistics{
		BaseHinter:   hint.NewBaseHinter(UserStatisticsHint),
		hp:           hp,
		strength:     strength,
		agility:      agility,
		dexterity:    dexterity,
		charisma:     charisma,
		intelligence: intelligence,
		vital:        vital,
	}
	return us
}

func MustNewUserStatistics(hp, strength, agility, dexterity, charisma, intelligence, vital uint) UserStatistics {
	us := NewUserStatistics(hp, strength, agility, dexterity, charisma, intelligence, vital)
	if err := us.IsValid(nil); err != nil {
		panic(err)
	}
	return us
}

func (us UserStatistics) Bytes() []byte {
	bs := make([][]byte, 7)

	bs[0] = util.UintToBytes(us.hp)
	bs[1] = util.UintToBytes(us.strength)
	bs[2] = util.UintToBytes(us.agility)
	bs[3] = util.UintToBytes(us.dexterity)
	bs[4] = util.UintToBytes(us.charisma)
	bs[5] = util.UintToBytes(us.intelligence)
	bs[6] = util.UintToBytes(us.vital)
	return util.ConcatBytesSlice(bs...)
}

func (us UserStatistics) Hash() valuehash.Hash {
	return us.GenerateHash()
}

func (us UserStatistics) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(us.Bytes())
}

func (us UserStatistics) IsValid([]byte) error { // nolint:revive
	return nil
}

func (us UserStatistics) Equal(b UserStatistics) bool {
	if us.hp != b.hp {
		return false
	}

	if us.strength != b.strength {
		return false
	}

	if us.agility != b.agility {
		return false
	}

	if us.dexterity != b.dexterity {
		return false
	}

	if us.charisma != b.charisma {
		return false
	}

	if us.intelligence != b.intelligence {
		return false
	}

	if us.vital != b.vital {
		return false
	}

	return true
}

var (
	DocInfoType   = hint.Type("mitum-document-info")
	DocInfoHint   = hint.NewHint(DocInfoType, "v0.0.1")
	DocInfoHinter = DocInfo{BaseHinter: hint.NewBaseHinter(DocInfoHint)}
)

type DocInfo struct {
	hint.BaseHinter
	id      DocID
	docType hint.Type
}

func NewDocInfo(id string, docType hint.Type) DocInfo {
	var i DocID
	switch docType {
	case BSDocDataType:
		i = NewBSDocID(id)
	case BCUserDataType:
		i = NewBCUserDocID(id)
	case BCLandDataType:
		i = NewBCLandDocID(id)
	case BCVotingDataType:
		i = NewBCVotingDocID(id)
	case BCHistoryDataType:
		i = NewBCHistoryDocID(id)
	default:
		return DocInfo{}
	}

	docInfo := DocInfo{
		BaseHinter: hint.NewBaseHinter(DocInfoHint),
		id:         i,
		docType:    docType,
	}
	return docInfo
}

func MustNewDocInfo(id string, docType hint.Type) DocInfo {
	docInfo := NewDocInfo(id, docType)
	if err := docInfo.IsValid(nil); err != nil {
		panic(err)
	}
	return docInfo
}

func (di DocInfo) DocumentID() string {
	return di.id.String()
}

func (di DocInfo) DocType() hint.Type {
	return di.docType
}

func (di DocInfo) Bytes() []byte {
	return util.ConcatBytesSlice(di.id.Bytes(), di.docType.Bytes())
}

func (di DocInfo) Hash() valuehash.Hash {
	return di.GenerateHash()
}

func (di DocInfo) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(di.Bytes())
}

func (di DocInfo) IsValid([]byte) error {
	if di.id == nil {
		return isvalid.InvalidError.Errorf("DocID in Docinfo is empty")
	}
	if di.docType == hint.Type("") {
		return isvalid.InvalidError.Errorf("DocType in Docinfo is empty")
	}

	if err := isvalid.Check(nil, false,
		di.id,
		di.BaseHinter,
		di.docType,
	); err != nil {
		return isvalid.InvalidError.Errorf("invalid Docinfo: %w", err)
	}
	return nil
}

func (di DocInfo) String() string {
	return fmt.Sprintf("%s:%s", di.id.String(), di.docType.String())
}

func (di DocInfo) Equal(b DocInfo) bool {
	return bytes.Equal(di.id.Bytes(), b.id.Bytes()) && di.docType == b.docType
}

type FileHash string

func (fh FileHash) Bytes() []byte {
	return []byte(fh)
}

func (fh FileHash) String() string {
	return string(fh)
}

func (fh FileHash) IsValid([]byte) error {
	if len(fh) < 1 {
		return errors.Errorf("empty fileHash")
	}
	return nil
}

func (fh FileHash) Equal(b FileHash) bool {
	return fh == b
}

var (
	DocSignType   = hint.Type("mitum-blocksign-docsign")
	DocSignHint   = hint.NewHint(DocSignType, "v0.0.1")
	DocSignHinter = DocSign{BaseHinter: hint.NewBaseHinter(DocSignHint)}
)

type DocSign struct {
	hint.BaseHinter
	address  base.Address
	signcode string
	signed   bool
}

func NewDocSign(address base.Address, signcode string, signed bool) DocSign {
	doc := DocSign{
		BaseHinter: hint.NewBaseHinter(DocSignHint),
		address:    address,
		signcode:   signcode,
		signed:     signed,
	}
	return doc
}

func MustNewDocSign(address base.Address, signcode string, signed bool) DocSign {
	doc := NewDocSign(address, signcode, signed)
	if err := doc.IsValid(nil); err != nil {
		panic(err)
	}
	return doc
}

func (ds DocSign) Address() base.Address {
	return ds.address
}

func (ds DocSign) Bytes() []byte {
	bs := make([][]byte, 3)
	bs[0] = ds.address.Bytes()
	var v int8
	if ds.signed {
		v = 1
	}
	bs[1] = []byte(ds.signcode)
	bs[2] = []byte{byte(v)}
	return util.ConcatBytesSlice(bs...)
}

func (ds DocSign) Hash() valuehash.Hash {
	return ds.GenerateHash()
}

func (ds DocSign) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(ds.Bytes())
}

func (ds DocSign) IsValid([]byte) error { // nolint:revive
	return nil
}

func (ds DocSign) IsEmpty() bool {
	return len(ds.address.String()) < 1
}

func (ds DocSign) String() string {
	v := fmt.Sprintf("%v", ds.signed)
	return fmt.Sprintf("%s:%s", ds.address.String(), v)
}

func (ds DocSign) Equal(b DocSign) bool {
	if !ds.address.Equal(b.address) {
		return false
	}

	if ds.signcode != b.signcode {
		return false
	}

	if ds.signed != b.signed {
		return false
	}

	return true
}

func (ds *DocSign) Signed() bool {
	return ds.signed
}

func (ds *DocSign) SetSigned() {
	ds.signed = true
}

var MaxManifest = 100

var (
	VotingCandidateType   = hint.Type("mitum-blockcity-voting-candidate")
	VotingCandidateHint   = hint.NewHint(VotingCandidateType, "v0.0.1")
	VotingCandidateHinter = VotingCandidate{BaseHinter: hint.NewBaseHinter(VotingCandidateHint)}
)

type VotingCandidate struct {
	hint.BaseHinter
	address  base.Address
	nickname string
	manifest string
	count    uint
}

func NewVotingCandidate(address base.Address, nickname, manifest string, count uint) VotingCandidate {
	votingCandidate := VotingCandidate{
		BaseHinter: hint.NewBaseHinter(VotingCandidateHint),
		address:    address,
		nickname:   nickname,
		manifest:   manifest,
		count:      count,
	}
	return votingCandidate
}

func MustNewVotingCandidate(address base.Address, nickname, manifest string, count uint) VotingCandidate {
	votingCandidate := NewVotingCandidate(address, nickname, manifest, count)
	if err := votingCandidate.IsValid(nil); err != nil {
		panic(err)
	}
	return votingCandidate
}

func (vc VotingCandidate) Address() base.Address { // nolint:stylecheck
	return vc.address
}

func (vc VotingCandidate) Bytes() []byte {
	return util.ConcatBytesSlice(vc.address.Bytes(), []byte(vc.nickname), []byte(vc.manifest), util.UintToBytes(vc.count))
}

func (vc VotingCandidate) Hash() valuehash.Hash {
	return vc.GenerateHash()
}

func (vc VotingCandidate) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(vc.Bytes())
}

func (vc VotingCandidate) IsValid([]byte) error {
	if len(vc.manifest) > MaxManifest {
		return isvalid.InvalidError.Errorf("Over candidate max manifest")
	}

	if err := isvalid.Check(nil, false,
		vc.address,
	); err != nil {
		return isvalid.InvalidError.Errorf("invalid VotingCandidate: %w", err)
	}
	return nil
}

func (vc VotingCandidate) String() string {
	return fmt.Sprintf("%s:%s:%s:%s", vc.address.String(), vc.nickname, vc.manifest, fmt.Sprint(vc.count))
}

func (vc VotingCandidate) Equal(b VotingCandidate) bool {
	return vc.address.Equal(b.address) && vc.nickname == b.nickname && vc.manifest == b.manifest && vc.count == b.count
}
