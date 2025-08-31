package pincer

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/grafana/dskit/grpcutil"
	"github.com/grafana/dskit/services"
	"github.com/milsim-tools/pincer/internal/modules"
	"github.com/milsim-tools/pincer/internal/signals"
	"github.com/milsim-tools/pincer/pkg/db"
	"github.com/milsim-tools/pincer/pkg/server"
	"github.com/milsim-tools/pincer/pkg/units"
	"github.com/milsim-tools/pincer/pkg/users"
	"github.com/urfave/cli/v2"
	"go.uber.org/atomic"
	"google.golang.org/grpc/health/grpc_health_v1"
)

const (
	FlagTarget        = "target"
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
	Flags = append(Flags, db.Flags...)
	Flags = append(Flags, units.Flags...)
	Flags = append(Flags, users.Flags...)
}

type Config struct {
	Target        []string
	ShutdownDelay time.Duration
	Version       string

	Server server.Config
	Db     db.Config

	Units units.Config
	Users users.Config
}

func ConfigFromFlags(version string, ctx *cli.Context) Config {
	var config Config

	config.Target = ctx.StringSlice(FlagTarget)
	config.ShutdownDelay = ctx.Duration(FlagShutdownDelay)
	config.Version = version

	config.Server = server.ConfigFromFlags(ctx)
	config.Db = db.ConfigFromFlags(ctx)
	config.Units = units.ConfigFromFlags(ctx)
	config.Users = users.ConfigFromFlags(ctx)

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
	Db     *db.Db

	Units *units.Units
	Users *users.Users
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
	mm.RegisterModule(Db, p.initDb, modules.UserInvisibleModule)

	mm.RegisterModule(Units, p.initUnits)
	mm.RegisterModule(Users, p.initUsers)

	mm.RegisterModule(All, nil)
	mm.RegisterModule(Backend, nil)

	deps := map[string][]string{
		Units: {Db, Server},
		Users: {Db, Server},

		// Groups
		All:     {Units, Users},
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

	p.Server.HTTP.Path("/ready").Methods("GET").Handler(p.readyHandler(sm, shutdownRequested))
	grpc_health_v1.RegisterHealthServer(p.Server.GRPCServer, grpcutil.NewHealthCheck(sm))

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

func (p *Pincer) readyHandler(sm *services.Manager, shutdownRequested *atomic.Bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if shutdownRequested.Load() {
			p.logger.Debug("application is stopping")
			http.Error(w, "application is stopping", http.StatusServiceUnavailable)
			return
		}
		if !sm.IsHealthy() {
			msg := bytes.Buffer{}
			msg.WriteString("services not running:\n")

			byState := sm.ServicesByState()
			for state, servs := range byState {
				msg.WriteString(fmt.Sprintf("%v: %d\n", state, len(servs)))
			}

			http.Error(w, msg.String(), http.StatusServiceUnavailable)
			return
		}

		http.Error(w, "ready", http.StatusOK)
	}
}
