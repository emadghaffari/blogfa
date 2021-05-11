package grpc

import (
	"blogfa/auth/domain/jwt"
	"blogfa/auth/pkg/jtrace"
	pb "blogfa/auth/proto"
	"blogfa/auth/service/grpc"
	"context"
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UpdateUser method for update users
func (a *Auth) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.Response, error) {
	span := jtrace.Tracer.StartSpan("update-user")
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

	response, err := grpc.Service.UpdateUser(jtrace.Tracer.ContextWithSpan(context.Background(), span), req)
	if err != nil {
		return response, err
	}

	return response, nil
}
