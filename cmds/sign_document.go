package cmds

import (
	"github.com/pkg/errors"

	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/util"

	"github.com/soonkuk/mitum-blocksign/document"
	currencycmds "github.com/spikeekips/mitum-currency/cmds"
	mitumcmds "github.com/spikeekips/mitum/launch/cmds"
)

type SignDocumentCommand struct {
	*BaseCommand
	currencycmds.OperationFlags
	Sender   currencycmds.AddressFlag    `arg:"" name:"sender" help:"sender address" required:""`
	DocId    string                      `arg:"" name:"documentid" help:"document id" required:""`
	Owner    currencycmds.AddressFlag    `arg:"" name:"owner" help:"owner address" required:""`
	Currency currencycmds.CurrencyIDFlag `arg:"" name:"currency" help:"currency id" required:""`
	Seal     mitumcmds.FileLoad          `help:"seal" optional:""`
	sender   base.Address
	owner    base.Address
}

func NewSignDocumentCommand() SignDocumentCommand {
	return SignDocumentCommand{
		BaseCommand: NewBaseCommand("sign-document-operation"),
	}
}

func (cmd *SignDocumentCommand) Run(version util.Version) error { // nolint:dupl
	if err := cmd.Initialize(cmd, version); err != nil {
		return errors.Errorf("failed to initialize command: %q", err)
	}

	if err := cmd.parseFlags(); err != nil {
		return err
	}

	var op operation.Operation
	if o, err := cmd.createOperation(); err != nil {
		return err
	} else {
		op = o
	}

	if sl, err := LoadSealAndAddOperation(
		cmd.Seal.Bytes(),
		cmd.Privatekey,
		cmd.NetworkID.NetworkID(),
		op,
	); err != nil {
		return err
	} else {
		currencycmds.PrettyPrint(cmd.Out, cmd.Pretty, sl)
	}

	return nil
}

func (cmd *SignDocumentCommand) parseFlags() error {
	if err := cmd.OperationFlags.IsValid(nil); err != nil {
		return err
	}

	if a, err := cmd.Sender.Encode(jenc); err != nil {
		return errors.Errorf("invalid sender format, %q: %q", cmd.Sender.String(), err)
	} else {
		cmd.sender = a
	}
	if a, err := cmd.Owner.Encode(jenc); err != nil {
		return errors.Errorf("invalid receiver format, %q: %q", cmd.Owner.String(), err)
	} else {
		cmd.owner = a
	}

	return nil
}

func (cmd *SignDocumentCommand) createOperation() (operation.Operation, error) { // nolint:dupl
	var items []document.SignDocumentItem
	if i, err := loadOperations(cmd.Seal.Bytes(), cmd.NetworkID.NetworkID()); err != nil {
		return nil, err
	} else {
		for j := range i {
			if t, ok := i[j].(document.SignDocuments); ok {
				items = t.Fact().(document.SignDocumentsFact).Items()
			}
		}
	}

	item := document.NewSignDocumentsItemSingleFile(cmd.DocId, cmd.owner, cmd.Currency.CID)

	if err := item.IsValid(nil); err != nil {
		return nil, err
	} else {
		items = append(items, item)
	}

	fact := document.NewSignDocumentsFact([]byte(cmd.Token), cmd.sender, items)

	var fs []base.FactSign
	sig, err := base.NewFactSignature(cmd.Privatekey, fact, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, err
	}
	fs = append(fs, base.NewBaseFactSign(cmd.Privatekey.Publickey(), sig))

	op, err := document.NewSignDocuments(fact, fs, cmd.Memo)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create sign-document operation")
	}
	return op, nil
}
