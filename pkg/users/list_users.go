package users

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/milsim-tools/pincer/internal/helpers"
	usersv1 "github.com/milsim-tools/pincer/pkg/api/gen/milsimtools/users/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func (s *Users) ListUsers(ctx context.Context, req *usersv1.ListUsersRequest) (*usersv1.ListUsersResponse, error) {
	limit := int(req.PageSize)
	if limit <= 0 || limit > 100 {
		limit = 100
	}

	qb := gorm.G[UsersUser](s.db.Db).Limit(limit)

	// TODO: Cursor pagination
	if req.PageToken != "" {
		// The user is trying to paginate
		cursor, err := helpers.CursorFromString(req.PageToken)
		if err != nil {
			return &usersv1.ListUsersResponse{}, status.Error(
				codes.InvalidArgument,
				"invalid page_token format",
			)
		}

		if cursor.OrderBy != req.OrderBy {
			return &usersv1.ListUsersResponse{}, status.Error(
				codes.InvalidArgument,
				"invalid page_token, order_by mismatch",
			)
		}

		qb.Where("created_at > ?", cursor.CreatedAt)
	}

	if req.OrderBy != "" {
		parts := strings.SplitSeq(req.OrderBy, ",")
		for part := range parts {
			order := strings.Split(strings.TrimSpace(part), " ")
			if len(order) != 2 || order[0] != "created_at" || slices.Contains([]string{"asc", "desc", "ASC", "DESC"}, order[1]) {
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

	var pageToken string

	if len(users) > 0 {
		pc := helpers.PaginationCursor{
			OrderBy:   req.OrderBy,
			CreatedAt: users[len(users)-1].CreatedAt.Format(time.RFC3339),
		}

		pageToken, err = pc.String()
		if err != nil {
			return &usersv1.ListUsersResponse{}, status.Error(
				codes.Internal,
				"failed to create page_token: "+err.Error(),
			)
		}
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
		NextPageToken: pageToken,
	}

	return resp, nil
}
