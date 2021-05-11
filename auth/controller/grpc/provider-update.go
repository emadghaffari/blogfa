package grpc

import (
	"blogfa/auth/domain/jwt"
	"blogfa/auth/pkg/jtrace"
	pb "blogfa/auth/proto"
	"blogfa/auth/service/grpc"
	"context"
	"fmt"
	"net/http"
	"strconv"

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

	// convert string to int for userID
	ID, err := strconv.Atoi(req.GetID())
	if err != nil {
		return &pb.Response{
			Message: fmt.Sprintf("invalid userID: %s", err.Error()),
			Status: &pb.Status{
				Code:    http.StatusUnauthorized,
				Message: "FAILED",
			},
		}, status.Errorf(codes.Internal, "invalid userID")
	}

	response, err := grpc.Service.UpdateProvider(jtrace.Tracer.ContextWithSpan(ctx, span), req, ID)
	if err != nil {
		return response, err
	}

	// response if updated successfully
	return response, nil

}
