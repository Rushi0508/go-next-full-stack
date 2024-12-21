package main

import (
	"log"
	"os"

	"crud_app/db"
	"crud_app/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	db.InitDB()

	port := os.Getenv("PORT")

	r := gin.Default()

	routes.SetupRoutes(r)

	log.Printf("Server is running on port %s", port)
	r.Run(":" + port)
}
