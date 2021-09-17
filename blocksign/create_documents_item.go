package blocksign

import (
	"github.com/pkg/errors"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
)

type BaseCreateDocumentsItem struct {
	hint       hint.Hint
	fileHash   FileHash
	documentid currency.Big
	signcode   string //creator signcode
	title      string
	size       currency.Big
	signers    []base.Address
	signcodes  []string //signers signcode
	cid        currency.CurrencyID
}

func NewBaseCreateDocumentsItem(ht hint.Hint,
	filehash FileHash,
	documentid currency.Big,
	signcode, title string,
	size currency.Big,
	signers []base.Address,
	signcodes []string,
	cid currency.CurrencyID) BaseCreateDocumentsItem {
	return BaseCreateDocumentsItem{
		hint:       ht,
		fileHash:   filehash,
		documentid: documentid,
		signcode:   signcode,
		title:      title,
		size:       size,
		signers:    signers,
		signcodes:  signcodes,
		cid:        cid,
	}
}

func (it BaseCreateDocumentsItem) Hint() hint.Hint {
	return it.hint
}

func (it BaseCreateDocumentsItem) Bytes() []byte {
	bs := make([][]byte, len(it.signers)+len(it.signcodes)+6)
	bs[0] = it.fileHash.Bytes()
	bs[1] = it.documentid.Bytes()
	bs[2] = []byte(it.signcode)
	bs[3] = []byte(it.title)
	bs[4] = it.size.Bytes()
	bs[5] = it.cid.Bytes()
	for i := range it.signers {
		bs[i+6] = it.signers[i].Bytes()
	}
	for i := range it.signcodes {
		bs[i+len(it.signers)+6] = []byte(it.signcodes[i])
	}

	return util.ConcatBytesSlice(bs...)
}

func (it BaseCreateDocumentsItem) IsValid([]byte) error {
	if len(it.fileHash) < 1 {
		return errors.Errorf("empty fileHash")
	}
	if (it.documentid == currency.Big{}) {
		return errors.Errorf("empty documentid")
	}
	if !it.documentid.OverZero() {
		return errors.Errorf("documentid is negative number")
	}
	if len(it.signcode) < 1 {
		return errors.Errorf("empty creator signcode")
	}
	if err := it.cid.IsValid(nil); err != nil {
		return err
	}
	if len(it.signers) != len(it.signcodes) {
		return errors.Errorf("length of signers array is not same with length of signcodes array")
	}
	return nil
}

// FileHash return BaseCreateDocumetsItem's owner address.
func (it BaseCreateDocumentsItem) FileHash() FileHash {
	return it.fileHash
}

func (it BaseCreateDocumentsItem) DocumentId() currency.Big {
	return it.documentid
}

func (it BaseCreateDocumentsItem) Signcode() string {
	return it.signcode
}

func (it BaseCreateDocumentsItem) Title() string {
	return it.title
}

func (it BaseCreateDocumentsItem) Size() currency.Big {
	return it.size
}

func (it BaseCreateDocumentsItem) Signers() []base.Address {
	return it.signers
}

func (it BaseCreateDocumentsItem) Signcodes() []string {
	return it.signcodes
}

// FileData return BaseCreateDocumentsItem's fileData.
func (it BaseCreateDocumentsItem) Currency() currency.CurrencyID {
	return it.cid
}

func (it BaseCreateDocumentsItem) Rebuild() CreateDocumentsItem {
	return it
}
