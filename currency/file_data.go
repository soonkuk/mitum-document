package currency

import (
	"fmt"

	"golang.org/x/xerrors"

	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
	"github.com/spikeekips/mitum/util/valuehash"
)

var (
	FileDataType  = hint.Type("mitum-blocksign-file-data")
	FileDatatHint = hint.NewHint(FileDataType, "v0.0.1")
)

type FileData struct {
	// fid   FileID
	signcode SignCode
	owner    base.Address
}

func NewEmptyFileData() FileData {
	fd := FileData{signcode: "", owner: EmptyAddress}

	return fd
}

func NewFileData(signcode SignCode, owner base.Address) FileData {
	fd := FileData{signcode: signcode, owner: owner}

	return fd
}

func MustNewFileData(signcode SignCode, owner base.Address) FileData {
	fd := NewFileData(signcode, owner)
	if err := fd.IsValid(nil); err != nil {
		panic(err)
	}

	return fd
}

func (fd FileData) Hint() hint.Hint {
	return FileDatatHint
}

func (fd FileData) Bytes() []byte {
	return util.ConcatBytesSlice(
		fd.signcode.Bytes(),
		fd.owner.Bytes(),
	)
}

func (fd FileData) Hash() valuehash.Hash {
	return fd.GenerateHash()
}

func (fd FileData) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(fd.Bytes())
}

func (fd FileData) IsEmpty() bool {
	return len(fd.signcode) < 1 || fd.owner == EmptyAddress
}

func (fd FileData) IsValid([]byte) error {
	if err := isvalid.Check([]isvalid.IsValider{
		fd.signcode,
		fd.owner,
	}, nil, false); err != nil {
		return xerrors.Errorf("invalid file data: %w", err)
	}

	return nil
}

func (fd FileData) SignCode() SignCode {
	return fd.signcode
}

func (fd FileData) Owner() base.Address {
	return fd.owner
}

func (fd FileData) String() string {
	return fmt.Sprintf("%s:%s", fd.signcode.String(), fd.owner.String())
}

func (fd FileData) Equal(b FileData) bool {
	switch {
	case fd.signcode != b.signcode:
		return false
	case fd.owner != b.owner:
		return false
	default:
		return true
	}
}

func (fd FileData) WithData(signcode SignCode, owner base.Address) FileData {
	fd.signcode = signcode
	fd.owner = owner
	return fd
}

var (
	FileIDType = hint.Type("mitum-blocksign-file-id")
	FileIDHint = hint.NewHint(FileIDType, "v0.0.1")
)

type FileID string

func (fid FileID) Bytes() []byte {
	return []byte(fid)
}

func (fid FileID) String() string {
	return string(fid)
}

func (fid FileID) Hint() hint.Hint {
	return FileIDHint
}

func (fid FileID) Hash() valuehash.Hash {
	return fid.GenerateHash()
}

func (fid FileID) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(fid.Bytes())
}

func (fid FileID) IsValid([]byte) error {
	return nil
}

func (fid FileID) Equal(b FileID) bool {
	return fid == b
}

var (
	SignCodeType = hint.Type("mitum-blocksign-owner-signcode")
	SignCodeHint = hint.NewHint(SignCodeType, "v0.0.1")
)

type SignCode string

func (sc SignCode) Bytes() []byte {
	return []byte(sc)
}

func (sc SignCode) String() string {
	return string(sc)
}

func (sc SignCode) Hint() hint.Hint {
	return SignCodeHint
}

func (sc SignCode) Hash() valuehash.Hash {
	return sc.GenerateHash()
}

func (sc SignCode) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(sc.Bytes())
}

func (sc SignCode) IsValid([]byte) error {
	return nil
}

func (sc SignCode) Equal(b SignCode) bool {
	return sc == b
}
