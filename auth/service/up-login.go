package service

import (
	"blogfa/auth/model/user"
	"blogfa/auth/pkg/cript"
	"blogfa/auth/pkg/jtrace"
	pb "blogfa/auth/proto"
	"context"
	"fmt"
)

func (a *Auth) UPLogin(ctx context.Context, req *pb.UPLoginRequest) (*pb.UPLoginResponse, error) {
	span := jtrace.Tracer.StartSpan("up-login")
	defer span.Finish()
	span.SetTag("register", "username password login")

	password, err := cript.Hash(req.GetPassword())
	if err != nil {
		return &pb.UPLoginResponse{Message: fmt.Sprintf("ERROR: %s", err.Error()), Status: &pb.Response{Code: 400, Message: "FAILED"}}, fmt.Errorf("error in hash password: %s", err.Error())
	}

	user, err := user.Model.Get(jtrace.Tracer.ContextWithSpan(ctx, span), "users", "username = ? OR email = ?", req.GetUsername())
	if err != nil || user.Password != &password {
		return &pb.UPLoginResponse{
			Message: "username or password not matched! ",
			Status: &pb.Response{
				Code:    403,
				Message: "invalid username or password",
			},
		}, fmt.Errorf("invalid username or password")
	}

	return &pb.UPLoginResponse{
		Message: "user loggedin successfully",
		Status: &pb.Response{
			Code:    200,
			Message: "user loggedin successfully",
		},
	}, nil
}
