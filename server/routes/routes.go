package routes

import (
	"crud_app/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.GET("/health", controllers.HealthCheck())
}

