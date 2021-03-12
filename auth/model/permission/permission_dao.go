package permission

import (
	"time"
)

// Permission struct
type Permission struct {
	ID        uint64    `json:"-" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"unique;not null;type:varchar(30);"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
