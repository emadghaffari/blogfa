package service

import (
	"blogfa/auth/pkg/jtrace"
	pb "blogfa/auth/proto"
	"context"
)

// Auth service
type Auth struct{}

func (a *Auth) PLogin(ctx context.Context, req *pb.PLoginRequest) (*pb.PLoginResponse, error) {
	span := jtrace.Tracer.StartSpan("p-login")
	defer span.Finish()

	return &pb.PLoginResponse{}, nil
}
