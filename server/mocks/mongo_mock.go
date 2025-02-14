package mocks

import (
	"errors"
	"url-shortener/models"

	"github.com/stretchr/testify/mock"
)

// MockMongoDBClient is a mock implementation of MongoDBClient.
type MockMongoDBClient struct {
	mock.Mock
	ForceError bool
}

// LogAccess stores an access log or returns an error if forced.
func (m *MockMongoDBClient) LogAccess(log models.AccessLog) error {
	if m.ForceError {
		return errors.New("failed to log access")
	}
	args := m.Called(log)
	return args.Error(0)
}
