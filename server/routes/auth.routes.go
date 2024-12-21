package routes

import (
	"crud_app/controllers"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(rg *gin.RouterGroup, authController *controllers.AuthController) {
	auth := rg.Group("/auth")
	{
		auth.POST("/register", authController.RegisterUser)
		auth.POST("/login", authController.LoginUser)
	}
}
