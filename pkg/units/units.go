package units

import (
	"log/slog"

	"github.com/grafana/dskit/services"
	unitsv1 "github.com/milsim-tools/pincer/pkg/api/gen/milsimtools/units/v1"
	usersv1 "github.com/milsim-tools/pincer/pkg/api/gen/milsimtools/users/v1"
	"github.com/milsim-tools/pincer/pkg/db"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	FlagUsersGrpcAddr = "users-grpc-addr"
)

var Flags = []cli.Flag{
	&cli.StringFlag{
		Name:    FlagUsersGrpcAddr,
		Value:   "localhost:9000",
		EnvVars: []string{"PINCER_UNITS_USERS_GRPC_ADDR"},
	},
}

type Config struct {
	UsersGrpcAddr string
}

func ConfigFromFlags(ctx *cli.Context) Config {
	var config Config

	config.UsersGrpcAddr = ctx.String(FlagUsersGrpcAddr)

	return config
}

type Units struct {
	unitsv1.UnitsServiceServer
	services.Service

	cfg    Config
	logger *slog.Logger

	db *db.Db

	users usersv1.UsersServiceClient
}

func New(
	logger *slog.Logger,
	cfg Config,
	db *db.Db,
) (*Units, error) {
	u := &Units{
		cfg:    cfg,
		logger: logger,
		db:     db,
	}

	if err := db.Db.AutoMigrate(&UnitsUnit{}); err != nil {
		return nil, err
	}

	u.Service = services.NewIdleService(nil, nil)

	return u, nil
}

func (u *Units) UsersClient() (usersv1.UsersServiceClient, error) {
	if u.users != nil {
		return u.users, nil
	}

	usersConn, err := grpc.NewClient(u.cfg.UsersGrpcAddr, grpc.WithTransportCredentials(
		insecure.NewCredentials(),
	))
	if err != nil {
		return nil, err
	}
	users := usersv1.NewUsersServiceClient(usersConn)

	u.users = users
	return u.users, nil
}
