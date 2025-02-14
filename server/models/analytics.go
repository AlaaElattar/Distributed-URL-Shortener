package models

import "time"

// AccessLog type for saving logs in mongoDB.
type AccessLog struct {
	ShortID   string    `bson:"shortID"`
	UserIP    string    `bson:"userIP"`
	Timestamp time.Time `bson:"timestamp"`
}
