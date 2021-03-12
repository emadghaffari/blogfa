package service

import (
	pb "blogfa/auth/proto"
	"context"
	"fmt"
)

// Auth service
type Auth struct{}

func (a *Auth) RegisterUser(ctx context.Context, req *pb.UserRegisterRequest) (*pb.UserRegisterResponse, error) {
	fmt.Println("Hello user")
	return &pb.UserRegisterResponse{}, nil
}
func (a *Auth) RegisterProvider(ctx context.Context, req *pb.ProviderRegisterRequest) (*pb.ProviderRegisterResponse, error) {
	return &pb.ProviderRegisterResponse{}, nil
}
func (a *Auth) UPLogin(ctx context.Context, req *pb.UPLoginRequest) (*pb.UPLoginResponse, error) {
	return &pb.UPLoginResponse{}, nil
}
func (a *Auth) PLogin(ctx context.Context, req *pb.PLoginRequest) (*pb.PLoginResponse, error) {
	return &pb.PLoginResponse{}, nil
}
