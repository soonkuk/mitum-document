package document

import (
	"encoding/json"

	"github.com/spikeekips/mitum-currency/currency"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
	"github.com/spikeekips/mitum/util/hint"
)

type CreateDocumentsItemImplJSONPacker struct {
	jsonenc.HintedHead
	DT hint.Type           `json:"doctype"`
	DD DocumentData        `json:"doc"`
	CI currency.CurrencyID `json:"currency"`
}

func (it CreateDocumentsItemImpl) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(CreateDocumentsItemImplJSONPacker{
		HintedHead: jsonenc.NewHintedHead(it.Hint()),
		DT:         it.doctype,
		DD:         it.doc,
		CI:         it.cid,
	})
}

type CreateDocumentsItemImplJSONUnpacker struct {
	DT string          `json:"doctype"`
	DD json.RawMessage `json:"doc"`
	CI string          `json:"currency"`
}

func (it *CreateDocumentsItemImpl) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var ucd CreateDocumentsItemImplJSONUnpacker
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
