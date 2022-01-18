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
	AccountType   = hint.Type("mitum-currency-extended-account")
	AccountHint   = hint.NewHint(AccountType, "v0.0.1")
	AccountHinter = Account{BaseHinter: hint.NewBaseHinter(AccountHint)}
)

type Account struct {
	hint.BaseHinter
	currency.Account
	h          valuehash.Hash
	permission AccountPermission
}

func NewAccount(address base.Address, keys currency.AccountKeys) (Account, error) {
	if err := address.IsValid(nil); err != nil {
		return Account{}, err
	}
	if keys != nil {
		if err := keys.IsValid(nil); err != nil {
			return Account{}, err
		}
	}

	cac, err := currency.NewAccount(address, keys)
	if err != nil {
		return Account{}, err
	}
	ac := Account{BaseHinter: hint.NewBaseHinter(AccountHint), Account: cac}
	ac.h = ac.GenerateHash()

	return ac, nil
}

func NewAccountFromKeys(keys currency.AccountKeys) (Account, error) {
	if a, err := currency.NewAddressFromKeys(keys); err != nil {
		return Account{}, err
	} else if ac, err := NewAccount(a, keys); err != nil {
		return Account{}, err
	} else {
		return ac, nil
	}
}

func (ac Account) Bytes() []byte {
	bs := make([][]byte, 3)
	bs[0] = ac.Account.Address().Bytes()

	if ac.Keys() != nil {
		bs[1] = ac.Keys().Bytes()
	}

	bs[2] = ac.permission.Bytes()

	return util.ConcatBytesSlice(bs...)
}

func (ac Account) Hash() valuehash.Hash {
	return ac.h
}

func (ac Account) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(ac.Bytes())
}

func (ac Account) SetKeys(keys currency.AccountKeys) (Account, error) {
	if err := keys.IsValid(nil); err != nil {
		return Account{}, err
	}

	cac, err := ac.Account.SetKeys(keys)
	if err != nil {
		return Account{}, err
	}

	ac.Account = cac

	return ac, nil
}

func (ac Account) SetPermission(ap AccountPermission) (Account, error) {
	if err := ap.IsValid(nil); err != nil {
		return Account{}, err
	}

	ac.permission = ap
	ac.h = ac.GenerateHash()

	return ac, nil
}

func (ac Account) IsEmpty() bool {
	return ac.Account.IsEmpty() || ac.h.IsEmpty()
}

type AccountPermission uint

const (
	blockcityAdmin AccountPermission = 1 + iota
	blocksignAdmin
)

var accountPermissions = map[AccountPermission]string{
	blockcityAdmin: "blockcityAdmin",
	blocksignAdmin: "blocksignAdmin",
}

func (ap AccountPermission) String() string {
	return accountPermissions[ap]
}

func (ap AccountPermission) Bytes() []byte {
	return util.UintToBytes(uint(ap))
}

func (ap AccountPermission) IsValid([]byte) error {
	if uint(ap) > uint(len(accountPermissions)) {
		return isvalid.InvalidError.Errorf("invalid AccountPermission")
	}
	return nil
}
