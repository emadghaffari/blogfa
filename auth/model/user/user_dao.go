package user

import (
	"blogfa/auth/model/role"
	"context"
	"time"
)

var (
	Model UserInterface = &User{}
)

type UserInterface interface {
	Register(ctx context.Context, user User) (*User, error)
	Get(ctx context.Context, table string, query interface{}, args ...interface{}) (*User, error)
	Update(ctx context.Context, user User) error
}

// User model
type User struct {
	ID        uint64    `gorm:"primaryKey"`
	Username  string    `json:"username" validate:"required" gorm:"unique;not null;type:varchar(100);"`
	Password  *string   `json:"-" validate:"required,gte=7" gorm:"type:varchar(100);"`
	Name      string    `json:"name" validate:"required" gorm:"type:varchar(100);"`
	LastName  string    `json:"lastName" validate:"required" gorm:"type:varchar(100);"`
	Phone     string    `json:"phone" validate:"required" gorm:"type:varchar(25);"`
	Email     string    `json:"email" gorm:"type:varchar(100);"`
	BirthDate string    `json:"birthDate" gorm:"type:varchar(50);"`
	Gender    string    `json:"gender" gorm:"type:varchar(20);"`
	RoleID    uint64    `json:"-"`
	Role      role.Role `json:"role" gorm:"foreignKey:RoleID;references:ID"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
