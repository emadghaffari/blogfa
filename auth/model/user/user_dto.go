package user

import (
	"blogfa/auth/pkg/jtrace"
	"context"
)

func (u *User) Register(ctx context.Context, user User) error {
	// mysql.Storage.Create(ctx, u)
	span, _ := jtrace.Tracer.SpanFromContext(ctx, "Register_user")
	defer span.Finish()
	span.SetTag("register", "register user model")

	return nil
}

func (u *User) Get() {

}
