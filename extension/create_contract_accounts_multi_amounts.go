package extension

import (
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
)

var maxCurenciesCreateContractAccountsItemMultiAmounts = 10

var (
	CreateContractAccountsItemMultiAmountsType   = hint.Type("mitum-currency-create-contract-accounts-multiple-amounts")
	CreateContractAccountsItemMultiAmountsHint   = hint.NewHint(CreateContractAccountsItemMultiAmountsType, "v0.0.1")
	CreateContractAccountsItemMultiAmountsHinter = CreateContractAccountsItemMultiAmounts{
		BaseCreateContractAccountsItem: BaseCreateContractAccountsItem{
			BaseHinter: hint.NewBaseHinter(CreateContractAccountsItemMultiAmountsHint),
		},
	}
)

type CreateContractAccountsItemMultiAmounts struct {
	BaseCreateContractAccountsItem
}

func NewCreateContractAccountsItemMultiAmounts(keys currency.AccountKeys, amounts []currency.Amount) CreateContractAccountsItemMultiAmounts {
	return CreateContractAccountsItemMultiAmounts{
		BaseCreateContractAccountsItem: NewBaseCreateContractAccountsItem(CreateContractAccountsItemMultiAmountsHint, keys, amounts),
	}
}

func (it CreateContractAccountsItemMultiAmounts) IsValid([]byte) error {
	if err := it.BaseCreateContractAccountsItem.IsValid(nil); err != nil {
		return err
	}

	if n := len(it.amounts); n > maxCurenciesCreateContractAccountsItemMultiAmounts {
		return isvalid.InvalidError.Errorf("amounts over allowed; %d > %d", n, maxCurenciesCreateContractAccountsItemMultiAmounts)
	}

	return nil
}

func (it CreateContractAccountsItemMultiAmounts) Rebuild() CreateContractAccountsItem {
	it.BaseCreateContractAccountsItem = it.BaseCreateContractAccountsItem.Rebuild().(BaseCreateContractAccountsItem)

	return it
}
