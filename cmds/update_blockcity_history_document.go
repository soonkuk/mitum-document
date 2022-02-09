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

type UpdateBlockcityHistoryDocumentCommand struct {
	*BaseCommand
	currencycmds.OperationFlags
	Sender      currencycmds.AddressFlag    `arg:"" name:"sender" help:"sender address" required:""`
	Name        string                      `arg:"" name:"name" help:"name" required:""`
	Account     currencycmds.AddressFlag    `arg:"" name:"renteraccount" help:"renter account address" required:""`
	Date        string                      `arg:"" name:"date" help:"sdate" required:""`
	Usage       string                      `arg:"" name:"usage" help:"usage" required:""`
	Application string                      `arg:"" name:"application" help:"application" required:""`
	DocumentId  string                      `arg:"" name:"documentid" help:"document id" required:""`
	Currency    currencycmds.CurrencyIDFlag `arg:"" name:"currency" help:"currency id" required:""`
	Seal        mitumcmds.FileLoad          `help:"seal" optional:""`
	sender      base.Address
	account     base.Address
}

func NewUpdateBlockcityHistoryDocumentCommand() UpdateBlockcityHistoryDocumentCommand {
	return UpdateBlockcityHistoryDocumentCommand{
		BaseCommand: NewBaseCommand("update-blockcity-history-document-operation"),
	}
}

func (cmd *UpdateBlockcityHistoryDocumentCommand) Run(version util.Version) error {
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

func (cmd *UpdateBlockcityHistoryDocumentCommand) parseFlags() error {
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
		return errors.Wrapf(err, "invalid account format, %q", cmd.Account.String())
	}
	cmd.account = ba

	return nil
}

func (cmd *UpdateBlockcityHistoryDocumentCommand) createOperation() (operation.Operation, error) {
	i, err := loadOperations(cmd.Seal.Bytes(), cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, err
	}
	var items []document.CreateDocumentsItem
	for j := range i {
		if t, ok := i[j].(document.CreateDocuments); ok {
			items = t.Fact().(document.CreateDocumentsFact).Items()
		}
	}

	info := document.NewDocInfo(cmd.DocumentId, document.BCHistoryDataType)
	doc := document.NewBCHistoryData(info, cmd.sender, cmd.Name, cmd.account, cmd.Date, cmd.Usage, cmd.Application)

	item := document.NewCreateDocumentsItemImpl(
		doc,
		cmd.Currency.CID,
	)

	if err := item.IsValid(nil); err != nil {
		return nil, err
	}
	items = append(items, item)

	fact := document.NewCreateDocumentsFact([]byte(cmd.Token), cmd.sender, items)

	sig, err := base.NewFactSignature(cmd.Privatekey, fact, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, err
	}
	fs := []base.FactSign{
		base.NewBaseFactSign(cmd.Privatekey.Publickey(), sig),
	}

	op, err := document.NewCreateDocuments(fact, fs, cmd.Memo)
	if err != nil {
		return nil, errors.Errorf("failed to create create-blockcity-land-document operation operation: %q", err)
	}
	return op, nil
}
