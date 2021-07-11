package blocksign

import (
	"bytes"
	"fmt"
	"sort"

	"golang.org/x/xerrors"

	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
	"github.com/spikeekips/mitum/util/valuehash"
)

var (
	DocumentDataType = hint.Type("mitum-blocksign-document-data")
	DocumentDataHint = hint.NewHint(DocumentDataType, "v0.0.1")
)

type DocumentData struct {
	fileHash FileHash
	creator  DocSign
	owner    base.Address
	signers  []DocSign
}

func NewDocumentData(fileHash FileHash, creator DocSign, owner base.Address, signers []DocSign) DocumentData {
	doc := DocumentData{
		fileHash: fileHash,
		creator:  creator,
		owner:    owner,
		signers:  signers,
	}

	return doc
}

func MustNewDocumentData(fileHash FileHash, creator DocSign, owner base.Address, signers []DocSign) DocumentData {
	doc := NewDocumentData(fileHash, creator, owner, signers)
	if err := doc.IsValid(nil); err != nil {
		panic(err)
	}

	return doc
}

func (doc DocumentData) Hint() hint.Hint {
	return DocumentDataHint
}

func (doc DocumentData) Bytes() []byte {
	bs := make([][]byte, len(doc.signers)+3)

	sort.Slice(doc.signers, func(i, j int) bool {
		return bytes.Compare(doc.signers[i].Bytes(), doc.signers[j].Bytes()) < 0
	})

	bs[0] = doc.fileHash.Bytes()
	bs[1] = doc.creator.Bytes()
	bs[2] = doc.owner.Bytes()

	for i := range doc.signers {
		bs[i+3] = doc.signers[i].Bytes()
	}

	return util.ConcatBytesSlice(bs...)
}

func (doc DocumentData) Hash() valuehash.Hash {
	return doc.GenerateHash()
}

func (doc DocumentData) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(doc.Bytes())
}

func (doc DocumentData) IsEmpty() bool {
	return len(doc.fileHash) < 1 || doc.creator.IsEmpty() || len(doc.owner.Raw()) < 1 || len(doc.signers) < 1
}

func (doc DocumentData) IsValid([]byte) error {
	if err := isvalid.Check([]isvalid.IsValider{
		doc.fileHash,
		doc.creator,
		doc.owner,
	}, nil, false); err != nil {
		return xerrors.Errorf("invalid document data: %w", err)
	}

	for i := range doc.signers {
		c := doc.signers[i]
		if err := c.IsValid(nil); err != nil {
			return err
		}
	}

	// TODO : check owner and signer are not same

	return nil
}

func (doc DocumentData) FileHash() FileHash {
	return doc.fileHash
}

func (doc DocumentData) Creator() DocSign {
	return doc.creator
}

func (doc DocumentData) Owner() base.Address {
	return doc.owner
}

func (doc DocumentData) Signers() []DocSign {
	return doc.signers
}

func (doc DocumentData) String() string {

	/*
		var signers string
		signers = "signers("
		for i := range doc.signers {
			signers = signers + "" + doc.signers[i].String()
		}
		signers = signers + ")"

		var signedBy string
		signedBy = "signedBy("
		for i := range doc.signedBy {
			signers = signers + ":" + doc.signedBy[i].String()
		}
		signedBy = signedBy + ")"
	*/

	return fmt.Sprintf("%s:%s:%s", doc.String(), doc.fileHash.String(), doc.creator.String())
}

func (doc DocumentData) Equal(b DocumentData) bool {

	if doc.fileHash != b.fileHash {
		return false
	}

	if !doc.creator.Equal(b.creator) {
		return false
	}

	if !doc.owner.Equal(b.owner) {
		return false
	}

	sort.Slice(doc.signers, func(i, j int) bool {
		return bytes.Compare(doc.signers[i].Bytes(), doc.signers[j].Bytes()) < 0
	})
	sort.Slice(b.signers, func(i, j int) bool {
		return bytes.Compare(b.signers[i].Bytes(), b.signers[j].Bytes()) < 0
	})

	for i := range doc.signers {
		if !doc.signers[i].Equal(b.signers[i]) {
			return false
		}
	}

	return true
}

func (doc DocumentData) WithData(fileHash FileHash, creator DocSign, owner base.Address, signers []DocSign) DocumentData {
	doc.fileHash = fileHash
	doc.creator = creator
	doc.owner = owner
	doc.signers = signers
	return doc
}

var (
	FileHashType = hint.Type("mitum-blocksign-filehash")
	FileHashHint = hint.NewHint(FileHashType, "v0.0.1")
)

type FileHash string

func (fh FileHash) Bytes() []byte {
	return []byte(fh)
}

func (fh FileHash) String() string {
	return string(fh)
}

func (fh FileHash) Hint() hint.Hint {
	return FileHashHint
}

func (fh FileHash) Hash() valuehash.Hash {
	return fh.GenerateHash()
}

func (fh FileHash) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(fh.Bytes())
}

func (fh FileHash) IsValid([]byte) error {
	return nil
}

func (fh FileHash) Equal(b FileHash) bool {
	return fh == b
}

/*
var (
	DocSignType = hint.Type("mitum-blocksign-docsign")
	DocSignHint = hint.NewHint(DocSignType, "v0.0.1")
)

type DocSign struct {
	key       key.Publickey
	signature key.Signature
	signedAt  time.Time
}

func NewDocSign(key key.Publickey, signature key.Signature, signedAt time.Time) DocSign {
	docsign := DocSign{
		key:       key,
		signature: signature,
		signedAt:  signedAt,
	}
	return docsign
}

func MustNewDocSign(key key.Publickey, signature key.Signature, signedAt time.Time) DocSign {
	docsign := NewDocSign(key, signature, signedAt)
	if err := docsign.IsValid(nil); err != nil {
		panic(err)
	}
	return docsign
}

func (ds DocSign) Bytes() []byte {
	return util.ConcatBytesSlice(
		ds.key.Bytes(),
		ds.signature.Bytes(),
		localtime.NewTime(ds.signedAt).Bytes(),
	)
}

func (ds DocSign) String() string {
	return fmt.Sprintf("%s:%s:%s", ds.key.String(), ds.signature.String(), ds.signedAt.String())
}

func (ds DocSign) Hint() hint.Hint {
	return DocSignHint
}

func (ds DocSign) Hash() valuehash.Hash {
	return ds.GenerateHash()
}

func (ds DocSign) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(ds.Bytes())
}

func (ds DocSign) IsValid([]byte) error {
	return nil
}

func (ds DocSign) IsEmpty() bool {
	return ds.signedAt.IsZero() || len(ds.signature) < 1 || ds.key.Bytes() == nil
}

func (ds DocSign) Equal(b DocSign) bool {
	switch {
	case ds.key != b.key:
		return false
	case !ds.signature.Equal(b.signature):
		return false
	case !ds.signedAt.Equal(b.signedAt):
		return false
	default:
		return true
	}
}
*/

var (
	DocSignsType = hint.Type("mitum-blocksign-docsign")
	DocSignsHint = hint.NewHint(DocSignsType, "v0.0.1")
)

type DocSign struct {
	address base.Address
	signed  bool
}

func NewDocSign(address base.Address, signed bool) DocSign {
	doc := DocSign{
		address: address,
		signed:  signed,
	}
	return doc
}

func MustNewDocSign(address base.Address, signed bool) DocSign {
	doc := NewDocSign(address, signed)
	if err := doc.IsValid(nil); err != nil {
		panic(err)
	}
	return doc
}

func (doc DocSign) Address() base.Address {
	return doc.address
}

func (ds DocSign) Bytes() []byte {
	bs := make([][]byte, 2)

	bs[0] = ds.address.Bytes()
	var v int8
	if ds.signed {
		v = 1
	}
	bs[1] = []byte{byte(v)}
	return util.ConcatBytesSlice(bs...)
}

func (ds DocSign) Hash() valuehash.Hash {
	return ds.GenerateHash()
}

func (ds DocSign) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(ds.Bytes())
}

func (ds DocSign) Hint() hint.Hint {
	return DocSignsHint
}

func (ds DocSign) IsValid([]byte) error {
	return nil
}

func (ds DocSign) IsEmpty() bool {
	return len(ds.address.Raw()) < 1
}

func (ds DocSign) String() string {
	v := fmt.Sprintf("%v", ds.signed)
	return fmt.Sprintf("%s:%s", ds.address.Raw(), v)
}

func (ds DocSign) Equal(b DocSign) bool {

	if !ds.address.Equal(b.address) {
		return false
	}

	if ds.signed != b.signed {
		return false
	}

	return true
}
