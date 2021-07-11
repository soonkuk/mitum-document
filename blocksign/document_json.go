package blocksign

import (
	"encoding/json"

	"github.com/spikeekips/mitum/base"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
)

type DocumentJSONPacker struct {
	jsonenc.HintedHead
	FH FileHash     `json:"filehash"`
	CR DocSign      `json:"creator"`
	OW base.Address `json:"owner"`
	SG []DocSign    `json:"signers"`
}

func (doc DocumentData) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(DocumentJSONPacker{
		HintedHead: jsonenc.NewHintedHead(doc.Hint()),
		FH:         doc.fileHash,
		CR:         doc.creator,
		OW:         doc.owner,
		SG:         doc.signers,
	})
}

type DocumentJSONUnpacker struct {
	FH string              `json:"filehash"`
	CR json.RawMessage     `json:"creator"`
	OW base.AddressDecoder `json:"owner"`
	SG json.RawMessage     `json:"signers"`
}

func (doc *DocumentData) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var udoc DocumentJSONUnpacker
	if err := enc.Unmarshal(b, &udoc); err != nil {
		return err
	}

	return doc.unpack(enc, udoc.FH, udoc.CR, udoc.OW, udoc.SG)
}
