package cmds

import (
	"context"

	"github.com/pkg/errors"

	currencycmds "github.com/spikeekips/mitum-currency/cmds"
	"github.com/spikeekips/mitum/launch/config"
	"github.com/spikeekips/mitum/launch/pm"
	"github.com/spikeekips/mitum/launch/process"
	mongodbstorage "github.com/spikeekips/mitum/storage/mongodb"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/logging"

	"github.com/protoconNet/mitum-document/digest"
)

var ProcessorDigestDatabase pm.Process

func init() {
	if i, err := pm.NewProcess(
		currencycmds.ProcessNameDigestDatabase,
		[]string{process.ProcessNameBlockdata},
		ProcessDigestDatabase,
	); err != nil {
		panic(err)
	} else {
		ProcessorDigestDatabase = i
	}
}

func ProcessDigestDatabase(ctx context.Context) (context.Context, error) {
	var design currencycmds.DigestDesign
	if err := currencycmds.LoadDigestDesignContextValue(ctx, &design); err != nil {
		if errors.Is(err, util.ContextValueNotFoundError) {
			return ctx, nil
		}

		return nil, err
	}

	var mst *mongodbstorage.Database
	if err := currencycmds.LoadDatabaseContextValue(ctx, &mst); err != nil {
		return ctx, err
	}

	st, err := loadDigestDatabase(mst, false)
	if err != nil {
		return ctx, err
	}
	var log *logging.Logging
	if err := config.LoadLogContextValue(ctx, &log); err != nil {
		return ctx, err
	}

	_ = st.SetLogging(log)

	return context.WithValue(ctx, currencycmds.ContextValueDigestDatabase, st), nil
}

func loadDigestDatabase(st *mongodbstorage.Database, readonly bool) (*digest.Database, error) {
	mst := st
	ost, err := st.New()
	if err != nil {
		return nil, err
	}

	var dst *digest.Database
	if readonly {
		s, err := digest.NewReadonlyDatabase(mst, ost)
		if err != nil {
			return nil, err
		}
		dst = s
	} else {
		s, err := digest.NewDatabase(mst, ost)
		if err != nil {
			return nil, err
		}
		dst = s
	}

	if err := dst.Initialize(); err != nil {
		return nil, err
	}

	return dst, nil
}
