package blocksign

import (
	"fmt"
	"strings"

	"github.com/soonkuk/mitum-data/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/state"
	"github.com/spikeekips/mitum/util"
	"golang.org/x/xerrors"
)

var (
	StateKeyDocumentsSuffix    = ":documents"
	StateKeyDocumentDataSuffix = ":documentData"
	StateKeyLastDocumentId     = "lastdocumentId"
)

func StateLastDocumentIdValue(st state.State) (DocInfo, error) {
	v := st.Value()
	if v == nil {
		return DocInfo{}, util.NotFoundError.Errorf("document id not found in State")
	}

	if s, ok := v.Interface().(DocInfo); !ok {
		return DocInfo{}, xerrors.Errorf("invalid document id value found, %T", v.Interface())
	} else {
		return s, nil
	}
}

func SetStateLastDocumentIdValue(st state.State, v DocInfo) (state.State, error) {
	if uv, err := state.NewHintedValue(v); err != nil {
		return nil, err
	} else {
		return st.SetValue(uv)
	}
}

func StateKeyDocumentData(fh FileHash) string {
	return fmt.Sprintf("%s%s", fh.String(), StateKeyDocumentDataSuffix)
}

func IsStateDocumentDataKey(key string) bool {
	return strings.HasSuffix(key, StateKeyDocumentDataSuffix)
}

func StateDocumentDataValue(st state.State) (DocumentData, error) {
	v := st.Value()
	if v == nil {
		return DocumentData{}, util.NotFoundError.Errorf("document data not found in State")
	}

	if s, ok := v.Interface().(DocumentData); !ok {
		return DocumentData{}, xerrors.Errorf("invalid document data value found, %T", v.Interface())
	} else {
		return s, nil
	}
}

func SetStateDocumentDataValue(st state.State, v DocumentData) (state.State, error) {
	if uv, err := state.NewHintedValue(v); err != nil {
		return nil, err
	} else {
		return st.SetValue(uv)
	}
}

func StateKeyDocuments(a base.Address) string {
	return fmt.Sprintf("%s%s", currency.StateAddressKeyPrefix(a), StateKeyDocumentsSuffix)
}

func IsStateDocumentsKey(key string) bool {
	return strings.HasSuffix(key, StateKeyDocumentsSuffix)
}

func StateDocumentsValue(st state.State) (DocumentInventory, error) {
	v := st.Value()
	if v == nil {
		return DocumentInventory{}, util.NotFoundError.Errorf("document inventory not found in State")
	}

	if s, ok := v.Interface().(DocumentInventory); !ok {
		return DocumentInventory{}, xerrors.Errorf("invalid document inventory value found, %T", v.Interface())
	} else {
		return s, nil
	}
}

func SetStateDocumentsValue(st state.State, v DocumentInventory) (state.State, error) {
	if uv, err := state.NewHintedValue(v); err != nil {
		return nil, err
	} else {
		return st.SetValue(uv)
	}
}
