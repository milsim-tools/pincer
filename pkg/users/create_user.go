package users

import (
	"context"

	"github.com/milsim-tools/pincer/internal/models"
	usersv1 "github.com/milsim-tools/pincer/pkg/api/gen/milsimtools/users/v1"
	"github.com/oklog/ulid/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func (s *Users) CreateUser(ctx context.Context, req *usersv1.CreateUserRequest) (*usersv1.User, error) {
	user := &UsersUser{
		Model: models.Model{
			ID: ulid.Make().String(),
		},
		SSOID: req.User.SsoId,
		DisplayName: req.User.DisplayName,
		Username: req.User.Username,
		Email: req.User.Email,
		Bio: req.User.Bio,
		AvatarURL: req.User.AvatarUrl,
	}

	if err := gorm.G[UsersUser](s.db.Db).Create(ctx, user); err != nil {
		return &usersv1.User{}, status.Error(
			codes.Internal,
			"failed to create user: "+ err.Error(),
		)
	}

	return user.Proto(), nil
}
