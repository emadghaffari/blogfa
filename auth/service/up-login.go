package service

import (
	"blogfa/auth/pkg/jtrace"
	pb "blogfa/auth/proto"
	"context"
)

func (a *Auth) UPLogin(ctx context.Context, req *pb.UPLoginRequest) (*pb.UPLoginResponse, error) {
	span := jtrace.Tracer.StartSpan("up-login")
	defer span.Finish()

	return &pb.UPLoginResponse{}, nil
}
