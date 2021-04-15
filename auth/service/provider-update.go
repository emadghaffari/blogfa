package service

import (
	pb "blogfa/auth/proto"
	"context"
)

// UpdateProvider method for update provider
func (a *Auth) UpdateProvider(ctx context.Context, req *pb.UpdateProviderRequest) (*pb.Response, error) {
	return &pb.Response{}, nil
}
