// Package middleware implements all middleware functions.
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const userKey = "user"

func (m *Middleware) AuthTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Logica del middleware (esto ocurrira antes de concluir la petición)
		// 1. Empezamos obteniendo todos los parametros requeridos para las validaciones
		authHeader := c.GetHeader("Authorization")
		// 2. Dividimos el token para sacar el bearer y el tokn
		// Se hace uso de strings.Split y se pasa tanto el header como la separacion " "
		// " " eso es imporante para poder separar "Bearer" del "Token"
		parts := strings.Split(authHeader, " ")
		// Empezamos con las validaciones
		// Las validaciones seran entorno a tener "Bearer" y "Token"
		// El bearer sera validado con el indice [0]
		// El token sera validado con el indice [1] y bajo la funcion validate
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "authorization header is malformed"})
			c.Abort()
			return
		}
		// Obtenemos el token que se aloja en el indice [1]
		tokenString := parts[1]
		// Usamos la funcion Validate para validar el token
		authToken, err := m.authenticator.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}
		// Extraemos el claim del token validado
		claims, ok := authToken.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			c.Abort()
			return
		}
		// Con el token validado debemos obtener el userID del claim sub
		userIDFloat, ok := claims["sub"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id in token"})
			c.Abort()
			return
		}
		userID := uint(userIDFloat)
		// Con el userID, ahora consultamos el usuario en la base de datos
		ctx := c.Request.Context()
		user, err := m.userRepo.GetByID(ctx, userID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			c.Abort()
			return
		}
		// Guardamos el usuario en el contexto para que el handler
		c.Set(userKey, user)
		// y otros middlewares puedan acceder a él.
		// Llamar a c.Next() para que continúe la cadena de middlewares/handlers.
		c.Next()
	}
}
