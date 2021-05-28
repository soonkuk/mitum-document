package currency

import (
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"golang.org/x/xerrors"
)

type BaseCreateDocumentsItem struct {
	hint  hint.Hint
	keys  Keys
	sc    SignCode
	owner base.Address
	cid   CurrencyID
}

func NewBaseCreateDocumentsItem(ht hint.Hint, keys Keys, sc SignCode, owner base.Address, cid CurrencyID) BaseCreateDocumentsItem {
	return BaseCreateDocumentsItem{
		hint:  ht,
		keys:  keys,
		sc:    sc,
		owner: owner,
		cid:   cid,
	}
}

func (it BaseCreateDocumentsItem) Hint() hint.Hint {
	return it.hint
}

func (it BaseCreateDocumentsItem) Bytes() []byte {
	bs := make([][]byte, 4)
	bs[0] = it.keys.Bytes()
	bs[1] = it.sc.Bytes()
	bs[2] = it.owner.Bytes()
	bs[3] = it.cid.Bytes()

	return util.ConcatBytesSlice(bs...)
}

func (it BaseCreateDocumentsItem) IsValid([]byte) error {

	// empty key, duplicated key, threshold check
	if err := it.keys.IsValid(nil); err != nil {
		return err
	}

	if len(it.sc) < 1 || it.owner == EmptyAddress {
		return xerrors.Errorf("empty filedata")
	}

	// compare owner and signers
	// for i := range it.signers {
	// 	if it.signers[i] == it.owner {
	//		return xerrors.Errorf("document owner also found in signers, %v", it.signers[i])
	//	  }
	// }

	if err := it.sc.IsValid(nil); err != nil {
		return err
	} else if err := it.owner.IsValid(nil); err != nil {
		return err
	}

	return nil
}

// Keys return BaseCreateDocumentsItem's keys.
func (it BaseCreateDocumentsItem) Keys() Keys {
	return it.keys
}

// Address get address from BaseCreateDocumentsItem's keys and return it.
func (it BaseCreateDocumentsItem) Address() (base.Address, error) {
	return NewAddressFromKeys(it.keys)
}

// Owner return BaseCreateDocumetsItem's owner address.
func (it BaseCreateDocumentsItem) SignCode() SignCode {
	return it.sc
}

// Owner return BaseCreateDocumetsItem's owner address.
func (it BaseCreateDocumentsItem) Owner() base.Address {
	return it.owner
}

// FileData return BaseCreateDocumentsItem's fileData.
func (it BaseCreateDocumentsItem) Currency() CurrencyID {
	return it.cid
}

func (it BaseCreateDocumentsItem) Rebuild() CreateDocumentsItem {

	return it
}
