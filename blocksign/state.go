package blocksign

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/base/state"
	"github.com/spikeekips/mitum/util"
)

var (
	StateKeyDocumentsSuffix    = ":documents"
	StateKeyDocumentDataSuffix = ":documentData"
	StateKeyLastDocumentId     = "lastdocumentId"
)

func StateKeyDocumentData(documentid DocId) string {
	return fmt.Sprintf("%s%s", documentid.String(), StateKeyDocumentDataSuffix)
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
		return DocumentData{}, errors.Errorf("invalid document data value found, %T", v.Interface())
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
	return fmt.Sprintf("%s%s", a.String(), StateKeyDocumentsSuffix)
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
		return DocumentInventory{}, errors.Errorf("invalid document inventory value found, %T", v.Interface())
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

func checkExistsState(
	key string,
	getState func(key string) (state.State, bool, error),
) error {
	switch _, found, err := getState(key); {
	case err != nil:
		return err
	case !found:
		return operation.NewBaseReasonError("state, %q does not exist", key)
	default:
		return nil
	}
}

func existsState(
	k,
	name string,
	getState func(key string) (state.State, bool, error),
) (state.State, error) {
	switch st, found, err := getState(k); {
	case err != nil:
		return nil, err
	case !found:
		return nil, operation.NewBaseReasonError("%s does not exist", name)
	default:
		return st, nil
	}
}
