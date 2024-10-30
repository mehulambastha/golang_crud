package controllers

import (
	"net/http"
	"schooglink-task/database"
	"schooglink-task/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func CreateBlog(c *gin.Context) {
	var blog models.Blog

	if err := c.ShouldBindJSON(&blog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validate.Struct(blog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	blog.CreatedAt = time.Now()
	blog.UpdatedAt = time.Now()
	userid, _ := c.Get("userID")
	blog.AuthorID = uint(userid.(float64))

	if err := database.DB.Create(&blog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create blog."})
		return
	}
	// Retrieve the created blog with the associated author details
	var blogWithAuthor models.Blog
	if err := database.DB.Preload("Author").First(&blogWithAuthor, blog.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve blog with author"})
		return
	}

	// Respond with the created blog including author details
	c.JSON(http.StatusOK, blogWithAuthor)
}

func GetBlogs(c *gin.Context) {
	var blogs []models.Blog

	if err := database.DB.Find(&blogs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve all blogs."})
		return
	}

	c.JSON(http.StatusOK, blogs)
}

func GetBlogsById(c *gin.Context) {
	id := c.Param("id")

	var blog models.Blog

	if err := database.DB.First(&blog, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found."})
		return
	}

	var blogWithAuthor models.Blog

	if err := database.DB.Preload("Author").First(&blogWithAuthor, blog.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not load Author details."})
		return
	}

	c.JSON(http.StatusOK, blogWithAuthor)
}

func UpdateBlog(c *gin.Context) {
	id := c.Param("id")
	var blog models.Blog
	if err := database.DB.First(&blog, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found."})
		return
	}

	if err := c.ShouldBindJSON(&blog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	blog.UpdatedAt = time.Now()

	if err := database.DB.Save(&blog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Blog"})
		return
	}

	c.JSON(http.StatusOK, blog)
}

func DeleteBlog(c *gin.Context) {
	id := c.Param("id")

	if err := database.DB.Delete(&models.Blog{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete blog."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Blog Deleted."})
}
