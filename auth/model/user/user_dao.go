package user

import "time"

// User model
type User struct {
	ID        uint64    `gorm:"primaryKey"`
	Username  string    `json:"username" validate:"required" gorm:"unique;not null;type:varchar(100);"`
	Password  *string   `json:"password" validate:"required,gte=7" gorm:"type:varchar(100);"`
	Name      string    `json:"name" validate:"required" gorm:"type:varchar(100);"`
	LastName  string    `json:"lastName" validate:"required" gorm:"type:varchar(100);"`
	Phone     string    `json:"phone" validate:"required" gorm:"type:varchar(25);"`
	Email     string    `json:"email" gorm:"type:varchar(100);"`
	BirthDate string    `json:"birthDate" gorm:"type:varchar(50);"`
	Gender    string    `json:"gender" gorm:"type:varchar(20);"`
	RoleID    uint64    `json:"-"`
	Role      Role      `json:"role" gorm:"foreignKey:RoleID;references:ID"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// Role struct
type Role struct {
	ID          uint64        `json:"-" gorm:"primaryKey"`
	Name        string        `json:"name" gorm:"unique;not null;type:varchar(30);"`
	Permissions []*Permission `json:"permissions" gorm:"many2many:roles_permissions;association_foreignkey:ID;foreignkey:ID"`
	CreatedAt   time.Time     `json:"-"`
	UpdatedAt   time.Time     `json:"-"`
}

// Permission struct
type Permission struct {
	ID        uint64    `json:"-" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"unique;not null;type:varchar(30);"`
	Role      []*Role   `json:"-" gorm:"many2many:roles_permissions;association_foreignkey:ID;foreignkey:ID"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
