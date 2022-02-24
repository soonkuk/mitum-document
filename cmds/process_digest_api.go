package cmds

import (
	"context"
	"crypto/tls"

	"github.com/pkg/errors"

	currencycmds "github.com/spikeekips/mitum-currency/cmds"
	currencydigest "github.com/spikeekips/mitum-currency/digest"
	"github.com/spikeekips/mitum/launch/config"
	"github.com/spikeekips/mitum/launch/pm"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/logging"

	"github.com/protoconNet/mitum-document/digest"
)

var (
	ProcessorDigestAPI      pm.Process
	ProcessorStartDigestAPI pm.Process
)

func init() {
	if i, err := pm.NewProcess(currencycmds.ProcessNameDigestAPI, []string{currencycmds.ProcessNameDigestDatabase}, ProcessDigestAPI); err != nil {
		panic(err)
	} else {
		ProcessorDigestAPI = i
	}

	if i, err := pm.NewProcess(
		currencycmds.ProcessNameStartDigestAPI,
		[]string{currencycmds.ProcessNameDigestDatabase, currencycmds.ProcessNameDigestAPI},
		ProcessStartDigestAPI,
	); err != nil {
		panic(err)
	} else {
		ProcessorStartDigestAPI = i
	}
}

func ProcessStartDigestAPI(ctx context.Context) (context.Context, error) {
	var nt *currencydigest.HTTP2Server
	if err := currencycmds.LoadDigestNetworkContextValue(ctx, &nt); err != nil {
		if errors.Is(err, util.ContextValueNotFoundError) {
			return ctx, nil
		}

		return ctx, err
	}

	return ctx, nt.Start()
}

func ProcessDigestAPI(ctx context.Context) (context.Context, error) {
	var design currencycmds.DigestDesign
	if err := currencycmds.LoadDigestDesignContextValue(ctx, &design); err != nil {
		if errors.Is(err, util.ContextValueNotFoundError) {
			return ctx, nil
		}

		return ctx, err
	}

	var log *logging.Logging
	if err := config.LoadLogContextValue(ctx, &log); err != nil {
		return ctx, err
	}

	var networkLog *logging.Logging
	if err := config.LoadNetworkLogContextValue(ctx, &networkLog); err != nil {
		return ctx, err
	}

	if design.Network() == nil {
		log.Log().Debug().Msg("digest api disabled; empty network")

		return ctx, nil
	}

	var st *digest.Database
	if err := LoadDigestDatabaseContextValue(ctx, &st); err != nil {
		log.Log().Debug().Err(err).Msg("digest api disabled; empty database")

		return ctx, nil
	} else if st == nil {
		log.Log().Debug().Msg("digest api disabled; empty database")

		return ctx, nil
	}

	log.Log().Info().
		Str("bind", design.Network().Bind().String()).
		Str("publish", design.Network().ConnInfo().String()).
		Msg("trying to start http2 server for digest API")

	var nt *currencydigest.HTTP2Server
	var certs []tls.Certificate
	if design.Network().Bind().Scheme == "https" {
		certs = design.Network().Certs()
	}

	if sv, err := currencydigest.NewHTTP2Server(
		design.Network().Bind().Host,
		design.Network().ConnInfo().URL().Host,
		certs,
	); err != nil {
		return ctx, err
	} else if err := sv.Initialize(); err != nil {
		return ctx, err
	} else {
		_ = sv.SetLogging(networkLog)

		nt = sv
	}

	return context.WithValue(ctx, currencycmds.ContextValueDigestNetwork, nt), nil
}
