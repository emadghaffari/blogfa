package provider

import (
	"blogfa/auth/model/user"
	"context"

	"gorm.io/gorm"
)

var (
	Model ProviderInterface = &Provider{}
)

// ProviderInterface interface
type ProviderInterface interface {
	Register(ctx context.Context, user Provider) error
	Update(ctx context.Context, prov Provider) error
	Get(ctx context.Context, table string, query interface{}, args ...interface{}) (Provider, error)
	Search(ctx context.Context, from, to int, search string) ([]Provider, error)
}

// Provider struct
type Provider struct {
	gorm.Model
	UserID      uint      `json:"-"`
	User        user.User `json:"user" gorm:"foreignKey:UserID;references:ID"`
	FixedNumber string    `json:"fixedNumber" validate:"required" gorm:"type:varchar(100);"`
	Company     string    `json:"company" validate:"required" gorm:"type:varchar(100);"`
	Card        string    `json:"card" validate:"required" gorm:"type:varchar(100);"`
	CardNumber  string    `json:"cardNumber" validate:"required" gorm:"type:varchar(25);"`
	ShebaNumber string    `json:"shebaNumber" validate:"required" gorm:"type:varchar(100);"`
	Address     string    `json:"address" validate:"required" gorm:"type:varchar(250);"`
}
