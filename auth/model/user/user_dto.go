package user

import (
	"blogfa/auth/database/mysql"
	"context"
)

func (u *User) Register(ctx context.Context,user User) error {
	mysql.Storage.Create(ctx, u)

	return nil
}

func (u *User) Get() {

}
