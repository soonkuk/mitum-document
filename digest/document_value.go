package digest

import (
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util/hint"

	"github.com/soonkuk/mitum-data/blocksign"
)

var (
	DocumentValueType = hint.Type("mitum-blocksign-document-value")
	DocumentValueHint = hint.NewHint(DocumentValueType, "v0.0.1")
)

type DocumentValue struct {
	doc    blocksign.DocumentData
	height base.Height
}

func NewDocumentValue(
	doc blocksign.DocumentData,
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

func (dv DocumentValue) Document() blocksign.DocumentData {
	return dv.doc
}

func (dv DocumentValue) Height() base.Height {
	return dv.height
}
