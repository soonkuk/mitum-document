package cmds

import (
	"golang.org/x/xerrors"

	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/util"

	"github.com/soonkuk/mitum-data/blocksign"
)

type CreateDocumentCommand struct {
	*BaseCommand
	OperationFlags
	Sender   AddressFlag    `arg:"" name:"sender" help:"sender address" required:""`
	FileHash string         `arg:"" name:"filehash" help:"filehash" required:""`
	Currency CurrencyIDFlag `arg:"" name:"currency" help:"currency id" required:""`

	Signers []AddressFlag `name:"signers" help:"signers for document"`
	// Keys      []KeyFlag `name:"key" help:"key for new document account (ex: \"<public key>,<weight>\")" sep:"@"`
	Seal   FileLoad `help:"seal" optional:""`
	sender base.Address
}

func NewCreateDocumentCommand() CreateDocumentCommand {
	return CreateDocumentCommand{
		BaseCommand: NewBaseCommand("create-document-operation"),
	}
}

func (cmd *CreateDocumentCommand) Run(version util.Version) error { // nolint:dupl
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

func (cmd *CreateDocumentCommand) parseFlags() error {
	if err := cmd.OperationFlags.IsValid(nil); err != nil {
		return err
	}

	if a, err := cmd.Sender.Encode(jenc); err != nil {
		return xerrors.Errorf("invalid sender format, %q: %w", cmd.Sender.String(), err)
	} else {
		cmd.sender = a
	}

	return nil
}

func (cmd *CreateDocumentCommand) createOperation() (operation.Operation, error) { // nolint:dupl
	var items []blocksign.CreateDocumentsItem
	if i, err := loadOperations(cmd.Seal.Bytes(), cmd.NetworkID.NetworkID()); err != nil {
		return nil, err
	} else {
		for j := range i {
			if t, ok := i[j].(blocksign.CreateDocuments); ok {
				items = t.Fact().(blocksign.CreateDocumentsFact).Items()
			}
		}
	}

	//TODO : Signers 추가
	var signers []base.Address
	for i := range cmd.Signers {
		if signer, err := cmd.Signers[i].Encode(jenc); err != nil {
			return nil, err
		} else {
			signers = append(signers, signer)
		}
	}

	item := blocksign.NewCreateDocumentsItemSingleFile(blocksign.FileHash(cmd.FileHash), signers, cmd.Currency.CID)

	if err := item.IsValid(nil); err != nil {
		return nil, err
	} else {
		items = append(items, item)
	}

	fact := blocksign.NewCreateDocumentsFact([]byte(cmd.Token), cmd.sender, items)

	var fs []operation.FactSign
	if sig, err := operation.NewFactSignature(cmd.Privatekey, fact, cmd.NetworkID.NetworkID()); err != nil {
		return nil, err
	} else {
		fs = append(fs, operation.NewBaseFactSign(cmd.Privatekey.Publickey(), sig))
	}

	if op, err := blocksign.NewCreateDocuments(fact, fs, cmd.Memo); err != nil {
		return nil, xerrors.Errorf("failed to create create-account operation: %w", err)
	} else {
		return op, nil
	}
}
