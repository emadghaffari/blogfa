package user

import (
	"blogfa/auth/database/mysql"
	"blogfa/auth/pkg/jtrace"
	"context"
)

// Register method for register new user
func (u *User) Register(ctx context.Context, user User) (*User, error) {
	span, _ := jtrace.Tracer.SpanFromContext(ctx, "register_user")
	defer span.Finish()
	span.SetTag("register", "register user model")

	tx := mysql.Storage.GetDatabase().Begin()

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	defer tx.Commit()

	return &user, nil
}

// Get method, get a user with table name, query and args for search
func (u *User) Get(ctx context.Context, table string, query interface{}, args ...interface{}) (*User, error) {
	span, _ := jtrace.Tracer.SpanFromContext(ctx, "get user model")
	defer span.Finish()
	span.SetTag("get", "get user model")

	tx := mysql.Storage.GetDatabase().Begin()

	var user = User{}
	if err := tx.Preload("Role").Preload("Role.Permissions").Table(table).Where(query, args...).First(&user); err.Error != nil {
		tx.Rollback()
		return nil, err.Error
	}
	defer tx.Commit()

	return &user, nil
}

func (u User) Update(ctx context.Context, user User) error {

	return nil
}
