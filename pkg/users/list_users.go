package users

import (
	"context"

	"github.com/milsim-tools/pincer/internal/helpers"
	usersv1 "github.com/milsim-tools/pincer/pkg/api/gen/milsimtools/users/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func (s *Users) ListUsers(ctx context.Context, req *usersv1.ListUsersRequest) (*usersv1.ListUsersResponse, error) {
	qb := gorm.G[UsersUser](s.db.Db).Select("*")
	qb = helpers.ApplyPageLimit(qb, int(req.PageSize))

	cursor, err := helpers.CursorFromString(req.PageToken)
	if err != nil {
		return &usersv1.ListUsersResponse{}, status.Error(
			codes.InvalidArgument,
			"invalid page_token format",
		)
	}

	if cursor != nil {
		qb = helpers.ApplyCursor(qb, cursor)
	}

	users, err := qb.Find(ctx)
	if err != nil {
		return &usersv1.ListUsersResponse{}, status.Error(
			codes.Internal,
			"failed to query unit: "+err.Error(),
		)
	}

	var userViews []*usersv1.UserView
	for _, user := range users {
		userViews = append(userViews, &usersv1.UserView{
			User:      user.Proto(),
			UnitCount: 0,
		})
	}

	resp := &usersv1.ListUsersResponse{
		Users:         userViews,
		NextPageToken: "",
	}

	return resp, nil
}
