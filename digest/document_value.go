package digest

/*
import (
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/state"
	"github.com/spikeekips/mitum/util/hint"
	"golang.org/x/xerrors"

	"github.com/soonkuk/mitum-data/blocksign"
	"github.com/soonkuk/mitum-data/currency"
)

var (
	DocumentValueType = hint.Type("mitum-blocksign-document-value")
	DocumentValueHint = hint.NewHint(DocumentValueType, "v0.0.1")
)

type DocumentValue struct {
	ac             currency.Account
	filedata       blocksign.FileData
	height         base.Height
	previousHeight base.Height
}

func NewDocumentValue(st state.State) (DocumentValue, error) {
	var ac currency.Account
	switch a, ok, err := IsDocumentState(st); {
	case err != nil:
		return DocumentValue{}, err
	case !ok:
		return DocumentValue{}, xerrors.Errorf("not state for currency.Account, %T", st.Value().Interface())
	default:
		ac = a
	}

	return DocumentValue{
		ac:             ac,
		height:         st.Height(),
		previousHeight: st.PreviousHeight(),
	}, nil
}

func (dv DocumentValue) Hint() hint.Hint {
	return DocumentValueHint
}

func (dv DocumentValue) Account() currency.Account {
	return dv.ac
}

func (dv DocumentValue) FileData() blocksign.FileData {
	return dv.filedata
}

func (dv DocumentValue) Height() base.Height {
	return dv.height
}

func (dv DocumentValue) SetHeight(height base.Height) DocumentValue {
	dv.height = height

	return dv
}

func (dv DocumentValue) PreviousHeight() base.Height {
	return dv.previousHeight
}

func (dv DocumentValue) SetPreviousHeight(height base.Height) DocumentValue {
	dv.previousHeight = height

	return dv
}

func (dv DocumentValue) SetFileData(filedata blocksign.FileData) DocumentValue {
	dv.filedata = filedata

	return dv
}
*/
