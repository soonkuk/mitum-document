package cmds

import (
	"github.com/spikeekips/mitum/launch"
	"github.com/spikeekips/mitum/util/hint"

	"github.com/soonkuk/mitum-blocksign/blockcity"
	"github.com/soonkuk/mitum-blocksign/blocksign"
	"github.com/soonkuk/mitum-blocksign/digest"
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
	blockcity.CreateDocumentsItemImplType,
	blockcity.CreateDocumentsFactType,
	blockcity.CreateDocumentsType,
	blockcity.UpdateDocumentsItemImplType,
	blockcity.UpdateDocumentsFactType,
	blockcity.UpdateDocumentsType,
	blockcity.DocumentType,
	blockcity.CityUserDataType,
	blockcity.CityLandDataType,
	blockcity.CityVotingDataType,
	blockcity.UserStatisticsType,
	blockcity.UserDocIdType,
	blockcity.LandDocIdType,
	blockcity.VotingDocIdType,
	blockcity.DocInfoType,
	blockcity.VotingCandidateType,
	blockcity.DocumentInventoryType,
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
	blockcity.CreateDocumentsFactHinter,
	blockcity.CreateDocumentsHinter,
	blockcity.CreateDocumentsItemImplHinter,
	blockcity.UpdateDocumentsFactHinter,
	blockcity.UpdateDocumentsHinter,
	blockcity.UpdateDocumentsItemImplHinter,
	blockcity.DocumentHinter,
	blockcity.CityUserDataHinter,
	blockcity.CityLandDataHinter,
	blockcity.CityVotingDataHinter,
	blockcity.UserStatisticsHinter,
	blockcity.DocInfoHinter,
	blockcity.VotingCandidateHinter,
	blockcity.UserDocIdHinter,
	blockcity.LandDocIdHinter,
	blockcity.VotingDocIdHinter,
	blockcity.DocumentInventoryHinter,
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
