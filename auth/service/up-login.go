package service

import (
	"blogfa/auth/model/jwt"
	"blogfa/auth/model/permission"
	"blogfa/auth/model/user"
	"blogfa/auth/pkg/cript"
	"blogfa/auth/pkg/jtrace"
	pb "blogfa/auth/proto"
	"context"
	"fmt"
)

// login with username or password
func (a *Auth) UPLogin(ctx context.Context, req *pb.UPLoginRequest) (*pb.UPLoginResponse, error) {
	span := jtrace.Tracer.StartSpan("up-login")
	defer span.Finish()
	span.SetTag("login", "username password login")

	// get user with email or username
	user, err := user.Model.Get(jtrace.Tracer.ContextWithSpan(ctx, span), "users", "username = ? OR email = ?", req.GetUsername(), req.GetUsername())

	// check password and errors
	if err != nil || !cript.CheckHash(req.GetPassword(), *user.Password) {
		return &pb.UPLoginResponse{
			Message: "username or password not matched! ",
			Status: &pb.Status{
				Code:    403,
				Message: "invalid username or password",
			},
		}, fmt.Errorf("invalid username or password")
	}

	// generate jwt token
	jwt, err := jwt.Model.Generate(ctx, user)
	if err != nil {
		return &pb.UPLoginResponse{
			Message: "error in generate accessToken try after 10 seconds!",
			Status: &pb.Status{
				Code:    403,
				Message: "error in generate accessToken try after 10 seconds!",
			},
		}, fmt.Errorf("error in generate accessToken try after 10 seconds!")
	}

	// return jwt,user
	return &pb.UPLoginResponse{
		Message: "user loggedin successfully",
		Token:   jwt.AccessToken,
		User: &pb.User{
			Username:  user.Username,
			Name:      user.Name,
			LastName:  user.LastName,
			Phone:     user.Phone,
			Email:     user.Email,
			BirthDate: user.BirthDate,
			Gender:    pb.User_Gender(pb.User_Gender_value[user.Gender]),
			Role: &pb.Role{
				Name:         user.Role.Name,
				Permissions: permission.ToList(user.Role.Permissions),
			},
		},
		Status: &pb.Status{
			Code:    200,
			Message: "user loggedin successfully",
		},
	}, nil
}
