package currency

import (
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/base/state"
	"github.com/spikeekips/mitum/util/valuehash"
	"golang.org/x/xerrors"
)

func (op CreateDocuments) Process(
	func(key string) (state.State, bool, error),
	func(valuehash.Hash, ...state.State) error,
) error {
	return nil
}

type CreateDocumentsItemProcessor struct {
	cp   *CurrencyPool
	h    valuehash.Hash
	item CreateDocumentsItem
	nas  state.State // document의 new account state
	nfs  state.State // document의 new FileID state
	keys Keys        // owner의 keys
}

func (opp *CreateDocumentsItemProcessor) PreProcess(
	getState func(key string) (state.State, bool, error),
	_ func(valuehash.Hash, ...state.State) error,
) error {

	if err := opp.item.IsValid(nil); err != nil {
		return err
	}

	// get account address of new target document
	var dadr base.Address
	if a, err := opp.item.Address(); err != nil {
		return err
	} else {
		dadr = a
	}

	// check the existence of owner account state key
	if st, err := existsState(StateKeyAccount(opp.item.Owner()), "key of Owner", getState); err != nil {
		return err
	} else {
		// stateAccount에서 owner keys 값 가져오기
		if ks, err := StateKeysValue(st); err != nil {
			return operation.NewBaseReasonErrorFromError(err)
		} else {
			opp.keys = ks
		}
	}

	// check the existence of document account state by keyAccount and prepare new account state
	// new account state generated with keyAccount(addr+hintType+:account)
	if st, err := notExistsState(StateKeyDocument(dadr), "key of Document", getState); err != nil {
		return err
	} else if st, err := SetStateKeysValue(st, opp.item.Keys()); err != nil {
		return err
	} else {
		opp.nas = st
	}

	// FileData state를 만든다.
	switch st, found, err := getState(StateKeyFileData(dadr)); {
	case err != nil:
		return err
	case found:
		return xerrors.Errorf(" already registered, %q", dadr)
	default:
		opp.nfs = st
	}

	return nil
}

func (opp *CreateDocumentsItemProcessor) Process(
	_ func(key string) (state.State, bool, error),
	_ func(valuehash.Hash, ...state.State) error,
) ([]state.State, error) {

	/*
		// new generated document account
		var nac Account
		// make account for new generated document from owner's keys
		// 여기에서 사용하는 keys 값이 생성되는 document account의 key 값이 된다.
		// 따라서 document의 key값으로 사용하려고 하는 owner의 key값을 넣는 것이 맞다.
		if ac, err := NewAccountFromKeys(opp.keys); err != nil {
			return nil, err
		} else {
			nac = ac
		}
	*/

	sts := make([]state.State, 2)
	// opp.ns has already key in the form of keyAccount.
	// So SetStateAccountValue set keys to state.
	// setStateAccountValue는 account 전체를 value로 사용하는 것이고
	// setStateKeysValue는 사실은 account를 받아서 account의 keys값을 업데이트 하고 account를 다시 setValue하는 것이다.
	// Keys로 만들어진 위의 account를 사용하여 StateAccount의 Value를 set한다.
	// StateAccount의 Value는 Account 자체가 모두 들어간다. 결국 address와 Keys가 모두 저장된다.
	// 나중에 이 StateAccount의 Value를 사용할 때는 LoadStateAccountValue라는 함수에서
	// 불러와서 사용하는데 이 때 keys 값만 사용하기는 한다.
	// 하지만 결국 document의 addr과 state로 저장한 addr값이 일치 하지 않기 때문에 문제발생의 소지가 있다.
	// 방법은 document라는 객체를 따로 만들어서 StateDocument를 저장하는 것이 옳은 방법일 것 같다.
	// 하지만 또 생각해보니 괜찮을 것 같은 것은, 어차피 나중에는 account의 address와 key들은 연관성이
	// 없어지게 되므로 account에서 address와 key들이 서로 맞지 않는다고 문제는 없을 것 같다.

	// 새로운 document account에 owner의 key를 갖도록 set value를 한다.
	if st, err := SetStateKeysValue(opp.nas, opp.keys); err != nil {
		return nil, err
	} else {
		sts[0] = st
	}
	filedata := NewFileData(opp.item.SignCode(), opp.item.Owner())

	if f, err := SetStateFileDataValue(opp.nfs, filedata); err != nil {
		return nil, err
	} else {
		sts[1] = f
	}

	return sts, nil
}

type CreateDocumentsProcessor struct {
	cp *CurrencyPool
	CreateDocuments
	sb       map[CurrencyID]AmountState      // sender의 StateBalance
	ns       []*CreateDocumentsItemProcessor // ItemProcessor
	required map[CurrencyID][2]Big           // Fee
}

func NewCreateDocumentsProcessor(cp *CurrencyPool) GetNewProcessor {
	return func(op state.Processor) (state.Processor, error) {
		if i, ok := op.(CreateDocuments); !ok {
			return nil, xerrors.Errorf("not CreateDocuments, %T", op)
		} else {
			return &CreateDocumentsProcessor{
				cp:              cp,
				CreateDocuments: i,
			}, nil
		}
	}
}

// TODO : Sender의 StateBalance 준비하기
func (opp *CreateDocumentsProcessor) PreProcess(
	getState func(key string) (state.State, bool, error),
	setState func(valuehash.Hash, ...state.State) error,
) (state.Processor, error) {
	fact := opp.Fact().(CreateDocumentsFact)

	// sender의 State Account 확인
	if err := checkExistsState(StateKeyAccount(fact.sender), getState); err != nil {
		return nil, err
	}

	// CreateDocument를 처리하기 위한 fee를 계산한다.
	// required는 k : currencyID v : Big의 map이다.
	if required, err := opp.calculateItemsFee(); err != nil {
		return nil, operation.NewBaseReasonError("failed to calculate fee: %w", err)
		// sender의 balance가 required를 처리하기에 충분한지 확인하고 sender의 현재 StateBalance(처리전)를 가져온다.
	} else if sb, err := CheckDocumentOwnerEnoughBalance(fact.sender, required, getState); err != nil {
		return nil, err
	} else {
		opp.required = required
		opp.sb = sb
	}

	ns := make([]*CreateDocumentsItemProcessor, len(fact.items))
	for i := range fact.items {
		c := &CreateDocumentsItemProcessor{cp: opp.cp, h: opp.Hash(), item: fact.items[i]}
		if err := c.PreProcess(getState, setState); err != nil {
			return nil, operation.NewBaseReasonErrorFromError(err)
		}

		ns[i] = c
	}

	if err := checkFactSignsByState(fact.sender, opp.Signs(), getState); err != nil {
		return nil, operation.NewBaseReasonError("invalid signing: %w", err)
	}

	opp.ns = ns

	return opp, nil
}

func (opp *CreateDocumentsProcessor) Process( // nolint:dupl
	getState func(key string) (state.State, bool, error),
	setState func(valuehash.Hash, ...state.State) error,
) error {
	fact := opp.Fact().(CreateDocumentsFact)

	var sts []state.State // nolint:prealloc
	// itemProcessor의 본처리를 한다. 처리후 나오는 state 배열은 새로 생성되는 target document의 state들(stateAccount, stateFileData)이다.
	for i := range opp.ns {
		if s, err := opp.ns[i].Process(getState, setState); err != nil {
			return operation.NewBaseReasonError("failed to process create account item: %w", err)
		} else {
			sts = append(sts, s...)
		}
	}

	// sender의 currencyID별 stateBalance를 변경하여 배열에 모은다.
	for k := range opp.required {
		rq := opp.required[k]
		sts = append(sts, opp.sb[k].Sub(rq[0]).AddFee(rq[1]))
	}

	return setState(fact.Hash(), sts...)
}

func (opp *CreateDocumentsProcessor) calculateItemsFee() (map[CurrencyID][2]Big, error) {
	fact := opp.Fact().(CreateDocumentsFact)

	items := make([]CreateDocumentsItem, len(fact.items))
	for i := range fact.items {
		items[i] = fact.items[i]
	}

	return CalculateDocumentItemsFee(opp.cp, items)
}

// CalculateItemsFee는 item별로 fee를 계산하여 currencyID별 fee값의 map을 반환한다.
// fee는 currencyPool에서 currencyID를 사용하여 찾아내고 적용한다.
// currencyPool이 nil일 때는 fee없이 생성한다.
func CalculateDocumentItemsFee(cp *CurrencyPool, items []CreateDocumentsItem) (map[CurrencyID][2]Big, error) {
	required := map[CurrencyID][2]Big{}

	for i := range items {
		it := items[i]

		rq := [2]Big{ZeroBig, ZeroBig}

		if k, found := required[it.Currency()]; found {
			rq = k
		}

		if cp == nil {
			required[it.Currency()] = [2]Big{rq[0], rq[1]}

			continue
		}

		// currencyID로 fee 계산
		feeer, found := cp.Feeer(it.Currency())
		if !found {
			return nil, xerrors.Errorf("unknown currency id found, %q", it.Currency())
		}
		switch k, err := feeer.Fee(ZeroBig); {
		case err != nil:
			return nil, err
		case !k.OverZero():
			required[it.Currency()] = [2]Big{rq[0], rq[1]}
		default:
			required[it.Currency()] = [2]Big{rq[0].Add(k), rq[1].Add(k)}
		}

	}

	return required, nil
}

// CheckEnoughBalance는 holder의 현재 상태 Balance가 Fee를 지불하기에 충분한지 확인한다.
// 충분하다면 현재상태의 StateBalance를 반환한다.
func CheckDocumentOwnerEnoughBalance(
	holder base.Address,
	required map[CurrencyID][2]Big,
	getState func(key string) (state.State, bool, error),
) (map[CurrencyID]AmountState, error) {
	sb := map[CurrencyID]AmountState{}

	for cid := range required {
		rq := required[cid]

		st, err := existsState(StateKeyBalance(holder, cid), "currency of holder", getState)
		if err != nil {
			return nil, err
		}

		// sender의 StateBalance를 가져와서
		am, err := StateBalanceValue(st)
		if err != nil {
			return nil, operation.NewBaseReasonError("insufficient balance of sender: %w", err)
		}

		// fee와 비교하여 모자라면 에러 발생
		if am.Big().Compare(rq[0]) < 0 {
			return nil, operation.NewBaseReasonError(
				"insufficient balance of sender, %s; %d !> %d", holder.String(), am.Big(), rq[0])
		} else {
			// 잔고가 충분하면 현재상태의 잔고 amount의 AmountState 생성
			sb[cid] = NewAmountState(st, cid)
		}
	}
	// 현재 상태의 AmountState 반환
	return sb, nil
}
