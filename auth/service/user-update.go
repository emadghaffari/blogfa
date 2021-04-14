package service

import (
	pb "blogfa/auth/proto"
	"context"
)

// UpdateUser method for update users
func (a *Auth) UpdateUser(ctx context.Context, stream *pb.UpdateUserRequest) (*pb.Response, error) {
	return &pb.Response{}, nil
}
