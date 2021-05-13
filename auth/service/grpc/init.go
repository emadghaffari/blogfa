package grpc

import (
	"blogfa/auth/model"
	pb "blogfa/auth/proto"
	"context"
)

// Service var
var (
	Service auth = &Auth{}
)

// Auth service
type Auth struct{}

type auth interface {
	PLogin(ctx context.Context, req *pb.PLoginRequest) (*pb.PLoginResponse, error)
	CreateProvider(ctx context.Context, req *pb.CreateProviderRequest, userID int) (*pb.Response, error)
	RegisterProvider(ctx context.Context, req *pb.ProviderRegisterRequest) (*pb.Response, error)
	SearchProvider(ctx context.Context, req *pb.SearchRequest) ([]model.Provider, error)
	UpdateProvider(ctx context.Context, req *pb.UpdateProviderRequest, id int) (*pb.Response, error)
	UPLogin(ctx context.Context, req *pb.UPLoginRequest) (*pb.UPLoginResponse, error)
	RegisterUser(ctx context.Context, req *pb.UserRegisterRequest, password string) (*pb.Response, error)
	SearchUser(ctx context.Context, from, to int, search string) ([]model.User, error)
	UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.Response, error)
}
