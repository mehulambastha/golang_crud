package main

import (
	"log"
	"schooglink-task/config"
	"schooglink-task/database"
	"schooglink-task/routes"
	"schooglink-task/utils"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"golang.org/x/time/rate"
)

func main() {
	config.LoadConfig()

	database.InitDB()

	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	rateLimitMiddleware := utils.RateLimiter(rate.Limit(1), 5)

	r.Use(rateLimitMiddleware)

	routes.RegisterBlogRoutes(r)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start the server.")
	}
}
