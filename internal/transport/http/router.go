// Package http
package http

import "github.com/gin-gonic/gin"

func NewRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return router
}
