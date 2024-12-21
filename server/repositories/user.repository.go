package repositories

import (
	"context"
	"crud_app/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	BaseRepository
}

type IUserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	FindUserByEmail(ctx context.Context, email string) (*models.User, error)
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		BaseRepository{
			DB:         db,
			Collection: db.Collection("users"),
		},
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	_, err := r.Collection.InsertOne(ctx, user)
	return err
}

func (r *UserRepository) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.Collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
} 