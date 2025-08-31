package users

import (
	"log/slog"

	"github.com/grafana/dskit/services"
	usersv1 "github.com/milsim-tools/pincer/pkg/api/gen/milsimtools/users/v1"
	"github.com/milsim-tools/pincer/pkg/db"
	"github.com/urfave/cli/v2"
)

var Flags = []cli.Flag{}

type Config struct{}

func ConfigFromFlags(ctx *cli.Context) Config {
	var config Config
	return config
}

type Users struct {
	usersv1.UsersServiceServer
	services.Service

	cfg    Config
	logger *slog.Logger

	db *db.Db
}

func New(
	logger *slog.Logger,
	cfg Config,
	db *db.Db,
) (*Users, error) {
	u := &Users{
		cfg:    cfg,
		logger: logger,
		db:     db,
	}

	if err := db.Db.AutoMigrate(&UsersUser{}); err != nil {
		return nil, err
	}

	u.Service = services.NewIdleService(nil, nil)

	return u, nil
}
