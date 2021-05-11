package grpc

import (
	"blogfa/auth/domain/user"
	"blogfa/auth/model"
	"blogfa/auth/pkg/jtrace"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SearchUser method for search user by users fields
func (a *Auth) SearchUser(ctx context.Context, from, to int, search string) ([]model.User, error) {
	span := jtrace.Tracer.StartSpan("search-user")
	defer span.Finish()
	span.SetTag("service", "start to search")

	// search users
	users, err := user.Model.Search(jtrace.Tracer.ContextWithSpan(context.Background(), span), from, to, search)
	if err != nil {
		return []model.User{}, status.Errorf(codes.Internal, "internal error for search users")
	}

	return users, nil
}
