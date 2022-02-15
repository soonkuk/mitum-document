package cmds

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/soonkuk/mitum-blocksign/document"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util/encoder"
)

type AddressFlag struct {
	s string
}

func (v *AddressFlag) UnmarshalText(b []byte) error {
	v.s = string(b)

	return nil
}

func (v *AddressFlag) String() string {
	return v.s
}

func (v *AddressFlag) Encode(enc encoder.Encoder) (base.Address, error) {
	return base.DecodeAddressFromString(v.s, enc)
}

type FileHashFlag struct {
	FH document.FileHash
}

func (v *FileHashFlag) UnmarshalText(b []byte) error {
	fh := document.FileHash(string(b))
	if err := fh.IsValid(nil); err != nil {
		return err
	}
	v.FH = fh

	return nil
}

func (v *FileHashFlag) String() string {
	return v.FH.String()
}

type DocSignFlag struct {
	AD AddressFlag
	SC string
}

func (v *DocSignFlag) UnmarshalText(b []byte) error {

	docSign := strings.SplitN(string(b), ",", 2)
	if len(docSign) != 2 {
		return errors.Errorf(`wrong formatted; "<string address>,<string signcode>"`)
	}

	v.AD = AddressFlag{
		s: docSign[0],
	}

	v.SC = docSign[1]

	return nil
}

func (v *DocSignFlag) String() string {
	return v.AD.String()
}
