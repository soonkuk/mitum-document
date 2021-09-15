package cmds

import (
	"bytes"

	"github.com/pkg/errors"

	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/key"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/base/seal"
	"github.com/spikeekips/mitum/util"

	"github.com/soonkuk/mitum-blocksign/blocksign"
	currencycmds "github.com/spikeekips/mitum-currency/cmds"
)

type CreateDocumentCommand struct {
	*BaseCommand
	currencycmds.OperationFlags
	Sender     currencycmds.AddressFlag    `arg:"" name:"sender" help:"sender address" required:""`
	FileHash   string                      `arg:"" name:"filehash" help:"filehash" required:""`
	Signcode   string                      `arg:"" name:"signcode" help:"signcode" required:""`
	DocumentId currencycmds.BigFlag        `arg:"" name:"documentid" help:"document id" required:""`
	Title      string                      `arg:"" name:"title" help:"title" required:""`
	Size       currencycmds.BigFlag        `arg:"" name:"size" help:"size" required:""`
	Currency   currencycmds.CurrencyIDFlag `arg:"" name:"currency" help:"currency id" required:""`
	Signers    []DocSignFlag               `name:"signers" help:"signers for document (ex: \"<address>,<signcode>\")" sep:"@"`
	Seal       currencycmds.FileLoad       `help:"seal" optional:""`
	sender     base.Address
	signers    []base.Address
	signcodes  []string
}

func NewCreateDocumentCommand() CreateDocumentCommand {
	return CreateDocumentCommand{
		BaseCommand: NewBaseCommand("create-document-operation"),
	}
}

func (cmd *CreateDocumentCommand) Run(version util.Version) error { // nolint:dupl
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

	if sl, err := loadSealAndAddOperation(
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

func (cmd *CreateDocumentCommand) parseFlags() error {
	if err := cmd.OperationFlags.IsValid(nil); err != nil {
		return err
	}

	if a, err := cmd.Sender.Encode(jenc); err != nil {
		return errors.Errorf("invalid sender format, %q: %q", cmd.Sender.String(), err)
	} else {
		cmd.sender = a
	}

	{
		signers := make([]base.Address, len(cmd.Signers))
		signcodes := make([]string, len(cmd.Signers))
		for i := range cmd.Signers {
			if a, err := cmd.Signers[i].AD.Encode(jenc); err != nil {
				return errors.Errorf("invalid sender format, %q: %q", cmd.Signers[i].String(), err)
			} else {
				signers[i] = a
				signcodes[i] = cmd.Signers[i].SC

			}
		}
		cmd.signers = signers
		cmd.signcodes = signcodes
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

	item := blocksign.NewCreateDocumentsItemSingleFile(
		blocksign.FileHash(cmd.FileHash),
		cmd.DocumentId.Big,
		cmd.Signcode,
		cmd.Title,
		cmd.Size.Big,
		cmd.signers,
		cmd.signcodes,
		cmd.Currency.CID,
	)

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
		return nil, errors.Errorf("failed to create create-account operation: %q", err)
	} else {
		return op, nil
	}
}

func loadSeal(b []byte, networkID base.NetworkID) (seal.Seal, error) {
	if len(bytes.TrimSpace(b)) < 1 {
		return nil, errors.Errorf("empty input")
	}

	if sl, err := seal.DecodeSeal(b, jenc); err != nil {
		return nil, err
	} else if err := sl.IsValid(networkID); err != nil {
		return nil, errors.Wrap(err, "invalid seal")
	} else {
		return sl, nil
	}
}

func loadSealAndAddOperation(
	b []byte,
	privatekey key.Privatekey,
	networkID base.NetworkID,
	op operation.Operation,
) (operation.Seal, error) {
	if b == nil {
		bs, err := operation.NewBaseSeal(
			privatekey,
			[]operation.Operation{op},
			networkID,
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create operation.Seal")
		}
		return bs, nil
	}

	var sl operation.Seal
	if s, err := loadSeal(b, networkID); err != nil {
		return nil, err
	} else if so, ok := s.(operation.Seal); !ok {
		return nil, errors.Errorf("seal is not operation.Seal, %T", s)
	} else if _, ok := so.(operation.SealUpdater); !ok {
		return nil, errors.Errorf("seal is not operation.SealUpdater, %T", s)
	} else {
		sl = so
	}

	// NOTE add operation to existing seal
	sl = sl.(operation.SealUpdater).SetOperations([]operation.Operation{op}).(operation.Seal)

	s, err := currencycmds.SignSeal(sl, privatekey, networkID)
	if err != nil {
		return nil, err
	}
	sl = s.(operation.Seal)

	return sl, nil
}

func loadOperations(b []byte, networkID base.NetworkID) ([]operation.Operation, error) {
	if len(bytes.TrimSpace(b)) < 1 {
		return nil, nil
	}

	var sl seal.Seal
	if s, err := loadSeal(b, networkID); err != nil {
		return nil, err
	} else if so, ok := s.(operation.Seal); !ok {
		return nil, errors.Errorf("seal is not operation.Seal, %T", s)
	} else if _, ok := so.(operation.SealUpdater); !ok {
		return nil, errors.Errorf("seal is not operation.SealUpdater, %T", s)
	} else {
		sl = so
	}

	return sl.(operation.Seal).Operations(), nil
}
