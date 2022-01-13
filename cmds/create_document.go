package cmds

import (
	"bytes"

	"github.com/pkg/errors"

	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/key"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/base/seal"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"

	"github.com/soonkuk/mitum-blocksign/blocksign"
	currencycmds "github.com/spikeekips/mitum-currency/cmds"
	mitumcmds "github.com/spikeekips/mitum/launch/cmds"
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
	Seal       mitumcmds.FileLoad          `help:"seal" optional:""`
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

func (cmd *CreateDocumentCommand) parseFlags() error {
	if err := cmd.OperationFlags.IsValid(nil); err != nil {
		return err
	}

	a, err := cmd.Sender.Encode(jenc)
	if err != nil {
		return errors.Wrapf(err, "invalid sender format, %q", cmd.Sender.String())
	}
	cmd.sender = a

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
	i, err := loadOperations(cmd.Seal.Bytes(), cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, err
	}
	var items []blocksign.CreateDocumentsItem
	for j := range i {
		if t, ok := i[j].(blocksign.CreateDocuments); ok {
			items = t.Fact().(blocksign.CreateDocumentsFact).Items()
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
	}
	items = append(items, item)

	fact := blocksign.NewCreateDocumentsFact([]byte(cmd.Token), cmd.sender, items)

	sig, err := base.NewFactSignature(cmd.Privatekey, fact, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, err
	}
	fs := []base.FactSign{
		base.NewBaseFactSign(cmd.Privatekey.Publickey(), sig),
	}

	op, err := blocksign.NewCreateDocuments(fact, fs, cmd.Memo)
	if err != nil {
		return nil, errors.Errorf("failed to create create-account operation: %q", err)
	}
	return op, nil
}

func loadOperations(b []byte, networkID base.NetworkID) ([]operation.Operation, error) {
	if len(bytes.TrimSpace(b)) < 1 {
		return nil, nil
	}

	var sl seal.Seal
	if s, err := LoadSeal(b, networkID); err != nil {
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

func LoadSeal(b []byte, networkID base.NetworkID) (seal.Seal, error) {
	if len(bytes.TrimSpace(b)) < 1 {
		return nil, errors.Errorf("empty input")
	}

	var sl seal.Seal
	if err := encoder.Decode(b, jenc, &sl); err != nil {
		return nil, err
	}

	if err := sl.IsValid(networkID); err != nil {
		return nil, errors.Wrap(err, "invalid seal")
	}

	return sl, nil
}

func LoadSealAndAddOperation(
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
	if s, err := LoadSeal(b, networkID); err != nil {
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
