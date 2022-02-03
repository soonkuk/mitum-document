package digest

import (
	"github.com/soonkuk/mitum-blocksign/blocksign"
	"github.com/soonkuk/mitum-blocksign/document"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util/hint"
)

var (
	BlocksignDocumentValueType = hint.Type("mitum-blocksign-document-value")
	BlocksignDocumentValueHint = hint.NewHint(BlocksignDocumentValueType, "v0.0.1")
)

type BlocksignDocumentValue struct {
	doc    blocksign.DocumentData
	height base.Height
}

func NewBlocksignDocumentValue(
	doc blocksign.DocumentData,
	height base.Height,
) BlocksignDocumentValue {

	return BlocksignDocumentValue{
		doc:    doc,
		height: height,
	}
}

func (dv BlocksignDocumentValue) Hint() hint.Hint {
	return BlocksignDocumentValueHint
}

func (dv BlocksignDocumentValue) Document() blocksign.DocumentData {
	return dv.doc
}

func (dv BlocksignDocumentValue) Height() base.Height {
	return dv.height
}

var (
	BlockcityDocumentValueType = hint.Type("mitum-blockcity-document-value")
	BlockcityDocumentValueHint = hint.NewHint(BlockcityDocumentValueType, "v0.0.1")
)

type BlockcityDocumentValue struct {
	doc    document.DocumentData
	height base.Height
}

func NewBlockcityDocumentValue(
	doc document.DocumentData,
	height base.Height,
) BlockcityDocumentValue {

	return BlockcityDocumentValue{
		doc:    doc,
		height: height,
	}
}

func (dv BlockcityDocumentValue) Hint() hint.Hint {
	return BlockcityDocumentValueHint
}

func (dv BlockcityDocumentValue) Document() document.DocumentData {
	return dv.doc
}

func (dv BlockcityDocumentValue) Height() base.Height {
	return dv.height
}
