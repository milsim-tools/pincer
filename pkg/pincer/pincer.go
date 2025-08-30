package pincer

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/grafana/dskit/services"
	"github.com/milsim-tools/pincer/internal/modules"
	"github.com/milsim-tools/pincer/internal/signals"
	"github.com/milsim-tools/pincer/pkg/server"
	"github.com/urfave/cli/v2"
	"go.uber.org/atomic"
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

func init() {
	Flags = append(Flags, server.Flags...)
}

type Config struct {
	Target        []string
	ShutdownDelay time.Duration
	Version       string

	Server server.Config
}

func ConfigFromFlags(version string, ctx *cli.Context) Config {
	var config Config

	config.Target = ctx.StringSlice(FlagTarget)
	config.ShutdownDelay = ctx.Duration(FlagShutdownDelay)
	config.Version = version

	config.Server = server.ConfigFromFlags(ctx)

	return config
}

type Pincer struct {
	Config Config

	logger *slog.Logger

	ModuleManager *modules.Manager
	serviceMap    map[string]services.Service
	deps          map[string][]string
	SignalHandler *signals.Handler

	Server *server.Server
}

func New(logger *slog.Logger, cfg Config) (*Pincer, error) {
	pincer := &Pincer{
		Config: cfg,
		logger: logger,
	}

	if err := pincer.setupModuleManager(); err != nil {
		return nil, err
	}

	return pincer, nil
}

func (p *Pincer) setupModuleManager() error {
	mm := modules.NewManager(p.logger.WithGroup("module-manager"))

	mm.RegisterModule(Server, p.initServer, modules.UserInvisibleModule)

	mm.RegisterModule(All, nil)
	mm.RegisterModule(Backend, nil)

	deps := map[string][]string{
		// Groups
		All:     {},
		Backend: {},
	}

	for mod, targets := range deps {
		if err := mm.AddDependency(mod, targets...); err != nil {
			return err
		}
	}

	p.deps = deps
	p.ModuleManager = mm

	return nil
}

func (p *Pincer) Run(ctx context.Context) error {
	startTime := time.Now()

	serviceMap, err := p.ModuleManager.InitModuleServices(p.Config.Target...)
	if err != nil {
		return err
	}

	p.serviceMap = serviceMap

	// get all services, create service manager and tell it to start
	var servs []services.Service
	for _, s := range serviceMap {
		servs = append(servs, s)
	}

	sm, err := services.NewManager(servs...)
	if err != nil {
		return err
	}

	shutdownRequested := atomic.NewBool(false)

	// Let's listen for events from this manager, and log them.
	logHook := func(msg, key string) func() {
		return func() {
			started := startTime
			p.logger.Info(msg, key, time.Since(started))
		}
	}
	healthy := logHook("hearth started", "startup_time")
	stopped := logHook("hearth stopped", "running_time")
	serviceFailed := func(service services.Service) {
		// if any service fails, stop entire hearth
		sm.StopAsync()

		// let's find out which module failed
		for m, s := range serviceMap {
			if s == service {
				if service.FailureCase() == modules.ErrStopProcess {
					p.logger.Info("received stop signal via return error", "module", m, "error", service.FailureCase())
				} else {
					p.logger.Error("module failed", "module", m, "error", service.FailureCase())
				}
				return
			}
		}

		p.logger.Error("module failed", "module", "unknown", "error", service.FailureCase())
	}

	sm.AddListener(services.NewManagerListener(healthy, stopped, serviceFailed))

	p.SignalHandler = signals.NewHandler(p.logger)
	go func() {
		p.SignalHandler.Loop()
		shutdownRequested.Store(true)

		if p.Config.ShutdownDelay > 0 {
			p.logger.Info(fmt.Sprintf("waiting %v before shutting down services", p.Config.ShutdownDelay))
			time.Sleep(p.Config.ShutdownDelay)
		}

		sm.StopAsync()
	}()

	// Start all services. This can really only fail if some service is already
	// in other state than New, which should not be the case.
	err = sm.StartAsync(context.Background())
	if err == nil {
		// Wait until service manager stops. It can stop in two ways:
		// 1) Signal is received and manager is stopped.
		// 2) Any service fails.
		err = sm.AwaitStopped(context.Background())
	}

	// If there is no error yet (= service manager started and then stopped without problems),
	// but any service failed, report that failure as an error to caller.
	if err == nil {
		if failed := sm.ServicesByState()[services.Failed]; len(failed) > 0 {
			for _, f := range failed {
				if f.FailureCase() != modules.ErrStopProcess {
					// Details were reported via failure listener before
					err = errors.New("failed services")
					break
				}
			}
		}
	}
	return err
}
