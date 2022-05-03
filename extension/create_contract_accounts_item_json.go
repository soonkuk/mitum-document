package extension

import (
	"encoding/json"

	"github.com/spikeekips/mitum-currency/currency"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
)

type CreateContractAccountsItemJSONPacker struct {
	jsonenc.HintedHead
	KS currency.AccountKeys `json:"keys"`
	AS []currency.Amount    `json:"amounts"`
}

func (it BaseCreateContractAccountsItem) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(CreateContractAccountsItemJSONPacker{
		HintedHead: jsonenc.NewHintedHead(it.Hint()),
		KS:         it.keys,
		AS:         it.amounts,
	})
}

type CreateContractAccountsItemJSONUnpacker struct {
	KS json.RawMessage `json:"keys"`
	AM json.RawMessage `json:"amounts"`
}

func (it *BaseCreateContractAccountsItem) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var uca CreateContractAccountsItemJSONUnpacker
	if err := jsonenc.Unmarshal(b, &uca); err != nil {
		return err
	}

	return it.unpack(enc, uca.KS, uca.AM)
}
