package controllers

import (
	"crud_app/models"
	"crud_app/services"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlogController struct {
	blogService *services.BlogService
}

func NewBlogController(blogService *services.BlogService) *BlogController {
	return &BlogController{
		blogService: blogService,
	}
}

func (c *BlogController) CreateBlog(ctx *gin.Context) {
	var req models.CreateBlogRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Validation failed",
			"details": err.Error(),
		})
		return
	}

	// Get author ID from the authenticated user context
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Create new blog with validated data
	blog := &models.Blog{
		ID: primitive.NewObjectID(),
		Title:     req.Title,
		Content:   req.Content,
		AuthorID:  userID.(primitive.ObjectID),
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	err := c.blogService.CreateBlog(ctx, blog)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, blog)
}

func (c *BlogController) GetAllBlogs(ctx *gin.Context) {
	blogs, err := c.blogService.GetAllBlogs(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, blogs)
}

func (c *BlogController) GetBlogById(ctx *gin.Context) {
	id := ctx.Param("id")
	blog, err := c.blogService.GetBlogById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, blog)
}

func (c *BlogController) UpdateBlog(ctx *gin.Context) {
	id := ctx.Param("id")
	var req models.CreateBlogRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get author ID from the authenticated user context
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Create updated blog with validated data
	blog := &models.Blog{
		Title:     req.Title,
		Content:   req.Content,
		AuthorID:  userID.(primitive.ObjectID),
	}

	err := c.blogService.UpdateBlog(ctx, id, blog)
	if err != nil {
		if strings.Contains(err.Error(), "unauthorized") {
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Blog updated successfully"})
}

func (c *BlogController) DeleteBlog(ctx *gin.Context) {
	id := ctx.Param("id")
	
	// Get author ID from the authenticated user context
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Add debug logging
	log.Printf("Attempting to delete blog ID: %s by user ID: %v", id, userID)

	err := c.blogService.DeleteBlog(ctx, id, userID.(primitive.ObjectID))
	if err != nil {
		log.Printf("Error deleting blog: %v", err) // Add error logging
		if strings.Contains(err.Error(), "unauthorized") {
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Blog deleted successfully"})
}
