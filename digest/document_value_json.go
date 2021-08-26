package digest

import (
	"encoding/json"

	"github.com/soonkuk/mitum-blocksign/blocksign"
	"github.com/spikeekips/mitum/base"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
)

type DocumentValueJSONPacker struct {
	jsonenc.HintedHead
	DM blocksign.DocumentData `json:"document"`
	HT base.Height            `json:"height"`
}

func (va DocumentValue) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(DocumentValueJSONPacker{
		HintedHead: jsonenc.NewHintedHead(va.Hint()),
		DM:         va.doc,
		HT:         va.height,
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

	if err := dv.unpack(enc, uva.DM, uva.HT); err != nil {
		return err
	}
	return nil
}
