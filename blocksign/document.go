package blocksign

import (
	"bytes"
	"fmt"
	"math/big"
	"sort"

	"github.com/pkg/errors"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
	"github.com/spikeekips/mitum/util/valuehash"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	DocumentDataType = hint.Type("mitum-blocksign-document-data")
	DocumentDataHint = hint.NewHint(DocumentDataType, "v0.0.1")
)

type DocumentData struct {
	info    DocInfo
	creator DocSign
	title   string
	size    currency.Big
	signers []DocSign
}

func NewDocumentData(info DocInfo,
	creator base.Address,
	signcode string,
	title string,
	size currency.Big,
	signers []DocSign) DocumentData {
	doc := DocumentData{
		info: info,
		creator: DocSign{
			address:  creator,
			signcode: signcode,
			signed:   true,
		},
		title:   title,
		size:    size,
		signers: signers,
	}

	return doc
}

func MustNewDocumentData(info DocInfo, creator base.Address, signcode string, title string, size currency.Big, signers []DocSign) DocumentData {
	doc := NewDocumentData(info, creator, signcode, title, size, signers)
	if err := doc.IsValid(nil); err != nil {
		panic(err)
	}

	return doc
}

func (doc DocumentData) Hint() hint.Hint {
	return DocumentDataHint
}

func (doc DocumentData) Bytes() []byte {
	bs := make([][]byte, len(doc.signers)+4)

	sort.Slice(doc.signers, func(i, j int) bool {
		return bytes.Compare(doc.signers[i].Bytes(), doc.signers[j].Bytes()) < 0
	})

	bs[0] = doc.info.Bytes()
	bs[1] = doc.creator.Bytes()
	bs[2] = []byte(doc.title)
	bs[3] = doc.size.Bytes()
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
	return len(doc.info.FileHash()) < 1 || len(doc.signers) < 1 || len(doc.title) < 1 || !doc.size.OverZero()
}

func (doc DocumentData) IsValid([]byte) error {
	if err := isvalid.Check([]isvalid.IsValider{
		doc.info.FileHash(),
		doc.creator,
	}, nil, false); err != nil {
		return errors.Wrap(err, "invalid document data")
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
	return doc.info.FileHash()
}

func (doc DocumentData) SignCode() string {
	return doc.creator.signcode
}

func (doc DocumentData) Title() string {
	return doc.title
}

func (doc DocumentData) Size() currency.Big {
	return doc.size
}

func (doc DocumentData) Info() DocInfo {
	return doc.info
}

func (doc DocumentData) Creator() base.Address {
	return doc.creator.address
}

func (doc DocumentData) Signers() []DocSign {
	return doc.signers
}

func (doc DocumentData) Addresses() ([]base.Address, error) {
	addresses := make(map[base.Address]bool)
	addresses[doc.creator.Address()] = true
	for i := range doc.Signers() {
		_, found := addresses[doc.Signers()[i].Address()]
		if !found {
			addresses[doc.Signers()[i].Address()] = true
		}
	}
	result := make([]base.Address, len(addresses))
	i := 0
	for k := range addresses {
		result[i] = k
		i = i + 1
	}
	return result, nil
}

func (doc DocumentData) String() string {

	return fmt.Sprintf("%s:%s:%s:%s:%s",
		doc.FileHash().String(),
		doc.info.String(),
		doc.creator.String(),
		doc.title,
		doc.size)
}

func (doc DocumentData) Equal(b DocumentData) bool {

	if doc.FileHash() != b.FileHash() {
		return false
	}

	if !doc.creator.Equal(b.creator) {
		return false
	}

	if doc.title != (b.title) {
		return false
	}

	if !doc.size.Equal(b.size) {
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

func (doc DocumentData) WithData(info DocInfo, creator DocSign, signcode string, title string, size currency.Big, signers []DocSign) DocumentData {
	doc.info = info
	doc.creator = creator
	doc.title = title
	doc.size = size
	doc.signers = signers
	return doc
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
	DocSignType = hint.Type("mitum-blocksign-docsign")
	DocSignHint = hint.NewHint(DocSignType, "v0.0.1")
)

type DocSign struct {
	address  base.Address
	signcode string
	signed   bool
}

func NewDocSign(address base.Address, signcode string, signed bool) DocSign {
	doc := DocSign{
		address:  address,
		signcode: signcode,
		signed:   signed,
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
	return DocSignHint
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

type DocSignJSONPacker struct {
	jsonenc.HintedHead
	AD base.Address `json:"address"`
	SC string       `json:"signcode"`
	SG bool         `json:"signed"`
}

func (ds DocSign) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(DocSignJSONPacker{
		HintedHead: jsonenc.NewHintedHead(ds.Hint()),
		AD:         ds.address,
		SC:         ds.signcode,
		SG:         ds.signed,
	})
}

type DocSignJSONUnpacker struct {
	AD base.AddressDecoder `json:"address"`
	SC string              `json:"signcode"`
	SG bool                `json:"signed"`
}

func (ds *DocSign) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var uds DocSignJSONUnpacker
	if err := enc.Unmarshal(b, &uds); err != nil {
		return err
	}

	return ds.unpack(enc, uds.AD, uds.SC, uds.SG)
}

type DocSignBSONPacker struct {
	AD base.Address `bson:"address"`
	SC string       `bson:"signcode"`
	SG bool         `bson:"signed"`
}

func (ds DocSign) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bsonenc.MergeBSONM(
		bsonenc.NewHintedDoc(ds.Hint()),
		bson.M{
			"address":  ds.address,
			"signcode": ds.signcode,
			"signed":   ds.signed,
		}),
	)
}

type DocSignBSONUnpacker struct {
	AD base.AddressDecoder `bson:"address"`
	SC string              `bson:"signcode"`
	SG bool                `bson:"signed"`
}

func (ds *DocSign) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var uds DocSignBSONUnpacker
	if err := bsonenc.Unmarshal(b, &uds); err != nil {
		return err
	}

	return ds.unpack(enc, uds.AD, uds.SC, uds.SG)
}

var (
	DocInfoType = hint.Type("mitum-blocksign-document-info")
	DocInfoHint = hint.NewHint(DocInfoType, "v0.0.1")
)

type DocInfo struct {
	idx      currency.Big
	filehash FileHash
}

func NewDocInfo(idx int64, fh FileHash) DocInfo {
	id := currency.NewBig(idx)
	if !id.OverNil() {
		return DocInfo{}
	}
	docInfo := DocInfo{
		idx:      id,
		filehash: fh,
	}
	return docInfo
}

func MustNewDocInfo(idx int64, fh FileHash) DocInfo {
	docInfo := NewDocInfo(idx, fh)
	if err := docInfo.IsValid(nil); err != nil {
		panic(err)
	}
	return docInfo
}

func NewDocInfoFromString(id string, fh string) (DocInfo, error) {
	i, ok := new(big.Int).SetString(id, 10)
	if !ok {
		return DocInfo{}, errors.Errorf("not proper DocInfo string, %q", id)
	}
	idx := currency.NewBigFromBigInt(i)
	if !idx.OverNil() {
		return DocInfo{}, nil
	}
	docInfo := DocInfo{
		idx:      idx,
		filehash: FileHash(fh),
	}
	return docInfo, nil
}

func (di DocInfo) Index() currency.Big {
	return di.idx
}

func (di DocInfo) FileHash() FileHash {
	return di.filehash
}

func (di DocInfo) Bytes() []byte {

	return util.ConcatBytesSlice(di.idx.Bytes(), di.filehash.Bytes())
}

func (di DocInfo) Hash() valuehash.Hash {
	return di.GenerateHash()
}

func (di DocInfo) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(di.Bytes())
}

func (di DocInfo) Hint() hint.Hint {
	return DocInfoHint
}

func (di DocInfo) IsValid([]byte) error {
	if err := di.idx.IsValid(nil); err != nil {
		return err
	} else if err := di.filehash.IsValid(nil); err != nil {
		return err
	}

	return nil
}

func (di DocInfo) IsEmpty() bool {
	return !di.idx.OverNil() || len(di.filehash) < 1
}

func (di DocInfo) String() string {
	return fmt.Sprintf("%s:%s", di.idx.String(), di.filehash.String())
}

func (di DocInfo) Equal(b DocInfo) bool {
	return di.idx.Equal(b.idx) && di.filehash.Equal(b.filehash)
}

func (di DocInfo) WithData(idx currency.Big, fh FileHash) DocInfo {
	di.idx = idx
	di.filehash = fh
	return di
}

type DocInfoJSONPacker struct {
	jsonenc.HintedHead
	ID currency.Big `json:"documentid"`
	FH FileHash     `json:"filehash"`
}

func (di DocInfo) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(DocInfoJSONPacker{
		HintedHead: jsonenc.NewHintedHead(di.Hint()),
		ID:         di.idx,
		FH:         di.filehash,
	})
}

type DocInfoJSONUnpacker struct {
	ID currency.Big `json:"documentid"`
	FH string       `json:"filehash"`
}

func (di *DocInfo) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var udi DocInfoJSONUnpacker
	if err := enc.Unmarshal(b, &udi); err != nil {
		return err
	}

	di.idx = udi.ID
	di.filehash = FileHash(udi.FH)

	return nil
}

type DocInfoBSONPacker struct {
	ID currency.Big `bson:"documentid"`
	FH FileHash     `bson:"filehash"`
}

func (di DocInfo) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bsonenc.MergeBSONM(
		bsonenc.NewHintedDoc(di.Hint()),
		bson.M{
			"documentid": di.idx,
			"filehash":   di.filehash,
		}),
	)
}

type DocInfoBSONUnpacker struct {
	ID currency.Big `bson:"documentid"`
	FH string       `bson:"filehash"`
}

func (di *DocInfo) UnmarshalBSON(b []byte) error {
	var udi DocInfoBSONUnpacker
	if err := bsonenc.Unmarshal(b, &udi); err != nil {
		return err
	}

	di.idx = udi.ID
	di.filehash = FileHash(udi.FH)

	return nil
}

type DocId currency.Big

func NewDocId(idx int64) DocId {
	id := currency.NewBig(idx)
	if !id.OverNil() {
		return DocId{}
	}

	return DocId(id)
}

func MustNewDocId(idx int64) DocId {
	docId := NewDocId(idx)
	if err := docId.IsValid(nil); err != nil {
		panic(err)
	}
	return docId
}

func NewDocIdFromString(id string) (DocId, error) {
	i, ok := new(big.Int).SetString(id, 10)
	if !ok {
		return DocId{}, errors.Errorf("not proper DocId string, %q", id)
	}
	idx := currency.NewBigFromBigInt(i)
	if !idx.OverNil() {
		return DocId{}, nil
	}

	return DocId(idx), nil
}

func (di DocId) Index() currency.Big {
	return currency.Big(di)
}

func (di DocId) Bytes() []byte {
	return currency.Big(di).Bytes()
}

func (di DocId) Hash() valuehash.Hash {
	return di.GenerateHash()
}

func (di DocId) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(di.Bytes())
}

func (di DocId) IsValid([]byte) error {
	if err := currency.Big(di).IsValid(nil); err != nil {
		return err
	}

	return nil
}

func (di DocId) IsEmpty() bool {
	return !currency.Big(di).OverNil()
}

func (di DocId) String() string {
	return currency.Big(di).String()
}

func (di DocId) Equal(b DocId) bool {
	return currency.Big(di).Equal(currency.Big(b))
}

func (di DocId) WithData(idx currency.Big) DocId {
	return DocId(idx)
}

type SignCode string

func (sc SignCode) Bytes() []byte {
	return []byte(sc)
}

func (sc SignCode) String() string {
	return string(sc)
}

func (sc SignCode) IsValid([]byte) error {
	return nil
}
