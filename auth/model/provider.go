package model

import (
	"gorm.io/gorm"
)

// Provider struct
type Provider struct {
	gorm.Model
	UserID      uint   `json:"-"`
	User        User   `json:"user" gorm:"foreignKey:UserID;references:ID"`
	FixedNumber string `json:"fixedNumber" validate:"required" gorm:"type:varchar(100);"`
	Company     string `json:"company" validate:"required" gorm:"type:varchar(100);"`
	Card        string `json:"card" validate:"required" gorm:"type:varchar(100);"`
	CardNumber  string `json:"cardNumber" validate:"required" gorm:"type:varchar(25);"`
	ShebaNumber string `json:"shebaNumber" validate:"required" gorm:"type:varchar(100);"`
	Address     string `json:"address" validate:"required" gorm:"type:varchar(250);"`
}
