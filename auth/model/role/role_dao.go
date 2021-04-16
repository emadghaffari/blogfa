package role

import (
	"blogfa/auth/model/permission"

	"gorm.io/gorm"
)

// Role struct
type Role struct {
	gorm.Model
	Name        string                   `json:"name" gorm:"unique;not null;type:varchar(30);"`
	Permissions []*permission.Permission `json:"permissions" gorm:"many2many:roles_permissions;association_foreignkey:ID;foreignkey:ID"`
}
