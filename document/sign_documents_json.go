package document // nolint:dupl

import (
	"encoding/json"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
	"github.com/spikeekips/mitum/util/valuehash"
)

type SignDocumentsFactJSONPacker struct {
	jsonenc.HintedHead
	H  valuehash.Hash      `json:"hash"`
	TK []byte              `json:"token"`
	SD base.Address        `json:"sender"`
	IT []SignDocumentsItem `json:"items"`
}

func (fact SignDocumentsFact) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(SignDocumentsFactJSONPacker{
		HintedHead: jsonenc.NewHintedHead(fact.Hint()),
		H:          fact.h,
		TK:         fact.token,
		SD:         fact.sender,
		IT:         fact.items,
	})
}

type SignDocumentsFactJSONUnpacker struct {
	H  valuehash.Bytes     `json:"hash"`
	TK []byte              `json:"token"`
	SD base.AddressDecoder `json:"sender"`
	IT json.RawMessage     `json:"items"`
}

func (fact *SignDocumentsFact) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var uda SignDocumentsFactJSONUnpacker
	if err := jsonenc.Unmarshal(b, &uda); err != nil {
		return err
	}

	return fact.unpack(enc, uda.H, uda.TK, uda.SD, uda.IT)
}

func (op *SignDocuments) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var ubo currency.BaseOperation
	if err := ubo.UnpackJSON(b, enc); err != nil {
		return err
	}

	op.BaseOperation = ubo

	return nil
}
