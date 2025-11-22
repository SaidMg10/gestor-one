package middleware

import (
	"net/http"
	"slices"

	"github.com/SaidMg10/gestor-one/internal/domain"
	"github.com/gin-gonic/gin"
)

func (m *Middleware) CheckRole(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtenemos el usuario guardado en el context
		userCtx, ok := c.Get(userKey)
		// Error que diga que no hay un usuario en el context
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user not found in the context"})
			c.Abort()
			return
		}
		// verificamos si la interfaz obtenida del context es del tipo User
		user, ok := userCtx.(*domain.User)
		// En caso de no ser del tipo user mandar un error diciendo que no es del tipo
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "the type of User not exist"})
			c.Abort()
			return
		}
		// Validar si el rol del usuario está en los roles requeridos
		if !slices.Contains(requiredRoles, user.Role) {
			// Respuesta en caso de no cumplir con el rol
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			c.Abort()
			return
		}
		c.Next()
	}
}
