package service

import (
	"blogfa/auth/model/jwt"
	"blogfa/auth/pkg/jtrace"
	pb "blogfa/auth/proto"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AssignProviderToUser method, Provider to existing user
func (a *Auth) AssignProviderToUser(ctx context.Context, req *pb.AssignProviderToUserRequest) (*pb.Response, error) {
	span := jtrace.Tracer.StartSpan("assign-provider-to-user")
	defer span.Finish()
	span.SetTag("service", "assign a provider to user")

	// verify the jwt token
	if _, err := jwt.Model.Verify(req.GetToken()); err != nil {
		return &pb.Response{
			Message: "user not verified",
			Status: &pb.Status{
				Code:    uint32(codes.Internal),
				Message: "invalid user",
			},
		}, status.Errorf(codes.Internal, "user not verified: %s", err.Error())
	}

	return &pb.Response{}, nil
}
