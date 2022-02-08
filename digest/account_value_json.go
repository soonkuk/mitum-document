package digest

import (
	"encoding/json"

	"github.com/soonkuk/mitum-blocksign/blocksign"
	"github.com/soonkuk/mitum-blocksign/document"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
)

type AccountValueJSONPacker struct {
	jsonenc.HintedHead
	currency.AccountPackerJSON
	BL []currency.Amount           `json:"balance,omitempty"`
	SD blocksign.DocumentInventory `json:"blocksign_documents"`
	CD document.DocumentInventory  `json:"blockcity_documents"`
	HT base.Height                 `json:"height"`
	PT base.Height                 `json:"previous_height"`
}

func (va AccountValue) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(AccountValueJSONPacker{
		HintedHead:        jsonenc.NewHintedHead(va.Hint()),
		AccountPackerJSON: va.ac.PackerJSON(),
		BL:                va.balance,
		SD:                va.bsDocument,
		CD:                va.bcDocument,
		HT:                va.height,
		PT:                va.previousHeight,
	})
}

type AccountValueJSONUnpacker struct {
	BL json.RawMessage `json:"balance"`
	SD json.RawMessage `json:"blocksign_documents"`
	CD json.RawMessage `json:"blockcity_documents"`
	HT base.Height     `json:"height"`
	PT base.Height     `json:"previous_height"`
}

func (va *AccountValue) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var uva AccountValueJSONUnpacker
	if err := enc.Unmarshal(b, &uva); err != nil {
		return err
	}

	ac := new(currency.Account)
	if err := va.unpack(enc, nil, uva.BL, uva.SD, uva.CD, uva.HT, uva.PT); err != nil {
		return err
	} else if err := ac.UnpackJSON(b, enc); err != nil {
		return err
	} else {
		va.ac = *ac

		return nil
	}
}
