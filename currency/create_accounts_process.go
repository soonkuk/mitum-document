package currency

import (
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/base/state"
	"github.com/spikeekips/mitum/util/valuehash"
	"golang.org/x/xerrors"
)

func (CreateAccounts) Process(
	func(key string) (state.State, bool, error),
	func(valuehash.Hash, ...state.State) error,
) error {
	return nil
}

type CreateAccountsItemProcessor struct {
	cp   *CurrencyPool
	h    valuehash.Hash
	item CreateAccountsItem
	ns   state.State                // target acccount의 new account state. publicKey값을 다루는 state
	nb   map[CurrencyID]AmountState // target account의 new balance state 배열. Amount값을 다루는 state
}

// PreProcess는 item에 대한 전처리를 하며 target에 대한 empty account state와 amount state를 만든다.
func (opp *CreateAccountsItemProcessor) PreProcess(
	getState func(key string) (state.State, bool, error),
	_ func(valuehash.Hash, ...state.State) error,
) error {
	for i := range opp.item.Amounts() {
		am := opp.item.Amounts()[i]

		var policy CurrencyPolicy
		if opp.cp != nil {
			i, found := opp.cp.Policy(am.Currency())
			if !found {
				return xerrors.Errorf("currency not registered, %q", am.Currency())
			}
			policy = i
		}

		if am.Big().Compare(policy.NewAccountMinBalance()) < 0 {
			return xerrors.Errorf(
				"amount should be over minimum balance, %v < %v", am.Big(), policy.NewAccountMinBalance())
		}
	}

	target, err := opp.item.Address()
	if err != nil {
		return err
	}
	// Account Key 값으로 target account의 state를 가져온다. 없으면 empty state가 온다.
	st, err := notExistsState(StateKeyAccount(target), "keys of target", getState)
	if err != nil {
		return err
	}
	opp.ns = st

	nb := map[CurrencyID]AmountState{}
	// item의 amount들을 순환하며
	for i := range opp.item.Amounts() {
		am := opp.item.Amounts()[i]
		// target account의 Balance key로 state를 가져온다. 없으면 empyty state가 온다.
		b, _, err := getState(StateKeyBalance(target, am.Currency()))
		if err != nil {
			return err
		}
		// target account의 new balance에 currency별로 AmountState를 만든다.
		nb[am.Currency()] = NewAmountState(b, am.Currency())
	}

	opp.nb = nb

	return nil
}

// Process는 item에 대한 본처리를 하며 target account에 대한 처리후의 account state와 amount state를 만들어 배열로 반환한다.
func (opp *CreateAccountsItemProcessor) Process(
	_ func(key string) (state.State, bool, error),
	_ func(valuehash.Hash, ...state.State) error,
) ([]state.State, error) {
	nac, err := NewAccountFromKeys(opp.item.Keys())
	if err != nil {
		return nil, err
	}
	// state의 array인데 amountState 배열의 길이 + 1만큼 만든다.
	sts := make([]state.State, len(opp.item.Amounts())+1)
	st, err := SetStateAccountValue(opp.ns, nac)
	if err != nil {
		return nil, err
	}
	sts[0] = st

	// 나머지는 각각의 amount에 대한 state이다.
	// 중복되는 currencyID가 있으면 안될 것 같다.
	for i := range opp.item.Amounts() {
		am := opp.item.Amounts()[i]
		sts[i+1] = opp.nb[am.Currency()].Add(am.Big())
	}
	// target account에 대해 새롭게 생성되는 state를 만든다.
	return sts, nil
}

type CreateAccountsProcessor struct {
	cp *CurrencyPool
	CreateAccounts
	sb       map[CurrencyID]AmountState     // sender의 balance state
	ns       []*CreateAccountsItemProcessor //
	required map[CurrencyID][2]Big          // sender의 required amount(amount + fee)와 fee
}

func NewCreateAccountsProcessor(cp *CurrencyPool) GetNewProcessor {
	return func(op state.Processor) (state.Processor, error) {
		i, ok := op.(CreateAccounts)
		if !ok {
			return nil, xerrors.Errorf("not CreateAccounts, %T", op)
		}
		return &CreateAccountsProcessor{
			cp:             cp,
			CreateAccounts: i,
		}, nil
	}
}

// PreProcess는 Sender의 account state확인하고 fact의 item들로부터 required를 계산하고 sender의 처리전 amount state를 가져온다.
// item들에 대한 전처리를 통해서
func (opp *CreateAccountsProcessor) PreProcess(
	getState func(key string) (state.State, bool, error),
	setState func(valuehash.Hash, ...state.State) error,
) (state.Processor, error) {
	fact := opp.Fact().(CreateAccountsFact)

	// Sender account에 대한 Key값(address-HintType-"account")으로 account state 확인
	if err := checkExistsState(StateKeyAccount(fact.sender), getState); err != nil {
		return nil, err
	}

	// calculateItemsFee를 통해서 required를 계산함.
	// 함수이름이 calculateItemsFee이지만 사실 amount와 fee를 둘 다 계산함.
	if required, err := opp.calculateItemsFee(); err != nil {
		return nil, operation.NewBaseReasonError("failed to calculate fee: %w", err)
		// required amount와 required fee에 대하여 sender의 balance가 충분한지 검증만하고 sender의 현재 amount state를 돌려준다.
	} else if sb, err := CheckEnoughBalance(fact.sender, required, getState); err != nil {
		return nil, err
	} else {
		// 현재 processor에서 처리해야 하는 required amount와 required fee
		opp.required = required
		// sender의 처리 전 Amount state
		opp.sb = sb
	}

	// item에 대한 Processor를 생성한다.
	ns := make([]*CreateAccountsItemProcessor, len(fact.items))
	for i := range fact.items {
		c := &CreateAccountsItemProcessor{cp: opp.cp, h: opp.Hash(), item: fact.items[i]}
		// itemProcessor의 전처리는 먼저 실행한다.
		if err := c.PreProcess(getState, setState); err != nil {
			return nil, operation.NewBaseReasonErrorFromError(err)
		}
		// 전처리가 끝난 Processor를 ns에 저장한다.
		ns[i] = c
	}

	// sender의 fact sign을 확인한다.
	if err := checkFactSignsByState(fact.sender, opp.Signs(), getState); err != nil {
		return nil, operation.NewBaseReasonError("invalid signing: %w", err)
	}

	opp.ns = ns

	return opp, nil
}

func (opp *CreateAccountsProcessor) Process( // nolint:dupl
	getState func(key string) (state.State, bool, error),
	setState func(valuehash.Hash, ...state.State) error,
) error {
	fact := opp.Fact().(CreateAccountsFact)

	var sts []state.State // nolint:prealloc
	for i := range opp.ns {
		// itemProcessor의 본처리를 한다. 처리후 나오는 state 배열은 새로 생성되는 target의 state들(account state, amount state)이다.
		s, err := opp.ns[i].Process(getState, setState)
		if err != nil {
			return operation.NewBaseReasonError("failed to process create account item: %w", err)
		}
		sts = append(sts, s...)
	}

	for k := range opp.required {
		rq := opp.required[k]
		sts = append(sts, opp.sb[k].Sub(rq[0]).AddFee(rq[1]))
	}

	return setState(fact.Hash(), sts...)
}

func (opp *CreateAccountsProcessor) calculateItemsFee() (map[CurrencyID][2]Big, error) {
	fact := opp.Fact().(CreateAccountsFact)

	items := make([]AmountsItem, len(fact.items))
	for i := range fact.items {
		items[i] = fact.items[i]
	}

	return CalculateItemsFee(opp.cp, items)
}

func CalculateItemsFee(cp *CurrencyPool, items []AmountsItem) (map[CurrencyID][2]Big, error) {

	// 각 CurrencyID별로 length 2의 Big 배열을 만듬.
	required := map[CurrencyID][2]Big{}

	for i := range items {
		it := items[i]

		for j := range it.Amounts() {
			am := it.Amounts()[j]

			rq := [2]Big{ZeroBig, ZeroBig}
			// required에서 각 CurrencyID 별 [2]Big{}을 가져온다.
			// [2]Big{} 는 각각 {required amount, required fee} 이다.
			// amount에 대한 합산은 feer계산할 때 한꺼번에
			if k, found := required[am.Currency()]; found {
				rq = k
			}

			if cp == nil {
				required[am.Currency()] = [2]Big{rq[0].Add(am.Big()), rq[1]}

				continue
			}
			// CurrencyID에 대한 feeer를 가져와서
			feeer, found := cp.Feeer(am.Currency())
			if !found {
				return nil, xerrors.Errorf("unknown currency id found, %q", am.Currency())
			}
			// Fee 값을 가져온다. am.Big()을 넘기는 이유는 ratio일 경우가 있어서
			switch k, err := feeer.Fee(am.Big()); {
			case err != nil:
				return nil, err
			// Fee가 없거나 0보다 작으면 required amount는 더하고 required fee는 변동이 없다.
			case !k.OverZero():
				required[am.Currency()] = [2]Big{rq[0].Add(am.Big()), rq[1]}
			// Fee가 있으면 required amount 누적하고 required fee도 누적한다.
			default:
				required[am.Currency()] = [2]Big{rq[0].Add(am.Big()).Add(k), rq[1].Add(k)}
			}
		}
	}

	// Item 전체에 대하여 계산한 CurrencyID별 required amount와 required fee 값을 반환한다.
	return required, nil
}

func CheckEnoughBalance(
	holder base.Address,
	required map[CurrencyID][2]Big,
	getState func(key string) (state.State, bool, error),
) (map[CurrencyID]AmountState, error) {
	sb := map[CurrencyID]AmountState{}

	for cid := range required {
		rq := required[cid]

		// state.State는 기본적으로 mitum에 있는 state.StateV0를 말한다고 보면 됨.
		// type StateV0 struct {
		// 	  h              valuehash.Hash
		// 	  key            string
		// 	  value          Value
		// 	  height         base.Height
		// 	  previousHeight base.Height
		// 	  operations     []valuehash.Hash
		// }s
		st, err := existsState(StateKeyBalance(holder, cid), "currency of holder", getState)
		if err != nil {
			return nil, err
		}
		// state에 있는 Balance Value를 가져옴. Balance Value는 Amount를 말함.
		am, err := StateBalanceValue(st)
		if err != nil {
			return nil, operation.NewBaseReasonError("insufficient balance of sender: %w", err)
		}

		// Balance가 required Amount를 감당하기에 충분한지 확인한다.
		if am.Big().Compare(rq[0]) < 0 {
			return nil, operation.NewBaseReasonError(
				"insufficient balance of sender, %s; %d !> %d", holder.String(), am.Big(), rq[0])
		}
		// 현재 상태의 Amount state를 반환한다.
		sb[cid] = NewAmountState(st, cid)
	}

	return sb, nil
}
