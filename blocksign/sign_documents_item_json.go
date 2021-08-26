package blocksign

import (
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
)

type SignDocumentsItemJSONPacker struct {
	jsonenc.HintedHead
	DI currency.Big        `json:"documentid"`
	OW base.Address        `json:"owner"`
	CI currency.CurrencyID `json:"currency"`
}

func (it BaseSignDocumentsItem) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(SignDocumentsItemJSONPacker{
		HintedHead: jsonenc.NewHintedHead(it.Hint()),
		DI:         it.id,
		OW:         it.owner,
		CI:         it.cid,
	})
}

type SignDocumentsItemJSONUnpacker struct {
	DI currency.Big        `json:"documentid"`
	OW base.AddressDecoder `json:"owner"`
	CI string              `json:"currency"`
}

func (it *BaseSignDocumentsItem) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var ht jsonenc.HintedHead
	if err := enc.Unmarshal(b, &ht); err != nil {
		return err
	}

	var ucd SignDocumentsItemJSONUnpacker
	if err := jsonenc.Unmarshal(b, &ucd); err != nil {
		return err
	}

	return it.unpack(enc, ht.H, ucd.DI, ucd.OW, ucd.CI)
}
