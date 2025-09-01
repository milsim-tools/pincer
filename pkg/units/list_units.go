package units

import (
	"context"
	"fmt"
	"strings"

	"github.com/milsim-tools/pincer/internal/helpers"
	unitsv1 "github.com/milsim-tools/pincer/pkg/api/gen/milsimtools/units/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func (s *Units) ListUnits(ctx context.Context, req *unitsv1.ListUnitsRequest) (*unitsv1.ListUnitsResponse, error) {
	qb := gorm.G[UnitsUnit](s.db.Db).Limit(int(req.PageSize))

	if cursor, err := helpers.CursorFromString(req.PageToken); err != nil && cursor != nil {
		qb = helpers.ApplyCursor(qb, cursor)
	}

	// TODO: Cursor pagination

	if req.OrderBy != "" {
		parts := strings.SplitSeq(req.OrderBy, ",")
		for part := range parts {
			order := strings.Split(strings.TrimSpace(part), " ")
			if len(order) != 2 {
				return &unitsv1.ListUnitsResponse{}, status.Error(
					codes.InvalidArgument,
					"invalid order_by format",
				)
			}
			qb = qb.Order(fmt.Sprintf("%s %s", order[0], order[1]))
		}
	}

	units, err := qb.Find(ctx)
	if err != nil {
		return &unitsv1.ListUnitsResponse{}, status.Error(
			codes.Internal,
			"failed to query unit: "+ err.Error(),
		)
	}

	var unitViews []*unitsv1.UnitView
	for _, unit := range units {
		unitViews = append(unitViews, &unitsv1.UnitView{
			Unit:        unit.Proto(),
			MemberCount: 0,
			RankCount:   0,
		})
	}

	resp := &unitsv1.ListUnitsResponse{
		Units: unitViews,
		NextPageToken: "",
	}

	return resp, nil
}
