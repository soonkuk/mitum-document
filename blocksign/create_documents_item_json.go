package blocksign

import (
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
)

type CreateDocumentsItemJSONPacker struct {
	jsonenc.HintedHead
	FH FileHash            `json:"filehash"`
	SG []base.Address      `json:"signers"`
	CI currency.CurrencyID `json:"currency"`
}

func (it BaseCreateDocumentsItem) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(CreateDocumentsItemJSONPacker{
		HintedHead: jsonenc.NewHintedHead(it.Hint()),
		FH:         it.fileHash,
		SG:         it.signers,
		CI:         it.cid,
	})
}

type CreateDocumentsItemJSONUnpacker struct {
	FH string                `json:"filehash"`
	SG []base.AddressDecoder `json:"signers"`
	CI string                `json:"currency"`
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

	return it.unpack(enc, ht.H, ucd.FH, ucd.SG, ucd.CI)
}
