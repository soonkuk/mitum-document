package blocksign

import (
	"encoding/json"

	"github.com/spikeekips/mitum/base"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
)

type DocumentJSONPacker struct {
	jsonenc.HintedHead
	FH FileHash     `json:"filehash"`
	ID DocId        `json:"documentid"`
	CR base.Address `json:"creator"`
	SG []DocSign    `json:"signers"`
}

func (doc DocumentData) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(DocumentJSONPacker{
		HintedHead: jsonenc.NewHintedHead(doc.Hint()),
		FH:         doc.fileHash,
		ID:         doc.id,
		CR:         doc.creator,
		SG:         doc.signers,
	})
}

type DocumentJSONUnpacker struct {
	FH string              `json:"filehash"`
	ID json.RawMessage     `json:"documentid"`
	CR base.AddressDecoder `json:"creator"`
	SG json.RawMessage     `json:"signers"`
}

func (doc *DocumentData) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var udoc DocumentJSONUnpacker
	if err := enc.Unmarshal(b, &udoc); err != nil {
		return err
	}

	return doc.unpack(enc, udoc.FH, udoc.ID, udoc.CR, udoc.SG)
}
