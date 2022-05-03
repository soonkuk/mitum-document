package extension

import (
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
)

type BaseCreateContractAccountsItem struct {
	hint.BaseHinter
	keys    currency.AccountKeys
	amounts []currency.Amount
}

func NewBaseCreateContractAccountsItem(ht hint.Hint, keys currency.AccountKeys, amounts []currency.Amount) BaseCreateContractAccountsItem {
	return BaseCreateContractAccountsItem{
		BaseHinter: hint.NewBaseHinter(ht),
		keys:       keys,
		amounts:    amounts,
	}
}

func (it BaseCreateContractAccountsItem) Bytes() []byte {
	bs := make([][]byte, len(it.amounts)+1)
	bs[0] = it.keys.Bytes()

	for i := range it.amounts {
		bs[i+1] = it.amounts[i].Bytes()
	}

	return util.ConcatBytesSlice(bs...)
}

func (it BaseCreateContractAccountsItem) IsValid([]byte) error {
	if n := len(it.amounts); n == 0 {
		return isvalid.InvalidError.Errorf("empty amounts")
	}

	if err := isvalid.Check(nil, false, it.BaseHinter, it.keys); err != nil {
		return err
	}

	founds := map[currency.CurrencyID]struct{}{}
	for i := range it.amounts {
		am := it.amounts[i]
		if _, found := founds[am.Currency()]; found {
			return isvalid.InvalidError.Errorf("duplicated currency found, %q", am.Currency())
		}
		founds[am.Currency()] = struct{}{}

		if err := am.IsValid(nil); err != nil {
			return err
		} else if !am.Big().OverZero() {
			return isvalid.InvalidError.Errorf("amount should be over zero")
		}
	}

	return nil
}

func (it BaseCreateContractAccountsItem) Keys() currency.AccountKeys {
	return it.keys
}

func (it BaseCreateContractAccountsItem) Address() (base.Address, error) {
	return currency.NewAddressFromKeys(it.keys)
}

func (it BaseCreateContractAccountsItem) Amounts() []currency.Amount {
	return it.amounts
}

func (it BaseCreateContractAccountsItem) Rebuild() CreateContractAccountsItem {
	ams := make([]currency.Amount, len(it.amounts))
	for i := range it.amounts {
		am := it.amounts[i]
		ams[i] = am.WithBig(am.Big())
	}

	it.amounts = ams

	return it
}
