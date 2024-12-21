package routes

import (
	"crud_app/controllers"
	"crud_app/repositories"
	"crud_app/services"

	"fmt"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRoutes(router *gin.Engine, db *mongo.Database) {
	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo)

	// Initialize controllers
	authController := controllers.NewAuthController(authService)

	router.GET("/health", controllers.HealthCheck())

	// API v1 group
	v1 := router.Group("/api/v1")

	// Setup route groups
	SetupAuthRoutes(v1, authController)

	// Add this after all routes are set up
	routes := router.Routes()
	for _, route := range routes {
		fmt.Printf("Method: %v, Path: %v\n", route.Method, route.Path)
	}
}

