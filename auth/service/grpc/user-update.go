package grpc

import (
	"blogfa/auth/domain/user"
	"blogfa/auth/model"
	"blogfa/auth/pkg/jtrace"
	pb "blogfa/auth/proto"
	"context"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UpdateUser method for update users
func (a *Auth) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.Response, error) {
	span := jtrace.Tracer.StartSpan("update-user")
	defer span.Finish()
	span.SetTag("service", "get details in service")

	// update spesific user with userID
	if err := user.Model.Update(jtrace.Tracer.ContextWithSpan(ctx, span), model.User{
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
		}, status.Errorf(codes.Internal, "user not updated successfully")
	}

	return &pb.Response{
		Message: "user successfully updated",
		Status: &pb.Status{
			Code:    http.StatusOK,
			Message: "SUCCESS",
		},
	}, nil
}
