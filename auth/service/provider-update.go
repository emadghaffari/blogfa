package service

import (
	"blogfa/auth/model/jwt"
	"blogfa/auth/model/provider"
	"blogfa/auth/pkg/jtrace"
	pb "blogfa/auth/proto"
	"context"
	"fmt"
	"net/http"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// UpdateProvider method for update provider
func (a *Auth) UpdateProvider(ctx context.Context, req *pb.UpdateProviderRequest) (*pb.Response, error) {
	span := jtrace.Tracer.StartSpan("update-provider")
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
		}, status.Errorf(codes.Internal, "user not verified")
	}

	// convert string to int for userID
	userID, err := strconv.Atoi(req.GetUserID())
	if err != nil {
		return &pb.Response{
			Message: fmt.Sprintf("invalid userID: %s", err.Error()),
			Status: &pb.Status{
				Code:    http.StatusUnauthorized,
				Message: "FAILED",
			},
		}, status.Errorf(codes.Internal, "invalid userID")
	}

	// convert string to int for userID
	ID, err := strconv.Atoi(req.GetID())
	if err != nil {
		return &pb.Response{
			Message: fmt.Sprintf("invalid userID: %s", err.Error()),
			Status: &pb.Status{
				Code:    http.StatusUnauthorized,
				Message: "FAILED",
			},
		}, status.Errorf(codes.Internal, "invalid userID")
	}

	// try to update provider
	if err := provider.Model.Update(
		jtrace.Tracer.ContextWithSpan(ctx, span),
		provider.Provider{
			Model:       gorm.Model{ID: uint(ID)},
			UserID:      uint(userID),
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
