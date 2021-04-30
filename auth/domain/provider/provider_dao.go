package provider

import (
	"blogfa/auth/model"
	pb "blogfa/auth/proto"
	"context"
)

var (
	Model ProviderInterface = &Provider{}
)

// ProviderInterface interface
type ProviderInterface interface {
	Register(ctx context.Context, user model.Provider) error
	Update(ctx context.Context, prov model.Provider) error
	Get(ctx context.Context, table string, query interface{}, args ...interface{}) (model.Provider, error)
	Search(ctx context.Context, from, to int, search string) ([]model.Provider, error)
	ToProto(prvs []*model.Provider) (resp []*pb.Providers)
}

// Provider struct
type Provider struct{}
