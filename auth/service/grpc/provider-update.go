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
	"gorm.io/gorm"
)

// UpdateProvider method for update provider
func (a *Auth) UpdateProvider(ctx context.Context, req *pb.UpdateProviderRequest, ID int) (*pb.Response, error) {
	span := jtrace.Tracer.StartSpan("update-provider")
	defer span.Finish()
	span.SetTag("service", "get details in service")

	// try to update provider
	if err := provider.Model.Update(
		jtrace.Tracer.ContextWithSpan(ctx, span),
		model.Provider{
			Model:       gorm.Model{ID: uint(ID)},
			FixedNumber: req.GetFixedNumber(),
			Company:     req.GetCompany(),
			CardNumber:  req.GetCardNumber(),
			ShebaNumber: req.GetShebaNumber(),
			Card:        req.GetCard(),
			Address:     req.GetAddress(),
		},
	); err != nil {
		return &pb.Response{
			Message: "provider not updated successfully",
			Status: &pb.Status{
				Code:    http.StatusInternalServerError,
				Message: "FAILED",
			},
		}, status.Errorf(codes.Internal, "provider not updated successfully")
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
