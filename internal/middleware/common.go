package middleware

import (
	"github.com/SaidMg10/gestor-one/internal/auth"
	"github.com/SaidMg10/gestor-one/internal/domain"
	"github.com/SaidMg10/gestor-one/internal/service"
)

// Middleware mantiene lo necesario para autenticar
type Middleware struct {
	authService   *service.AuthService
	authenticator auth.Authenticator
	userRepo      domain.UserRepo
}

// NewMiddleware inicializa el middleware con lo necesario
func NewMiddleware(authService *service.AuthService, authenticator auth.Authenticator, userRepo domain.UserRepo) *Middleware {
	return &Middleware{
		authService:   authService,
		authenticator: authenticator,
		userRepo:      userRepo,
	}
}
