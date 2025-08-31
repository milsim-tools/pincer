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
	}, nil
}
