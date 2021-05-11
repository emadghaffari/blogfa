package grpc

import (
	pb "blogfa/auth/proto"
	"context"
)

var (
	Service auth = &Auth{}
)

// Auth service
type Auth struct{}

type auth interface {
	PLogin(ctx context.Context, req *pb.PLoginRequest) (*pb.PLoginResponse, error)
}
