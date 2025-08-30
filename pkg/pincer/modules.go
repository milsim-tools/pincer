package pincer

import (
	"github.com/grafana/dskit/services"
	"github.com/milsim-tools/pincer/pkg/server"
)

const (
	Server = "server"

	All     = "all"
	Backend = "backend"
)

func (p *Pincer) initServer() (services.Service, error) {
	s, err := server.New(p.logger.With("module", Server), p.Config.Server)
	if err != nil {
		return nil, err
	}
	p.Server = s
	return s, nil
}
