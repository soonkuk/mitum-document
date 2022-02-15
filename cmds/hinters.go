package cmds

import (
	"github.com/spikeekips/mitum/launch"
	"github.com/spikeekips/mitum/util/hint"

	"github.com/soonkuk/mitum-blocksign/digest"
	"github.com/soonkuk/mitum-blocksign/document"
	"github.com/spikeekips/mitum-currency/currency"
)

var (
	Hinters []hint.Hinter
	Types   []hint.Type
)

var types = []hint.Type{
	currency.AccountType,
	currency.AddressType,
	currency.AmountType,
	currency.CreateAccountsFactType,
	currency.CreateAccountsItemMultiAmountsType,
	currency.CreateAccountsItemSingleAmountType,
	currency.CreateAccountsType,
	currency.CurrencyDesignType,
	currency.CurrencyPolicyType,
	currency.CurrencyPolicyUpdaterFactType,
	currency.CurrencyPolicyUpdaterType,
	currency.CurrencyRegisterFactType,
	currency.CurrencyRegisterType,
	currency.FeeOperationFactType,
	currency.FeeOperationType,
	currency.FixedFeeerType,
	currency.GenesisCurrenciesFactType,
	currency.GenesisCurrenciesType,
	currency.AccountKeyType,
	currency.KeyUpdaterFactType,
	currency.KeyUpdaterType,
	currency.AccountKeysType,
	currency.NilFeeerType,
	currency.RatioFeeerType,
	currency.SuffrageInflationFactType,
	currency.SuffrageInflationType,
	currency.TransfersFactType,
	currency.TransfersItemMultiAmountsType,
	currency.TransfersItemSingleAmountType,
	currency.TransfersType,
	document.SignItemSingleDocumentType,
	document.SignDocumentsFactType,
	document.SignDocumentsType,
	document.DocSignType,
	document.CreateDocumentsItemImplType,
	document.CreateDocumentsFactType,
	document.CreateDocumentsType,
	document.UpdateDocumentsItemImplType,
	document.UpdateDocumentsFactType,
	document.UpdateDocumentsType,
	document.DocumentType,
	document.BSDocDataType,
	document.BCUserDataType,
	document.BCLandDataType,
	document.BCVotingDataType,
	document.BCHistoryDataType,
	document.UserStatisticsType,
	document.BSDocIdType,
	document.UserDocIdType,
	document.LandDocIdType,
	document.VotingDocIdType,
	document.HistoryDocIdType,
	document.DocInfoType,
	document.VotingCandidateType,
	document.DocumentInventoryType,
	digest.ProblemType,
	digest.NodeInfoType,
	digest.BaseHalType,
	digest.AccountValueType,
	// digest.BSDocumentValueType,
	digest.DocumentValueType,
	digest.OperationValueType,
}

var hinters = []hint.Hinter{
	currency.AccountHinter,
	currency.AddressHinter,
	currency.AmountHinter,
	currency.CreateAccountsFactHinter,
	currency.CreateAccountsItemMultiAmountsHinter,
	currency.CreateAccountsItemSingleAmountHinter,
	currency.CreateAccountsHinter,
	currency.CurrencyDesignHinter,
	currency.CurrencyPolicyUpdaterFactHinter,
	currency.CurrencyPolicyUpdaterHinter,
	currency.CurrencyPolicyHinter,
	currency.CurrencyRegisterFactHinter,
	currency.CurrencyRegisterHinter,
	currency.FeeOperationFactHinter,
	currency.FeeOperationHinter,
	currency.FixedFeeerHinter,
	currency.GenesisCurrenciesFactHinter,
	currency.GenesisCurrenciesHinter,
	currency.KeyUpdaterFactHinter,
	currency.KeyUpdaterHinter,
	currency.AccountKeysHinter,
	currency.AccountKeyHinter,
	currency.NilFeeerHinter,
	currency.RatioFeeerHinter,
	currency.SuffrageInflationFactHinter,
	currency.SuffrageInflationHinter,
	currency.TransfersFactHinter,
	currency.TransfersItemMultiAmountsHinter,
	currency.TransfersItemSingleAmountHinter,
	currency.TransfersHinter,
	document.SignDocumentsFactHinter,
	document.SignDocumentsHinter,
	document.SignItemSingleDocumentHinter,
	document.DocSignHinter,
	document.CreateDocumentsFactHinter,
	document.CreateDocumentsHinter,
	document.CreateDocumentsItemImplHinter,
	document.UpdateDocumentsFactHinter,
	document.UpdateDocumentsHinter,
	document.UpdateDocumentsItemImplHinter,
	document.DocumentHinter,
	document.BSDocDataHinter,
	document.BCUserDataHinter,
	document.BCLandDataHinter,
	document.BCVotingDataHinter,
	document.BCHistoryDataHinter,
	document.UserStatisticsHinter,
	document.DocInfoHinter,
	document.VotingCandidateHinter,
	document.BSDocIdHinter,
	document.UserDocIdHinter,
	document.LandDocIdHinter,
	document.VotingDocIdHinter,
	document.HistoryDocIdHinter,
	document.DocumentInventoryHinter,
	digest.AccountValue{},
	digest.DocumentValue{},
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
