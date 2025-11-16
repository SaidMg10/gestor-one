package service

import (
	"context"
	"fmt"

	"github.com/SaidMg10/gestor-one/internal/domain"
)

type UserService struct {
	userRepo domain.UserRepo
}

func NewUserService(u domain.UserRepo) *UserService {
	return &UserService{
		userRepo: u,
	}
}

func (s *UserService) Create(ctx context.Context, user *domain.User) error {
	// Validaciones de logica en el create
	// Regla de negocio:
	// 1. No pueden crear un superadmin
	if user.Role == domain.RoleSuperAdmin {
		return fmt.Errorf("Cannot create a superadmin")
	}
	// 2. Validar que el email no exista en la base de datos
	exists, err := s.userRepo.ExistsByEmail(ctx, user.Email)
	if err != nil {
		return err
	}
	if exists {
		return domain.ErrEmailExists
	}
	if !domain.IsValidRole(user.Role) {
		return fmt.Errorf("invalid role")
	}
	// 3. If role is ""
	if user.Role == "" {
		user.Role = domain.RoleEmployee
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

func (s *UserService) Update(ctx context.Context, id uint, input *domain.User) error {
	existing, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// 1. No tocar superadmin
	if existing.Role == domain.RoleSuperAdmin {
		return fmt.Errorf("Cannot update a superadmin")
	}

	if !domain.IsValidRole(input.Role) {
		return fmt.Errorf("invalid role")
	}

	// 2. Email único (si se envía)
	if input.Email != "" && input.Email != existing.Email {
		exists, err := s.userRepo.ExistsByEmail(ctx, input.Email)
		if err != nil {
			return err
		}
		if exists {
			return domain.ErrEmailExists
		}
		existing.Email = input.Email
	}

	// 3. Validar role (si se envía)
	if input.Role != "" {
		if input.Role == domain.RoleSuperAdmin {
			return fmt.Errorf("Cannot assign superadmin role")
		}
		existing.Role = input.Role
	}

	// 4. Actualizar campos simples
	if input.Name != "" {
		existing.Name = input.Name
	}
	if input.LastName != "" {
		existing.LastName = input.LastName
	}
	if input.Phone != "" {
		existing.Phone = input.Phone
	}

	// 5. Nueva contraseña (si viene)
	if input.Password.HasHash() {
		existing.Password = input.Password
	}

	if input.Active != nil {
		existing.Active = input.Active
	}

	return s.userRepo.Update(ctx, existing)
}

func (s *UserService) Delete(ctx context.Context, id uint) error {
	return s.userRepo.Delete(ctx, id)
}
