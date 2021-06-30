package blocksign

import (
	"encoding/json"

	"github.com/soonkuk/mitum-data/currency"
	"github.com/spikeekips/mitum/base"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
)

type CreateDocumentsItemJSONPacker struct {
	jsonenc.HintedHead
	KS currency.Keys       `json:"keys"`
	SC SignCode            `json:"signcode"`
	OW base.Address        `json:"owner"`
	CI currency.CurrencyID `json:"currency"`
}

func (it BaseCreateDocumentsItem) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(CreateDocumentsItemJSONPacker{
		HintedHead: jsonenc.NewHintedHead(it.Hint()),
		KS:         it.keys,
		SC:         it.sc,
		OW:         it.owner,
		CI:         it.cid,
	})
}

type CreateDocumentsItemJSONUnpacker struct {
	KS json.RawMessage     `json:"keys"`
	SC string              `json:"signcode"`
	OW base.AddressDecoder `json:"owner"`
	CI string              `json:"currency"`
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

	return it.unpack(enc, ht.H, ucd.KS, ucd.SC, ucd.OW, ucd.CI)
}
