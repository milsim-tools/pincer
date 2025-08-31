package users

import (
	"context"
	"errors"

	usersv1 "github.com/milsim-tools/pincer/pkg/api/gen/milsimtools/users/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

func (s *Users) DeleteUser(ctx context.Context, req *usersv1.DeleteUserRequest) (*emptypb.Empty, error) {
	if _, err := gorm.G[UsersUser](s.db.Db).Where("id = ?", req.UserId).First(ctx); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &emptypb.Empty{}, status.Error(
				codes.NotFound,
				"user not found",
			)
		}

		return &emptypb.Empty{}, status.Error(
			codes.Internal,
			"failed to query user: "+ err.Error(),
		)
	}

	if _, err := gorm.G[UsersUser](s.db.Db).Where("id = ?", req.UserId).Delete(ctx); err != nil {
		return &emptypb.Empty{}, status.Error(
			codes.Internal,
			"failed to delete user: "+ err.Error(),
		)
	}

	return &emptypb.Empty{}, nil
}
