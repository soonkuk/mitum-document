package digest

import (
	"encoding/json"

	"github.com/soonkuk/mitum-blocksign/blockcity"
	"github.com/soonkuk/mitum-blocksign/blocksign"
	"github.com/spikeekips/mitum/base"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
)

type BlocksignDocumentValueJSONPacker struct {
	jsonenc.HintedHead
	DM blocksign.DocumentData `json:"document"`
	HT base.Height            `json:"height"`
}

func (va BlocksignDocumentValue) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(BlocksignDocumentValueJSONPacker{
		HintedHead: jsonenc.NewHintedHead(va.Hint()),
		DM:         va.doc,
		HT:         va.height,
	})
}

type BlocksignDocumentValueJSONUnpacker struct {
	DM json.RawMessage `json:"document"`
	HT base.Height     `json:"height"`
}

func (dv *BlocksignDocumentValue) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var uva BlocksignDocumentValueJSONUnpacker
	if err := enc.Unmarshal(b, &uva); err != nil {
		return err
	}

	if err := dv.unpack(enc, uva.DM, uva.HT); err != nil {
		return err
	}
	return nil
}

type BlockcityDocumentValueJSONPacker struct {
	jsonenc.HintedHead
	DM blockcity.Document `json:"document"`
	HT base.Height        `json:"height"`
}

func (va BlockcityDocumentValue) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(BlockcityDocumentValueJSONPacker{
		HintedHead: jsonenc.NewHintedHead(va.Hint()),
		DM:         va.doc,
		HT:         va.height,
	})
}

type BlockcityDocumentValueJSONUnpacker struct {
	DM json.RawMessage `json:"document"`
	HT base.Height     `json:"height"`
}

func (dv *BlockcityDocumentValue) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var uva BlockcityDocumentValueJSONUnpacker
	if err := enc.Unmarshal(b, &uva); err != nil {
		return err
	}

	if err := dv.unpack(enc, uva.DM, uva.HT); err != nil {
		return err
	}
	return nil
}
