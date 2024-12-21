package main

import (
	"context"
	config "crud_app/configs"
	"crud_app/db"
	"crud_app/routes"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	mongoDB, err := db.InitDB(cfg.MongoURI, cfg.DBName)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	router := gin.Default()
	routes.SetupRoutes(router, mongoDB.Database)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := mongoDB.Close(ctx); err != nil {
		log.Printf("Failed to close MongoDB connection: %v", err)
	}

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
