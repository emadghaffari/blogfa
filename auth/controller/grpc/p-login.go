package grpc

import (
	"blogfa/auth/pkg/jtrace"
	pb "blogfa/auth/proto"
	"blogfa/auth/service/grpc"
	"context"
	"net/http"
)

// PLogin, login user with phone number with sms code
func (a *Auth) PLogin(ctx context.Context, req *pb.PLoginRequest) (*pb.PLoginResponse, error) {
	span := jtrace.Tracer.StartSpan("p-login")
	defer span.Finish()
	span.SetTag("login", "phone login")

	response, err := grpc.Service.PLogin(jtrace.Tracer.ContextWithSpan(ctx, span), req)
	if err != nil {
		return &pb.PLoginResponse{
			Message: "ERROR",
			Status: &pb.Status{
				Message: err.Error(),
				Code:    http.StatusInternalServerError,
			},
		}, err
	}

	// return response for check the phone
	return response, nil
}
