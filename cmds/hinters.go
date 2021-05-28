package cmds

import (
	"github.com/spikeekips/mitum/launch"
	"github.com/spikeekips/mitum/util/hint"

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
	currency.TransferDocumentsFactType,
	currency.TransferDocumentsType,
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
	currency.CreateDocumentsFactType,
	currency.CreateDocumentsType,
	currency.CreateAccountsItemSingleAmountType,
	currency.CreateDocumentsItemSingleFileType,
	currency.TransfersItemMultiAmountsType,
	currency.TransfersItemSingleDocumentType,
	currency.CurrencyPolicyType,
	currency.AddressType,
	currency.CreateAccountsItemMultiAmountsType,
	currency.TransfersItemSingleAmountType,
	currency.KeyUpdaterFactType,
	currency.KeyUpdaterType,
	currency.FileDataType,
	currency.SignCodeType,
	digest.ProblemType,
	digest.NodeInfoType,
	digest.BaseHalType,
	digest.AccountValueType,
	digest.DocumentValueType,
	digest.OperationValueType,
}

var hinters = []hint.Hinter{
	currency.Account{},                            // a014
	currency.Address(""),                          // a000
	currency.AmountState{},                        // a023
	currency.Amount{},                             // a022
	currency.CreateAccountsFact{},                 // a005
	currency.CreateAccountsItemMultiAmountsHinter, // a024
	currency.CreateAccountsItemSingleAmountHinter, // a025
	currency.CreateAccounts{},                     // a006
	currency.CreateDocumentsFact{},                // a042
	currency.CreateDocumentsItemSingleFileHinter,  // a041
	currency.CreateDocuments{},                    // a043
	currency.CurrencyDesign{},                     // a035
	currency.CurrencyPolicyUpdaterFact{},          // a034
	currency.CurrencyPolicyUpdater{},              // a035
	currency.CurrencyPolicy{},                     // a036
	currency.CurrencyRegisterFact{},               // a028
	currency.CurrencyRegister{},                   // a029
	currency.FeeOperationFact{},                   // a012
	currency.FeeOperation{},                       // a013
	currency.FileData{},
	currency.FixedFeeer{},            // a032
	currency.GenesisCurrenciesFact{}, // a020
	currency.GenesisCurrencies{},     // a021
	currency.KeyUpdaterFact{},        // a009
	currency.KeyUpdater{},            // a010
	currency.Keys{},                  // a004
	currency.Key{},                   // a003
	currency.NilFeeer{},              // a031
	//currency.Owner(""),                         // a046
	currency.RatioFeeer{},                      // a033
	currency.SignCode(""),                      // a045
	currency.TransfersFact{},                   // a001
	currency.TransfersItemMultiAmountsHinter,   // a026
	currency.TransfersItemSingleAmountHinter,   // a027
	currency.Transfers{},                       // a002
	currency.TransferDocumentsFact{},           // a047
	currency.TransferDocuments{},               // a048
	currency.TransfersItemSingleDocumentHinter, // a049
	digest.AccountValue{},                      // a018
	digest.DocumentValue{},
	digest.BaseHal{},        // a016
	digest.NodeInfo{},       // a015
	digest.OperationValue{}, // a019
	digest.Problem{},        // a017
}

func init() {
	Hinters = make([]hint.Hinter, len(launch.EncoderHinters)+len(hinters))
	copy(Hinters, launch.EncoderHinters)
	copy(Hinters[len(launch.EncoderHinters):], hinters)

	Types = make([]hint.Type, len(launch.EncoderTypes)+len(types))
	copy(Types, launch.EncoderTypes)
	copy(Types[len(launch.EncoderTypes):], types)
}
