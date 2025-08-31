package units

import (
	"context"
	"errors"

	"github.com/milsim-tools/pincer/internal/models"
	unitsv1 "github.com/milsim-tools/pincer/pkg/api/gen/milsimtools/units/v1"
	usersv1 "github.com/milsim-tools/pincer/pkg/api/gen/milsimtools/users/v1"
	"github.com/oklog/ulid/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func (s *Units) CreateUnit(ctx context.Context, req *unitsv1.CreateUnitRequest) (*unitsv1.Unit, error) {
	client, err := s.UsersClient()
	if err != nil {
		return &unitsv1.Unit{}, status.Error(
			codes.Internal,
			"failed to connect to users service: "+err.Error(),
		)
	}

	if _, err := client.GetUser(ctx, &usersv1.GetUserRequest{
		Value: &usersv1.GetUserRequest_Id{Id: req.Unit.OwnerId},
	}); err != nil {
		if errors.Is(err, status.Error(codes.NotFound, "user not found")) {
			return &unitsv1.Unit{}, status.Error(
				codes.InvalidArgument,
				"owner_id does not correspond to an existing user",
			)
		}
		return &unitsv1.Unit{}, status.Error(
			codes.Internal,
			"failed to call users service: "+err.Error(),
		)
	}

	unit := &UnitsUnit{
		Model: models.Model{
			ID: ulid.Make().String(),
		},
		DisplayName: req.Unit.DisplayName,
		Slug:        req.Unit.Slug,
		Description: req.Unit.Description,
		OwnerID:     req.Unit.OwnerId,
	}

	if err := gorm.G[UnitsUnit](s.db.Db).Create(ctx, unit); err != nil {
		return &unitsv1.Unit{}, status.Error(
			codes.Internal,
			"failed to create unit: "+err.Error(),
		)
	}

	return unit.Proto(), nil
}
