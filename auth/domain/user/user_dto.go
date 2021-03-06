package user

import (
	"blogfa/auth/client/mysql"
	"blogfa/auth/model"
	"blogfa/auth/pkg/jtrace"
	"blogfa/auth/pkg/logger"
	"context"
	"fmt"

	"go.uber.org/zap"
)

// Register method for register new user
func (u *User) Register(ctx context.Context, user model.User) (*model.User, error) {
	span, _ := jtrace.Tracer.SpanFromContext(ctx, "register_user")
	defer span.Finish()
	span.SetTag("register", "register user model")

	tx := mysql.Storage.GetDatabase().Begin()

	if err := tx.Create(&user).Error; err != nil {
		log := logger.GetZapLogger(false)
		logger.Prepare(log).
			Append(zap.Any("error", fmt.Sprintf("register user: %s", err))).
			Level(zap.ErrorLevel).
			Development().
			Commit("env")
		tx.Rollback()
		return nil, err
	}
	defer tx.Commit()

	return &user, nil
}

// Get method, get a user with table name, query and args for search
func (u *User) Get(ctx context.Context, table string, query interface{}, args ...interface{}) (model.User, error) {
	span, _ := jtrace.Tracer.SpanFromContext(ctx, "get user model")
	defer span.Finish()
	span.SetTag("model", "get user model")

	tx := mysql.Storage.GetDatabase().Begin()

	var user = model.User{}
	if err := tx.Preload("Role").Preload("Role.Permissions").Table(table).Where(query, args...).First(&user).Error; err != nil {
		log := logger.GetZapLogger(false)
		logger.Prepare(log).
			Append(zap.Any("error", fmt.Sprintf("get user: %s", err))).
			Level(zap.ErrorLevel).
			Development().
			Commit("env")
		tx.Rollback()
		return user, err
	}
	defer tx.Commit()

	return user, nil
}

// Update method, for update users
func (u User) Update(ctx context.Context, user model.User) error {
	span, _ := jtrace.Tracer.SpanFromContext(ctx, "get user model")
	defer span.Finish()
	span.SetTag("model", fmt.Sprintf("update user with id: %d", user.ID))

	tx := mysql.Storage.GetDatabase().Begin()

	usr, err := u.Get(ctx, "users", "username = ?", user.ID)
	if err != nil {
		return err
	}

	// update usr fileds
	usr.Name = user.Name
	usr.LastName = user.LastName
	usr.Phone = user.Phone
	usr.Email = user.Email
	usr.BirthDate = user.BirthDate
	usr.Gender = user.Gender
	usr.RoleID = user.RoleID

	if err := tx.Table("users").Where("username = ?", user.Username).Select("*").Updates(&usr).Error; err != nil {
		log := logger.GetZapLogger(false)
		logger.Prepare(log).
			Append(zap.Any("error", fmt.Sprintf("update user error: %s", err))).
			Level(zap.ErrorLevel).
			Development().
			Commit("env")
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

// Search method for search users
func (u User) Search(ctx context.Context, from, to int, search string) ([]model.User, error) {
	span, _ := jtrace.Tracer.SpanFromContext(ctx, "search user model")
	defer span.Finish()
	span.SetTag("model", "search users")

	tx := mysql.Storage.GetDatabase().Begin()

	var users []model.User
	err := tx.
		Preload("Role").
		Preload("Provider").
		Preload("Role.Permissions").
		Table("users").
		Where("username LIKE ?", "%"+search+"%").
		Or("name LIKE ?", "%"+search+"%").
		Or("last_name LIKE ?", "%"+search+"%").
		Or("phone LIKE ?", "%"+search+"%").
		Or("email LIKE ?", "%"+search+"%").
		Or("gender LIKE ?", "%"+search+"%").
		Or("role_id LIKE ?", "%"+search+"%").
		Limit(to - from).
		Offset(from).
		Select("*").
		Find(&users).Error
	if err != nil {
		log := logger.GetZapLogger(false)
		logger.Prepare(log).
			Append(zap.Any("error", fmt.Sprintf("update user error: %s", err))).
			Level(zap.ErrorLevel).
			Development().
			Commit("env")
		tx.Rollback()
		return nil, err
	}
	tx.Commit()

	return users, nil
}
