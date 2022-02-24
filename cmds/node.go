package cmds

import (
	currencycmds "github.com/spikeekips/mitum-currency/cmds"
	mitumcmds "github.com/spikeekips/mitum/launch/cmds"
)

type NodeCommand struct {
	Init          currencycmds.InitCommand       `cmd:"" help:"initialize node"`
	Run           RunCommand                     `cmd:"" help:"run node"`
	Info          currencycmds.NodeInfoCommand   `cmd:"" help:"node information"`
	StartHandover mitumcmds.StartHandoverCommand `cmd:"" name:"start-handover" help:"start handover"`
}

func NewNodeCommand() (NodeCommand, error) {
	initCommand, err := currencycmds.NewInitCommand(false)
	if err != nil {
		return NodeCommand{}, err
	}

	runCommand, err := NewRunCommand(false)
	if err != nil {
		return NodeCommand{}, err
	}

	return NodeCommand{
		Init:          initCommand,
		Run:           runCommand,
		Info:          currencycmds.NewNodeInfoCommand(),
		StartHandover: mitumcmds.NewStartHandoverCommand(),
	}, nil
}
