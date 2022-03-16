package document

import (
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
)

type DocId interface {
	String() string
	Bytes() []byte
	Hint() hint.Hint
}

func NewDocId(id string) (did DocId) {
	_, ty, err := ParseDocId(id)
	if err != nil {
		did = nil
		return did
	}

	switch ty {
	case BSDocIdType:
		did = NewBSDocId(id)
	case BCUserDocIdType:
		did = NewBCUserDocId(id)
	case BCLandDocIdType:
		did = NewBCLandDocId(id)
	case BCVotingDocIdType:
		did = NewBCVotingDocId(id)
	case BCHistoryDocIdType:
		did = NewBCHistoryDocId(id)
	default:
		did = nil
	}
	return did
}

var DocIdShortTypeSize = 3

var DocIdShortTypeMap = map[string]hint.Type{
	"sdi": BSDocIdType,
	"cui": BCUserDocIdType,
	"cli": BCLandDocIdType,
	"cvi": BCVotingDocIdType,
	"chi": BCHistoryDocIdType,
}

var (
	BSDocIdType   = hint.Type("mitum-document-id")
	BSDocIdHint   = hint.NewHint(BSDocIdType, "v0.0.1")
	BSDocIdHinter = BSDocId{BaseHinter: hint.NewBaseHinter(BSDocIdHint)}
)

type BSDocId struct {
	hint.BaseHinter
	s string
}

func NewBSDocId(id string) BSDocId {
	return NewBSDocIdWithHint(BSDocIdHint, id)
}

func NewBSDocIdWithHint(ht hint.Hint, id string) BSDocId {

	return BSDocId{BaseHinter: hint.NewBaseHinter(ht), s: id}
}

func MustNewBSDocId(id string) BSDocId {
	uid := NewBSDocId(id)
	if err := uid.IsValid(nil); err != nil {
		panic(err)
	}

	return uid
}

func (ui BSDocId) IsValid([]byte) error {
	if _, _, err := ParseDocId(string(ui.s)); err != nil {
		return err
	}
	return nil
}

func (ui BSDocId) String() string {
	return ui.s
}

func (ui BSDocId) Hint() hint.Hint {
	return ui.BaseHinter.Hint()
}

func (ui BSDocId) Bytes() []byte {
	return []byte(ui.s)
}

func (ui BSDocId) Equal(b BCUserDocId) bool {
	if (b == BCUserDocId{}) {
		return false
	}

	if ui.Hint().Type() != b.Hint().Type() {
		return false
	}

	if err := b.IsValid(nil); err != nil {
		return false
	}

	return string(ui.s) == b.String()
}

var (
	BCUserDocIdType   = hint.Type("mitum-blockcity-user-document-id")
	BCUserDocIdHint   = hint.NewHint(BCUserDocIdType, "v0.0.1")
	BCUserDocIdHinter = BCUserDocId{BaseHinter: hint.NewBaseHinter(BCUserDocIdHint)}
)

type BCUserDocId struct {
	hint.BaseHinter
	s string
}

func NewBCUserDocId(id string) BCUserDocId {
	return NewBCUserDocIdWithHint(BCUserDocIdHint, id)
}

func NewBCUserDocIdWithHint(ht hint.Hint, id string) BCUserDocId {

	return BCUserDocId{BaseHinter: hint.NewBaseHinter(ht), s: id}
}

func MustNewBCUserDocId(id string) BCUserDocId {
	uid := NewBCUserDocId(id)
	if err := uid.IsValid(nil); err != nil {
		panic(err)
	}

	return uid
}

func (ui BCUserDocId) IsValid([]byte) error {
	if _, _, err := ParseDocId(string(ui.s)); err != nil {
		return err
	}
	return nil
}

func (ui BCUserDocId) String() string {
	return ui.s
}

func (ui BCUserDocId) Hint() hint.Hint {
	return ui.BaseHinter.Hint()
}

func (ui BCUserDocId) Bytes() []byte {
	return []byte(ui.s)
}

func (ui BCUserDocId) Equal(b BCUserDocId) bool {
	if (b == BCUserDocId{}) {
		return false
	}

	if ui.Hint().Type() != b.Hint().Type() {
		return false
	}

	if err := b.IsValid(nil); err != nil {
		return false
	}

	return string(ui.s) == b.String()
}

var (
	BCLandDocIdType   = hint.Type("mitum-blockcity-land-document-id")
	BCLandDocIdHint   = hint.NewHint(BCLandDocIdType, "v0.0.1")
	BCLandDocIdHinter = BCLandDocId{BaseHinter: hint.NewBaseHinter(BCLandDocIdHint)}
)

type BCLandDocId struct {
	hint.BaseHinter
	s string
}

func NewBCLandDocId(id string) BCLandDocId {
	return NewBCLandDocIdWithHint(BCLandDocIdHint, id)
}

func NewBCLandDocIdWithHint(ht hint.Hint, id string) BCLandDocId {

	return BCLandDocId{BaseHinter: hint.NewBaseHinter(ht), s: id}
}

func MustNewBCLandDocId(id string) BCLandDocId {
	uid := NewBCLandDocId(id)
	if err := uid.IsValid(nil); err != nil {
		panic(err)
	}

	return uid
}

func (ui BCLandDocId) IsValid([]byte) error {
	if _, _, err := ParseDocId(ui.s); err != nil {
		return err
	}
	return nil
}

func (ui BCLandDocId) String() string {
	return ui.s
}

func (ui BCLandDocId) Hint() hint.Hint {
	return ui.BaseHinter.Hint()
}

func (ui BCLandDocId) Bytes() []byte {
	return []byte(ui.s)
}

func (ui BCLandDocId) Equal(b BCLandDocId) bool {
	if (b == BCLandDocId{}) {
		return false
	}

	if ui.Hint().Type() != b.Hint().Type() {
		return false
	}

	if err := b.IsValid(nil); err != nil {
		return false
	}

	return ui.s == b.String()
}

var (
	BCVotingDocIdType   = hint.Type("mitum-blockcity-voting-document-id")
	BCVotingDocIdHint   = hint.NewHint(BCVotingDocIdType, "v0.0.1")
	BCVotingDocIdHinter = BCVotingDocId{BaseHinter: hint.NewBaseHinter(BCVotingDocIdHint)}
)

type BCVotingDocId struct {
	hint.BaseHinter
	s string
}

func NewBCVotingDocId(id string) BCVotingDocId {
	return NewBCVotingDocIdWithHint(BCVotingDocIdHint, id)
}

func NewBCVotingDocIdWithHint(ht hint.Hint, id string) BCVotingDocId {

	return BCVotingDocId{BaseHinter: hint.NewBaseHinter(ht), s: id}
}

func MustNewBCVotingDocId(id string) BCVotingDocId {
	uid := NewBCVotingDocId(id)
	if err := uid.IsValid(nil); err != nil {
		panic(err)
	}

	return uid
}

func (ui BCVotingDocId) IsValid([]byte) error {
	if _, _, err := ParseDocId(ui.s); err != nil {
		return err
	}
	return nil
}

func (ui BCVotingDocId) String() string {
	return ui.s
}

func (ui BCVotingDocId) Hint() hint.Hint {
	return ui.BaseHinter.Hint()
}

func (ui BCVotingDocId) Bytes() []byte {
	return []byte(ui.s)
}

func (ui BCVotingDocId) Equal(b BCVotingDocId) bool {
	if (b == BCVotingDocId{}) {
		return false
	}

	if ui.Hint().Type() != b.Hint().Type() {
		return false
	}

	if err := b.IsValid(nil); err != nil {
		return false
	}

	return ui.s == b.String()
}

var (
	BCHistoryDocIdType   = hint.Type("mitum-blockcity-history-document-id")
	BCHistoryDocIdHint   = hint.NewHint(BCHistoryDocIdType, "v0.0.1")
	BCHistoryDocIdHinter = BCHistoryDocId{BaseHinter: hint.NewBaseHinter(BCHistoryDocIdHint)}
)

type BCHistoryDocId struct {
	hint.BaseHinter
	s string
}

func NewBCHistoryDocId(id string) BCHistoryDocId {
	return NewBCHistoryDocIdWithHint(BCHistoryDocIdHint, id)
}

func NewBCHistoryDocIdWithHint(ht hint.Hint, id string) BCHistoryDocId {

	return BCHistoryDocId{BaseHinter: hint.NewBaseHinter(ht), s: id}
}

func MustNewBCHistoryDocId(id string) BCHistoryDocId {
	uid := NewBCHistoryDocId(id)
	if err := uid.IsValid(nil); err != nil {
		panic(err)
	}

	return uid
}

func (ui BCHistoryDocId) IsValid([]byte) error {
	if _, _, err := ParseDocId(ui.s); err != nil {
		return err
	}
	return nil
}

func (ui BCHistoryDocId) String() string {
	return ui.s
}

func (ui BCHistoryDocId) Hint() hint.Hint {
	return ui.BaseHinter.Hint()
}

func (ui BCHistoryDocId) Bytes() []byte {
	return []byte(ui.s)
}

func (ui BCHistoryDocId) Equal(b BCHistoryDocId) bool {
	if (b == BCHistoryDocId{}) {
		return false
	}

	if ui.Hint().Type() != b.Hint().Type() {
		return false
	}

	if err := b.IsValid(nil); err != nil {
		return false
	}

	return ui.s == b.String()
}

func ParseDocId(s string) (string, hint.Type, error) {
	if len(s) <= DocIdShortTypeSize {
		return "", hint.Type(""), isvalid.InvalidError.Errorf("invalid DocId, %q", s)
	}

	shortType := s[len(s)-DocIdShortTypeSize:]

	if len(shortType) != DocIdShortTypeSize {
		return "", hint.Type(""), isvalid.InvalidError.Errorf("invalid ShortType for DocId, %q", shortType)
	}

	v, ok := DocIdShortTypeMap[shortType]
	if !ok {
		return "", hint.Type(""), isvalid.InvalidError.Errorf("wrong ShortType of DocId : %q", shortType)
	}

	return s[:len(s)-DocIdShortTypeSize], v, nil
}
