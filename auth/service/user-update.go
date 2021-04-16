package service

import (
	"blogfa/auth/model/jwt"
	"blogfa/auth/model/user"
	"blogfa/auth/pkg/jtrace"
	pb "blogfa/auth/proto"
	"context"
	"fmt"
	"net/http"
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
		}, nil
	}

	// update spesific user with userID
	if err := user.Model.Update(jtrace.Tracer.ContextWithSpan(ctx, span), user.User{
		Username:  req.GetID(),
		Name:      req.GetName(),
		LastName:  req.GetLastName(),
		Phone:     req.GetPhone(),
		Email:     req.GetEmail(),
		BirthDate: req.GetBirthDate(),
		Gender:    req.GetGender().String(),
		RoleID:    req.GetRole(),
	}); err != nil {
		return &pb.Response{
			Message: "user not updated successfully",
			Status: &pb.Status{
				Code:    http.StatusInternalServerError,
				Message: "FAILED",
			},
		}, nil
	}

	return &pb.Response{
		Message: "user successfully updated",
		Status: &pb.Status{
			Code:    http.StatusOK,
			Message: "SUCCESS",
		},
	}, nil
}
