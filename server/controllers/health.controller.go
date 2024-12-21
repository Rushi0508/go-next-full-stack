package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Server is running fine", "status": "healthy"})
	}
}
