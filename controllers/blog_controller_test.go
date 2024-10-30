package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"schooglink-task/controllers"
	"schooglink-task/models"
	"schooglink-task/utils/middleware"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateUserWithInvalidData(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(middleware.AuthMiddleware())
	router.POST("/users", controllers.CreateUser)

	token, _ := middleware.CreateToken(1)

	jsonData := `{"name": "", "email": "inValidEmail", "password": "short"}`

	req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(jsonData))

	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var responseBlog models.Blog

	errUnmarshalling := json.Unmarshal(w.Body.Bytes(), &responseBlog)

	if errUnmarshalling != nil {
		t.Errorf("Unamrshalling failed")
	}

	assert.Equal(t, "New Blog Post", responseBlog.Title)
	assert.Equal(t, "Blog Subtitle", responseBlog.Subtitle)
	assert.Equal(t, "This is the content of the blog post.", responseBlog.Content)
	assert.WithinDuration(t, time.Now(), responseBlog.CreatedAt, time.Second*2)
}

// Similarly we will write unit tests for other functions (Modifying, Deleting Users)
// Those route with have the BlogOwnership middleware too, after Authorization middleware.
//
// And similary we write unit test for user controller.
