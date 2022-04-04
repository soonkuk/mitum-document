package document // nolint: dupl, revive

import (
	"bytes"
	"encoding/json"
	"sort"

	"github.com/pkg/errors"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/valuehash"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	DocumentInventoryType   = hint.Type("mitum-document-inventory")
	DocumentInventoryHint   = hint.NewHint(DocumentInventoryType, "v0.0.1")
	DocumentInventoryHinter = DocumentInventory{BaseHinter: hint.NewBaseHinter(DocumentInventoryHint)}
)

type DocumentInventory struct {
	hint.BaseHinter
	docInfos []DocInfo
}

func NewDocumentInventory(docInfos []DocInfo) DocumentInventory {
	if docInfos == nil {
		return DocumentInventory{
			BaseHinter: hint.NewBaseHinter(DocumentInventoryHint),
			docInfos:   []DocInfo{},
		}
	}
	return DocumentInventory{
		BaseHinter: hint.NewBaseHinter(DocumentInventoryHint),
		docInfos:   docInfos,
	}
}

func MustNewDocumentInventory(docInfos []DocInfo) DocumentInventory {
	d := NewDocumentInventory(docInfos)
	if err := d.IsValid(nil); err != nil {
		panic(err)
	}

	return d
}

func (div DocumentInventory) Bytes() []byte {
	bs := make([][]byte, len(div.docInfos))
	for i := range div.docInfos {
		bs[i] = div.docInfos[i].Bytes()
	}

	return util.ConcatBytesSlice(bs...)
}

func (div DocumentInventory) Hint() hint.Hint {
	return DocumentInventoryHint
}

func (div DocumentInventory) Hash() valuehash.Hash {
	return div.GenerateHash()
}

func (div DocumentInventory) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(div.Bytes())
}

func (div DocumentInventory) IsEmpty() bool {
	return len(div.docInfos) < 1
}

func (div DocumentInventory) IsValid([]byte) error {
	for i := range div.docInfos {
		if err := div.docInfos[i].IsValid(nil); err != nil {
			return err
		}
	}
	return nil
}

func (div DocumentInventory) Equal(b DocumentInventory) bool {
	if len(div.docInfos) != len(b.docInfos) {
		return false
	}
	if len(div.docInfos) < 1 && len(b.docInfos) < 1 {
		return true
	}
	div.Sort(true)
	b.Sort(true)
	for i := range div.docInfos {
		if !div.docInfos[i].Equal(b.docInfos[i]) {
			return false
		}
	}
	return true
}

func (div *DocumentInventory) Sort(ascending bool) {
	if len(div.docInfos) < 1 {
		return
	}
	sort.Slice(div.docInfos, func(i, j int) bool {
		if ascending {
			return bytes.Compare(div.docInfos[j].id.Bytes(), div.docInfos[i].id.Bytes()) > 0
		}
		return bytes.Compare(div.docInfos[j].id.Bytes(), div.docInfos[i].id.Bytes()) < 0
	})
}

func (div DocumentInventory) Exists(id string) bool {
	if len(div.docInfos) < 1 {
		return false
	}
	for i := range div.docInfos {
		if id == div.docInfos[i].id.String() {
			return true
		}
	}
	return false
}

func (div DocumentInventory) Get(id string) (DocInfo, error) {
	for i := range div.docInfos {
		if div.docInfos[i].id.String() == id {
			return div.docInfos[i], nil
		}
	}
	return DocInfo{}, errors.Errorf("document not found in Owner's DocumentInventory, %v", id)
}

func (div *DocumentInventory) Append(d DocInfo) error {
	if err := d.IsValid(nil); err != nil {
		return err
	}
	if div.Exists(d.id.String()) {
		return errors.Errorf("document id %v already exists in document inventory", d.id.String())
	}
	div.docInfos = append(div.docInfos, d)
	return nil
}

func (div *DocumentInventory) Romove(d DocInfo) error {
	if !div.Exists(d.id.String()) {
		return errors.Errorf("document id %v not found in document inventory", d.id.String())
	}
	for i := range div.docInfos {
		if d.id.String() == div.docInfos[i].id.String() {
			div.docInfos[i] = div.docInfos[len(div.docInfos)-1]
			div.docInfos[len(div.docInfos)-1] = DocInfo{}
			div.docInfos = div.docInfos[:len(div.docInfos)-1]
			return nil
		}
	}
	return nil
}

func (div DocumentInventory) Documents() []DocInfo {
	return div.docInfos
}

type DocumentInventoryJSONPacker struct {
	jsonenc.HintedHead
	DI []DocInfo `json:"documents"`
}

func (div DocumentInventory) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(DocumentInventoryJSONPacker{
		HintedHead: jsonenc.NewHintedHead(div.Hint()),
		DI:         div.docInfos,
	})
}

type DocumentInventoryJSONUnpacker struct {
	DI json.RawMessage `json:"address"`
}

func (div *DocumentInventory) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var udi DocumentInventoryJSONUnpacker
	if err := enc.Unmarshal(b, &udi); err != nil {
		return err
	}

	return div.unpack(enc, udi.DI)
}

type DocumentInventoryBSONPacker struct {
	DI []DocInfo `bson:"documents"`
}

func (div DocumentInventory) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bsonenc.MergeBSONM(
		bsonenc.NewHintedDoc(div.Hint()),
		bson.M{
			"documents": div.docInfos,
		}),
	)
}

type DocumentInventoryBSONUnpacker struct {
	DI bson.Raw `bson:"documents"`
}

func (div *DocumentInventory) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var udi DocumentInventoryBSONUnpacker
	if err := bsonenc.Unmarshal(b, &udi); err != nil {
		return err
	}

	return div.unpack(enc, udi.DI)
}

func (div *DocumentInventory) unpack(
	enc encoder.Encoder,
	dis []byte, // DocInfos
) error {
	hits, err := enc.DecodeSlice(dis)
	if err != nil {
		return err
	}

	docInfos := make([]DocInfo, len(hits))
	for i := range hits {
		j, ok := hits[i].(DocInfo)
		if !ok {
			return util.WrongTypeError.Errorf("expected DocInfo, not %T", hits[i])
		}

		docInfos[i] = j
	}

	div.docInfos = docInfos

	return nil
}
