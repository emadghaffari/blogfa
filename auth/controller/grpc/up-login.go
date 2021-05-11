package grpc

import (
	"blogfa/auth/pkg/jtrace"
	pb "blogfa/auth/proto"
	"blogfa/auth/service/grpc"
	"context"
)

// login with username or password
// ID is UserName
func (a *Auth) UPLogin(ctx context.Context, req *pb.UPLoginRequest) (*pb.UPLoginResponse, error) {
	span := jtrace.Tracer.StartSpan("up-login")
	defer span.Finish()
	span.SetTag("login", "username password login")

	response, err := grpc.Service.UPLogin(jtrace.Tracer.ContextWithSpan(ctx, span), req)
	if err != nil {
		return response, err
	}

	// return jwt,user
	return response, nil
}
