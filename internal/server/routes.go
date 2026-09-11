package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func registerRoutes(r *gin.Engine) {
	r.GET("hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"hello": "world",
		})
	})

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{
			"code":    "not_found",
			"message": "no such endpoint",
		}})
	})
}
