// +build test

package blocksign

import "github.com/soonkuk/mitum-data/currency"

func MustAddress(s string) currency.Address {
	a, err := currency.NewAddress(s)
	if err != nil {
		panic(err)
	}

	return a
}
