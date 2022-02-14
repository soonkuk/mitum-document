package cmds

import (
	"github.com/pkg/errors"
	"github.com/soonkuk/mitum-blocksign/document"

	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/util"

	currencycmds "github.com/spikeekips/mitum-currency/cmds"
	mitumcmds "github.com/spikeekips/mitum/launch/cmds"
)

type UpdateBlockcityVotingDocumentCommand struct {
	*BaseCommand
	currencycmds.OperationFlags
	Sender      currencycmds.AddressFlag    `arg:"" name:"sender" help:"sender address" required:""`
	Round       uint                        `arg:"" name:"round" help:"voting round" required:""`
	EndVoteTime string                      `arg:"" name:"endvotetime" help:"end vote time" required:""`
	Candidates  []currencycmds.AddressFlag  `name:"candidates" help:"candidates address" required:""`
	Nicknames   []string                    `name:"nicknames" help:"candidates nicknames" required:""`
	BossName    string                      `arg:"" name:"bossname" help:"boss name" required:""`
	Account     currencycmds.AddressFlag    `arg:"" name:"bossaccount" help:"boss account address" required:""`
	Term        string                      `arg:"" name:"termofoffice" help:"term of office" required:""`
	DocumentId  string                      `arg:"" name:"documentid" help:"document id" required:""`
	Currency    currencycmds.CurrencyIDFlag `arg:"" name:"currency" help:"currency id" required:""`
	Seal        mitumcmds.FileLoad          `help:"seal" optional:""`
	sender      base.Address
	candidates  []document.VotingCandidate
	account     base.Address
}

func NewUpdateBlockcityVotingDocumentCommand() UpdateBlockcityVotingDocumentCommand {
	return UpdateBlockcityVotingDocumentCommand{
		BaseCommand: NewBaseCommand("update-blockcity-voting-document-operation"),
	}
}

func (cmd *UpdateBlockcityVotingDocumentCommand) Run(version util.Version) error {
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

func (cmd *UpdateBlockcityVotingDocumentCommand) parseFlags() error {
	if err := cmd.OperationFlags.IsValid(nil); err != nil {
		return err
	}

	sa, err := cmd.Sender.Encode(jenc)
	if err != nil {
		return errors.Wrapf(err, "invalid sender format, %q", cmd.Sender.String())
	}
	cmd.sender = sa

	ba, err := cmd.Account.Encode(jenc)
	if err != nil {
		return errors.Wrapf(err, "invalid boss account format, %q", cmd.Account.String())
	}
	cmd.account = ba

	if len(cmd.Candidates) < 1 {
		return errors.Errorf("empty candidates, must be given at least one")
	}

	{
		candidates := make([]document.VotingCandidate, len(cmd.Candidates))
		for i := range cmd.Candidates {
			ca, err := cmd.Candidates[i].Encode(jenc)
			if err != nil {
				return errors.Wrapf(err, "invalid address format, %q", cmd.Candidates[i].String())
			}
			candidates[i] = document.MustNewVotingCandidate(ca, cmd.Nicknames[i], "")
		}
		cmd.candidates = candidates
	}

	return nil
}

func (cmd *UpdateBlockcityVotingDocumentCommand) createOperation() (operation.Operation, error) {
	i, err := loadOperations(cmd.Seal.Bytes(), cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, err
	}
	var items []document.UpdateDocumentsItem
	for j := range i {
		if t, ok := i[j].(document.UpdateDocuments); ok {
			items = t.Fact().(document.UpdateDocumentsFact).Items()
		}
	}

	info := document.NewDocInfo(cmd.DocumentId, document.BCVotingDataType)
	doc := document.NewBCVotingData(info, cmd.sender, cmd.Round, cmd.EndVoteTime, cmd.candidates, cmd.BossName, cmd.account, cmd.Term)

	item := document.NewUpdateDocumentsItemImpl(
		doc,
		cmd.Currency.CID,
	)

	if err := item.IsValid(nil); err != nil {
		return nil, err
	}
	items = append(items, item)

	fact := document.NewUpdateDocumentsFact([]byte(cmd.Token), cmd.sender, items)

	sig, err := base.NewFactSignature(cmd.Privatekey, fact, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, err
	}
	fs := []base.FactSign{
		base.NewBaseFactSign(cmd.Privatekey.Publickey(), sig),
	}

	op, err := document.NewUpdateDocuments(fact, fs, cmd.Memo)
	if err != nil {
		return nil, errors.Errorf("failed to create update-blockcity-voting-document operation operation: %q", err)
	}
	return op, nil
}
