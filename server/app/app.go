// Package app for url-shortener backend
package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"url-shortener/models"
	"url-shortener/storage"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// App holds all dependencies for the backend server.
type App struct {
	router        *gin.Engine
	redisClient   storage.RedisClient
	mongoDBClient storage.MongoDBClient
	logsChannel   chan models.AccessLog
}

// NewApp creates new app with all configurations.
func NewApp() (*App, error) {
	// Load environment variables if running app locally
	if os.Getenv("DOCKER_ENV") == "" {
		err := godotenv.Load("../.env")
		if err != nil {
			log.Println(".env file not found")
			return nil, err
		}
	}

	// Initialize Gin
	router := gin.Default()

	// Initialize redis client
	redis, err := storage.NewRedisClient()
	if err != nil {
		return nil, err
	}

	// Initialize mongodb client
	mongodb, err := storage.NewMongoDBClient()
	if err != nil {
		return nil, err
	}

	// Create a buffered logs channel of size 100.
	channel := make(chan models.AccessLog, 100)

	app := &App{
		router,
		redis,
		mongodb,
		channel,
	}

	// Start a background worker to log analytics in MongoDB.
	go app.startWorker()

	router.Use(app.rateLimiterMiddleware())

	app.setupRouters()

	return app, nil
}

// rateLimiterMiddleware limits users to 10 requests per minute
func (app *App) rateLimiterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		key := fmt.Sprintf("rate_limit:%s", clientIP)

		count, err := app.redisClient.IncrementRequests(key)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to check rate limit",
			})
			c.Abort()
			return
		}

		if count == 1 {
			err := app.redisClient.SetExpiration(key, 60*time.Second)
			if err != nil {
				log.Println("Failed to set expiration time")
			}

		}

		if count > 10 {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Please wait a minute before trying again.",
			})

			c.Abort()
			return
		}

		c.Next()
	}
}

// startWorker sends logs from buffered channel to mongodb
func (app *App) startWorker() {
	for logEntry := range app.logsChannel {
		err := app.mongoDBClient.LogAccess(logEntry)
		if err != nil {
			log.Printf("Failed to save analytics log for %v: %v:", logEntry, err)
		}
	}
}

// setupRouters setup the routers
func (app *App) setupRouters() {
	app.router.POST("/shorten", app.ShortenURL)
	app.router.GET("/:shortID", app.GetURL)
}

// Run starts the server on the specified port.
func (app *App) Run() {
	port := os.Getenv("SERVER_PORT")
	log.Printf("Server running on port %v", port)

	app.router.Run(":" + port)
}
