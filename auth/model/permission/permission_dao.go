package permission

import "gorm.io/gorm"

// Permission struct
type Permission struct {
	gorm.Model
	Name string `json:"name" gorm:"unique;not null;type:varchar(30);"`
}
