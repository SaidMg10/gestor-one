// Package domain contains the core domain models for the application.
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

// Income represents an income record in the system.
type Income struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Amount      float64    `gorm:"not null" json:"amount"`
	Description string     `gorm:"size:255" json:"description"`
	Date        time.Time  `gorm:"not null" json:"date"`
	Type        IncomeType `gorm:"size:50;not null" json:"type"`
	CreatedBy   uint       `gorm:"not null" json:"created_by"`
	Receipt     Receipt    `gorm:"constraint:OnDelete:CASCADE;foreignKey:IncomeID" json:"receipt"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

// Expense represents an expense record in the system.
type Expense struct {
	ID          uint        `gorm:"primaryKey" json:"id"`
	Amount      float64     `gorm:"not null" json:"amount"`
	Description string      `gorm:"size:255" json:"description"`
	Date        time.Time   `gorm:"not null" json:"date"`
	Type        ExpenseType `gorm:"size:50;not null" json:"type"`
	CreatedBy   uint        `gorm:"not null" json:"created_by"`
	Receipt     Receipt     `gorm:"constraint:OnDelete:CASCADE;foreignKey:ExpenseID" json:"receipt"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	DeletedAt   *time.Time  `gorm:"index" json:"deleted_at,omitempty"`
}

// Receipt represents a receipt associated with an income.
type Receipt struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	IncomeID   *uint     `gorm:"index"`
	ExpenseID  *uint     `gorm:"index"`
	FileName   string    `gorm:"size:255;not null" json:"file_name"`
	RelPath    string    `gorm:"size:255;not null" json:"relPath_url"`
	MimeType   string    `gorm:"size:50;not null" json:"mime_type"`
	UploadedBy uint      `gorm:"not null" json:"uploaded_by"`
	Checksum   string    `gorm:"size:255" json:"checksum,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
