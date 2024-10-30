package middleware

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"schooglink-task/database"
	"schooglink-task/models"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

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
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized. Invalid login token. Unable to fetch user associated with this."})
			c.Abort()
			return
		}

		c.Set("userID", userId)
		c.Next()
	}
}

func BlogOwnership() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		postID := c.Param("id")
		post, err := GetBlogPostByID(postID)

		if userIdFloat, ok := userId.(float64); ok {
			// Convert userId from float64 to uint, then to string
			userIdStr := fmt.Sprintf("%.0f", userIdFloat) // Removes any decimal point
			if err != nil || strconv.FormatUint(uint64(post.AuthorID), 10) != userIdStr {
				log.Println("Error is: ", err)
				log.Println("AuthorID is: ", post.AuthorID, "current userID is: ", userId)
				c.JSON(http.StatusForbidden, gin.H{"error": `You don't have permission to modify this post`})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

func GetBlogPostByID(postID string) (models.Blog, error) {
	var blog models.Blog

	if err := database.DB.First(&blog, postID).Error; err != nil {
		return models.Blog{}, err
	}

	return blog, nil
}

func isValidToken(tokenString string) (*jwt.Token, bool) {
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method.")
		}
		return jwtSecret, nil
	})
	return token, err == nil && token.Valid
}

func GetUserIDFromToken(tokenString string) (float64, error) {
	token, isValid := isValidToken(tokenString)

	if !isValid {
		return 0, errors.New("invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if userID, ok := claims["sub"].(float64); ok {
			return userID, nil
		}
	}

	return 0, errors.New("user ID not found in token.")
}

func CreateToken(userID uint) (string, error) {
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))

	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 72).Unix(), // Token expires in 72 hours
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}
