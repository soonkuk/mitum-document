// +build test

package blocksign

import "github.com/spikeekips/mitum-currency/currency"

func MustAddress(s string) currency.Address {
	a, err := currency.NewAddress(s)
	if err != nil {
		panic(err)
	}

	return a
}
