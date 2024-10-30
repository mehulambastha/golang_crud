package controllers

import (
	"log"
	"net/http"
	"schooglink-task/database"
	"schooglink-task/models"
	"schooglink-task/utils"
	"schooglink-task/utils/middleware"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateUser godoc
// @Summary Create a new user
// @Description Register a new user with a name, email, and password
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "User data"
// @Success 200 {object} models.User
// @Failure 400 {object} gin.H{"error": "error message"}
// @Failure 500 {object} gin.H{"error": "Internal server error"}
// @Router /users [post]
func CreateUser(c *gin.Context) {
	var user models.User

	// Bind JSON input to the User struct
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validate.Struct(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the user's password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Set the hashed password to the user model
	user.Password = hashedPassword
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Save the user to the database
	if err := database.DB.Create(&user).Error; err != nil {
		log.Printf("Failed to create user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Do not return the password in the response
	user.Password = ""

	// Respond with the created user
	c.JSON(http.StatusOK, user)
}

func GetAllUsers(c *gin.Context) {
	var users []models.User

	// Preload the Blogs field to load associated blogs for each user
	if err := database.DB.Preload("Blogs").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}

	// Hide the password field from the response
	for i := range users {
		users[i].Password = ""
		for j := range users[i].Blogs {
			users[i].Blogs[j].Author = nil
		}
	}

	// Respond with the list of users in JSON format
	c.JSON(http.StatusOK, users)
}

func LoginUser(c *gin.Context) {
	var input models.User
	var user models.User

	// Bind JSON input to the User struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find the user by email
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Check if the provided password matches the hashed password
	if !utils.CheckPassword(user.Password, input.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, _ := middleware.CreateToken(user.ID)

	// Do not return the password in the response
	user.Password = ""

	// Respond with the user information (e.g., for creating a session or token)
	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "user": user, "token": token})
}
