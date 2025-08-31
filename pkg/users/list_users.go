package users

import (
	"context"
	"fmt"
	"strings"

	usersv1 "github.com/milsim-tools/pincer/pkg/api/gen/milsimtools/users/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func (s *Users) ListUsers(ctx context.Context, req *usersv1.ListUsersRequest) (*usersv1.ListUsersResponse, error) {
	qb := gorm.G[UsersUser](s.db.Db).Limit(int(req.PageSize))

	// TODO: Cursor pagination

	if req.OrderBy != "" {
		parts := strings.SplitSeq(req.OrderBy, ",")
		for part := range parts {
			order := strings.Split(strings.TrimSpace(part), " ")
			if len(order) != 2 {
				return &usersv1.ListUsersResponse{}, status.Error(
					codes.InvalidArgument,
					"invalid order_by format",
				)
			}
			qb = qb.Order(fmt.Sprintf("%s %s", order[0], order[1]))
		}
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
