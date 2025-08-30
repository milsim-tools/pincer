package server

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/grafana/dskit/services"
	"github.com/milsim-tools/pincer/internal/middleware"
	"github.com/milsim-tools/pincer/internal/signals"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	FlagHTTPBindAddr     = "http-bind-addr"
	FlagHTTPReadTimeout  = "http-read-timeout"
	FlagHTTPWriteTimeout = "http-write-timeout"
	FlagHTTPIdleTimeout  = "http-idle-timeout"

	FlagGRPCBindAddr     = "grpc-bind-addr"
	FlagGRPCReadTimeout  = "grpc-read-timeout"
	FlagGRPCWriteTimeout = "grpc-write-timeout"
	FlagGRPCIdleTimeout  = "grpc-idle-timeout"
)

// SignalHandler used by Server.
type SignalHandler interface {
	// Starts the signals handler. This method is blocking, and returns only after signal is received,
	// or "Stop" is called.
	Loop()

	// Stop blocked "Loop" method.
	Stop()
}

var Flags = []cli.Flag{
	&cli.StringFlag{
		Name:    FlagHTTPBindAddr,
		Value:   ":8080",
		Usage:   "The HTTP address to bind the server to.",
		EnvVars: []string{"HEARTH_HTTP_BIND_ADDR"},
	},

	&cli.DurationFlag{
		Name:    FlagHTTPReadTimeout,
		Value:   30 * time.Second,
		Usage:   "Maximum time for reading the request body.",
		EnvVars: []string{"HEARTH_HTTP_READ_TIMEOUT"},
	},

	&cli.DurationFlag{
		Name:    FlagHTTPWriteTimeout,
		Value:   30 * time.Second,
		Usage:   "Maximum time for writing the response body.",
		EnvVars: []string{"HEARTH_HTTP_WRITE_TIMEOUT"},
	},

	&cli.DurationFlag{
		Name:    FlagHTTPIdleTimeout,
		Value:   30 * time.Second,
		Usage:   "Maximum time to wait for another request when keep-alives are used.",
		EnvVars: []string{"HEARTH_HTTP_IDLE_TIMEOUT"},
	},

	&cli.StringFlag{
		Name:    FlagGRPCBindAddr,
		Value:   ":9000",
		Usage:   "The gRPC address to bind the server to.",
		EnvVars: []string{"HEARTH_GRPC_BIND_ADDR"},
	},

	&cli.DurationFlag{
		Name:    FlagGRPCReadTimeout,
		Value:   30 * time.Second,
		Usage:   "Maximum time for reading the request body.",
		EnvVars: []string{"HEARTH_GRPC_READ_TIMEOUT"},
	},

	&cli.DurationFlag{
		Name:    FlagGRPCWriteTimeout,
		Value:   30 * time.Second,
		Usage:   "Maximum time for writing the response body.",
		EnvVars: []string{"HEARTH_GRPC_WRITE_TIMEOUT"},
	},

	&cli.DurationFlag{
		Name:    FlagGRPCIdleTimeout,
		Value:   30 * time.Second,
		Usage:   "Maximum time to wait for another request when keep-alives are used.",
		EnvVars: []string{"HEARTH_GRPC_IDLE_TIMEOUT"},
	},
}

type Config struct {
	HTTPBindAddr     string
	HTTPReadTimeout  time.Duration
	HTTPWriteTimeout time.Duration
	HTTPIdleTimeout  time.Duration

	GRPCBindAddr     string
	GRPCReadTimeout  time.Duration
	GRPCWriteTimeout time.Duration
	GRPCIdleTimeout  time.Duration
}

func ConfigFromFlags(ctx *cli.Context) Config {
	var config Config

	config.HTTPBindAddr = ctx.String(FlagHTTPBindAddr)
	config.HTTPReadTimeout = ctx.Duration(FlagHTTPReadTimeout)
	config.HTTPWriteTimeout = ctx.Duration(FlagHTTPWriteTimeout)
	config.HTTPIdleTimeout = ctx.Duration(FlagHTTPIdleTimeout)

	config.GRPCBindAddr = ctx.String(FlagGRPCBindAddr)
	config.GRPCReadTimeout = ctx.Duration(FlagGRPCReadTimeout)
	config.GRPCWriteTimeout = ctx.Duration(FlagGRPCWriteTimeout)
	config.GRPCIdleTimeout = ctx.Duration(FlagGRPCIdleTimeout)

	return config
}

type Server struct {
	services.Service

	config       Config
	handler      SignalHandler
	httpListener net.Listener
	grpcListener net.Listener

	HTTPServer *http.Server
	GRPCServer *grpc.Server

	logger *slog.Logger
}

func New(logger *slog.Logger, config Config) (*Server, error) {
	httpListener, err := net.Listen("tcp", config.HTTPBindAddr)
	if err != nil {
		return nil, err
	}

	grpcListener, err := net.Listen("tcp", config.GRPCBindAddr)
	if err != nil {
		return nil, err
	}

	logger.Info("server listening on addr", "http", httpListener.Addr(), "grpc", grpcListener.Addr())

	httpServer := http.Server{
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	serverLog := middleware.GRPCServerLog{
		Log:                      logger,
		WithRequest:              true,
		DisableRequestSuccessLog: false,
	}

	grpcMiddleware := []grpc.UnaryServerInterceptor{serverLog.UnaryServerInterceptor}
	grpcStreamMiddleware := []grpc.StreamServerInterceptor{serverLog.StreamServerInterceptor}

	grpcOptions := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(grpcMiddleware...),
		grpc.ChainStreamInterceptor(grpcStreamMiddleware...),
	}

	grpcServer := grpc.NewServer(grpcOptions...)

	srv := Server{
		config:       config,
		HTTPServer:   &httpServer,
		GRPCServer:   grpcServer,
		logger:       logger,
		httpListener: httpListener,
		grpcListener: grpcListener,
		handler:      signals.NewHandler(logger),
	}

	return &srv, nil
}

func (s *Server) Run() error {
	errChan := make(chan error, 1)

	// Wait for a signal
	go func() {
		s.handler.Loop()
		select {
		case errChan <- nil:
		default:
		}
	}()

	go func() {
		err := s.HTTPServer.Serve(s.httpListener)
		if err == http.ErrServerClosed {
			err = nil
		}

		select {
		case errChan <- err:
		default:
		}
	}()

	reflection.Register(s.GRPCServer)

	go func() {
		err := s.GRPCServer.Serve(s.grpcListener)
		handleGRPCError(err, errChan)
	}()

	return <-errChan
}

func (s *Server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_ = s.HTTPServer.Shutdown(ctx)
	s.GRPCServer.GracefulStop()
}

// handleGRPCError consolidates GRPC Server error handling by sending
// any error to errChan except for grpc.ErrServerStopped which is ignored.
func handleGRPCError(err error, errChan chan error) {
	if err == grpc.ErrServerStopped {
		err = nil
	}

	select {
	case errChan <- err:
	default:
	}
}
