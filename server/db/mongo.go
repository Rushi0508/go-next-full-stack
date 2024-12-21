package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Client *mongo.Client
	Database *mongo.Database
}

func InitDB(uri, dbName string) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	log.Println("Connected to MongoDB")

	log.Println("Connected to MongoDB")
    return &MongoDB{
        Client:   client,
        Database: client.Database(dbName),
    }, nil
}

func (m *MongoDB) Close(ctx context.Context) error {
	return m.Client.Disconnect(ctx)
}
