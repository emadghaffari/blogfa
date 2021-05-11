package grpc

import (
	"blogfa/auth/domain/provider"
	"blogfa/auth/model"
	"blogfa/auth/pkg/jtrace"
	pb "blogfa/auth/proto"
	"context"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateProvider method for create provider and assign to user
func (a *Auth) CreateProvider(ctx context.Context, req *pb.CreateProviderRequest, userID int) (*pb.Response, error) {
	span := jtrace.Tracer.StartSpan("create-provider")
	defer span.Finish()
	span.SetTag("service", "start to create")

	// create provider
	if err := provider.Model.Register(jtrace.Tracer.ContextWithSpan(ctx, span), model.Provider{
		UserID:      uint(userID),
		FixedNumber: req.GetFixedNumber(),
		Company:     req.GetCompany(),
		Card:        req.GetCard(),
		CardNumber:  req.GetCardNumber(),
		ShebaNumber: req.GetShebaNumber(),
		Address:     req.GetAddress(),
	}); err != nil {
		return &pb.Response{Message: "internal error: invalid data", Status: &pb.Status{Code: http.StatusInternalServerError, Message: "FAILED"}}, status.Errorf(codes.Internal, "error in store new provider")
	}

	// response if updated successfully
	return &pb.Response{
		Message: "provider successfully updated",
		Status: &pb.Status{
			Code:    http.StatusOK,
			Message: "SUCCESS",
		},
	}, nil
}
