package cmds

import (
	"context"

	"github.com/pkg/errors"

	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/launch/config"
	"github.com/spikeekips/mitum/launch/pm"
	"github.com/spikeekips/mitum/launch/process"
	"github.com/spikeekips/mitum/storage"
	"github.com/spikeekips/mitum/storage/blockdata/localfs"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/logging"

	"github.com/protoconNet/mitum-document/digest"
	currencycmds "github.com/spikeekips/mitum-currency/cmds"
	"github.com/spikeekips/mitum-currency/currency"
)

var (
	ProcessorDigester      pm.Process
	ProcessorStartDigester pm.Process
)

func init() {
	if i, err := pm.NewProcess(
		currencycmds.ProcessNameDigester,
		[]string{currencycmds.ProcessNameDigestDatabase},
		ProcessDigester,
	); err != nil {
		panic(err)
	} else {
		ProcessorDigester = i
	}

	if i, err := pm.NewProcess(
		currencycmds.ProcessNameStartDigester,
		[]string{currencycmds.ProcessNameDigester},
		ProcessStartDigester,
	); err != nil {
		panic(err)
	} else {
		ProcessorStartDigester = i
	}
}

func ProcessDigester(ctx context.Context) (context.Context, error) {
	var log *logging.Logging
	if err := config.LoadLogContextValue(ctx, &log); err != nil {
		return ctx, err
	}

	var st *digest.Database
	if err := LoadDigestDatabaseContextValue(ctx, &st); err != nil {
		if errors.Is(err, util.ContextValueNotFoundError) {
			return ctx, nil
		}

		return ctx, err
	}

	di := digest.NewDigester(st, nil)
	_ = di.SetLogging(log)

	return context.WithValue(ctx, currencycmds.ContextValueDigester, di), nil
}

func ProcessStartDigester(ctx context.Context) (context.Context, error) {
	var di *digest.Digester
	if err := LoadDigesterContextValue(ctx, &di); err != nil {
		if errors.Is(err, util.ContextValueNotFoundError) {
			return ctx, nil
		}

		return ctx, err
	}

	return ctx, di.Start()
}

func HookDigesterFollowUp(ctx context.Context) (context.Context, error) {
	var log *logging.Logging
	if err := config.LoadLogContextValue(ctx, &log); err != nil {
		return ctx, err
	}

	log.Log().Debug().Msg("digester trying to follow up")

	var mst storage.Database
	if err := process.LoadDatabaseContextValue(ctx, &mst); err != nil {
		return ctx, err
	}

	var st *digest.Database
	if err := LoadDigestDatabaseContextValue(ctx, &st); err != nil {
		if errors.Is(err, util.ContextValueNotFoundError) {
			return ctx, nil
		}

		return ctx, err
	}

	switch m, found, err := mst.LastManifest(); {
	case err != nil:
		return ctx, err
	case !found:
		log.Log().Debug().Msg("last manifest not found")
	case m.Height() > st.LastBlock():
		log.Log().Info().
			Int64("last_manifest", m.Height().Int64()).
			Int64("last_block", st.LastBlock().Int64()).
			Msg("new blocks found to digest")

		if err := digestFollowup(ctx, m.Height()); err != nil {
			log.Log().Error().Err(err).Msg("failed to follow up")

			return ctx, err
		}
		log.Log().Info().Msg("digested new blocks")
	default:
		log.Log().Info().Msg("digested blocks is up-to-dated")
	}

	return ctx, nil
}

func digestFollowup(ctx context.Context, height base.Height) error {
	var st *digest.Database
	if err := LoadDigestDatabaseContextValue(ctx, &st); err != nil {
		return err
	}

	var bd *localfs.Blockdata
	if err := util.LoadFromContextValue(ctx, process.ContextValueBlockdata, &bd); err != nil {
		return err
	}

	var cp *currency.CurrencyPool
	if err := currencycmds.LoadCurrencyPoolContextValue(ctx, &cp); err != nil {
		return err
	}

	if height <= st.LastBlock() {
		return nil
	}

	lastBlock := st.LastBlock()
	if lastBlock < base.PreGenesisHeight {
		lastBlock = base.PreGenesisHeight
	}

	for i := lastBlock; i <= height; i++ {
		_, blk, err := localfs.LoadBlock(bd, i)
		if err != nil {
			return err
		}

		if err := digest.DigestBlock(ctx, st, blk); err != nil {
			return err
		}

		if err := st.SetLastBlock(blk.Height()); err != nil {
			return err
		}
	}

	return nil
}
