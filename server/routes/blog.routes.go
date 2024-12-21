package routes

import (
	"crud_app/controllers"
	"crud_app/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupBlogRoutes(rg *gin.RouterGroup, blogController *controllers.BlogController) {
	blog := rg.Group("/blog")

	// Public routes	
	{
		blog.GET("/", blogController.GetAllBlogs)
		blog.GET("/:id", blogController.GetBlogById)
	}

	// Private routes
	blog.Use(middlewares.AuthMiddleware())
	{
		blog.POST("/", blogController.CreateBlog)
		blog.PUT("/:id", blogController.UpdateBlog)
		blog.DELETE("/:id", blogController.DeleteBlog)
	}
}
