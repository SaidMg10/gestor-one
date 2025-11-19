// Package service implements the business logic of the application.
package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/SaidMg10/gestor-one/internal/domain"
	"github.com/SaidMg10/gestor-one/internal/validator"
)

type UserService struct {
	userRepo domain.UserRepo
}

func NewUserService(u domain.UserRepo) *UserService {
	return &UserService{
		userRepo: u,
	}
}

func (s *UserService) Create(ctx context.Context, user *domain.User, pwd string) error {
	if err := validator.ValidatePassword(pwd); err != nil {
		return err
	}

	if !validator.IsValidEmail(user.Email) {
		return domain.ErrInvalidEmail
	}

	exists, err := s.userRepo.ExistsByEmail(ctx, user.Email)
	if err != nil {
		return err
	}
	if exists {
		return domain.ErrEmailExists
	}

	if user.Role == domain.RoleSuperAdmin {
		return errors.New("no se puede crear un superadmin")
	}
	if user.Role == "" {
		user.Role = domain.RoleEmployee
	}

	if user.Active == nil {
		active := true
		user.Active = &active
	}

	if err := user.Password.Set(pwd); err != nil {
		return fmt.Errorf("error al hashear la contraseña: %w", err)
	}
	return s.userRepo.Create(ctx, user)
}

func (s *UserService) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *UserService) List(ctx context.Context) ([]domain.User, error) {
	return s.userRepo.List(ctx)
}

func (s *UserService) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return s.userRepo.GetByEmail(ctx, email)
}

func (s *UserService) Update(ctx context.Context, id uint, updates *domain.User, pwd *string) error {
	existing, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return domain.ErrNotFound
	}

	if updates.Email != "" && updates.Email != existing.Email {
		if !validator.IsValidEmail(updates.Email) {
			return domain.ErrInvalidEmail
		}

		exists, err := s.userRepo.ExistsByEmail(ctx, updates.Email)
		if err != nil {
			return err
		}
		if exists {
			return domain.ErrEmailExists
		}

		existing.Email = updates.Email
	}

	if pwd != nil && *pwd != "" {
		if err := validator.ValidatePassword(*pwd); err != nil {
			return err
		}
		if err := existing.Password.Set(*pwd); err != nil {
			return fmt.Errorf("error al hashear la contraseña: %w", err)
		}
	}

	if updates.Name != "" {
		existing.Name = updates.Name
	}
	if updates.LastName != "" {
		existing.LastName = updates.LastName
	}
	if updates.Phone != "" {
		existing.Phone = updates.Phone
	}
	if updates.Role != "" {
		if updates.Role == domain.RoleSuperAdmin {
			return errors.New("no se puede asignar rol superadmin")
		}
		existing.Role = updates.Role
	}
	if updates.Active != nil {
		existing.Active = updates.Active
	}

	// 6️⃣ Guardar
	return s.userRepo.Update(ctx, existing)
}

func (s *UserService) Delete(ctx context.Context, id uint) error {
	return s.userRepo.Delete(ctx, id)
}
