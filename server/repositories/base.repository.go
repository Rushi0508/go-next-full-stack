package repositories

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type BaseRepository struct {
    DB         *mongo.Database
    Collection *mongo.Collection
} 