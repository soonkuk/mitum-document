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
	StateKeyDocumentSuffix         = ":document"
	StateKeyDocumentDataExclSuffix = ":documentDataExcl"
	StateKeyDocumentDataSuffix     = ":documentData"
	//StateKeyFileIDSuffix   = ":fileid"
	//StateKeySignCodeSuffix = ":signcode"
	//StateKeyOwnerSuffix    = ":owner"
	StateKeyLastDocumentId = "lastId:document"
)

/*
func StateKeyDocument() string {
	return fmt.Sprintf("%s%s", "lastestId", StateKeyDocumentSuffix)
}
*/
func StateLastDocumentIdValue(st state.State) (DocId, error) {
	v := st.Value()
	if v == nil {
		return DocId{}, util.NotFoundError.Errorf("document id not found in State")
	}

	if s, ok := v.Interface().(DocId); !ok {
		return DocId{}, xerrors.Errorf("invalid document id value found, %T", v.Interface())
	} else {
		return s, nil
	}
}

func SetStateLastDocumentIdValue(st state.State, v DocId) (state.State, error) {
	if uv, err := state.NewHintedValue(v); err != nil {
		return nil, err
	} else {
		return st.SetValue(uv)
	}
}

func IsStateDocumentDataExclKey(key string) bool {
	return strings.HasSuffix(key, StateKeyDocumentDataExclSuffix)
}

func StateKeyDocumentDataExcl(fh FileHash) string {
	return fmt.Sprintf("%s%s", fh.String(), StateKeyDocumentDataSuffix)
}

func StateDocumentDataKeyPrefix(a base.Address, di DocId) string {
	return fmt.Sprintf("%s-%s", currency.StateAddressKeyPrefix(a), di)
}

func StateKeyDocumentData(a base.Address, di DocId) string {
	return fmt.Sprintf("%s%s", StateDocumentDataKeyPrefix(a, di), StateKeyDocumentDataSuffix)
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

/*
func StateFileIDKeyPrefix(a base.Address, fid FileID) string {
	return fmt.Sprintf("%s-%s", currency.StateAddressKeyPrefix(a), fid)
}

func StateKeyFileID(a base.Address, fid FileID) string {
	return fmt.Sprintf("%s%s", StateFileIDKeyPrefix(a, fid), StateKeyFileIDSuffix)
}

func IsStateFileIDKey(key string) bool {
	return strings.HasSuffix(key, StateKeyFileIDSuffix)
}

func StateFileIDValue(st state.State) (FileID, error) {
	v := st.Value()
	if v == nil {
		return FileID(""), util.NotFoundError.Errorf("filedata not found in State")
	}

	if s, ok := v.Interface().(FileID); !ok {
		return FileID(""), xerrors.Errorf("invalid filedata value found, %T", v.Interface())
	} else {
		return s, nil
	}
}

func SetStateFileIDValue(st state.State, v FileID) (state.State, error) {
	if uv, err := state.NewHintedValue(v); err != nil {
		return nil, err
	} else {
		return st.SetValue(uv)
	}
}

func StateSignCodeKeyPrefix(a base.Address, sc SignCode) string {
	return fmt.Sprintf("%s-%s", currency.StateAddressKeyPrefix(a), sc)
}

func StateKeySignCode(a base.Address, sc SignCode) string {
	return fmt.Sprintf("%s%s", StateSignCodeKeyPrefix(a, sc), StateKeySignCodeSuffix)
}

func IsStateSignCodeKey(key string) bool {
	return strings.HasSuffix(key, StateKeySignCodeSuffix)
}

func StateSignCodeValue(st state.State) (SignCode, error) {
	v := st.Value()
	if v == nil {
		return SignCode(""), util.NotFoundError.Errorf("filedata not found in State")
	}

	if s, ok := v.Interface().(SignCode); !ok {
		return SignCode(""), xerrors.Errorf("invalid filedata value found, %T", v.Interface())
	} else {
		return s, nil
	}
}

func SetStateSignCodeValue(st state.State, v SignCode) (state.State, error) {
	if uv, err := state.NewHintedValue(v); err != nil {
		return nil, err
	} else {
		return st.SetValue(uv)
	}
}

func StateOwnerKeyPrefix(a base.Address, owner base.Address) string {
	return fmt.Sprintf("%s-%s", currency.StateAddressKeyPrefix(a), owner)
}

func StateKeyOwner(a base.Address, owner base.Address) string {
	return fmt.Sprintf("%s%s", StateOwnerKeyPrefix(a, owner), StateKeyOwnerSuffix)
}

func IsStateOwnerKey(key string) bool {
	return strings.HasSuffix(key, StateKeyOwnerSuffix)
}

func StateOwnerValue(st state.State) (base.Address, error) {
	v := st.Value()
	if v == nil {
		return currency.EmptyAddress, util.NotFoundError.Errorf("filedata not fousnd in State")
	}

	if s, ok := v.Interface().(base.Address); !ok {
		return currency.EmptyAddress, xerrors.Errorf("invalid filedata value found, %T", v.Interface())
	} else {
		return s, nil
	}
}

func SetStateOwnerValue(st state.State, v base.Address) (state.State, error) {
	if uv, err := state.NewHintedValue(v); err != nil {
		return nil, err
	} else {
		return st.SetValue(uv)
	}
}
*/
