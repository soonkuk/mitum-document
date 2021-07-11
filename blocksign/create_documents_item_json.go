package blocksign

import (
	"encoding/json"

	"github.com/soonkuk/mitum-data/currency"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
)

type CreateDocumentsItemJSONPacker struct {
	jsonenc.HintedHead
	KS currency.Keys       `json:"keys"`
	DC DocumentData        `json:"document"`
	CI currency.CurrencyID `json:"currency"`
}

func (it BaseCreateDocumentsItem) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(CreateDocumentsItemJSONPacker{
		HintedHead: jsonenc.NewHintedHead(it.Hint()),
		KS:         it.keys,
		DC:         it.doc,
		CI:         it.cid,
	})
}

type CreateDocumentsItemJSONUnpacker struct {
	KS json.RawMessage `json:"keys"`
	DC json.RawMessage `json:"document"`
	CI string          `json:"currency"`
}

func (it *BaseCreateDocumentsItem) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var ht jsonenc.HintedHead
	if err := enc.Unmarshal(b, &ht); err != nil {
		return err
	}

	var ucd CreateDocumentsItemJSONUnpacker
	if err := jsonenc.Unmarshal(b, &ucd); err != nil {
		return err
	}

	return it.unpack(enc, ht.H, ucd.KS, ucd.DC, ucd.CI)
}
