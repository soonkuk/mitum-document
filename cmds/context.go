package cmds

import (
	"context"

	"github.com/protoconNet/mitum-document/digest"
	"github.com/spikeekips/mitum/util"

	currencycmds "github.com/spikeekips/mitum-currency/cmds"
)

func LoadDigestDatabaseContextValue(ctx context.Context, l **digest.Database) error {
	return util.LoadFromContextValue(ctx, currencycmds.ContextValueDigestDatabase, l)
}

func LoadDigesterContextValue(ctx context.Context, l **digest.Digester) error {
	return util.LoadFromContextValue(ctx, currencycmds.ContextValueDigester, l)
}
