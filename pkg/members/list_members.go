package members

import (
	"context"

	membersv1 "github.com/milsim-tools/pincer/pkg/api/gen/milsimtools/members/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (m *Members) ListMembers(ctx context.Context, req *membersv1.ListMembersRequest) (*membersv1.ListMembersResponse, error) {
	return &membersv1.ListMembersResponse{}, status.Errorf(codes.Unimplemented, "method ListMembers not implemented")
}
