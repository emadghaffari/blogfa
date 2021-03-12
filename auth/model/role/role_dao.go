package role

import (
	"blogfa/auth/model/permission"
	"time"
)

// Role struct
type Role struct {
	ID          uint64                   `json:"-" gorm:"primaryKey"`
	Name        string                   `json:"name" gorm:"unique;not null;type:varchar(30);"`
	Permissions []*permission.Permission `json:"permissions" gorm:"many2many:roles_permissions;association_foreignkey:ID;foreignkey:ID"`
	CreatedAt   time.Time                `json:"-"`
	UpdatedAt   time.Time                `json:"-"`
}
