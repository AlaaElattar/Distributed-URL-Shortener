package models

import "time"

// AccessLog type for saving logs in mongoDB.
type AccessLog struct {
	ShortID   string    `bson:"shortID" validate:"required"`
	UserIP    string    `bson:"userIP" validate:"required,ip"`
	Timestamp time.Time `bson:"timestamp" validate:"lte"`
}
