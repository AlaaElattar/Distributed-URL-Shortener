package mocks

import (
	"errors"
	"time"

	"github.com/stretchr/testify/mock"
)

// MockRedisClient is a mock implementation of Redis Client.
type MockRedisClient struct {
	mock.Mock
	ForceError bool
}

// SaveURL mocks the behavior of saving shortened URL in Redis.
func (m *MockRedisClient) SaveURL(shortID, originalURL string) error {
	if m.ForceError {
		return errors.New("failed to save url.")
	}
	args := m.Called(shortID, originalURL)
	return args.Error(0)
}

// GetURL mocks the behavior of retrieving url from Redis.
func (m *MockRedisClient) GetURL(shortID string) (string, error) {
	if m.ForceError {
		return "", errors.New("failed to get url")
	}
	args := m.Called(shortID)
	return args.String(0), args.Error(1)
}

// IncrementRequests mocks the behavior of incrementing requests count in Redis.
func (m *MockRedisClient) IncrementRequests(key string) (int64, error) {
	if m.ForceError {
		return 0, errors.New("failed to increment request")
	}
	args := m.Called(key)
	return int64(args.Int(0)), args.Error(1)
}

// SetExpiration mocks the behavior of setting an expiration time for a Redis key.
func (m *MockRedisClient) SetExpiration(key string, expiration time.Duration) error {
	if m.ForceError {
		return errors.New("failed to set expiration")
	}
	args := m.Called(key, expiration)
	return args.Error(0)
}
