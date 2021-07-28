package cmds

import (
	"golang.org/x/xerrors"

	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/util"

	"github.com/soonkuk/mitum-data/blocksign"
)

type SignDocumentCommand struct {
	*BaseCommand
	OperationFlags
	Sender   AddressFlag    `arg:"" name:"sender" help:"sender address" required:""`
	DocId    BigFlag        `arg:"" name:"documentid" help:"document id" required:""`
	Owner    AddressFlag    `arg:"" name:"owner" help:"owner address" required:""`
	Currency CurrencyIDFlag `arg:"" name:"currency" help:"currency id" required:""`
	Seal     FileLoad       `help:"seal" optional:""`
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
		return xerrors.Errorf("failed to initialize command: %w", err)
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

	if sl, err := loadSealAndAddOperation(
		cmd.Seal.Bytes(),
		cmd.Privatekey,
		cmd.NetworkID.NetworkID(),
		op,
	); err != nil {
		return err
	} else {
		cmd.pretty(cmd.Pretty, sl)
	}

	return nil
}

func (cmd *SignDocumentCommand) parseFlags() error {
	if err := cmd.OperationFlags.IsValid(nil); err != nil {
		return err
	}

	if a, err := cmd.Sender.Encode(jenc); err != nil {
		return xerrors.Errorf("invalid sender format, %q: %w", cmd.Sender.String(), err)
	} else {
		cmd.sender = a
	}
	if a, err := cmd.Owner.Encode(jenc); err != nil {
		return xerrors.Errorf("invalid receiver format, %q: %w", cmd.Owner.String(), err)
	} else {
		cmd.owner = a
	}

	return nil
}

func (cmd *SignDocumentCommand) createOperation() (operation.Operation, error) { // nolint:dupl
	var items []blocksign.SignDocumentItem
	if i, err := loadOperations(cmd.Seal.Bytes(), cmd.NetworkID.NetworkID()); err != nil {
		return nil, err
	} else {
		for j := range i {
			if t, ok := i[j].(blocksign.TransferDocuments); ok {
				items = t.Fact().(blocksign.SignDocumentsFact).Items()
			}
		}
	}

	item := blocksign.NewSignDocumentsItemSingleFile(cmd.DocId.Big, cmd.owner, cmd.Currency.CID)

	if err := item.IsValid(nil); err != nil {
		return nil, err
	} else {
		items = append(items, item)
	}

	fact := blocksign.NewSignDocumentsFact([]byte(cmd.Token), cmd.sender, items)

	var fs []operation.FactSign
	if sig, err := operation.NewFactSignature(cmd.Privatekey, fact, cmd.NetworkID.NetworkID()); err != nil {
		return nil, err
	} else {
		fs = append(fs, operation.NewBaseFactSign(cmd.Privatekey.Publickey(), sig))
	}

	if op, err := blocksign.NewSignDocuments(fact, fs, cmd.Memo); err != nil {
		return nil, xerrors.Errorf("failed to create sign-document operation: %w", err)
	} else {
		return op, nil
	}
}
