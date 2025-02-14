package app

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"url-shortener/mocks"
	"url-shortener/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestShortenURL(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRedis := new(mocks.MockRedisClient)

	testApp := &App{
		router:      gin.Default(),
		redisClient: mockRedis,
	}

	testApp.router.POST("/shorten", testApp.ShortenURL)

	tests := []struct {
		name           string
		requestBody    string
		forceRedisFail bool
		expectedStatus int
		expectedBody   string
		expectSaveURL  bool
	}{
		{
			name:           "Test with valid URL.",
			requestBody:    `{"long_url": "https://example.com"}`,
			forceRedisFail: false,
			expectedStatus: http.StatusOK,
			expectedBody:   `"short_url"`,
			expectSaveURL:  true,
		},
		{
			name:           "Test with Invalid URL.",
			requestBody:    `{"long_url": }`,
			forceRedisFail: false,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"error":"Failed to read url"`,
			expectSaveURL:  false,
		},
		{
			name:           "Test with failure of redis.",
			requestBody:    `{"long_url": "https://example.com"}`,
			forceRedisFail: true,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `"error":"Failed to save URL"`,
			expectSaveURL:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRedis.ForceError = tc.forceRedisFail

			if tc.expectSaveURL {
				var returnErr error
				if tc.forceRedisFail {
					returnErr = errors.New("Redis failure")
				}

				mockRedis.On("SaveURL", mock.Anything, mock.Anything).Return(returnErr).Once()
			}

			req, err := http.NewRequest("POST", "/shorten", bytes.NewBufferString(tc.requestBody))
			require.NoError(t, err, "Failed to create request")
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()
			testApp.router.ServeHTTP(recorder, req)

			assert.Equal(t, tc.expectedStatus, recorder.Code, "Unexpected status code")
			assert.Contains(t, recorder.Body.String(), tc.expectedBody, "Unexpected response body")

			if tc.expectSaveURL {
				mockRedis.AssertExpectations(t)
			}
		})
	}
}

func TestGetURL(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRedis := new(mocks.MockRedisClient)
	mockMongo := new(mocks.MockMongoDBClient)

	logsChannel := make(chan models.AccessLog, 10)

	testApp := &App{
		router:        gin.Default(),
		redisClient:   mockRedis,
		mongoDBClient: mockMongo,
		logsChannel:   logsChannel,
	}

	testApp.router.GET("/:shortID", testApp.GetURL)

	tests := []struct {
		name           string
		shortID        string
		mockRedisURL   string
		mockRedisError error
		expectedStatus int
		expectedBody   string
		expectLogEntry bool
	}{
		{
			name:           "Test with existing URL.",
			shortID:        "short123",
			mockRedisURL:   "https://example.com",
			mockRedisError: nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `"long_url":"https://example.com"`,
			expectLogEntry: true,
		},
		{
			name:           "Test with non-existing URL.",
			shortID:        "notfound",
			mockRedisURL:   "",
			mockRedisError: assert.AnError,
			expectedStatus: http.StatusNotFound,
			expectedBody:   `"error":"Short URL not found"`,
			expectLogEntry: false,
		},
		{
			name:           "Test with failure of redis.",
			shortID:        "error",
			mockRedisURL:   "",
			mockRedisError: assert.AnError,
			expectedStatus: http.StatusNotFound,
			expectedBody:   `"error":"Short URL not found"`,
			expectLogEntry: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRedis.On("GetURL", tc.shortID).Return(tc.mockRedisURL, tc.mockRedisError).Once()

			req, err := http.NewRequest("GET", "/"+tc.shortID, nil)
			require.NoError(t, err, "Failed to create request")

			recorder := httptest.NewRecorder()
			testApp.router.ServeHTTP(recorder, req)

			assert.Equal(t, tc.expectedStatus, recorder.Code, "Unexpected status code")
			assert.Contains(t, recorder.Body.String(), tc.expectedBody, "Unexpected response body")

			if tc.expectLogEntry {
				select {
				case logEntry := <-logsChannel:
					assert.Equal(t, tc.shortID, logEntry.ShortID, "Log entry has incorrect ShortID")
				default:
					assert.Fail(t, "Expected log entry, but none was sent")
				}
			} else {
				assert.Empty(t, logsChannel, "Expected no log entry, but one was sent")
			}

			mockRedis.AssertExpectations(t)
		})
	}
}
