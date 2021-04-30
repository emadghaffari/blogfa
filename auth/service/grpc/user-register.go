package grpc

import (
	"blogfa/auth/domain/user"
	"blogfa/auth/model"
	"blogfa/auth/pkg/cript"
	"blogfa/auth/pkg/jtrace"
	pb "blogfa/auth/proto"
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

	// create new user requested.
	_, err = user.Model.Register(jtrace.Tracer.ContextWithSpan(ctx, span), model.User{
		Username:  req.GetUsername(),
		Password:  &password,
		Name:      req.GetName(),
		LastName:  req.GetLastName(),
		Phone:     req.GetPhone(),
		Email:     req.GetEmail(),
		BirthDate: req.GetBirthDate(),
		Gender:    req.GetGender().String(),
		RoleID:    1, // USER
	})
	if err != nil {
		return &pb.Response{Message: fmt.Sprintf("ERROR: %s", err.Error()), Status: &pb.Status{Code: http.StatusInternalServerError, Message: "FAILED"}}, status.Errorf(codes.Internal, "error in store user: %s", err.Error())
	}

	child := jtrace.Tracer.ChildOf(span, "register")
	child.SetTag("register", "after register user")
	defer child.Finish()

	return &pb.Response{Message: "user created successfully", Status: &pb.Status{Code: http.StatusOK, Message: "SUCCESS"}}, nil
}
