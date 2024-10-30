package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"schooglink-task/controllers"
	"schooglink-task/database"
	"schooglink-task/models"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	// Assuming you have a route for creating a user
	router.POST("/users", controllers.CreateUser)
	return router
}

func TestCreateUserIntegration(t *testing.T) {
	// Setup the router
	router := setupRouter()

	// Sample user data for registration
	user := models.User{
		Name:     "Test User",
		Email:    "testuser@example.com",
		Password: "securepassword", // Make sure to hash this in your controller
	}

	// Convert user to JSON
	jsonUser, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("could not marshal user: %v", err)
	}

	// Create a request to create a user
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonUser))
	req.Header.Set("Content-Type", "application/json")

	// Record the response
	w := httptest.NewRecorder()

	// Serve the HTTP request
	router.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)

	// Deserialize response to verify response data
	var createdUser models.User
	err = json.Unmarshal(w.Body.Bytes(), &createdUser)
	assert.NoError(t, err)

	// Check that the user was created successfully
	assert.Equal(t, user.Name, createdUser.Name)
	assert.Equal(t, user.Email, createdUser.Email)

	// Cleanup - remove the user from the database if needed
	database.DB.Delete(&createdUser, createdUser.ID)
}
