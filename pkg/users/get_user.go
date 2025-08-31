package users

import (
	"context"
	"errors"

	usersv1 "github.com/milsim-tools/pincer/pkg/api/gen/milsimtools/users/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func (s *Users) GetUser(ctx context.Context, req *usersv1.GetUserRequest) (*usersv1.UserView, error) {
	qb := gorm.G[UsersUser](s.db.Db).Select("*")

	if id := req.GetId(); id != "" {
		qb = qb.Where("id = ?", id)
	} else if email := req.GetEmail(); email != "" {
		qb = qb.Where("email = ?", email)
	} else if username := req.GetUsername(); username != "" {
		qb = qb.Where("username = ?", username)
	}

	user, err := qb.First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &usersv1.UserView{}, status.Error(
				codes.NotFound,
				"unit not found",
			)
		}

		return &usersv1.UserView{}, status.Error(
			codes.Internal,
			"failed to query unit: "+ err.Error(),
		)
	}

	return &usersv1.UserView{
		User: user.Proto(),
		UnitCount: 0,
	}, nil
}
