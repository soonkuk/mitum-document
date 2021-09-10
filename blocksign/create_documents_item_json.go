package blocksign

import (
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
)

type CreateDocumentsItemJSONPacker struct {
	jsonenc.HintedHead
	FH FileHash            `json:"filehash"`
	DI currency.Big        `json:"documentid"`
	SC string              `json:"signcode"`
	TL string              `json:"title"`
	SZ currency.Big        `json:"size"`
	SG []base.Address      `json:"signers"`
	SD []string            `json:"signcodes"`
	CI currency.CurrencyID `json:"currency"`
}

func (it BaseCreateDocumentsItem) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(CreateDocumentsItemJSONPacker{
		HintedHead: jsonenc.NewHintedHead(it.Hint()),
		FH:         it.fileHash,
		DI:         it.documentid,
		SC:         it.signcode,
		TL:         it.title,
		SZ:         it.size,
		SG:         it.signers,
		SD:         it.signcodes,
		CI:         it.cid,
	})
}

type CreateDocumentsItemJSONUnpacker struct {
	FH string                `json:"filehash"`
	DI currency.Big          `json:"documentid"`
	SC string                `json:"signcode"`
	TL string                `json:"title"`
	SZ currency.Big          `json:"size"`
	SG []base.AddressDecoder `json:"signers"`
	SD []string              `json:"signcodes"`
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

	return it.unpack(enc, ht.H, ucd.FH, ucd.DI, ucd.SC, ucd.TL, ucd.SZ, ucd.SG, ucd.SD, ucd.CI)
}
