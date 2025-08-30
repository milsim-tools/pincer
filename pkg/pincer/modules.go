package pincer

import (
	"context"
	"fmt"

	"github.com/grafana/dskit/services"
	"github.com/milsim-tools/pincer/pkg/server"
)

const (
	Server = "server"

	All     = "all"
	Backend = "backend"
)

func (p *Pincer) initServer() (services.Service, error) {
	server, err := server.New(p.logger.With("module", Server), p.Config.Server)
	if err != nil {
		return nil, err
	}
	p.Server = server

	servicesToWaitFor := func() []services.Service {
		svs := []services.Service(nil)
		for m, s := range p.serviceMap {
			// Server should not wait for itself.
			if m != Server {
				svs = append(svs, s)
			}
		}
		return svs
	}

	s := p.newServerService(p.Server, servicesToWaitFor)

	return s, nil
}

func (p *Pincer) newServerService(serv *server.Server, servicesToWaitFor func() []services.Service) services.Service {
	serverDone := make(chan error, 1)

	runFn := func(ctx context.Context) error {
		go func() {
			defer close(serverDone)
			serverDone <- serv.Run()
		}()

		select {
		case <-ctx.Done():
			return nil
		case err := <-serverDone:
			if err != nil {
				return err
			}
			return fmt.Errorf("server stopped unexpectedly")
		}
	}

	stoppingFn := func(_ error) error {
		// wait until all modules are done, and then shutdown server.
		for _, s := range servicesToWaitFor() {
			_ = s.AwaitTerminated(context.Background())
		}

		// shutdown HTTP and gRPC servers (this also unblocks Run)
		serv.Shutdown()

		// if not closed yet, wait until server stops.
		<-serverDone
		p.logger.Info("server stopped")
		return nil
	}

	return services.NewBasicService(nil, runFn, stoppingFn)
}
