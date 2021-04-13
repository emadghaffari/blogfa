package service

import (
	pb "blogfa/auth/proto"
	"context"
)


func (a *Auth) UpdateProvider(ctx context.Context, req *pb.UpdateProviderRequest) (*pb.Response, error)  {
	return &pb.Response{},nil
}