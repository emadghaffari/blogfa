package model

import (
	"gorm.io/gorm"
)

// User model
type User struct {
	gorm.Model
	Username  string     `json:"username" validate:"required" gorm:"unique;not null;type:varchar(100);"`
	Password  *string    `json:"-" validate:"required,gte=7" gorm:"type:varchar(100);"`
	Name      string     `json:"name" validate:"required" gorm:"type:varchar(100);"`
	LastName  string     `json:"lastName" validate:"required" gorm:"type:varchar(100);"`
	Phone     string     `json:"phone" validate:"required" gorm:"type:varchar(25);"`
	Email     string     `json:"email" gorm:"type:varchar(100);"`
	BirthDate string     `json:"birthDate" gorm:"type:varchar(50);"`
	Gender    string     `json:"gender" gorm:"type:varchar(20);"`
	RoleID    uint64     `json:"-"`
	Role      Role       `json:"role" gorm:"foreignKey:RoleID;references:ID"`
	Provider  []Provider `json:"provides" gorm:"foreignKey:UserID;references:ID"`
}
