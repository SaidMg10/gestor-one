package service

import (
	"context"
	"fmt"
	"time"

	"github.com/SaidMg10/gestor-one/internal/auth"
	"github.com/SaidMg10/gestor-one/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	UserRepo      domain.UserRepo
	Authenticator auth.Authenticator
	AccessTTL     time.Duration
	RefreshTTL    time.Duration
	Issuer        string
}

func NewAuthService(
	userRepo domain.UserRepo,
	authenticator auth.Authenticator,
	accessTTL time.Duration,
	refreshTTL time.Duration,
	issuer string,
) *AuthService {
	return &AuthService{
		UserRepo:      userRepo,
		Authenticator: authenticator,
		AccessTTL:     accessTTL,
		RefreshTTL:    refreshTTL,
		Issuer:        issuer,
	}
}

// Login is the method to authenticate a user and generate tokens
func (s *AuthService) Login(ctx context.Context, email, password string) (string, string, error) {
	// Buscar usuario
	user, err := s.UserRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", "", err
	}

	// Comparar contraseña
	if err := user.Password.Compare(password); err != nil {
		return "", "", err
	}

	// === ACCESS TOKEN ===
	accessClaims := jwt.MapClaims{
		"sub": user.ID,
		"rol": user.Role,
		"exp": time.Now().Add(s.AccessTTL).Unix(),
		"iat": time.Now().Unix(),
		"iss": s.Issuer,
		"aud": s.Issuer,
	}

	accessToken, err := s.Authenticator.GenerateToken(accessClaims)
	if err != nil {
		return "", "", err
	}

	// === REFRESH TOKEN ===
	refreshClaims := jwt.MapClaims{
		"sub": user.ID,
		"typ": "refresh",
		"exp": time.Now().Add(s.RefreshTTL).Unix(),
		"iat": time.Now().Unix(),
		"iss": s.Issuer,
		"aud": s.Issuer,
	}

	refreshToken, err := s.Authenticator.GenerateToken(refreshClaims)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// Refresh is the method to refresh an access token using a valid refresh token
func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (string, error) {
	// Validar refresh
	token, err := s.Authenticator.ValidateToken(refreshToken)
	if err != nil {
		return "", err
	}

	claims := token.Claims.(jwt.MapClaims)

	if claims["typ"] != "refresh" {
		return "", fmt.Errorf("invalid token type")
	}

	userID := uint(claims["sub"].(float64))

	user, err := s.UserRepo.GetByID(ctx, userID)
	if err != nil {
		return "", err
	}

	// Generar nuevo access
	newAccessClaims := jwt.MapClaims{
		"sub": user.ID,
		"rol": user.Role,
		"exp": time.Now().Add(s.AccessTTL).Unix(),
		"iat": time.Now().Unix(),
		"iss": s.Issuer,
		"aud": s.Issuer,
	}

	accessToken, err := s.Authenticator.GenerateToken(newAccessClaims)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
