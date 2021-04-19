package service

import (
	"blogfa/auth/model/jwt"
	"blogfa/auth/pkg/jtrace"
	pb "blogfa/auth/proto"
	"context"
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UpdateProvider method for update provider
func (a *Auth) UpdateProvider(ctx context.Context, req *pb.UpdateProviderRequest) (*pb.Response, error) {
	span := jtrace.Tracer.StartSpan("update-provider")
	defer span.Finish()
	span.SetTag("service", "get details in service")

	// verify the jwt token
	if _, err := jwt.Model.Verify(req.GetToken()); err != nil {
		return &pb.Response{
			Message: fmt.Sprintf("user not verified: %s", err.Error()),
			Status: &pb.Status{
				Code:    http.StatusUnauthorized,
				Message: "FAILED",
			},
		}, status.Errorf(codes.Internal, "user not verified")
	}

	return &pb.Response{}, nil
}
