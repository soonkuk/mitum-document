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
	case UserDocIdType:
		did = NewUserDocId(id)
	case LandDocIdType:
		did = NewLandDocId(id)
	case VotingDocIdType:
		did = NewVotingDocId(id)
	case HistoryDocIdType:
		did = NewHistoryDocId(id)
	default:
		did = nil
	}
	return did
}

var DocIdShortTypeSize = 3

var DocIdShortTypeMap = map[string]hint.Type{
	"cui": UserDocIdType,
	"cli": LandDocIdType,
	"cvi": VotingDocIdType,
	"chi": HistoryDocIdType,
}

var (
	UserDocIdType   = hint.Type("mitum-blockcity-user-document-id")
	UserDocIdHint   = hint.NewHint(UserDocIdType, "v0.0.1")
	UserDocIdHinter = UserDocId{BaseHinter: hint.NewBaseHinter(UserDocIdHint)}
)

type UserDocId struct {
	hint.BaseHinter
	s string
}

func NewUserDocId(id string) UserDocId {
	return NewUserDocIdWithHint(UserDocIdHint, id)
}

func NewUserDocIdWithHint(ht hint.Hint, id string) UserDocId {

	return UserDocId{BaseHinter: hint.NewBaseHinter(ht), s: id}
}

func MustNewUserDocId(id string) UserDocId {
	uid := NewUserDocId(id)
	if err := uid.IsValid(nil); err != nil {
		panic(err)
	}

	return uid
}

func (ui UserDocId) IsValid([]byte) error {
	if _, _, err := ParseDocId(string(ui.s)); err != nil {
		return err
	}
	return nil
}

func (ui UserDocId) String() string {
	return ui.s
}

func (ui UserDocId) Hint() hint.Hint {
	return ui.BaseHinter.Hint()
}

func (ui UserDocId) Bytes() []byte {
	return []byte(ui.s)
}

func (ui UserDocId) Equal(b UserDocId) bool {
	if (b == UserDocId{}) {
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
	LandDocIdType   = hint.Type("mitum-blockcity-land-document-id")
	LandDocIdHint   = hint.NewHint(LandDocIdType, "v0.0.1")
	LandDocIdHinter = LandDocId{BaseHinter: hint.NewBaseHinter(LandDocIdHint)}
)

type LandDocId struct {
	hint.BaseHinter
	s string
}

func NewLandDocId(id string) LandDocId {
	return NewLandDocIdWithHint(LandDocIdHint, id)
}

func NewLandDocIdWithHint(ht hint.Hint, id string) LandDocId {

	return LandDocId{BaseHinter: hint.NewBaseHinter(ht), s: id}
}

func MustNewLandDocId(id string) LandDocId {
	uid := NewLandDocId(id)
	if err := uid.IsValid(nil); err != nil {
		panic(err)
	}

	return uid
}

func (ui LandDocId) IsValid([]byte) error {
	if _, _, err := ParseDocId(ui.s); err != nil {
		return err
	}
	return nil
}

func (ui LandDocId) String() string {
	return ui.s
}

func (ui LandDocId) Hint() hint.Hint {
	return ui.BaseHinter.Hint()
}

func (ui LandDocId) Bytes() []byte {
	return []byte(ui.s)
}

func (ui LandDocId) Equal(b LandDocId) bool {
	if (b == LandDocId{}) {
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
	VotingDocIdType   = hint.Type("mitum-blockcity-voting-document-id")
	VotingDocIdHint   = hint.NewHint(VotingDocIdType, "v0.0.1")
	VotingDocIdHinter = VotingDocId{BaseHinter: hint.NewBaseHinter(VotingDocIdHint)}
)

type VotingDocId struct {
	hint.BaseHinter
	s string
}

func NewVotingDocId(id string) VotingDocId {
	return NewVotingDocIdWithHint(VotingDocIdHint, id)
}

func NewVotingDocIdWithHint(ht hint.Hint, id string) VotingDocId {

	return VotingDocId{BaseHinter: hint.NewBaseHinter(ht), s: id}
}

func MustNewVotingDocId(id string) VotingDocId {
	uid := NewVotingDocId(id)
	if err := uid.IsValid(nil); err != nil {
		panic(err)
	}

	return uid
}

func (ui VotingDocId) IsValid([]byte) error {
	if _, _, err := ParseDocId(ui.s); err != nil {
		return err
	}
	return nil
}

func (ui VotingDocId) String() string {
	return ui.s
}

func (ui VotingDocId) Hint() hint.Hint {
	return ui.BaseHinter.Hint()
}

func (ui VotingDocId) Bytes() []byte {
	return []byte(ui.s)
}

func (ui VotingDocId) Equal(b VotingDocId) bool {
	if (b == VotingDocId{}) {
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
	HistoryDocIdType   = hint.Type("mitum-blockcity-history-document-id")
	HistoryDocIdHint   = hint.NewHint(HistoryDocIdType, "v0.0.1")
	HistoryDocIdHinter = HistoryDocId{BaseHinter: hint.NewBaseHinter(HistoryDocIdHint)}
)

type HistoryDocId struct {
	hint.BaseHinter
	s string
}

func NewHistoryDocId(id string) HistoryDocId {
	return NewHistoryDocIdWithHint(HistoryDocIdHint, id)
}

func NewHistoryDocIdWithHint(ht hint.Hint, id string) HistoryDocId {

	return HistoryDocId{BaseHinter: hint.NewBaseHinter(ht), s: id}
}

func MustNewHistoryDocId(id string) HistoryDocId {
	uid := NewHistoryDocId(id)
	if err := uid.IsValid(nil); err != nil {
		panic(err)
	}

	return uid
}

func (ui HistoryDocId) IsValid([]byte) error {
	if _, _, err := ParseDocId(ui.s); err != nil {
		return err
	}
	return nil
}

func (ui HistoryDocId) String() string {
	return ui.s
}

func (ui HistoryDocId) Hint() hint.Hint {
	return ui.BaseHinter.Hint()
}

func (ui HistoryDocId) Bytes() []byte {
	return []byte(ui.s)
}

func (ui HistoryDocId) Equal(b HistoryDocId) bool {
	if (b == HistoryDocId{}) {
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
