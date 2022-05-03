package extension

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/base/state"
	"github.com/spikeekips/mitum/util"
)

var (
	StateKeyContractAccountOwnerSuffix = ":contractaccountowner"
	// StateKeyBalanceSuffix              = ":exntensionbalance"
)

func StateKeyContractAccountOwner(a base.Address) string {
	return fmt.Sprintf("%s%s", a.String(), StateKeyContractAccountOwnerSuffix)
}

func IsStateContractAccountOwnerKey(key string) bool {
	return strings.HasSuffix(key, StateKeyContractAccountOwnerSuffix)
}

func StateContractAccountOwnerValue(st state.State) (currency.Account, error) {
	v := st.Value()
	if v == nil {
		return currency.Account{}, util.NotFoundError.Errorf("account not found in State")
	}

	s, ok := v.Interface().(currency.Account)
	if !ok {
		return currency.Account{}, errors.Errorf("invalid account value found, %T", v.Interface())
	}
	return s, nil
}

func SetStateContractAccountOwnerValue(st state.State, v currency.Account) (state.State, error) {
	uv, err := state.NewHintedValue(v)
	if err != nil {
		return nil, err
	}
	return st.SetValue(uv)
}

func SetStateContractAccountKeysValue(st state.State) (state.State, error) {
	var ac currency.Account

	v := NewContractAccountKeys()

	if a, err := currency.LoadStateAccountValue(st); err != nil {
		if !errors.Is(err, util.NotFoundError) {
			return nil, err
		}

		n, err := currency.NewAccountFromKeys(v)
		if err != nil {
			return nil, err
		}
		ac = n
	} else {
		ac = a
	}

	if uac, err := ac.SetKeys(v); err != nil {
		return nil, err
	} else if uv, err := state.NewHintedValue(uac); err != nil {
		return nil, err
	} else {
		return st.SetValue(uv)
	}
}

/*
func StateKeyBalance(a base.Address, cid currency.CurrencyID) string {
	return fmt.Sprintf("%s%s", currency.StateBalanceKeyPrefix(a, cid), StateKeyBalanceSuffix)
}

func IsStateBalanceKey(key string) bool {
	return strings.HasSuffix(key, StateKeyBalanceSuffix)
}
*/

func checkExistsState(
	key string,
	getState func(key string) (state.State, bool, error),
) error {
	switch _, found, err := getState(key); {
	case err != nil:
		return err
	case !found:
		return operation.NewBaseReasonError("state, %q does not exist", key)
	default:
		return nil
	}
}

func existsState(
	k,
	name string,
	getState func(key string) (state.State, bool, error),
) (state.State, error) {
	switch st, found, err := getState(k); {
	case err != nil:
		return nil, err
	case !found:
		return nil, operation.NewBaseReasonError("%s does not exist", name)
	default:
		return st, nil
	}
}

func notExistsState(
	k,
	name string,
	getState func(key string) (state.State, bool, error),
) (state.State, error) {
	switch st, found, err := getState(k); {
	case err != nil:
		return nil, err
	case found:
		return nil, operation.NewBaseReasonError("%s already exists", name)
	default:
		return st, nil
	}
}
