package domain

import (
	"context"
	"mime/multipart"
)

// UserRepo defines an interface with methods for managing User entities.
type UserRepo interface {
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id uint) (*User, error)
	List(ctx context.Context) ([]User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uint) error
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

// IncomeRepo defines an interface with methods for managing Income entities.
type IncomeRepo interface {
	GetByID(ctx context.Context, id uint) (*Income, error)
	List(ctx context.Context) ([]Income, error)
	CreateWithReceipt(ctx context.Context, income *Income, receipt *Receipt) error
	UpdateWithReceipt(ctx context.Context, income *Income, receipt *Receipt) error
	Delete(ctx context.Context, id uint) error
	SoftDelete(ctx context.Context, id uint) error
	Restore(ctx context.Context, id uint) error
}

// ExpenseRepo defines an interface with methods for managing Expense entities.
type ExpenseRepo interface {
	GetByID(ctx context.Context, id uint) (*Expense, error)
	List(ctx context.Context) ([]Expense, error)
	CreateWithReceipt(ctx context.Context, expense *Expense, receipt *Receipt) error
	UpdateWithReceipt(ctx context.Context, expense *Expense, receipt *Receipt) error
	Delete(ctx context.Context, id uint) error
	SoftDelete(ctx context.Context, id uint) error
	Restore(ctx context.Context, id uint) error
}

// ReceiptRepo defines an interface with methods for managing Receipt entities.
type ReceiptRepo interface {
	GetByID(ctx context.Context, id uint) (*Receipt, error)
	List(ctx context.Context) ([]Receipt, error)
	Create(ctx context.Context, receipt *Receipt) error
	Update(ctx context.Context, receipt *Receipt) error
	Delete(ctx context.Context, id uint) error
}

type FileStorage interface {
	SavePDF(fileHeader *multipart.FileHeader) (string, string, string, error)
	DeletePDF(filePath string) error
}
