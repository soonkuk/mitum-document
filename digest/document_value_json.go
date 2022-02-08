package digest

import (
	"encoding/json"

	"github.com/soonkuk/mitum-blocksign/blocksign"
	"github.com/soonkuk/mitum-blocksign/document"
	"github.com/spikeekips/mitum/base"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
)

type BSDocumentValueJSONPacker struct {
	jsonenc.HintedHead
	DM blocksign.DocumentData `json:"document"`
	HT base.Height            `json:"height"`
}

func (va BSDocumentValue) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(BSDocumentValueJSONPacker{
		HintedHead: jsonenc.NewHintedHead(va.Hint()),
		DM:         va.doc,
		HT:         va.height,
	})
}

type BSDocumentValueJSONUnpacker struct {
	DM json.RawMessage `json:"document"`
	HT base.Height     `json:"height"`
}

func (dv *BSDocumentValue) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var uva BSDocumentValueJSONUnpacker
	if err := enc.Unmarshal(b, &uva); err != nil {
		return err
	}

	if err := dv.unpack(enc, uva.DM, uva.HT); err != nil {
		return err
	}
	return nil
}

type BCDocumentValueJSONPacker struct {
	jsonenc.HintedHead
	DM document.DocumentData `json:"document"`
	HT base.Height           `json:"height"`
}

func (va BCDocumentValue) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(BCDocumentValueJSONPacker{
		HintedHead: jsonenc.NewHintedHead(va.Hint()),
		DM:         va.doc,
		HT:         va.height,
	})
}

type BCDocumentValueJSONUnpacker struct {
	DM json.RawMessage `json:"document"`
	HT base.Height     `json:"height"`
}

func (dv *BCDocumentValue) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var uva BCDocumentValueJSONUnpacker
	if err := enc.Unmarshal(b, &uva); err != nil {
		return err
	}

	if err := dv.unpack(enc, uva.DM, uva.HT); err != nil {
		return err
	}
	return nil
}
