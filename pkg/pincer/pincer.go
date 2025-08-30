package pincer

import (
	"context"
	"log/slog"
	"time"

	"github.com/grafana/dskit/services"
	"github.com/milsim-tools/pincer/internal/modules"
	"github.com/milsim-tools/pincer/internal/signals"
	"github.com/urfave/cli/v2"
)

const (
	FlagTarget = "target"
	FlagShutdownDelay = "shutdown-delay"
)

var Flags = []cli.Flag{
	&cli.StringSliceFlag{
		Name:    FlagTarget,
		Value:   cli.NewStringSlice(All),
		Usage:   "Which services to run when starting the app",
		EnvVars: []string{"PINCER_TARGET"},
	},

	&cli.DurationFlag{
		Name:    FlagShutdownDelay,
		Value:   0 * time.Second,
		Usage:   "How long to wait before shutting down services",
		EnvVars: []string{"PINCER_SHUTDOWN_DELAY"},
	},
}

type Config struct {
	Target        []string
	ShutdownDelay time.Duration
	Version       string
}

func ConfigFromFlags(version string, ctx *cli.Context) Config {
	var config Config

	config.Target = ctx.StringSlice(FlagTarget)
	config.ShutdownDelay = ctx.Duration(FlagShutdownDelay)
	config.Version = version

	return config
}

type Pincer struct {
	Config Config

	logger *slog.Logger

	ModuleManager *modules.Manager
	serviceMap    map[string]services.Service
	deps          map[string][]string
	SignalHandler *signals.Handler
}

func New(logger *slog.Logger, cfg Config) (*Pincer, error) {
	pincer := &Pincer{
		Config: cfg,
		logger: logger,
	}

	return pincer, nil
}

func (p *Pincer) Run(ctx context.Context) error {
	return nil
}
