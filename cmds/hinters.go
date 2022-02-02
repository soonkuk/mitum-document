package cmds

import (
	"github.com/spikeekips/mitum/launch"
	"github.com/spikeekips/mitum/util/hint"

	"github.com/soonkuk/mitum-blocksign/blocksign"
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
	blocksign.CreateDocumentsItemSingleFileType,
	blocksign.CreateDocumentsFactType,
	blocksign.CreateDocumentsType,
	blocksign.SignItemSingleDocumentType,
	blocksign.SignDocumentsFactType,
	blocksign.SignDocumentsType,
	blocksign.DocumentDataType,
	blocksign.DocInfoType,
	blocksign.DocSignType,
	blocksign.DocumentInventoryType,
	document.CreateDocumentsItemImplType,
	document.CreateDocumentsFactType,
	document.CreateDocumentsType,
	document.UpdateDocumentsItemImplType,
	document.UpdateDocumentsFactType,
	document.UpdateDocumentsType,
	document.DocumentType,
	document.CityUserDataType,
	document.CityLandDataType,
	document.CityVotingDataType,
	document.UserStatisticsType,
	document.UserDocIdType,
	document.LandDocIdType,
	document.VotingDocIdType,
	document.DocInfoType,
	document.VotingCandidateType,
	document.DocumentInventoryType,
	digest.ProblemType,
	digest.NodeInfoType,
	digest.BaseHalType,
	digest.AccountValueType,
	digest.BlocksignDocumentValueType,
	digest.BlockcityDocumentValueType,
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
	blocksign.CreateDocumentsFactHinter,
	blocksign.CreateDocumentsItemSingleFileHinter,
	blocksign.CreateDocumentsHinter,
	blocksign.SignDocumentsFactHinter,
	blocksign.SignDocumentsHinter,
	blocksign.SignItemSingleDocumentHinter,
	blocksign.DocumentDataHinter,
	blocksign.DocInfoHinter,
	blocksign.DocSignHinter,
	blocksign.DocumentInventoryHinter,
	document.CreateDocumentsFactHinter,
	document.CreateDocumentsHinter,
	document.CreateDocumentsItemImplHinter,
	document.UpdateDocumentsFactHinter,
	document.UpdateDocumentsHinter,
	document.UpdateDocumentsItemImplHinter,
	document.DocumentHinter,
	document.CityUserDataHinter,
	document.CityLandDataHinter,
	document.CityVotingDataHinter,
	document.UserStatisticsHinter,
	document.DocInfoHinter,
	document.VotingCandidateHinter,
	document.UserDocIdHinter,
	document.LandDocIdHinter,
	document.VotingDocIdHinter,
	document.DocumentInventoryHinter,
	digest.AccountValue{},
	digest.BlocksignDocumentValue{},
	digest.BlockcityDocumentValue{},
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
