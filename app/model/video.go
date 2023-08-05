package model

type Video struct {
	ID        int 			`bson:"vid"`
	Title     string        `bson:"title"`
	VideoData []byte        `bson:"video_data"`
}