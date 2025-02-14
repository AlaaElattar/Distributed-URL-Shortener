package storage

import (
	"context"
	"fmt"
	"os"
	"time"
	"url-shortener/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDBClient for storing analytics logs in MongoDB.
type MongoDBClient struct {
	Client *mongo.Client
	DB     *mongo.Database
}

// NewMongoDBClient creates new MongoDB client
func NewMongoDBClient() (MongoDBClient, error) {
	mongoURI := os.Getenv("MONGO_URI")

	//TODO: check timeout ??
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return MongoDBClient{}, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	db := client.Database("analytics")

	return MongoDBClient{
		Client: client,
		DB:     db,
	}, nil
}

// LogAccess saves an access log to MongoDB collection.
func (mongodb *MongoDBClient) LogAccess(log models.AccessLog) error {
	collection := mongodb.DB.Collection("access_logs")

	_, err := collection.InsertOne(context.TODO(), bson.M{
		"short_id": log.ShortID,
		"timestamp": log.Timestamp,
		"user_ip": log.UserIP,
	})
	return err

}
