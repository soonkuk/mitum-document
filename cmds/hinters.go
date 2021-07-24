package cmds

import (
	"github.com/spikeekips/mitum/launch"
	"github.com/spikeekips/mitum/util/hint"

	"github.com/soonkuk/mitum-data/blocksign"
	"github.com/soonkuk/mitum-data/currency"
	"github.com/soonkuk/mitum-data/digest"
)

var (
	Hinters []hint.Hinter
	Types   []hint.Type
)

var types = []hint.Type{
	currency.KeyType,
	currency.KeysType,
	currency.NilFeeerType,
	currency.FixedFeeerType,
	currency.RatioFeeerType,
	currency.TransfersFactType,
	currency.TransfersType,
	blocksign.TransferDocumentsFactType,
	blocksign.TransferDocumentsType,
	currency.AccountType,
	currency.AmountStateType,
	currency.GenesisCurrenciesFactType,
	currency.GenesisCurrenciesType,
	currency.AmountType,
	currency.FeeOperationFactType,
	currency.FeeOperationType,
	currency.CurrencyDesignType,
	currency.CurrencyRegisterFactType,
	currency.CurrencyRegisterType,
	currency.CurrencyPolicyUpdaterFactType,
	currency.CurrencyPolicyUpdaterType,
	currency.CreateAccountsFactType,
	currency.CreateAccountsType,
	blocksign.CreateDocumentsFactType,
	blocksign.CreateDocumentsType,
	currency.CreateAccountsItemSingleAmountType,
	blocksign.CreateDocumentsItemSingleFileType,
	currency.TransfersItemMultiAmountsType,
	blocksign.TransfersItemSingleDocumentType,
	currency.CurrencyPolicyType,
	currency.AddressType,
	currency.CreateAccountsItemMultiAmountsType,
	currency.TransfersItemSingleAmountType,
	currency.KeyUpdaterFactType,
	currency.KeyUpdaterType,
	blocksign.DocumentDataType,
	blocksign.DocInfoType,
	blocksign.DocSignType,
	blocksign.DocumentInventoryType,
	blocksign.FileHashType,
	digest.ProblemType,
	digest.NodeInfoType,
	digest.BaseHalType,
	digest.AccountValueType,
	// digest.DocumentValueType,
	digest.OperationValueType,
}

var hinters = []hint.Hinter{
	currency.Account{},
	currency.Address(""),
	currency.AmountState{},
	currency.Amount{},
	currency.CreateAccountsFact{},
	currency.CreateAccountsItemMultiAmountsHinter,
	currency.CreateAccountsItemSingleAmountHinter,
	currency.CreateAccounts{},
	blocksign.CreateDocumentsFact{},
	blocksign.CreateDocumentsItemSingleFileHinter,
	blocksign.CreateDocuments{},
	currency.CurrencyDesign{},
	currency.CurrencyPolicyUpdaterFact{},
	currency.CurrencyPolicyUpdater{},
	currency.CurrencyPolicy{},
	currency.CurrencyRegisterFact{},
	currency.CurrencyRegister{},
	currency.FeeOperationFact{},
	currency.FeeOperation{},
	blocksign.DocumentData{},
	blocksign.DocInfo{},
	blocksign.DocSign{},
	blocksign.DocumentInventory{},
	currency.FixedFeeer{},
	currency.GenesisCurrenciesFact{},
	currency.GenesisCurrencies{},
	currency.KeyUpdaterFact{},
	currency.KeyUpdater{},
	currency.Keys{},
	currency.Key{},
	currency.NilFeeer{},
	//currency.Owner(""),
	currency.RatioFeeer{},
	blocksign.FileHash(""),
	currency.TransfersFact{},
	currency.TransfersItemMultiAmountsHinter,
	currency.TransfersItemSingleAmountHinter,
	currency.Transfers{},
	blocksign.TransferDocumentsFact{},
	blocksign.TransferDocuments{},
	blocksign.TransfersItemSingleDocumentHinter,
	digest.AccountValue{},
	// digest.DocumentValue{},
	digest.BaseHal{},
	digest.NodeInfo{},
	digest.OperationValue{},
	digest.Problem{},
}

func init() {
	Hinters = make([]hint.Hinter, len(launch.EncoderHinters)+len(hinters))
	copy(Hinters, launch.EncoderHinters)
	copy(Hinters[len(launch.EncoderHinters):], hinters)

	Types = make([]hint.Type, len(launch.EncoderTypes)+len(types))
	copy(Types, launch.EncoderTypes)
	copy(Types[len(launch.EncoderTypes):], types)
}
