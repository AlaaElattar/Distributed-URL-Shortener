package storage

import (
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis"
)

// RedisClient for redis DB storing short url mapping with expiration.
type RedisClient struct {
	DB *redis.Client
}

// NewRedisClient creates a new Redis client.
func NewRedisClient() (RedisClient, error) {
	redisAddress := os.Getenv("REDIS_ADDRESS")

	client := redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: "",
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return RedisClient{}, fmt.Errorf("failed to connect to Redis: %v", err)
	}
	fmt.Println("Connected to Redis")

	return RedisClient{DB: client}, nil
}

// SaveURL stores a shortened URL with an expiration of 30 days.
func (redis *RedisClient) SaveURL(shortID, originalURL string) error {
	err := redis.DB.Set(shortID, originalURL, 30*24*time.Hour).Err()
	if err != nil {
		fmt.Printf("failed to save URL: %v", err)
	}
	return err
}

// GetURL retrieves the original url from the shortened ID.
func (redis *RedisClient) GetURL(shortID string) (string, error) {
	url, err := redis.DB.Get(shortID).Result()

	if err != nil {
		return "", fmt.Errorf("failed to get URL: %v", err)
	}

	return url, nil
}

// IncrementRequests increments a request counter for a given key.
func (redis *RedisClient) IncrementRequests(key string)(int64, error){
	count, err := redis.DB.Incr(key).Result()
	if err != nil{
		return 0, err
	}

	return count, nil
}

// SetExpiration sets an expiration time for a given key.
func (redis *RedisClient) SetExpiration(key string, expiration time.Duration) error{
	_, err:= redis.DB.Expire(key, expiration).Result()
	return err

}


