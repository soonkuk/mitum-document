package blocksign

import (
	"encoding/json"
	"sort"

	"github.com/soonkuk/mitum-data/currency"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/valuehash"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/xerrors"
)

var (
	DocumentInventoryType = hint.Type("mbs-document-inventory")
	DocumentInventoryHint = hint.NewHint(DocumentInventoryType, "v0.0.1")
)

type DocumentInventory struct {
	documents []DocInfo
}

func NewDocumentInventory(documents []DocInfo) DocumentInventory {
	if documents == nil {
		return DocumentInventory{documents: []DocInfo{}}
	}
	return DocumentInventory{documents: documents}
}

func (div DocumentInventory) Bytes() []byte {
	bs := make([][]byte, len(div.documents))
	for i := range div.documents {
		bs[i] = div.documents[i].Bytes()
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
	return len(div.documents) < 1
}

func (div DocumentInventory) IsValid([]byte) error {
	for i := range div.documents {
		if err := div.documents[i].IsValid(nil); err != nil {
			return err
		}
	}
	return nil
}

func (div DocumentInventory) Equal(b DocumentInventory) bool {
	div.Sort(true)
	b.Sort(true)
	for i := range div.documents {
		if !div.documents[i].Equal(b.documents[i]) {
			return false
		}
	}
	return true
}

func (div *DocumentInventory) Sort(ascending bool) {
	sort.Slice(div.documents, func(i, j int) bool {
		if ascending {
			return div.documents[j].idx.Sub(div.documents[i].idx).OverZero()
		}
		return div.documents[i].idx.Sub(div.documents[j].idx).OverZero()
	})
}

func (div DocumentInventory) Exists(id currency.Big) bool {
	if len(div.documents) < 1 {
		return false
	}
	for i := range div.documents {
		if id.Equal(div.documents[i].idx) {
			return true
		}
	}
	return false
}

func (div DocumentInventory) Get(id currency.Big) (DocInfo, error) {
	for i := range div.documents {
		if div.documents[i].idx.Equal(id) {
			return div.documents[i], nil
		}
	}
	return DocInfo{}, xerrors.Errorf("Document not found in Owner's DocumentInventory, %v", id)
}

func (div *DocumentInventory) Append(d DocInfo) error {
	if err := d.IsValid(nil); err != nil {
		return err
	}
	if div.Exists(d.Index()) {
		return xerrors.Errorf("document id %v already exists in document inventory", d.idx)
	}
	div.documents = append(div.documents, d)
	return nil
}

func (div *DocumentInventory) Romove(d DocInfo) error {
	if !div.Exists(d.Index()) {
		return xerrors.Errorf("document id %v not found in document inventory", d.idx)
	}
	for i := range div.documents {
		if d.idx.Equal(div.documents[i].idx) {
			div.documents[i] = div.documents[len(div.documents)-1]
			div.documents[len(div.documents)-1] = DocInfo{}
			div.documents = div.documents[:len(div.documents)-1]
			return nil
		}
	}
	return nil
}

func (div DocumentInventory) Documents() []DocInfo {
	return div.documents
}

type DocumentInventoryJSONPacker struct {
	jsonenc.HintedHead
	DI []DocInfo `json:"documents"`
}

func (div DocumentInventory) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(DocumentInventoryJSONPacker{
		HintedHead: jsonenc.NewHintedHead(div.Hint()),
		DI:         div.documents,
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
			"documents": div.documents,
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

	div.documents = docInfos

	return nil
}
