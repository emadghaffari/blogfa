package service

import (
	pb "blogfa/auth/proto"
	"context"
)

// AssignProviderToUser method, Provider to existing user
func (a *Auth) AssignProviderToUser(ctx context.Context, req *pb.AssignProviderToUserRequest) (*pb.Response, error) {
	return &pb.Response{}, nil
}
