package units

import (
	"log/slog"

	"github.com/grafana/dskit/services"
	unitsv1 "github.com/milsim-tools/pincer/pkg/api/gen/milsimtools/units/v1"
	"github.com/milsim-tools/pincer/pkg/server"
	"github.com/urfave/cli/v2"
)

var Flags = []cli.Flag{}

type Config struct{}

func ConfigFromFlags(ctx *cli.Context) Config {
	var config Config
	return config
}

type Units struct {
	unitsv1.UnitsServiceServer
	services.Service

	cfg    Config
	logger *slog.Logger

	server *server.Server
}

func New(
	logger *slog.Logger,
	cfg Config,
	serv *server.Server,
) (*Units, error) {
	u := &Units{
		cfg:    cfg,
		logger: logger,
		server: serv,
	}

	u.Service = services.NewIdleService(nil, nil)

	return u, nil
}
