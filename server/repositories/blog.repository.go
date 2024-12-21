package repositories

import (
	"context"
	"crud_app/models"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BlogRepository struct {
	BaseRepository
}

type IBlogRepository interface {
	CreateBlog(ctx context.Context, blog *models.Blog) error
	GetAllBlogs(ctx context.Context) ([]*models.Blog, error)
	GetBlogById(ctx context.Context, id string) (*models.Blog, error)
	UpdateBlog(ctx context.Context, id string, blog *models.Blog) error
	DeleteBlog(ctx context.Context, id string, authorID string) error
}

func NewBlogRepository(db *mongo.Database) *BlogRepository {
	return &BlogRepository{
		BaseRepository{
				DB:         db,
				Collection: db.Collection("blogs"),
		},
	}
}

func (r *BlogRepository) CreateBlog(ctx context.Context, blog *models.Blog) error {
	_, err := r.Collection.InsertOne(ctx, blog)
	return err
}

func (r *BlogRepository) GetAllBlogs(ctx context.Context) ([]*models.Blog, error) {
	cursor, err := r.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var blogs []*models.Blog
	for cursor.Next(ctx) {
		var blog models.Blog
		if err := cursor.Decode(&blog); err != nil {
			return nil, err
		}
		blogs = append(blogs, &blog)
	}
	return blogs, nil
}

func (r *BlogRepository) GetBlogById(ctx context.Context, id string) (*models.Blog, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var blog models.Blog
	err = r.Collection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&blog)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("blog not found")
		}
		return nil, err
	}
	return &blog, nil
}

func (r *BlogRepository) UpdateBlog(ctx context.Context, id string, blog *models.Blog) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid blog ID format: %v", err)
	}

	// First, verify the blog exists and belongs to the author
	var existingBlog models.Blog
	err = r.Collection.FindOne(ctx, bson.M{
		"_id": objectId,
	}).Decode(&existingBlog)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("blog not found")
		}
		return err
	}

	// Compare the author IDs
	if existingBlog.AuthorID.String() != blog.AuthorID.String() {
		return errors.New("unauthorized: only the author can update this blog")
	}

	// If verification passes, proceed with update
	blog.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	result, err := r.Collection.UpdateOne(
		ctx,
		bson.M{"_id": objectId},
		bson.M{"$set": blog},
	)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return errors.New("blog not found")
	}

	return nil
}

func (r *BlogRepository) DeleteBlog(ctx context.Context, id string, authorID string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid blog ID format: %v", err)
	}

	// First, verify the blog exists and belongs to the author
	var blog models.Blog
	err = r.Collection.FindOne(ctx, bson.M{
		"_id": objectId,
	}).Decode(&blog)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("blog not found")
		}
		return err
	}

	// Compare the author IDs
	if blog.AuthorID.String() != authorID {
		return errors.New("unauthorized: only the author can delete this blog")
	}

	// If verification passes, proceed with deletion
	result, err := r.Collection.DeleteOne(ctx, bson.M{"_id": objectId})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("blog not found")
	}

	return nil
}
