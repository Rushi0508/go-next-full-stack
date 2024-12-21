package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Blog struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title     string             `bson:"title" json:"title"`
	Content   string             `bson:"content" json:"content"`
	AuthorID  primitive.ObjectID          `bson:"author_id" json:"author_id"`
	CreatedAt primitive.DateTime `bson:"created_at" json:"created_at" default:"now()"`
	UpdatedAt primitive.DateTime `bson:"updated_at" json:"updated_at" default:"now()"`
}

type CreateBlogRequest struct {
	Title   string `json:"title" binding:"required,min=3,max=100"`
	Content string `json:"content" binding:"required,min=10"`
}

