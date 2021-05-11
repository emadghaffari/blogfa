package grpc

import (
	"blogfa/auth/domain/jwt"
	"blogfa/auth/pkg/jtrace"
	pb "blogfa/auth/proto"
	"blogfa/auth/service/grpc"
	"context"
	"net/http"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateProvider method for create provider and assign to user
func (a *Auth) CreateProvider(ctx context.Context, req *pb.CreateProviderRequest) (*pb.Response, error) {
	span := jtrace.Tracer.StartSpan("create-provider")
	defer span.Finish()
	span.SetTag("service", "start to create")

	// verify the jwt token
	if _, err := jwt.Model.Verify(req.GetToken()); err != nil {
		return &pb.Response{Message: "invalid user", Status: &pb.Status{Code: http.StatusUnauthorized, Message: "FAILED"}}, status.Errorf(codes.Internal, "user not verified")
	}

	// convert string to int
	userID, err := strconv.Atoi(req.GetUserID())
	if err != nil {
		return &pb.Response{Message: "invalid user id", Status: &pb.Status{Code: http.StatusInternalServerError, Message: "FAILED"}}, status.Errorf(codes.Internal, "invalid user id")
	}

	response, err := grpc.Service.CreateProvider(jtrace.Tracer.ContextWithSpan(ctx, span), req, userID)
	if err != nil {
		return response, err
	}

	child := jtrace.Tracer.ChildOf(span, "register")
	child.SetTag("register", "after create provider")
	defer child.Finish()

	// response if updated successfully
	return response, nil
}
