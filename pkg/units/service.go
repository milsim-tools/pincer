package units

import (
	"context"

	corev1 "github.com/milsim-tools/pincer/pkg/api/gen/milsimtools/core/v1"
	unitsv1 "github.com/milsim-tools/pincer/pkg/api/gen/milsimtools/units/v1"
	statusproto "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/status"
)

func (s Units) GetUnit(context.Context, *unitsv1.GetUnitRequest) (*unitsv1.UnitView, error) {
	return &unitsv1.UnitView{}, status.ErrorProto(&statusproto.Status{
		Code:    int32(corev1.Code_CODE_UNIMPLEMENTED),
		Message: "GetUnit is not implemented",
	})
}

func (s Units) ListUnits(context.Context, *unitsv1.ListUnitsRequest) (*unitsv1.ListUnitsResponse, error) {
	return &unitsv1.ListUnitsResponse{}, status.ErrorProto(&statusproto.Status{
		Code:    int32(corev1.Code_CODE_UNIMPLEMENTED),
		Message: "ListUnits is not implemented",
	})
}
