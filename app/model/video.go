package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Video struct {
	VID       primitive.ObjectID `bson:"vid"`
	VideoData []byte             `bson:"video_data"`
	UpdatedAt time.Time          `bson:"updated_at"`
}
