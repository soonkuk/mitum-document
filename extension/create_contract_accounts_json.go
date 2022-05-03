package extension

import (
	"encoding/json"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
	"github.com/spikeekips/mitum/util/valuehash"
)

type CreateContractAccountsFactJSONPacker struct {
	jsonenc.HintedHead
	H  valuehash.Hash               `json:"hash"`
	TK []byte                       `json:"token"`
	SD base.Address                 `json:"sender"`
	IT []CreateContractAccountsItem `json:"items"`
}

func (fact CreateContractAccountsFact) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(CreateContractAccountsFactJSONPacker{
		HintedHead: jsonenc.NewHintedHead(fact.Hint()),
		H:          fact.h,
		TK:         fact.token,
		SD:         fact.sender,
		IT:         fact.items,
	})
}

type CreateContractAccountsFactJSONUnpacker struct {
	H  valuehash.Bytes     `json:"hash"`
	TK []byte              `json:"token"`
	SD base.AddressDecoder `json:"sender"`
	IT json.RawMessage     `json:"items"`
}

func (fact *CreateContractAccountsFact) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var uca CreateContractAccountsFactJSONUnpacker
	if err := jsonenc.Unmarshal(b, &uca); err != nil {
		return err
	}

	return fact.unpack(enc, uca.H, uca.TK, uca.SD, uca.IT)
}

func (op *CreateContractAccounts) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var ubo currency.BaseOperation
	if err := ubo.UnpackJSON(b, enc); err != nil {
		return err
	}

	op.BaseOperation = ubo

	return nil
}
