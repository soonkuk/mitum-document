package cmds

import currencycmds "github.com/spikeekips/mitum-currency/cmds"

type NodeCommand struct {
	Init InitCommand                  `cmd:"" help:"initialize node"`
	Run  RunCommand                   `cmd:"" help:"run node"`
	Info currencycmds.NodeInfoCommand `cmd:"" help:"node information"`
}

func NewNodeCommand() (NodeCommand, error) {
	initCommand, err := NewInitCommand(false)
	if err != nil {
		return NodeCommand{}, err
	}

	runCommand, err := NewRunCommand(false)
	if err != nil {
		return NodeCommand{}, err
	}

	return NodeCommand{
		Init: initCommand,
		Run:  runCommand,
		Info: currencycmds.NewNodeInfoCommand(),
	}, nil
}
