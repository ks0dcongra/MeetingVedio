package repository

import (
	"MeetingVideoHelper/app/model"
	"MeetingVideoHelper/database"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VideoRepository struct {
}

func NewVideoRepository() *VideoRepository {
	return &VideoRepository{}
}

func (vr *VideoRepository) GetVideoData(videoID primitive.ObjectID) ([]byte, error) {
	var videoResult model.Video

	err := database.QmgoConnection.Find(context.Background(), bson.M{"vid": videoID}).One(&videoResult)
	if err != nil {
		return nil, err
	}

	return videoResult.VideoData, nil
}

func (vr *VideoRepository) SaveVideo(concatVideoBytes []byte) (*model.Video, error) {
	// save to mongodb
	insertVideo := model.Video{
		VID:       primitive.NewObjectID(),
		Title:     "concatVideo",
		VideoData: concatVideoBytes,
		UpdatedAt: time.Now(),
	}

	_, err := database.QmgoConnection.InsertOne(context.Background(), insertVideo)
	if err != nil {
		return nil, err
	}

	return &insertVideo, nil
}
