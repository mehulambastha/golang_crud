(the env should never be pushed upstream, however in this project the env did not contain anything significant, therefore it is pushed as it is.)
# Go Blog API

## Overview

The Go Blog API is a RESTful web service built with Go, utilizing the Gin framework for efficient routing and GORM for database interactions. This project allows users to create, read, update, and delete blog posts securely. It adheres to best practices in software development, such as code organization, testing, and security, making it a robust choice for web applications.

## Features

- **User Management:** Users can register and log in, with passwords securely hashed.
- **Blog Management:** Users can create, read, update, and delete blog posts.
- **JWT Authentication:** Routes are protected with JWT to ensure secure access.
- **Input Validation:** User input is validated and sanitized to prevent SQL injection and other attacks.
- **Testing:** The application includes unit and integration tests for core functionality.
- **API Documentation:** Swagger is integrated for easy API documentation.

## Installation

To get started with the Go Blog API, follow these steps:

### Prerequisites

- Go (version 1.16 or later)
- SQLite database

### Steps

1. **Clone the Repository:**

   ```bash
   git clone https://github.com/mehulambastha/golang_crud.git
   cd golang_crud
   ```

2. **Install Dependencies:**

   Use Go modules to install the necessary dependencies:

   ```bash
   go mod tidy
   ```

3. **Set Up Environment Variables:**

   Create a `.env` file in the root directory and define the following variables:

   ```bash
   JWT_SECRET=your_jwt_secret
   ```

4. **Run Database Migrations:**

   Ensure the database schema is up-to-date:

   ```bash
   go run main.go migrate
   ```

5. **Start the Application:**

   Run the application:

   ```bash
   go run main.go
   ```

6. **Access the API:**

   The API will be available at `http://localhost:8080`.

## Code Snippets

### User Registration

The `CreateUser` function handles user registration, hashing passwords, and saving users to the database:

```go
func CreateUser(c *gin.Context) {
    var user models.User

    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    hashedPassword, err := utils.HashPassword(user.Password)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        return
    }

    user.Password = hashedPassword
    user.CreatedAt = time.Now()
    user.UpdatedAt = time.Now()

    if err := database.DB.Create(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        return
    }

    user.Password = ""
    c.JSON(http.StatusOK, user)
}
```

### JWT Authentication Middleware

The authentication middleware verifies the JWT token before allowing access to protected routes:

```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if len(token) != 0 {
            token = strings.Split(token, " ")[1]
        }
        _, valid := isValidToken(token)
        if token == "" || !valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized. Invalid login token."})
            c.Abort()
            return
        }

        userId, err := GetUserIDFromToken(token)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized. Invalid login token."})
            c.Abort()
            return
        }

        c.Set("userID", userId)
        c.Next()
    }
}
```

### Testing

The application includes comprehensive tests for core functionalities:

```go
func TestCreateUserIntegration(t *testing.T) {
    router := setupRouter()
    user := models.User{Name: "Test User", Email: "testuser@example.com", Password: "securepassword"}

    jsonUser, _ := json.Marshal(user)
    req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonUser))
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
}
```

## Best Practices

This project follows best practices in software development:

- **Code Organization:** Code is structured with a clear separation of concerns, utilizing models, controllers, and middleware.
- **Security:** Passwords are hashed using bcrypt, and input validation prevents common vulnerabilities like SQL injection.
- **Testing:** Unit and integration tests ensure the application's reliability and robustness.
- **Documentation:** Swagger documentation provides clear guidelines on API usage.

## Conclusion

This project showcases my ability to develop secure and scalable web applications using Go. I am eager to continue learning and contributing as a Go intern, where I can further hone my skills and collaborate with a team of talented developers.

Feel free to explore the code, and I welcome any feedback or contributions!
