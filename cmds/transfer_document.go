package cmds

import (
	"golang.org/x/xerrors"

	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/util"

	"github.com/soonkuk/mitum-data/currency"
)

type TransferDocumentCommand struct {
	*BaseCommand
	OperationFlags
	Sender   AddressFlag    `arg:"" name:"sender" help:"sender address" required:""`
	Currency CurrencyIDFlag `arg:"" name:"currency" help:"currency id" required:""`
	Document AddressFlag    `arg:"" name:"document" help:"document address" required:""`
	Receiver AddressFlag    `arg:"" name:"reciever" help:"reciever address" required:""`
	Seal     FileLoad       `help:"seal" optional:""`
	sender   base.Address
	document base.Address
	receiver base.Address
}

func NewTransferDocumentCommand() TransferDocumentCommand {
	return TransferDocumentCommand{
		BaseCommand: NewBaseCommand("transfer-document-operation"),
	}
}

func (cmd *TransferDocumentCommand) Run(version util.Version) error { // nolint:dupl
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

func (cmd *TransferDocumentCommand) parseFlags() error {
	if err := cmd.OperationFlags.IsValid(nil); err != nil {
		return err
	}

	if a, err := cmd.Sender.Encode(jenc); err != nil {
		return xerrors.Errorf("invalid sender format, %q: %w", cmd.Sender.String(), err)
	} else {
		cmd.sender = a
	}
	if a, err := cmd.Document.Encode(jenc); err != nil {
		return xerrors.Errorf("invalid document format, %q: %w", cmd.Document.String(), err)
	} else {
		cmd.document = a
	}
	if a, err := cmd.Receiver.Encode(jenc); err != nil {
		return xerrors.Errorf("invalid receiver format, %q: %w", cmd.Receiver.String(), err)
	} else {
		cmd.receiver = a
	}

	return nil
}

func (cmd *TransferDocumentCommand) createOperation() (operation.Operation, error) { // nolint:dupl
	var items []currency.TransferDocumentsItem
	if i, err := loadOperations(cmd.Seal.Bytes(), cmd.NetworkID.NetworkID()); err != nil {
		return nil, err
	} else {
		for j := range i {
			if t, ok := i[j].(currency.TransferDocuments); ok {
				items = t.Fact().(currency.TransferDocumentsFact).Items()
			}
		}
	}

	item := currency.NewTransferDocumentsItemSingleFile(cmd.document, cmd.receiver, cmd.Currency.CID)

	if err := item.IsValid(nil); err != nil {
		return nil, err
	} else {
		items = append(items, item)
	}

	fact := currency.NewTransferDocumentsFact([]byte(cmd.Token), cmd.sender, items)

	var fs []operation.FactSign
	if sig, err := operation.NewFactSignature(cmd.Privatekey, fact, cmd.NetworkID.NetworkID()); err != nil {
		return nil, err
	} else {
		fs = append(fs, operation.NewBaseFactSign(cmd.Privatekey.Publickey(), sig))
	}

	if op, err := currency.NewTransferDocuments(fact, fs, cmd.Memo); err != nil {
		return nil, xerrors.Errorf("failed to create transfer-document operation: %w", err)
	} else {
		return op, nil
	}
}
