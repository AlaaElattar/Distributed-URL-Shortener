package app

import (
	"net/http"
	"time"
	"url-shortener/models"

	"github.com/gin-gonic/gin"
	"github.com/jaevor/go-nanoid"
)

// LongURLInput type for getting long url from user.
type LongURLInput struct {
	Url string `json:"long_url" binding:"required"`
}

// ShortenURL handles request parsing, generating shortened ID.
func (app *App) ShortenURL(c *gin.Context) {
	var longInput LongURLInput

	if err := c.ShouldBindJSON(&longInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read url",
		})
		return
	}

	canonicID, err := nanoid.Standard(6)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate short ID"})
		return
	}

	shortID := canonicID()

	if err := app.redisClient.SaveURL(shortID, longInput.Url); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save URL",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"short_url": shortID})
}

// GetURL retrieves the original url from shortened one.
func (app *App) GetURL(c *gin.Context) {
	shortID := c.Param("shortID")
	userIP := c.ClientIP()

	originalURL, err := app.redisClient.GetURL(shortID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Short URL not found",
		})
		return
	}

	logEntry := models.AccessLog{
		ShortID:   shortID,
		Timestamp: time.Now(),
		UserIP:    userIP,
	}

	app.logsChannel <- logEntry

	c.JSON(http.StatusOK, gin.H{
		"long_url": originalURL,
	})
}