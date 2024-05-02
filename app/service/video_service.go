package service

import (
	"MeetingVideoHelper/app/model"
	"MeetingVideoHelper/app/repository"
	"fmt"
	"os"
	"os/exec"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VideoService struct {
	videoRepository *repository.VideoRepository
}

func NewVideoService() *VideoService {
	return &VideoService{
		videoRepository: repository.NewVideoRepository(),
	}
}

func (vs *VideoService) GetVideoFile(videoID primitive.ObjectID) (string, error) {
	videoData, err := vs.videoRepository.GetVideoData(videoID)
	if err != nil {
		return "", err
	}

	tempFile, err := os.CreateTemp("./", "tempVideo*.mp4")
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	_, err = tempFile.Write(videoData)
	if err != nil {
		return "", err
	}

	return tempFile.Name(), nil
}

func (vs *VideoService) UploadedVideo(fileBytes []byte) (*model.Video, error) {
	// write uploaded files to temporary files
	tempVideoPath, err := createTempFile(fileBytes)
	if err != nil {
		return nil, err
	}
	defer os.Remove(tempVideoPath)

	// reverse video
	reversedVideoPath, err := reverseVideo(tempVideoPath)
	if err != nil {
		return nil, err
	}
	defer os.Remove(reversedVideoPath)

	// merge original video and reverse it
	concatVideoPath, err := concatVideos(tempVideoPath, reversedVideoPath)
	if err != nil {
		return nil, err
	}
	defer os.Remove(concatVideoPath)

	// read file data
	concatVideoBytes, err := os.ReadFile(concatVideoPath)
	if err != nil {
		return nil, err
	}

	// save data to database
	insertVideo, err := vs.videoRepository.SaveVideo(concatVideoBytes)
	if err != nil {
		return nil, err
	}

	return insertVideo, nil
}

func createTempFile(fileBytes []byte) (string, error) { 
	tempVideo, err := os.CreateTemp("./static/videos", "tempVideo*.mp4")
	if err != nil {
		fmt.Println("Error creating temporary file:", err)
		return "", err
	}
	defer tempVideo.Close()

	tempVideo.Write(fileBytes)
	
	return tempVideo.Name(), nil
}

func reverseVideo(tempVideoName string) (string, error) {
	tempVideoNameMp4 := tempVideoName 
	reversedVideoPath := "./reversedVideo.mp4"
	 // reverse vedio
	reverseCmd := exec.Command("ffmpeg", "-i", tempVideoNameMp4, "-vf", "reverse", "-af", "areverse", reversedVideoPath)
	reverseCmd.Stdout = os.Stdout
	reverseCmd.Stderr = os.Stderr

	err := reverseCmd.Run()
	if err != nil {
		fmt.Println("Error reversing video:", err)
		return "", err
	}
	return reversedVideoPath, nil
}

func concatVideos(video1Path, video2Path string) (string, error) { 
	concatVideoPath := "./flippedVideo.mp4"
	// ffmpeg -i 要被串接的影音檔案路徑1 -i 要被串接的影音檔案路徑2 -i 要被串接的影音檔案路徑3 -filter_complex "[0:v][0:a][1:v][1:a][2:v][2:a]concat=n=3:v=1:a=1[outv][outa]" -map "[outv]" -map "[outa]" 輸出的影片檔案路徑
	finalCmd := exec.Command("ffmpeg", "-i", video1Path, "-i", video2Path, "-filter_complex", "[0:v][0:a][1:v][1:a]concat=n=2:v=1:a=1[outv][outa]", "-map", "[outv]", "-map", "[outa]", "-n", concatVideoPath)
	finalCmd.Stdout = os.Stdout
	finalCmd.Stderr = os.Stderr

	err := finalCmd.Run()
	if err != nil {
		fmt.Println("Error concat video:", err)
		return "", err
	}

	return concatVideoPath, nil
}