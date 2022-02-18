package digest

import (
	"encoding/json"

	"github.com/protoconNet/mitum-document/document"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
)

type AccountValueJSONPacker struct {
	jsonenc.HintedHead
	currency.AccountPackerJSON
	BL []currency.Amount          `json:"balance,omitempty"`
	CD document.DocumentInventory `json:"documents"`
	HT base.Height                `json:"height"`
	PT base.Height                `json:"previous_height"`
}

func (va AccountValue) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(AccountValueJSONPacker{
		HintedHead:        jsonenc.NewHintedHead(va.Hint()),
		AccountPackerJSON: va.ac.PackerJSON(),
		BL:                va.balance,
		CD:                va.document,
		HT:                va.height,
		PT:                va.previousHeight,
	})
}

type AccountValueJSONUnpacker struct {
	BL json.RawMessage `json:"balance"`
	CD json.RawMessage `json:"documents"`
	HT base.Height     `json:"height"`
	PT base.Height     `json:"previous_height"`
}

func (va *AccountValue) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var uva AccountValueJSONUnpacker
	if err := enc.Unmarshal(b, &uva); err != nil {
		return err
	}

	ac := new(currency.Account)
	if err := va.unpack(enc, nil, uva.BL /*uva.SD, */, uva.CD, uva.HT, uva.PT); err != nil {
		return err
	} else if err := ac.UnpackJSON(b, enc); err != nil {
		return err
	} else {
		va.ac = *ac

		return nil
	}
}
