package members

import (
	"context"

	membersv1 "github.com/milsim-tools/pincer/pkg/api/gen/milsimtools/members/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (m *Members) DeleteMember(ctx context.Context, req *membersv1.DeleteMemberRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, status.Errorf(codes.Unimplemented, "method DeleteMember not implemented")
}
