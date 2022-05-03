package extension

import (
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/base/state"
)

func checkFactSignsByState(
	address base.Address,
	fs []base.FactSign,
	getState func(string) (state.State, bool, error),
) error {
	st, err := existsState(currency.StateKeyAccount(address), "keys of account", getState)
	if err != nil {
		return err
	}
	keys, err := currency.StateKeysValue(st)
	if err != nil {
		return operation.NewBaseReasonErrorFromError(err)
	}

	if err := checkThreshold(fs, keys); err != nil {
		return operation.NewBaseReasonErrorFromError(err)
	}

	return nil
}
