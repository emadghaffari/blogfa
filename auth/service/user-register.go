package service

import (
	"blogfa/auth/model/user"
	"blogfa/auth/pkg/cript"
	"blogfa/auth/pkg/jtrace"
	pb "blogfa/auth/proto"
	"context"
	"fmt"
)

// RegisterUser, for create a new user
func (a *Auth) RegisterUser(ctx context.Context, req *pb.UserRegisterRequest) (*pb.UserRegisterResponse, error) {
	span := jtrace.Tracer.StartSpan("register-user")
	defer span.Finish()
	span.SetTag("register", "register user")

	password, err := cript.Hash(req.GetPassword())
	if err != nil {
		return &pb.UserRegisterResponse{Message: "ERROR"}, fmt.Errorf("error in hash password: %s", err.Error())
	}

	// create new user requested.
	_, err = user.Model.Register(jtrace.Tracer.ContextWithSpan(ctx, span), user.User{
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
		return &pb.UserRegisterResponse{Message: "ERROR"}, fmt.Errorf("error in store user: %s", err.Error())
	}

	child := jtrace.Tracer.ChildOf(span, "register")
	child.SetTag("register", "after register user")
	defer child.Finish()

	return &pb.UserRegisterResponse{Message: "DONE"}, nil
}
