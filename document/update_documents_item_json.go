package document

import (
	"encoding/json"

	"github.com/spikeekips/mitum-currency/currency"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
	"github.com/spikeekips/mitum/util/hint"
)

type UpdateDocumentsItemImplJSONPacker struct {
	jsonenc.HintedHead
	DT hint.Type           `json:"doctype"`
	DD Document            `json:"doc"`
	CI currency.CurrencyID `json:"currency"`
}

func (it UpdateDocumentsItemImpl) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(UpdateDocumentsItemImplJSONPacker{
		HintedHead: jsonenc.NewHintedHead(it.Hint()),
		DT:         it.doctype,
		DD:         it.doc,
		CI:         it.cid,
	})
}

type UpdateDocumentsItemImplJSONUnpacker struct {
	DT string          `json:"doctype"`
	DD json.RawMessage `json:"doc"`
	CI string          `json:"currency"`
}

func (it *UpdateDocumentsItemImpl) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var ucd UpdateDocumentsItemImplJSONUnpacker
	if err := jsonenc.Unmarshal(b, &ucd); err != nil {
		return err
	}

	return it.unpack(
		enc,
		ucd.DT,
		ucd.DD,
		ucd.CI,
	)
}
