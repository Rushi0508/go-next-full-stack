package services

import (
	"context"
	"crud_app/models"
	"crud_app/repositories"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlogService struct {
	blogRepo repositories.IBlogRepository
}

type IBlogService interface {
	CreateBlog(ctx context.Context, blog *models.Blog) error
	GetAllBlogs(ctx context.Context) ([]*models.Blog, error)
	GetBlogById(ctx context.Context, id string) (*models.Blog, error)
	UpdateBlog(ctx context.Context, id string, blog *models.Blog) error
	DeleteBlog(ctx context.Context, id string, authorID primitive.ObjectID) error
}

func NewBlogService(blogRepo repositories.IBlogRepository) *BlogService {
	return &BlogService{
		blogRepo: blogRepo,
	}
}

func (s *BlogService) CreateBlog(ctx context.Context, blog *models.Blog) error {
	return s.blogRepo.CreateBlog(ctx, blog)
}

func (s *BlogService) GetAllBlogs(ctx context.Context) ([]*models.Blog, error) {
	return s.blogRepo.GetAllBlogs(ctx)
}

func (s *BlogService) GetBlogById(ctx context.Context, id string) (*models.Blog, error) {
	return s.blogRepo.GetBlogById(ctx, id)
}

func (s *BlogService) UpdateBlog(ctx context.Context, id string, blog *models.Blog) error {
	// No need for manual authorization check here anymore as it's handled in the repository
	return s.blogRepo.UpdateBlog(ctx, id, blog)
}

func (s *BlogService) DeleteBlog(ctx context.Context, id string, authorID primitive.ObjectID) error {
	// Get the existing blog first
	fmt.Println(id)
	fmt.Println(authorID)
	existingBlog, err := s.blogRepo.GetBlogById(ctx, id)
	fmt.Println(existingBlog.AuthorID.String())
	if err != nil {
		return err
	}

	// Check if the user is the author
	if existingBlog.AuthorID.String() != authorID.String() {
		return errors.New("unauthorized: only the author can delete this blog")
	}

	return s.blogRepo.DeleteBlog(ctx, id, authorID.String())
}
