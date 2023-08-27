package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Video struct {
	ID        primitive.ObjectID `bson:"vid"`
	Title     string             `bson:"title"`
	VideoData []byte             `bson:"video_data"`
}
