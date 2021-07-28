package cmds

type SealCommand struct {
	Send                  SendCommand                  `cmd:"" name:"send" help:"send seal to remote mitum node"`
	CreateAccount         CreateAccountCommand         `cmd:"" name:"create-account" help:"create new account"`
	CreateDocument        CreateDocumentCommand        `cmd:"" name:"create-document" help:"create new document"`
	SignDocument          SignDocumentCommand          `cmd:"" name:"sign-document" help:"sign document"`
	Transfer              TransferCommand              `cmd:"" name:"transfer" help:"transfer big"`
	TransferDocument      TransferDocumentCommand      `cmd:"" name:"transfer-document" help:"transfer document"`
	KeyUpdater            KeyUpdaterCommand            `cmd:"" name:"key-updater" help:"update keys"`
	CurrencyRegister      CurrencyRegisterCommand      `cmd:"" name:"currency-register" help:"register new currency"`
	CurrencyPolicyUpdater CurrencyPolicyUpdaterCommand `cmd:"" name:"currency-policy-updater" help:"update currency policy"` // revive:disable-line:line-length-limit
	Sign                  SignSealCommand              `cmd:"" name:"sign" help:"sign seal"`
	SignFact              SignFactCommand              `cmd:"" name:"sign-fact" help:"sign facts of operation seal"`
}

func NewSealCommand() SealCommand {
	return SealCommand{
		Send:                  NewSendCommand(),
		CreateAccount:         NewCreateAccountCommand(),
		CreateDocument:        NewCreateDocumentCommand(),
		SignDocument:          NewSignDocumentCommand(),
		Transfer:              NewTransferCommand(),
		TransferDocument:      NewTransferDocumentCommand(),
		KeyUpdater:            NewKeyUpdaterCommand(),
		CurrencyRegister:      NewCurrencyRegisterCommand(),
		CurrencyPolicyUpdater: NewCurrencyPolicyUpdaterCommand(),
		Sign:                  NewSignSealCommand(),
		SignFact:              NewSignFactCommand(),
	}
}
