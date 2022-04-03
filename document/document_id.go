package document // nolint: dupl, revive

import (
	"regexp"

	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
)

type DocID interface {
	String() string
	Bytes() []byte
	Hint() hint.Hint
	IsValid([]byte) error
}

func NewDocID(id string) (did DocID) {
	_, ty, err := ParseDocID(id)
	if err != nil {
		did = nil
		return did
	}

	switch ty {
	case BSDocIDType:
		did = NewBSDocID(id)
	case BCUserDocIDType:
		did = NewBCUserDocID(id)
	case BCLandDocIDType:
		did = NewBCLandDocID(id)
	case BCVotingDocIDType:
		did = NewBCVotingDocID(id)
	case BCHistoryDocIDType:
		did = NewBCHistoryDocID(id)
	default:
		did = nil
	}
	return did
}

var (
	MaxLengthDocID = 10
	ReValidDocID   = regexp.MustCompile(`[a-z0-9]+`)
)

var DocIDShortTypeSize = 3

var DocIDShortTypeMap = map[string]hint.Type{
	"sdi": BSDocIDType,
	"cui": BCUserDocIDType,
	"cli": BCLandDocIDType,
	"cvi": BCVotingDocIDType,
	"chi": BCHistoryDocIDType,
}

var (
	BSDocIDType   = hint.Type("mitum-document-id")
	BSDocIDHint   = hint.NewHint(BSDocIDType, "v0.0.1")
	BSDocIDHinter = BSDocID{BaseHinter: hint.NewBaseHinter(BSDocIDHint)}
)

type BSDocID struct {
	hint.BaseHinter
	s string
}

func NewBSDocID(id string) BSDocID {
	return NewBSDocIDWithHint(BSDocIDHint, id)
}

func NewBSDocIDWithHint(ht hint.Hint, id string) BSDocID {
	return BSDocID{BaseHinter: hint.NewBaseHinter(ht), s: id}
}

func MustNewBSDocID(id string) BSDocID {
	uid := NewBSDocID(id)
	if err := uid.IsValid(nil); err != nil {
		panic(err)
	}

	return uid
}

func (ui BSDocID) IsValid([]byte) error {
	s, _, err := ParseDocID(ui.s)
	if err != nil {
		return err
	} else if !ReValidDocID.MatchString(s) {
		return isvalid.InvalidError.Errorf("wrong doc id, %q", s)
	} else if l := len(s); l > MaxLengthDocID {
		return isvalid.InvalidError.Errorf(
			"invalid length of document id, %d <= %d", l, MaxLengthDocID)
	}
	return nil
}

func (ui BSDocID) String() string {
	return ui.s
}

func (ui BSDocID) Hint() hint.Hint {
	return ui.BaseHinter.Hint()
}

func (ui BSDocID) Bytes() []byte { // nolint:stylecheck
	return []byte(ui.s)
}

func (ui BSDocID) Equal(b BCUserDocID) bool {
	if (b == BCUserDocID{}) {
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
	BCUserDocIDType   = hint.Type("mitum-blockcity-user-document-id")
	BCUserDocIDHint   = hint.NewHint(BCUserDocIDType, "v0.0.1")
	BCUserDocIDHinter = BCUserDocID{BaseHinter: hint.NewBaseHinter(BCUserDocIDHint)}
)

type BCUserDocID struct {
	hint.BaseHinter
	s string
}

func NewBCUserDocID(id string) BCUserDocID {
	return NewBCUserDocIDWithHint(BCUserDocIDHint, id)
}

func NewBCUserDocIDWithHint(ht hint.Hint, id string) BCUserDocID {
	return BCUserDocID{BaseHinter: hint.NewBaseHinter(ht), s: id}
}

func MustNewBCUserDocID(id string) BCUserDocID {
	uid := NewBCUserDocID(id)
	if err := uid.IsValid(nil); err != nil {
		panic(err)
	}

	return uid
}

func (ui BCUserDocID) IsValid([]byte) error {
	s, _, err := ParseDocID(ui.s)
	if err != nil {
		return err
	} else if !ReValidDocID.MatchString(s) {
		return isvalid.InvalidError.Errorf("wrong doc id, %q", s)
	} else if l := len(s); l > MaxLengthDocID {
		return isvalid.InvalidError.Errorf(
			"invalid length of document id, %d <= %d", l, MaxLengthDocID)
	}
	return nil
}

func (ui BCUserDocID) String() string {
	return ui.s
}

func (ui BCUserDocID) Hint() hint.Hint {
	return ui.BaseHinter.Hint()
}

func (ui BCUserDocID) Bytes() []byte { // nolint:stylecheck
	return []byte(ui.s)
}

func (ui BCUserDocID) Equal(b BCUserDocID) bool {
	if (b == BCUserDocID{}) {
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
	BCLandDocIDType   = hint.Type("mitum-blockcity-land-document-id")
	BCLandDocIDHint   = hint.NewHint(BCLandDocIDType, "v0.0.1")
	BCLandDocIDHinter = BCLandDocID{BaseHinter: hint.NewBaseHinter(BCLandDocIDHint)}
)

type BCLandDocID struct {
	hint.BaseHinter
	s string
}

func NewBCLandDocID(id string) BCLandDocID {
	return NewBCLandDocIDWithHint(BCLandDocIDHint, id)
}

func NewBCLandDocIDWithHint(ht hint.Hint, id string) BCLandDocID {
	return BCLandDocID{BaseHinter: hint.NewBaseHinter(ht), s: id}
}

func MustNewBCLandDocID(id string) BCLandDocID {
	uid := NewBCLandDocID(id)
	if err := uid.IsValid(nil); err != nil {
		panic(err)
	}

	return uid
}

func (ui BCLandDocID) IsValid([]byte) error {
	s, _, err := ParseDocID(ui.s)
	if err != nil {
		return err
	} else if !ReValidDocID.MatchString(s) {
		return isvalid.InvalidError.Errorf("wrong doc id, %q", s)
	} else if l := len(s); l > MaxLengthDocID {
		return isvalid.InvalidError.Errorf(
			"invalid length of document id, %d <= %d", l, MaxLengthDocID)
	}
	return nil
}

func (ui BCLandDocID) String() string {
	return ui.s
}

func (ui BCLandDocID) Hint() hint.Hint {
	return ui.BaseHinter.Hint()
}

func (ui BCLandDocID) Bytes() []byte { // nolint:stylecheck
	return []byte(ui.s)
}

func (ui BCLandDocID) Equal(b BCLandDocID) bool {
	if (b == BCLandDocID{}) {
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
	BCVotingDocIDType   = hint.Type("mitum-blockcity-voting-document-id")
	BCVotingDocIDHint   = hint.NewHint(BCVotingDocIDType, "v0.0.1")
	BCVotingDocIDHinter = BCVotingDocID{BaseHinter: hint.NewBaseHinter(BCVotingDocIDHint)}
)

type BCVotingDocID struct {
	hint.BaseHinter
	s string
}

func NewBCVotingDocID(id string) BCVotingDocID {
	return NewBCVotingDocIDWithHint(BCVotingDocIDHint, id)
}

func NewBCVotingDocIDWithHint(ht hint.Hint, id string) BCVotingDocID {
	return BCVotingDocID{BaseHinter: hint.NewBaseHinter(ht), s: id}
}

func MustNewBCVotingDocID(id string) BCVotingDocID {
	uid := NewBCVotingDocID(id)
	if err := uid.IsValid(nil); err != nil {
		panic(err)
	}

	return uid
}

func (ui BCVotingDocID) IsValid([]byte) error {
	s, _, err := ParseDocID(ui.s)
	if err != nil {
		return err
	} else if !ReValidDocID.MatchString(s) {
		return isvalid.InvalidError.Errorf("wrong doc id, %q", s)
	} else if l := len(s); l > MaxLengthDocID {
		return isvalid.InvalidError.Errorf(
			"invalid length of document id, %d <= %d", l, MaxLengthDocID)
	}
	return nil
}

func (ui BCVotingDocID) String() string {
	return ui.s
}

func (ui BCVotingDocID) Hint() hint.Hint {
	return ui.BaseHinter.Hint()
}

func (ui BCVotingDocID) Bytes() []byte { // nolint:stylecheck
	return []byte(ui.s)
}

func (ui BCVotingDocID) Equal(b BCVotingDocID) bool {
	if (b == BCVotingDocID{}) {
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
	BCHistoryDocIDType   = hint.Type("mitum-blockcity-history-document-id")
	BCHistoryDocIDHint   = hint.NewHint(BCHistoryDocIDType, "v0.0.1")
	BCHistoryDocIDHinter = BCHistoryDocID{BaseHinter: hint.NewBaseHinter(BCHistoryDocIDHint)}
)

type BCHistoryDocID struct {
	hint.BaseHinter
	s string
}

func NewBCHistoryDocID(id string) BCHistoryDocID {
	return NewBCHistoryDocIDWithHint(BCHistoryDocIDHint, id)
}

func NewBCHistoryDocIDWithHint(ht hint.Hint, id string) BCHistoryDocID {
	return BCHistoryDocID{BaseHinter: hint.NewBaseHinter(ht), s: id}
}

func MustNewBCHistoryDocID(id string) BCHistoryDocID {
	uid := NewBCHistoryDocID(id)
	if err := uid.IsValid(nil); err != nil {
		panic(err)
	}

	return uid
}

func (ui BCHistoryDocID) IsValid([]byte) error {
	s, _, err := ParseDocID(ui.s)
	if err != nil {
		return err
	} else if !ReValidDocID.MatchString(s) {
		return isvalid.InvalidError.Errorf("wrong doc id, %q", s)
	} else if l := len(s); l > MaxLengthDocID {
		return isvalid.InvalidError.Errorf(
			"invalid length of document id, %d <= %d", l, MaxLengthDocID)
	}
	return nil
}

func (ui BCHistoryDocID) String() string {
	return ui.s
}

func (ui BCHistoryDocID) Hint() hint.Hint {
	return ui.BaseHinter.Hint()
}

func (ui BCHistoryDocID) Bytes() []byte { // nolint:stylecheck
	return []byte(ui.s)
}

func (ui BCHistoryDocID) Equal(b BCHistoryDocID) bool {
	if (b == BCHistoryDocID{}) {
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

// ParseDocID receive typed docID string and return untyped docID string and docID type
func ParseDocID(s string) (string, hint.Type, error) {
	if len(s) <= DocIDShortTypeSize {
		return "", hint.Type(""), isvalid.InvalidError.Errorf("invalid DocID, %q", s)
	}

	shortType := s[len(s)-DocIDShortTypeSize:]

	if len(shortType) != DocIDShortTypeSize {
		return "", hint.Type(""), isvalid.InvalidError.Errorf("invalid ShortType for DocID, %q", shortType)
	}

	v, ok := DocIDShortTypeMap[shortType]
	if !ok {
		return "", hint.Type(""), isvalid.InvalidError.Errorf("wrong ShortType of DocID : %q", shortType)
	}

	return s[:len(s)-DocIDShortTypeSize], v, nil
}
