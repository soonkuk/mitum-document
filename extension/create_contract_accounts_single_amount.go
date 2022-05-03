package extension

import (
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
)

var (
	CreateContractAccountsItemSingleAmountType   = hint.Type("mitum-currency-create-contract-accounts-single-amount")
	CreateContractAccountsItemSingleAmountHint   = hint.NewHint(CreateContractAccountsItemSingleAmountType, "v0.0.1")
	CreateContractAccountsItemSingleAmountHinter = CreateContractAccountsItemSingleAmount{
		BaseCreateContractAccountsItem: BaseCreateContractAccountsItem{
			BaseHinter: hint.NewBaseHinter(CreateContractAccountsItemSingleAmountHint),
		},
	}
)

type CreateContractAccountsItemSingleAmount struct {
	BaseCreateContractAccountsItem
}

func NewCreateContractAccountsItemSingleAmount(keys currency.AccountKeys, amount currency.Amount) CreateContractAccountsItemSingleAmount {
	return CreateContractAccountsItemSingleAmount{
		BaseCreateContractAccountsItem: NewBaseCreateContractAccountsItem(CreateContractAccountsItemSingleAmountHint, keys, []currency.Amount{amount}),
	}
}

func (it CreateContractAccountsItemSingleAmount) IsValid([]byte) error {
	if err := it.BaseCreateContractAccountsItem.IsValid(nil); err != nil {
		return err
	}

	if n := len(it.amounts); n != 1 {
		return isvalid.InvalidError.Errorf("only one amount allowed; %d", n)
	}

	return nil
}

func (it CreateContractAccountsItemSingleAmount) Rebuild() CreateContractAccountsItem {
	it.BaseCreateContractAccountsItem = it.BaseCreateContractAccountsItem.Rebuild().(BaseCreateContractAccountsItem)

	return it
}
