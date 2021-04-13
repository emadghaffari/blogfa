package service

import (
	pb "blogfa/auth/proto"
	"context"
)


func (a *Auth) CreateProvider(ctx context.Context, req *pb.CreateProviderRequest) (*pb.Response, error)  {
	return &pb.Response{},nil
}