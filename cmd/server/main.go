package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SaidMg10/gestor-one/internal/auth"
	"github.com/SaidMg10/gestor-one/internal/config"
	"github.com/SaidMg10/gestor-one/internal/db"
	"github.com/SaidMg10/gestor-one/internal/repository"
	"github.com/SaidMg10/gestor-one/internal/service"
	httpTransport "github.com/SaidMg10/gestor-one/internal/transport/http"
)

// Main initializes the application and starts the server.
func main() {
	// Cargar configuración desde ./config/config.yml
	// Inicializar configuración
	if err := config.Init("."); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	cfg := config.Cfg

	if err := db.Init(cfg.Database); err != nil {
		log.Fatalf("X error initializing database: %v", err)
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	auth := auth.NewJWTAuthenticatorFromConfig(cfg.JWT)

	// 3. Inicializar repositorios y servicios
	userRepo := repository.NewGormUserRepo(db.DB)
	userSvc := service.NewUserService(userRepo)
	authSvc := service.NewAuthService(
		userRepo, // repositorio de usuarios
		auth,     // Authenticator
		cfg.JWT.AccessTokenTTL,
		cfg.JWT.RefreshTokenTTL,
		cfg.JWT.Issuer,
	)

	r := httpTransport.NewRouter(userSvc, authSvc)

	// Mostrar que la config se cargó correctamente
	fmt.Println("=================================")
	fmt.Printf("App: %s v%s [%s]\n", cfg.App.Name, cfg.App.Version, cfg.App.Env)
	fmt.Printf("Server running on %s:%d\n", cfg.Server.Host, cfg.Server.Port)
	fmt.Printf("Debug mode: %v\n", cfg.App.Debug)
	fmt.Printf("Google Client ID: %s\n", cfg.Google.ClientID)
	fmt.Printf("Google Client Secret: %s\n", cfg.Google.ClientSecret)
	fmt.Printf("Google Redirect URL: %s\n", cfg.Google.RedirectURL)
	fmt.Println("=================================")

	s := &http.Server{
		Addr:    cfg.Server.Addr,
		Handler: r,
	}
	// Canal para graceful shutdown
	serverErr := make(chan error, 1)
	go func() {
		log.Printf("🚀 Servidor iniciado en %s (Modo: %s)", s.Addr, cfg.App.Env)
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	// Esperar señales de interrupción
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		log.Fatalf("❌ Error del servidor: %v", err)
	case <-quit:
		log.Println("🛑 Recibida señal de apagado, cerrando servidor...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := s.Shutdown(ctx); err != nil {
			log.Fatalf("❌ Error al apagar servidor: %v", err)
		}
		log.Println("✅ Servidor apagado correctamente")
	}
}
