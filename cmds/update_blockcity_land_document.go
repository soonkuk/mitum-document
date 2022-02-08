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

type UpdateBlockcityLandDocumentCommand struct {
	*BaseCommand
	currencycmds.OperationFlags
	Sender        currencycmds.AddressFlag    `arg:"" name:"sender" help:"sender address" required:""`
	Address       string                      `arg:"" name:"landaddress" help:"land address" required:""`
	Area          string                      `arg:"" name:"landarea" help:"land area" required:""`
	Renter        string                      `arg:"" name:"renter" help:"renter nickname" required:""`
	Account       currencycmds.AddressFlag    `arg:"" name:"renteraccount" help:"renter account address" required:""`
	Rentdate      string                      `arg:"" name:"rentdate" help:"rent date" required:""`
	Periodday     uint                        `arg:"" name:"periodday" help:"periodday" required:""`
	DocumentId    string                      `arg:"" name:"documentid" help:"document id" required:""`
	Currency      currencycmds.CurrencyIDFlag `arg:"" name:"currency" help:"currency id" required:""`
	Seal          mitumcmds.FileLoad          `help:"seal" optional:""`
	sender        base.Address
	renterAccount base.Address
}

func NewUpdateBlockcityLandDocumentCommand() UpdateBlockcityLandDocumentCommand {
	return UpdateBlockcityLandDocumentCommand{
		BaseCommand: NewBaseCommand("update-blockcity-land-document-operation"),
	}
}

func (cmd *UpdateBlockcityLandDocumentCommand) Run(version util.Version) error {
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

func (cmd *UpdateBlockcityLandDocumentCommand) parseFlags() error {
	if err := cmd.OperationFlags.IsValid(nil); err != nil {
		return err
	}

	sa, err := cmd.Sender.Encode(jenc)
	if err != nil {
		return errors.Wrapf(err, "invalid sender format, %q", cmd.Sender.String())
	}
	cmd.sender = sa

	ra, err := cmd.Account.Encode(jenc)
	if err != nil {
		return errors.Wrapf(err, "invalid renter account format, %q", cmd.Account.String())
	}
	cmd.renterAccount = ra

	return nil
}

func (cmd *UpdateBlockcityLandDocumentCommand) createOperation() (operation.Operation, error) {
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

	info := document.NewDocInfo(cmd.DocumentId, document.BCLandDataType)
	doc := document.NewBCLandData(info, cmd.sender, cmd.Address, cmd.Area, cmd.Renter, cmd.renterAccount, cmd.Rentdate, cmd.Periodday)

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
		return nil, errors.Errorf("failed to create update-blockcity-land-document operation: %q", err)
	}
	return op, nil
}
