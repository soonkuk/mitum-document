package cmds

import (
	"github.com/soonkuk/mitum-blocksign/blocksign"
)

type FileHashFlag struct {
	FH blocksign.FileHash
}

func (v *FileHashFlag) UnmarshalText(b []byte) error {
	fh := blocksign.FileHash(string(b))
	if err := fh.IsValid(nil); err != nil {
		return err
	}
	v.FH = fh

	return nil
}

func (v *FileHashFlag) String() string {
	return v.FH.String()
}
