// Package http
package http

import (
	"github.com/SaidMg10/gestor-one/internal/service"
	"github.com/gin-gonic/gin"
)

// NewRouter creates a new Gin router.
func NewRouter(userSvc *service.UserService) *gin.Engine {
	r := gin.Default()

	// Middleware de CORS básico
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "gestor-one-api",
		})
	})

	// API V1
	v1 := r.Group("/api/v1")
	{
		// Appointments routes
		users := v1.Group("/users")
		{
			userHandler := NewUserHandler(userSvc)
			users.POST("", userHandler.Create)
			users.GET("", userHandler.List)
			users.GET("/:id", userHandler.GetByID)
			users.PATCH("/:id", userHandler.Update)
			users.DELETE("/:id", userHandler.Delete)
		}

		// Products routes
		/*
			products := v1.Group("/products")
			{
				prodHandler := NewProductHandler(prodSvc)
				products.POST("", prodHandler.Create)
				products.GET("", prodHandler.List)
				products.GET(":id", prodHandler.Get)
				products.PUT(":id", prodHandler.Update)
				products.DELETE(":id", prodHandler.Delete)
			}
		*/
	}
	return r
}
