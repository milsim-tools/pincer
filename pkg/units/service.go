package units

import (
	"context"
	"errors"

	unitsv1 "github.com/milsim-tools/pincer/pkg/api/gen/milsimtools/units/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func (s *Units) GetUnit(ctx context.Context, req *unitsv1.GetUnitRequest) (*unitsv1.UnitView, error) {
	unit, err := gorm.G[UnitsUnit](s.db.Db).Where("id = ?", req.Id).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &unitsv1.UnitView{}, status.Error(
				codes.NotFound,
				"unit not found",
			)
		}

		return &unitsv1.UnitView{}, status.Error(
			codes.Internal,
			"failed to query unit: "+ err.Error(),
		)
	}

	return &unitsv1.UnitView{
		Unit: unit.Proto(),
		MemberCount: 0,
		RankCount: 0,
	}, status.Error(
		codes.Unimplemented,
		"GetUnit is not implemented",
	)
}

func (s *Units) ListUnits(context.Context, *unitsv1.ListUnitsRequest) (*unitsv1.ListUnitsResponse, error) {
	resp := &unitsv1.ListUnitsResponse{
		Units: []*unitsv1.UnitView{
			{
				Unit: &unitsv1.Unit{
					Id:          "unit-1",
					DisplayName: "Unit 1",
					Slug:        "unit-1",
					Description: "Hello, I am Unit 1",
					OwnerId:     "",
				},
			},
		},
		NextPageToken: "",
	}

	return resp, nil
}
