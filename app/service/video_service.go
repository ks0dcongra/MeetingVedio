package service

import (
	"MeetingVideoHelper/app/model"
	"MeetingVideoHelper/app/repository"
	"crypto/rand"
	"fmt"
	"io"
	"log"
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

func (vs *VideoService) GetVideoFile(videoID primitive.ObjectID, GIFTag string) (string, error) {
	videoData, err := vs.videoRepository.GetVideoData(videoID)
	if err != nil {
		log.Println("get video fail:", err)
		return "", err
	}

	tempFile, err := os.CreateTemp("./", "tempVideo*.mp4")
	if err != nil {
		log.Println("create file fail:", err)
		return "", err
	}
	defer tempFile.Close()

	_, err = tempFile.Write(videoData)
	if err != nil {
		log.Println("write file fail:", err)
		return "", err
	}
	
	if GIFTag != "" {
		randomFilePath, err := convertGIF(tempFile.Name())
		if err != nil {
			log.Println("covert file to GIF fail:", err)
			return "", err
		}
		tempFile.Close()
		os.Remove(tempFile.Name())
	
		return randomFilePath, nil
	}

	return tempFile.Name(), nil
}

func (vs *VideoService) UploadedVideo(fileBytes []byte) (*model.Video, error) {
	// write uploaded files to temporary files
	tempVideoPath, err := createTempFile(fileBytes)
	if err != nil {
		log.Println("create temporary file fail:", err)
		return nil, err
	}
	defer os.Remove(tempVideoPath)

	// reverse video
	reversedVideoPath, err := reverseVideo(tempVideoPath)
	if err != nil {
		log.Println("reverse video fail:", err)
		return nil, err
	}
	defer os.Remove(reversedVideoPath)

	// merge original video and reverse it
	concatVideoPath, err := concatVideos(tempVideoPath, reversedVideoPath)
	if err != nil {
		log.Println("concat video fail:", err)
		return nil, err
	}
	defer os.Remove(concatVideoPath)

	// read file data
	concatVideoBytes, err := os.ReadFile(concatVideoPath)
	if err != nil {
		log.Println("read file fail:", err)
		return nil, err
	}

	// save data to database
	insertVideo, err := vs.videoRepository.SaveVideo(concatVideoBytes)
	if err != nil {
		log.Println("insert video fail", err)
		return nil, err
	}

	return insertVideo, nil
}

func createTempFile(fileBytes []byte) (string, error) {
	tempVideo, err := os.CreateTemp("./static/videos", "tempVideo*.mp4")
	if err != nil {
		return "", err
	}
	defer tempVideo.Close()

	tempVideo.Write(fileBytes)

	return tempVideo.Name(), nil
}

func reverseVideo(tempVideoName string) (string, error) {
	tempVideoNameMp4 := tempVideoName

	// 生成隨機檔案名稱
	randomBytes := make([]byte, 8)
	_, err := io.ReadFull(rand.Reader, randomBytes)
	if err != nil {
		return "", err
	}
	randomFilePath := fmt.Sprintf("./%x.mp4", randomBytes)

	// reverse video
	reverseCmd := exec.Command("ffmpeg", "-i", tempVideoNameMp4, "-vf", "reverse", "-af", "areverse", randomFilePath)
	reverseCmd.Stdout = os.Stdout
	reverseCmd.Stderr = os.Stderr

	err = reverseCmd.Run()
	if err != nil {
		return "", err
	}
	return randomFilePath, nil
}

func concatVideos(video1Path, video2Path string) (string, error) {
	// generate random file name
	randomBytes := make([]byte, 8)
	_, err := io.ReadFull(rand.Reader, randomBytes)
	if err != nil {
		return "", err
	}
	randomFilePath := fmt.Sprintf("./%x.mp4", randomBytes)

	// ffmpeg -i 要被串接的影音檔案路徑1 -i 要被串接的影音檔案路徑2 -i 要被串接的影音檔案路徑3 -filter_complex "[0:v][0:a][1:v][1:a][2:v][2:a]concat=n=3:v=1:a=1[outv][outa]" -map "[outv]" -map "[outa]" 輸出的影片檔案路徑
	finalCmd := exec.Command("ffmpeg", "-i", video1Path, "-i", video2Path, "-filter_complex", "[0:v][0:a][1:v][1:a]concat=n=2:v=1:a=1[outv][outa]", "-map", "[outv]", "-map", "[outa]", "-n", randomFilePath)
	finalCmd.Stdout = os.Stdout
	finalCmd.Stderr = os.Stderr

	err = finalCmd.Run()
	if err != nil {
		return "", err
	}

	return randomFilePath, nil
}

func convertGIF(tempFileName string) (string, error) {
	tempFileNameGIF := tempFileName

	// generate random file name
	randomBytes := make([]byte, 8)
	_, err := io.ReadFull(rand.Reader, randomBytes)
	if err != nil {
		return "", err
	}
	randomFilePath := fmt.Sprintf("./tempGIF%x.gif", randomBytes)

	// reverse video
	reverseCmd := exec.Command("ffmpeg", "-i", tempFileNameGIF, "-vf", "fps=10,scale=320:-1:flags=lanczos", "-c:v", "gif", "-n", randomFilePath)
	reverseCmd.Stdout = os.Stdout
	reverseCmd.Stderr = os.Stderr

	err = reverseCmd.Run()
	if err != nil {
		return "", err
	}
	return randomFilePath, nil
}