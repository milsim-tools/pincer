package members

import (
	"log/slog"

	"github.com/grafana/dskit/services"
	membersv1 "github.com/milsim-tools/pincer/pkg/api/gen/milsimtools/members/v1"
	unitsv1 "github.com/milsim-tools/pincer/pkg/api/gen/milsimtools/units/v1"
	usersv1 "github.com/milsim-tools/pincer/pkg/api/gen/milsimtools/users/v1"
	"github.com/milsim-tools/pincer/pkg/db"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	FlagUsersGrpcAddr = "members-users-grpc-addr"
	FlagUnitsGrpcAddr = "members-units-grpc-addr"
)

var Flags = []cli.Flag{
	&cli.StringFlag{
		Name:    FlagUsersGrpcAddr,
		Value:   "localhost:9000",
		EnvVars: []string{"PINCER_MEMBERS_USERS_GRPC_ADDR"},
	},

	&cli.StringFlag{
		Name:    FlagUnitsGrpcAddr,
		Value:   "localhost:9000",
		EnvVars: []string{"PINCER_MEMBERS_UNITS_GRPC_ADDR"},
	},
}

type Config struct {
	UsersGrpcAddr string
	UnitsGrpcAddr string
}

func ConfigFromFlags(ctx *cli.Context) Config {
	var config Config

	config.UsersGrpcAddr = ctx.String(FlagUsersGrpcAddr)
	config.UnitsGrpcAddr = ctx.String(FlagUnitsGrpcAddr)

	return config
}

type Members struct {
	membersv1.MembersServiceServer
	services.Service

	cfg    Config
	logger *slog.Logger

	db *db.Db

	users usersv1.UsersServiceClient
	units unitsv1.UnitsServiceClient
}

func New(
	logger *slog.Logger,
	cfg Config,
	db *db.Db,
) (*Members, error) {
	u := &Members{
		cfg:    cfg,
		logger: logger,
		db:     db,
	}

	if err := db.Db.AutoMigrate(&MembersUnitMember{}); err != nil {
		return nil, err
	}

	u.Service = services.NewIdleService(nil, nil)

	return u, nil
}

func (u *Members) UsersClient() (usersv1.UsersServiceClient, error) {
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

func (u *Members) UnitsClient() (unitsv1.UnitsServiceClient, error) {
	if u.units != nil {
		return u.units, nil
	}

	unitsConn, err := grpc.NewClient(u.cfg.UnitsGrpcAddr, grpc.WithTransportCredentials(
		insecure.NewCredentials(),
	))
	if err != nil {
		return nil, err
	}
	units := unitsv1.NewUnitsServiceClient(unitsConn)

	u.units = units
	return u.units, nil
}
