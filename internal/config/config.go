// Package config provides configuration settings for the application
package config

import (
	"time"

	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// -----------------------
// Estructuras de Config |
// ----------------------

type Config struct {
	App      AppConfig          `mapstructure:"app"`
	Server   ServerConfig       `mapstructure:"server"`
	Database DBConfig           `mapstructure:"database"`
	JWT      JWTConfig          `mapstructure:"jwt"`
	Google   GoogleOAuth2Config `mapstructure:"google_oauth2"`
}

// AppConfig es la Configuración general de la aplicación
type AppConfig struct {
	Name        string `mapstructure:"name"`          // Nombre de la app
	Env         string `mapstructure:"env"`           // dev, prod, test
	Version     string `mapstructure:"version"`       // Versión de la app
	Debug       bool   `mapstructure:"debug"`         // true/false
	FEOriginURL string `mapstructure:"fe_origin_url"` // URL del frontend para CORS
}

// ServerConfig es la Configuración del servidor HTTP
type ServerConfig struct {
	Addr         string        `mapstructure:"addr"`          // ej: 0.0.0.0:8080
	Host         string        `mapstructure:"host"`          // ej: 0.0.0.0
	Port         int           `mapstructure:"port"`          // ej: 8080
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`  // ej: 5s
	WriteTimeout time.Duration `mapstructure:"write_timeout"` // ej: 10s
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`  // ej: 120s
}

// DBConfig es la Configuración de la base de datos
type DBConfig struct {
	DSN          string        `mapstructure:"dsn"`      // DSN directo (prod)
	Host         string        `mapstructure:"host"`     // DB host
	Port         int           `mapstructure:"port"`     // DB port
	User         string        `mapstructure:"user"`     // DB usuario
	Password     string        `mapstructure:"password"` // DB contraseña
	Name         string        `mapstructure:"name"`     // DB nombre
	SSLMode      string        `mapstructure:"sslmode"`  // disable/require/verify-full
	MaxOpenConns int           `mapstructure:"MAX_OPEN_CONNS"`
	MaxIdleConns int           `mapstructure:"MAX_IDLE_CONNS"`
	MaxIdleTime  time.Duration `mapstructure:"MAX_IDLE_TIME"`
}

// JWTConfig es la Configuración de JWTConfig
type JWTConfig struct {
	Secret           string        `mapstructure:"secret"`            // clave secreta
	AccessTokenTTL   time.Duration `mapstructure:"access_token_ttl"`  // ej: 15m
	RefreshTokenTTL  time.Duration `mapstructure:"refresh_token_ttl"` // ej: 7d
	Issuer           string        `mapstructure:"issuer"`            // emisor
	SigningAlgorithm string        `mapstructure:"signing_algorithm"` // ej: HS256
}

// GoogleOAuth2Config es la Configuración de Google OAuth2
type GoogleOAuth2Config struct {
	ClientID     string   `mapstructure:"client_id"`
	ClientSecret string   `mapstructure:"client_secret"`
	RedirectURL  string   `mapstructure:"redirect_url"`
	Scopes       []string `mapstructure:"scopes"`
}

// -----------------------
// Funcion LoadConfig    |
// ----------------------

func LoadConfig(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigName(".env") // sin extensión
	v.SetConfigType("yaml")
	v.AddConfigPath(path)

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) GetGoogleOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     c.Google.ClientID,
		ClientSecret: c.Google.ClientSecret,
		RedirectURL:  c.Google.RedirectURL,
		Scopes:       c.Google.Scopes,
		Endpoint:     google.Endpoint,
	}
}

var Cfg *Config

// Init inicializa la configuración y la deja accesible globalmente
func Init(path string) error {
	cfg, err := LoadConfig(path)
	if err != nil {
		return err
	}
	Cfg = cfg
	return nil
}
