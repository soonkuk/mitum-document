package digest

import (
	"encoding/json"

	"github.com/protoconNet/mitum-document/document"
	"github.com/spikeekips/mitum/base"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
)

type DocumentValueJSONPacker struct {
	jsonenc.HintedHead
	DM document.DocumentData `json:"document"`
	HT base.Height           `json:"height"`
}

func (dv DocumentValue) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(DocumentValueJSONPacker{
		HintedHead: jsonenc.NewHintedHead(dv.Hint()),
		DM:         dv.doc,
		HT:         dv.height,
	})
}

type DocumentValueJSONUnpacker struct {
	DM json.RawMessage `json:"document"`
	HT base.Height     `json:"height"`
}

func (dv *DocumentValue) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var uva DocumentValueJSONUnpacker
	if err := enc.Unmarshal(b, &uva); err != nil {
		return err
	}

	err := dv.unpack(enc, uva.DM, uva.HT)
	if err != nil {
		return err
	}
	return nil
}
