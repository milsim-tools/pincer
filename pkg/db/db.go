package db

import (
	"log/slog"

	"github.com/grafana/dskit/services"
	"github.com/urfave/cli/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	FlagsDSN = "db-dsn"
)

var Flags = []cli.Flag{
	&cli.StringFlag{
		Name:     FlagsDSN,
		Required: true,
		Usage:    "Database data source name",
		EnvVars:  []string{"PINCER_DB_DSN"},
	},
}

type Config struct {
	DSN string
}

func ConfigFromFlags(ctx *cli.Context) Config {
	var config Config
	config.DSN = ctx.String(FlagsDSN)
	return config
}

type Db struct {
	services.Service

	cfg    Config
	logger *slog.Logger

	db *gorm.DB
}

func New(
	logger *slog.Logger,
	cfg Config,
) (*Db, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// TODO: Query logging

	u := &Db{
		cfg:    cfg,
		logger: logger,
		db:     db,
	}

	u.Service = services.NewIdleService(nil, u.stopping)

	return u, nil
}

func (d *Db) stopping(_ error) error {
	return nil
}
