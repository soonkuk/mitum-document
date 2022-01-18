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

type UpdateBlockcityUserDocumentCommand struct {
	*BaseCommand
	currencycmds.OperationFlags
	Sender       currencycmds.AddressFlag    `arg:"" name:"sender" help:"sender address" required:""`
	Gold         currencycmds.BigFlag        `arg:"" name:"gold" help:"gold" required:""`
	Bankgold     currencycmds.BigFlag        `arg:"" name:"bankgold" help:"bankgold" required:""`
	Hp           uint                        `arg:"" name:"hp" help:"hp" required:""`
	Strength     uint                        `arg:"" name:"strength" help:"strength" required:""`
	Agility      uint                        `arg:"" name:"agility" help:"agility" required:""`
	Dexterity    uint                        `arg:"" name:"dexterity" help:"dexterity" required:""`
	Charisma     uint                        `arg:"" name:"charisma" help:"charisma" required:""`
	Intelligence uint                        `arg:"" name:"intelligence" help:"intelligence" required:""`
	Vital        uint                        `arg:"" name:"vital" help:"vital" required:""`
	DocumentId   string                      `arg:"" name:"documentid" help:"document id" required:""`
	Currency     currencycmds.CurrencyIDFlag `arg:"" name:"currency" help:"currency id" required:""`
	Seal         mitumcmds.FileLoad          `help:"seal" optional:""`
	sender       base.Address
}

func NewUpdateBlockcityUserDocumentCommand() UpdateBlockcityUserDocumentCommand {
	return UpdateBlockcityUserDocumentCommand{
		BaseCommand: NewBaseCommand("update-blockcity-user-document-operation"),
	}
}

func (cmd *UpdateBlockcityUserDocumentCommand) Run(version util.Version) error { // nolint:dupl
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

func (cmd *UpdateBlockcityUserDocumentCommand) parseFlags() error {
	if err := cmd.OperationFlags.IsValid(nil); err != nil {
		return err
	}

	a, err := cmd.Sender.Encode(jenc)
	if err != nil {
		return errors.Wrapf(err, "invalid sender format, %q", cmd.Sender.String())
	}
	cmd.sender = a

	return nil
}

func (cmd *UpdateBlockcityUserDocumentCommand) createOperation() (operation.Operation, error) { // nolint:dupl
	i, err := loadOperations(cmd.Seal.Bytes(), cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, err
	}
	var items []blockcity.UpdateDocumentsItem
	for j := range i {
		if t, ok := i[j].(blockcity.UpdateDocuments); ok {
			items = t.Fact().(blockcity.UpdateDocumentsFact).Items()
		}
	}
	info := blockcity.NewDocInfo(cmd.DocumentId, blockcity.CityUserDataType)
	statistics := blockcity.NewUserStatistics(cmd.Hp, cmd.Strength, cmd.Agility, cmd.Dexterity, cmd.Charisma, cmd.Intelligence, cmd.Vital)
	userDoc := blockcity.NewCityUserData(info, cmd.sender, cmd.Gold.Big, cmd.Bankgold.Big, statistics)
	doc := blockcity.NewDocument(userDoc)

	item := blockcity.NewUpdateDocumentsItemImpl(
		doc,
		cmd.Currency.CID,
	)

	if err := item.IsValid(nil); err != nil {
		return nil, err
	}
	items = append(items, item)

	fact := blockcity.NewUpdateDocumentsFact([]byte(cmd.Token), cmd.sender, items)

	sig, err := base.NewFactSignature(cmd.Privatekey, fact, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, err
	}
	fs := []base.FactSign{
		base.NewBaseFactSign(cmd.Privatekey.Publickey(), sig),
	}

	op, err := blockcity.NewUpdateDocuments(fact, fs, cmd.Memo)
	if err != nil {
		return nil, errors.Errorf("failed to create update-blockcity-user-document operation: %q", err)
	}
	return op, nil
}
