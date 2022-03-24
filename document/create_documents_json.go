package document // nolint: dupl, revive

import (
	"encoding/json"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
	"github.com/spikeekips/mitum/util/valuehash"
)

type CreateDocumentsFactJSONPacker struct {
	jsonenc.HintedHead
	H  valuehash.Hash        `json:"hash"`
	TK []byte                `json:"token"`
	SD base.Address          `json:"sender"`
	IT []CreateDocumentsItem `json:"items"`
}

func (fact CreateDocumentsFact) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(CreateDocumentsFactJSONPacker{
		HintedHead: jsonenc.NewHintedHead(fact.Hint()),
		H:          fact.h,
		TK:         fact.token,
		SD:         fact.sender,
		IT:         fact.items,
	})
}

type CreateDocumentsFactJSONUnpacker struct {
	H  valuehash.Bytes     `json:"hash"`
	TK []byte              `json:"token"`
	SD base.AddressDecoder `json:"sender"`
	IT json.RawMessage     `json:"items"`
}

func (fact *CreateDocumentsFact) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var uda CreateDocumentsFactJSONUnpacker
	if err := jsonenc.Unmarshal(b, &uda); err != nil {
		return err
	}

	return fact.unpack(enc, uda.H, uda.TK, uda.SD, uda.IT)
}

func (op *CreateDocuments) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var ubo currency.BaseOperation
	if err := ubo.UnpackJSON(b, enc); err != nil {
		return err
	}

	op.BaseOperation = ubo

	return nil
}
