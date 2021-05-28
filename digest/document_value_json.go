package digest

import (
	"encoding/json"

	"github.com/soonkuk/mitum-data/currency"
	"github.com/spikeekips/mitum/base"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
)

type DocumentValueJSONPacker struct {
	jsonenc.HintedHead
	currency.AccountPackerJSON
	BL []currency.Amount `json:"balance"`
	FD currency.FileData `json:"filedata"`
	HT base.Height       `json:"height"`
	PT base.Height       `json:"previous_height"`
}

func (va DocumentValue) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(DocumentValueJSONPacker{
		HintedHead:        jsonenc.NewHintedHead(va.Hint()),
		AccountPackerJSON: va.ac.PackerJSON(),
		FD:                va.filedata,
		HT:                va.height,
		PT:                va.previousHeight,
	})
}

type DocumentValueJSONUnpacker struct {
	BL []json.RawMessage `json:"balance"`
	FD json.RawMessage   `json:"filedata"`
	HT base.Height       `json:"height"`
	PT base.Height       `json:"previous_height"`
}

func (dv *DocumentValue) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var uva DocumentValueJSONUnpacker
	if err := enc.Unmarshal(b, &uva); err != nil {
		return err
	}

	ac := new(currency.Account)
	if err := dv.unpack(enc, nil, uva.FD, uva.HT, uva.PT); err != nil {
		return err
	} else if err := ac.UnpackJSON(b, enc); err != nil {
		return err
	} else {
		dv.ac = *ac

		return nil
	}
}
