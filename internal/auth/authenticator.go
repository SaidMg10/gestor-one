// Package auth provides a functionality to generate and validate JWT tokens.
package auth

import (
	"fmt"

	"github.com/SaidMg10/gestor-one/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

// JWTAuthenticator contiene la configuración necesaria para crear y validar tokens JWT.
// Incluye la clave secreta para firma, la audiencia esperada (aud) y el emisor esperado (iss).
type JWTAuthenticator struct {
	Secret string // Clave secreta para firmar y verificar los tokens JWT.
	Aud    string // Audiencia esperada en el token (claim "aud").
	Iss    string // Emisor esperado en el token (claim "iss").
}

// NewJWTAuthenticator construye una nueva instancia de JWTAuthenticator.
// Recibe explícitamente la clave secreta, audiencia e issuer para configurar la autenticación JWT.
// Esta función es útil cuando se tienen estos valores separados y se quieren pasar directamente.
func NewJWTAuthenticator(secret, aud, iss string) *JWTAuthenticator {
	return &JWTAuthenticator{
		Secret: secret,
		Aud:    aud,
		Iss:    iss,
	}
}

// NewJWTAuthenticatorFromConfig crea un JWTAuthenticator usando la configuración TokenConfig.
// En esta implementación se asume que Iss (issuer) y Aud (audience) son iguales,
// por eso ambos campos se asignan con authCfg.Iss.
// Si en el futuro quieres diferenciarlos, considera agregar Aud explícitamente en TokenConfig.
func NewJWTAuthenticatorFromConfig(authCfg config.JWTConfig) *JWTAuthenticator {
	return NewJWTAuthenticator(
		authCfg.Secret,
		authCfg.Issuer, // audiencia igual al issuer por defecto
		authCfg.Issuer,
	)
}

// GenerateToken genera un token JWT firmado usando HS256.
// Recibe los claims deseados y los firma con la clave secreta definida en el JWTAuthenticator.
// Devuelve el token como string o un error si la firma falla.
func (a *JWTAuthenticator) GenerateToken(claims jwt.Claims) (string, error) {
	// Crear token con el método de firma y los claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Firmar el token usando la clave secreta
	signedString, err := token.SignedString([]byte(a.Secret)) // Importante: se usa []byte para la firma HMAC
	if err != nil {
		return "", err
	}

	// Devolver el token firmado
	return signedString, nil
}

// ValidateToken parsea y valida un token JWT string.
// Verifica que:
// - El método de firma sea HMAC (HS256) para evitar ataques con 'alg' maliciosos.
// - El token tenga los claims requeridos como 'exp', 'aud' y 'iss' con los valores esperados.
// - La firma sea válida con la clave secreta.
// Devuelve el token validado o un error.
func (a *JWTAuthenticator) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (any, error) {
		// Validamos que el método usado sea de tipo HMAC (ej. HS256)
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: % v", t.Header["alg"])
		}
		// Devolvemos la clave con la que se firmó para que jwt la use en la verificación de firma
		return []byte(a.Secret), nil
	},
		// Verificamos los claims esperados para que el token no sea aceptado si no cumple las reglas
		jwt.WithExpirationRequired(),                                // debe tener 'exp' y no estar vencido
		jwt.WithAudience(a.Aud),                                     // debe tener 'aud' correcto
		jwt.WithIssuer(a.Iss),                                       // debe tener 'iss' correcto
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}), // aceptamos solo HS256
	)
}
