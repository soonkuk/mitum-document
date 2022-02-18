package digest

import (
	"github.com/protoconNet/mitum-document/document"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util/hint"
)

var (
	DocumentValueType = hint.Type("mitum-document-value")
	DocumentValueHint = hint.NewHint(DocumentValueType, "v0.0.1")
)

type DocumentValue struct {
	doc    document.DocumentData
	height base.Height
}

func NewDocumentValue(
	doc document.DocumentData,
	height base.Height,
) DocumentValue {

	return DocumentValue{
		doc:    doc,
		height: height,
	}
}

func (dv DocumentValue) Hint() hint.Hint {
	return DocumentValueHint
}

func (dv DocumentValue) Document() document.DocumentData {
	return dv.doc
}

func (dv DocumentValue) Height() base.Height {
	return dv.height
}
