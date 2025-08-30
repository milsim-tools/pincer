package units

import (
	"context"

	unitsv1 "github.com/milsim-tools/pincer/pkg/api/gen/milsimtools/units/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Units) GetUnit(context.Context, *unitsv1.GetUnitRequest) (*unitsv1.UnitView, error) {
	return &unitsv1.UnitView{}, status.Error(
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
