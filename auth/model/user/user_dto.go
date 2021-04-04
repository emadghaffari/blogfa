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

	tx := mysql.Storage.GetDatabase().Begin()

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return err
	}
	defer tx.Commit()

	return nil
}

func (u *User) Get() {

}
