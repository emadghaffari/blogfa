package user

import (
	"blogfa/auth/database/mysql"
	"blogfa/auth/pkg/jtrace"
	"context"
)

func (u *User) Register(ctx context.Context, user User) error {
	span, _ := jtrace.Tracer.SpanFromContext(ctx, "Register_user")
	defer span.Finish()
	span.SetTag("register", "register user model")

	if err := mysql.Storage.Create(jtrace.Tracer.ContextWithSpan(ctx, span), "users", user); err != nil {
		return err
	}

	return nil
}

func (u *User) Get() {

}
