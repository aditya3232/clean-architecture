package model

import (
	"time"
)

type UserRole struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	RoleID    int64     `gorm:"not null"`
	UserID    int64     `gorm:"not null"`
	CreatedAt time.Time `gorm:"type:timestamp;default:current_timestamp"`
	UpdatedAt *time.Time
	DeletedAt *time.Time `gorm:"index"`

	// Relasi ke User & Role
	User User `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	Role Role `gorm:"foreignKey:RoleID;references:ID;constraint:OnDelete:CASCADE"`
}

func (UserRole) TableName() string {
	return "user_role"
}
