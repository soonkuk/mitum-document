package cmds

import (
	"github.com/pkg/errors"

	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/util"

	"github.com/soonkuk/mitum-blocksign/blockcity"
	currencycmds "github.com/spikeekips/mitum-currency/cmds"
	mitumcmds "github.com/spikeekips/mitum/launch/cmds"
)

type CreateBlockcityVotingDocumentCommand struct {
	*BaseCommand
	currencycmds.OperationFlags
	Sender     currencycmds.AddressFlag    `arg:"" name:"sender" help:"sender address" required:""`
	Round      uint                        `arg:"" name:"round" help:"voting round" required:""`
	Candidates []currencycmds.AddressFlag  `name:"candidates" help:"candidates address" required:""`
	DocumentId string                      `arg:"" name:"documentid" help:"document id" required:""`
	Currency   currencycmds.CurrencyIDFlag `arg:"" name:"currency" help:"currency id" required:""`
	Seal       mitumcmds.FileLoad          `help:"seal" optional:""`
	sender     base.Address
	candidates []blockcity.VotingCandidate
}

func NewCreateBlockcityVotingDocumentCommand() CreateBlockcityVotingDocumentCommand {
	return CreateBlockcityVotingDocumentCommand{
		BaseCommand: NewBaseCommand("create-blockcity-voting-document-operation"),
	}
}

func (cmd *CreateBlockcityVotingDocumentCommand) Run(version util.Version) error {
	if err := cmd.Initialize(cmd, version); err != nil {
		return errors.Errorf("failed to initialize command: %q", err)
	}

	if err := cmd.parseFlags(); err != nil {
		return err
	}

	op, err := cmd.createOperation()
	if err != nil {
		return err
	}

	sl, err := LoadSealAndAddOperation(
		cmd.Seal.Bytes(),
		cmd.Privatekey,
		cmd.NetworkID.NetworkID(),
		op,
	)
	if err != nil {
		return err
	}
	currencycmds.PrettyPrint(cmd.Out, cmd.Pretty, sl)
	return nil
}

func (cmd *CreateBlockcityVotingDocumentCommand) parseFlags() error {
	if err := cmd.OperationFlags.IsValid(nil); err != nil {
		return err
	}

	sa, err := cmd.Sender.Encode(jenc)
	if err != nil {
		return errors.Wrapf(err, "invalid sender format, %q", cmd.Sender.String())
	}
	cmd.sender = sa

	if len(cmd.Candidates) < 1 {
		return errors.Errorf("empty candidates, must be given at least one")
	}

	{
		candidates := make([]blockcity.VotingCandidate, len(cmd.Candidates))
		for i := range cmd.Candidates {
			ca, err := cmd.Candidates[i].Encode(jenc)
			if err != nil {
				return errors.Wrapf(err, "invalid address format, %q", cmd.Candidates[i].String())
			}
			candidates[i] = blockcity.MustNewVotingCandidate(ca, "")
		}
		cmd.candidates = candidates
	}

	return nil
}

func (cmd *CreateBlockcityVotingDocumentCommand) createOperation() (operation.Operation, error) {
	i, err := loadOperations(cmd.Seal.Bytes(), cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, err
	}
	var items []blockcity.CreateDocumentsItem
	for j := range i {
		if t, ok := i[j].(blockcity.CreateDocuments); ok {
			items = t.Fact().(blockcity.CreateDocumentsFact).Items()
		}
	}

	info := blockcity.NewDocInfo(cmd.DocumentId, blockcity.CityVotingDataType)
	votingDoc := blockcity.NewCityVotingData(info, cmd.sender, cmd.Round, cmd.candidates)
	doc := blockcity.NewDocument(votingDoc)
	item := blockcity.NewCreateDocumentsItemImpl(
		doc,
		cmd.Currency.CID,
	)

	if err := item.IsValid(nil); err != nil {
		return nil, err
	}
	items = append(items, item)

	fact := blockcity.NewCreateDocumentsFact([]byte(cmd.Token), cmd.sender, items)

	sig, err := base.NewFactSignature(cmd.Privatekey, fact, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, err
	}
	fs := []base.FactSign{
		base.NewBaseFactSign(cmd.Privatekey.Publickey(), sig),
	}

	op, err := blockcity.NewCreateDocuments(fact, fs, cmd.Memo)
	if err != nil {
		return nil, errors.Errorf("failed to create create-blockcity-voting-document operation operation: %q", err)
	}
	return op, nil
}
