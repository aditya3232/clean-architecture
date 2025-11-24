package model

import (
	"time"
)

type Role struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"type:varchar(255);unique;not null"`
	CreatedAt time.Time `gorm:"type:timestamp;default:current_timestamp"`
	UpdatedAt *time.Time
	DeletedAt *time.Time `gorm:"index"`

	// Relasi many-to-many ke User lewat tabel pivot user_role
	// Walaupun di tabel roles tidak ada kolom user_id,
	// GORM otomatis menggunakan tabel user_role (join table)
	// untuk menghubungkan role_id <-> user_id.
	// Jadi field Users ini berfungsi agar kita bisa langsung ambil
	// semua user yang memiliki role tertentu, tanpa query join manual.
	Users []User `gorm:"many2many:user_role"`
}

func (Role) TableName() string {
	return "roles"
}
