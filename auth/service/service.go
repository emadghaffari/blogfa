package service

import (
	pb "blogfa/auth/proto"
	"context"

	"github.com/opentracing/opentracing-go"
)

// Auth service
type Auth struct{}

func (a *Auth) RegisterUser(ctx context.Context, req *pb.UserRegisterRequest) (*pb.UserRegisterResponse, error) {
	span := opentracing.GlobalTracer().StartSpan("register-user")
	defer span.Finish()

	return &pb.UserRegisterResponse{}, nil
}
func (a *Auth) RegisterProvider(ctx context.Context, req *pb.ProviderRegisterRequest) (*pb.ProviderRegisterResponse, error) {
	span := opentracing.GlobalTracer().StartSpan("register-provider")
	defer span.Finish()

	return &pb.ProviderRegisterResponse{}, nil
}
func (a *Auth) UPLogin(ctx context.Context, req *pb.UPLoginRequest) (*pb.UPLoginResponse, error) {
	span := opentracing.GlobalTracer().StartSpan("up-login")
	defer span.Finish()

	return &pb.UPLoginResponse{}, nil
}
func (a *Auth) PLogin(ctx context.Context, req *pb.PLoginRequest) (*pb.PLoginResponse, error) {
	span := opentracing.GlobalTracer().StartSpan("p-login")
	defer span.Finish()

	return &pb.PLoginResponse{}, nil
}
