package domain

import "context"

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
