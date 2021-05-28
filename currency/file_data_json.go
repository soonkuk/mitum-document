package currency

import (
	"github.com/spikeekips/mitum/base"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
)

type FileDataJSONPacker struct {
	jsonenc.HintedHead
	US SignCode     `json:"signcode"`
	OW base.Address `json:"owner"`
}

func (fd FileData) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(FileDataJSONPacker{
		HintedHead: jsonenc.NewHintedHead(fd.Hint()),
		US:         fd.signcode,
		OW:         fd.owner,
	})
}

type FileDataJSONUnpacker struct {
	US string              `json:"signcode"`
	OW base.AddressDecoder `json:"owner"`
}

func (fd *FileData) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var ufd FileDataJSONUnpacker
	if err := enc.Unmarshal(b, &ufd); err != nil {
		return err
	}

	return fd.unpack(enc, ufd.US, ufd.OW)
}
