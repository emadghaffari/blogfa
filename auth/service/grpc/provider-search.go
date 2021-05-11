package grpc

import (
	"blogfa/auth/domain/provider"
	"blogfa/auth/model"
	"blogfa/auth/pkg/jtrace"
	pb "blogfa/auth/proto"
	"context"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SearchProvider method for search user by provider fields
func (a *Auth) SearchProvider(ctx context.Context, req *pb.SearchRequest) ([]model.Provider, error) {
	span := jtrace.Tracer.StartSpan("search-provider")
	defer span.Finish()
	span.SetTag("service", "start to search")

	// convert string to int
	from, err := strconv.Atoi(req.GetFrom())
	if err != nil {
		return []model.Provider{}, status.Errorf(codes.Internal, "invalid from number")
	}

	// convert string to int
	to, err := strconv.Atoi(req.GetTo())
	if err != nil {
		return []model.Provider{}, status.Errorf(codes.Internal, "invalid to number")
	}

	// search providers
	providers, err := provider.Model.Search(jtrace.Tracer.ContextWithSpan(context.Background(), span), from, to, req.GetSearch())
	if err != nil {
		return []model.Provider{}, status.Errorf(codes.Internal, "internal error for search providers")
	}

	return providers, nil
}
