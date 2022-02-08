package digest

import (
	"github.com/soonkuk/mitum-blocksign/blocksign"
	"github.com/soonkuk/mitum-blocksign/document"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util/hint"
)

var (
	BSDocumentValueType = hint.Type("mitum-blocksign-document-value")
	BSDocumentValueHint = hint.NewHint(BSDocumentValueType, "v0.0.1")
)

type BSDocumentValue struct {
	doc    blocksign.DocumentData
	height base.Height
}

func NewBSDocumentValue(
	doc blocksign.DocumentData,
	height base.Height,
) BSDocumentValue {

	return BSDocumentValue{
		doc:    doc,
		height: height,
	}
}

func (dv BSDocumentValue) Hint() hint.Hint {
	return BSDocumentValueHint
}

func (dv BSDocumentValue) Document() blocksign.DocumentData {
	return dv.doc
}

func (dv BSDocumentValue) Height() base.Height {
	return dv.height
}

var (
	BCDocumentValueType = hint.Type("mitum-blockcity-document-value")
	BCDocumentValueHint = hint.NewHint(BCDocumentValueType, "v0.0.1")
)

type BCDocumentValue struct {
	doc    document.DocumentData
	height base.Height
}

func NewBCDocumentValue(
	doc document.DocumentData,
	height base.Height,
) BCDocumentValue {

	return BCDocumentValue{
		doc:    doc,
		height: height,
	}
}

func (dv BCDocumentValue) Hint() hint.Hint {
	return BCDocumentValueHint
}

func (dv BCDocumentValue) Document() document.DocumentData {
	return dv.doc
}

func (dv BCDocumentValue) Height() base.Height {
	return dv.height
}
