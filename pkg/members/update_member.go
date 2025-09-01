package members

import (
	"context"

	membersv1 "github.com/milsim-tools/pincer/pkg/api/gen/milsimtools/members/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (m *Members) UpdateMember(ctx context.Context, req *membersv1.UpdateMemberRequest) (*membersv1.UnitMember, error) {
	return &membersv1.UnitMember{}, status.Errorf(codes.Unimplemented, "method UpdateMember not implemented")
}
