package users

import (
	"context"
	"errors"
	"slices"

	usersv1 "github.com/milsim-tools/pincer/pkg/api/gen/milsimtools/users/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func (s *Users) UpdateUser(ctx context.Context, req *usersv1.UpdateUserRequest) (*usersv1.User, error) {
	user, err := gorm.G[UsersUser](s.db.Db).Where("id = ?", req.User.Id).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &usersv1.User{}, status.Error(
				codes.NotFound,
				"user not found",
			)
		}

		return &usersv1.User{}, status.Error(
			codes.Internal,
			"failed to query user: "+ err.Error(),
		)
	}

	if slices.Contains(req.UpdateMask.GetPaths(), "user.sso_id") {
		user.SSOID = req.User.SsoId
	}

	if slices.Contains(req.UpdateMask.GetPaths(), "user.display_name") {
		user.DisplayName = req.User.DisplayName
	}

	if slices.Contains(req.UpdateMask.GetPaths(), "user.email") {
		user.Email = req.User.Email
	}

	if slices.Contains(req.UpdateMask.GetPaths(), "user.bio") {
		user.Bio = req.User.Bio
	}

	if slices.Contains(req.UpdateMask.GetPaths(), "user.username") {
		user.Username = req.User.Username
	}

	if slices.Contains(req.UpdateMask.GetPaths(), "user.avatar_url") {
		user.AvatarURL = req.User.AvatarUrl
	}

	if _, err := gorm.G[UsersUser](s.db.Db).Updates(ctx, user); err != nil {
		return &usersv1.User{}, status.Error(
			codes.Internal,
			"failed to update user: "+ err.Error(),
		)
	}

	return user.Proto(), nil
}
