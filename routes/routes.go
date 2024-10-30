package routes

import (
	"schooglink-task/controllers"
	"schooglink-task/utils/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterBlogRoutes(r *gin.Engine) {
	blogGroup := r.Group("/blogs")
	{
		blogGroup.POST("/", middleware.AuthMiddleware(), controllers.CreateBlog)
		blogGroup.GET("/", controllers.GetBlogs)
		blogGroup.GET("/:id", controllers.GetBlogsById)
		blogGroup.PUT("/:id", middleware.AuthMiddleware(), middleware.BlogOwnership(), controllers.UpdateBlog)
		blogGroup.DELETE("/:id", middleware.AuthMiddleware(), middleware.BlogOwnership(), controllers.DeleteBlog)
	}

	userGroup := r.Group("/users")
	{
		userGroup.POST("/register", controllers.CreateUser)
		userGroup.GET("/", controllers.GetAllUsers)
		userGroup.POST("/login", controllers.LoginUser)
	}
}
