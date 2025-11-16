package domain

import (
	"time"
)

// User represents a user entity in the system.
type User struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	Name      string     `gorm:"size:100;not null" json:"name"`
	LastName  string     `gorm:"size:100;not null" json:"last_name"`
	Email     string     `gorm:"size:150;not null;uniqueIndex" json:"email"`
	Phone     string     `gorm:"size:20" json:"phone"`
	Password  Password   `gorm:"size:255" json:"-"`
	Role      string     `gorm:"size:30;not null" json:"role"`
	GoogleID  string     `gorm:"size:150" json:"google_id,omitempty"`
	Active    *bool      `gorm:"default:true" json:"active"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}
