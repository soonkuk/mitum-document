package cmds

type DocumentCommand struct {
	CreateBSDocument              CreateBSDocumentCommand              `cmd:"" name:"create-blocksign-document" help:"create new blocksign document"`
	CreateBlockcityUserDocument   CreateBlockcityUserDocumentCommand   `cmd:"" name:"create-blockcity-user-document" help:"create new blockcity user document"`
	CreateBlockcityLandDocument   CreateBlockcityLandDocumentCommand   `cmd:"" name:"create-blockcity-land-document" help:"create new blockcity land document"`
	CreateBlockcityVotingDocument CreateBlockcityVotingDocumentCommand `cmd:"" name:"create-blockcity-voting-document" help:"create new blockcity voting document"`
	UpdateBlockcityUserDocument   UpdateBlockcityUserDocumentCommand   `cmd:"" name:"update-blockcity-user-document" help:"update blockcity user document"`
	UpdateBlockcityLandDocument   UpdateBlockcityLandDocumentCommand   `cmd:"" name:"update-blockcity-land-document" help:"update blockcity land document"`
	UpdateBlockcityVotingDocument UpdateBlockcityVotingDocumentCommand `cmd:"" name:"update-blockcity-voting-document" help:"update blockcity voting document"`
}

func NewDocumentCommand() DocumentCommand {
	return DocumentCommand{
		CreateBSDocument:              NewCreateBSDocumentCommand(),
		CreateBlockcityUserDocument:   NewCreateBlockcityUserDocumentCommand(),
		CreateBlockcityLandDocument:   NewCreateBlockcityLandDocumentCommand(),
		CreateBlockcityVotingDocument: NewCreateBlockcityVotingDocumentCommand(),
		UpdateBlockcityUserDocument:   NewUpdateBlockcityUserDocumentCommand(),
		UpdateBlockcityLandDocument:   NewUpdateBlockcityLandDocumentCommand(),
		UpdateBlockcityVotingDocument: NewUpdateBlockcityVotingDocumentCommand(),
	}
}
