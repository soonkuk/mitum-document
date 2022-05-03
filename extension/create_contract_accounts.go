package extension

import (
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/isvalid"
	"github.com/spikeekips/mitum/util/valuehash"
)

var (
	CreateContractAccountsFactType   = hint.Type("mitum-currency-create-contract-accounts-operation-fact")
	CreateContractAccountsFactHint   = hint.NewHint(CreateContractAccountsFactType, "v0.0.1")
	CreateContractAccountsFactHinter = CreateContractAccountsFact{BaseHinter: hint.NewBaseHinter(CreateContractAccountsFactHint)}
	CreateContractAccountsType       = hint.Type("mitum-currency-create-contract-accounts-operation")
	CreateContractAccountsHint       = hint.NewHint(CreateContractAccountsType, "v0.0.1")
	CreateContractAccountsHinter     = CreateContractAccounts{BaseOperation: operationHinter(CreateContractAccountsHint)}
)

var MaxCreateContractAccountsItems uint = 10

type AmountsItem interface {
	Amounts() []currency.Amount
}

type CreateContractAccountsItem interface {
	hint.Hinter
	isvalid.IsValider
	AmountsItem
	Bytes() []byte
	Keys() currency.AccountKeys
	Address() (base.Address, error)
	Rebuild() CreateContractAccountsItem
}

type CreateContractAccountsFact struct {
	hint.BaseHinter
	h      valuehash.Hash
	token  []byte
	sender base.Address
	items  []CreateContractAccountsItem
}

func NewCreateContractAccountsFact(token []byte, owner base.Address, items []CreateContractAccountsItem) CreateContractAccountsFact {
	fact := CreateContractAccountsFact{
		BaseHinter: hint.NewBaseHinter(CreateContractAccountsFactHint),
		token:      token,
		sender:     owner,
		items:      items,
	}
	fact.h = fact.GenerateHash()

	return fact
}

func (fact CreateContractAccountsFact) Hash() valuehash.Hash {
	return fact.h
}

func (fact CreateContractAccountsFact) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact CreateContractAccountsFact) Bytes() []byte {
	is := make([][]byte, len(fact.items))
	for i := range fact.items {
		is[i] = fact.items[i].Bytes()
	}

	return util.ConcatBytesSlice(
		fact.token,
		fact.sender.Bytes(),
		util.ConcatBytesSlice(is...),
	)
}

func (fact CreateContractAccountsFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := currency.IsValidOperationFact(fact, b); err != nil {
		return err
	}

	if n := len(fact.items); n < 1 {
		return isvalid.InvalidError.Errorf("empty items")
	} else if n > int(MaxCreateContractAccountsItems) {
		return isvalid.InvalidError.Errorf("items, %d over max, %d", n, MaxCreateContractAccountsItems)
	}

	if err := isvalid.Check(nil, false, fact.sender); err != nil {
		return err
	}

	foundKeys := map[string]struct{}{}
	for i := range fact.items {
		if err := isvalid.Check(nil, false, fact.items[i]); err != nil {
			return err
		}

		it := fact.items[i]
		k := it.Keys().Hash().String()
		if _, found := foundKeys[k]; found {
			return isvalid.InvalidError.Errorf("duplicated acocunt Keys found, %s", k)
		}

		switch a, err := it.Address(); {
		case err != nil:
			return err
		case fact.sender.Equal(a):
			return isvalid.InvalidError.Errorf("target address is same with sender, %q", fact.sender)
		default:
			foundKeys[k] = struct{}{}
		}
	}

	return nil
}

func (fact CreateContractAccountsFact) Token() []byte {
	return fact.token
}

func (fact CreateContractAccountsFact) Sender() base.Address {
	return fact.sender
}

func (fact CreateContractAccountsFact) Items() []CreateContractAccountsItem {
	return fact.items
}

func (fact CreateContractAccountsFact) Targets() ([]base.Address, error) {
	as := make([]base.Address, len(fact.items))
	for i := range fact.items {
		a, err := fact.items[i].Address()
		if err != nil {
			return nil, err
		}
		as[i] = a
	}

	return as, nil
}

func (fact CreateContractAccountsFact) Addresses() ([]base.Address, error) {
	as := make([]base.Address, len(fact.items)+1)

	tas, err := fact.Targets()
	if err != nil {
		return nil, err
	}
	copy(as, tas)

	as[len(fact.items)] = fact.Sender()

	return as, nil
}

func (fact CreateContractAccountsFact) Rebuild() CreateContractAccountsFact {
	items := make([]CreateContractAccountsItem, len(fact.items))
	for i := range fact.items {
		it := fact.items[i]
		items[i] = it.Rebuild()
	}

	fact.items = items
	fact.h = fact.GenerateHash()

	return fact
}

type CreateContractAccounts struct {
	currency.BaseOperation
}

func NewCreateContractAccounts(fact CreateContractAccountsFact, fs []base.FactSign, memo string) (CreateContractAccounts, error) {
	bo, err := currency.NewBaseOperationFromFact(CreateContractAccountsHint, fact, fs, memo)
	if err != nil {
		return CreateContractAccounts{}, err
	}

	return CreateContractAccounts{BaseOperation: bo}, nil
}
