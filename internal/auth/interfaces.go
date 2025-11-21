package auth

import "github.com/golang-jwt/jwt/v5"

// Entendimiento de las interfaces
// La interfaz Authenticator permite desacoplar la lógica JWT,
// de modo que se pueda cambiar la implementación (por ejemplo, en tests o para usar otro sistema de auth),
// sin afectar al resto de la aplicación.

// Comentario 1:
// Estas funciones no son parte del paquete jwt original.
// Se crean como parte de nuestra implementación para abstraer
// la lógica de generación y validación de tokens en una interfaz reutilizable.

// Comentario 2:
// La idea de separar GenerateToken y ValidateToken en una interfaz
// es evitar repetir lógica en cada parte del sistema que necesita JWT.
// Así podemos centralizar la firma y verificación en una sola capa
// (ej. el middleware llama a ValidateToken, el login a GenerateToken),
// manteniendo el código limpio, reutilizable y fácil de testear.

// Authenticator define el contrato para cualquier componente capaz de
// generar y validar tokens JWT. Esta interfaz permite abstraer la lógica
// del JWT del resto de la aplicación, facilitando pruebas y cambios de implementación.
type Authenticator interface {
	// GenerateToken recibe un conjunto de claims (datos del usuario o contexto)
	// y devuelve un token JWT firmado como string. Usado típicamente al iniciar sesión.
	GenerateToken(claims jwt.Claims) (string, error)

	// ValidateToken recibe un token JWT como string y lo valida (firma, expiración, issuer, etc.).
	// Devuelve el token parseado si es válido o un error en caso contrario.
	// Usado en middleware u otros puntos de verificación.
	ValidateToken(token string) (*jwt.Token, error)
}
