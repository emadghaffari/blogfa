package user

import (
	"blogfa/auth/model"
	"context"
)

var (
	Model UserInterface = &User{}
)

type UserInterface interface {
	Register(ctx context.Context, user model.User) (*model.User, error)
	Get(ctx context.Context, table string, query interface{}, args ...interface{}) (model.User, error)
	Update(ctx context.Context, user model.User) error
	Search(ctx context.Context, from, to int, search string) ([]model.User, error)
}

// User model
type User struct{}
