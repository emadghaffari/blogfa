package grpc

import (
	"blogfa/auth/pkg/cript"
	"blogfa/auth/pkg/jtrace"
	pb "blogfa/auth/proto"
	"blogfa/auth/service/grpc"
	"context"
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RegisterUser, for create a new user
func (a *Auth) RegisterUser(ctx context.Context, req *pb.UserRegisterRequest) (*pb.Response, error) {
	span := jtrace.Tracer.StartSpan("register-user")
	defer span.Finish()
	span.SetTag("register", "register user")

	password, err := cript.Hash(req.GetPassword())
	if err != nil {
		return &pb.Response{Message: fmt.Sprintf("ERROR: %s", err.Error()), Status: &pb.Status{Code: http.StatusInternalServerError, Message: "FAILED"}}, status.Errorf(codes.Internal, "error in hash password: %s", err.Error())
	}

	response, err := grpc.Service.RegisterUser(jtrace.Tracer.ContextWithSpan(ctx, span), req, password)
	if err != nil {
		return response, err
	}

	child := jtrace.Tracer.ChildOf(span, "register")
	child.SetTag("register", "after register user")
	defer child.Finish()

	return response, nil
}
