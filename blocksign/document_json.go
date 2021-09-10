package blocksign

import (
	"encoding/json"

	"github.com/spikeekips/mitum-currency/currency"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
)

type DocumentJSONPacker struct {
	jsonenc.HintedHead
	DI DocInfo      `json:"documentinfo"`
	CR DocSign      `json:"creator"`
	TL string       `json:"title"`
	SZ currency.Big `json:"size"`
	SG []DocSign    `json:"signers"`
}

func (doc DocumentData) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(DocumentJSONPacker{
		HintedHead: jsonenc.NewHintedHead(doc.Hint()),
		DI:         doc.info,
		CR:         doc.creator,
		TL:         doc.title,
		SZ:         doc.size,
		SG:         doc.signers,
	})
}

type DocumentJSONUnpacker struct {
	DI json.RawMessage `json:"documentinfo"`
	CR json.RawMessage `json:"creator"`
	TL string          `json:"title"`
	SZ currency.Big    `json:"size"`
	SG json.RawMessage `json:"signers"`
}

func (doc *DocumentData) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var udoc DocumentJSONUnpacker
	if err := enc.Unmarshal(b, &udoc); err != nil {
		return err
	}

	return doc.unpack(enc, udoc.DI, udoc.CR, udoc.TL, udoc.SZ, udoc.SG)
}
