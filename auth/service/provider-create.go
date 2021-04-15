package service

import (
	pb "blogfa/auth/proto"
	"context"
)

// CreateProvider method for create provider and assign to user
func (a *Auth) CreateProvider(ctx context.Context, req *pb.CreateProviderRequest) (*pb.Response, error) {
	return &pb.Response{}, nil
}
