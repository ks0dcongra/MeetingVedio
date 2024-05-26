package controller

import (
	"MeetingVideoHelper/app/contract"
	"MeetingVideoHelper/app/service"
	"encoding/base64"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VideoController struct {
	videoService *service.VideoService
}

func NewVideoController() *VideoController {
	return &VideoController{
		videoService: service.NewVideoService(),
	}
}

func (vc *VideoController) DownloadVideo(c *gin.Context) {
	videoID := c.PostForm("videoID")
	GIFTag := c.PostForm("GIF")

	if videoID == "" {
		c.HTML(http.StatusBadRequest, "index.html", gin.H{
			"status": contract.ERROR_Download_Video_NotExist,
			"msg":    contract.Message[contract.ERROR_Download_Video_NotExist],
		})
		return
	}

	videoDownloadID, err := primitive.ObjectIDFromHex(videoID[10:34])
	if err != nil {
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{
			"status": contract.ERROR_Video_To_Hex,
			"msg":    contract.Message[contract.ERROR_Video_To_Hex],
		})
		return
	}

	filePath, err := vc.videoService.GetVideoFile(videoDownloadID, GIFTag)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{
			"status": contract.ERROR_Download_Video,
			"msg":    contract.Message[contract.ERROR_Download_Video],
		})
		return
	}
	defer os.Remove(filePath)

	filePath = filePath[2:]
	c.Header("Content-Disposition", "attachment; filename="+filePath)
	c.Header("Content-Type", "video/mp4")
	c.File(filePath)

	c.HTML(http.StatusOK, "index.html", gin.H{})
}

func (vc *VideoController) UploadVideo(c *gin.Context) {
	// get the upload file
	file, header, err := c.Request.FormFile("myFile")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{
			"status": contract.ERROR_Upload_Video_NotExist,
			"msg":    contract.Message[contract.ERROR_Upload_Video_NotExist],
		})
		return
	}

	maxSize := int64(30 << 20)

	if header.Size > maxSize {
		c.HTML(http.StatusBadRequest, "index.html", gin.H{
			"status": contract.ERROR_Upload_Video_FileTooLarge,
			"msg":    contract.Message[contract.ERROR_Upload_Video_FileTooLarge],
		})
		return
	}

	file.Seek(0, 0)

	// convert to bytes
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{
			"status": contract.ERROR_Video_Format,
			"msg":    contract.Message[contract.ERROR_Video_Format],
		})
		return
	}

	file.Close()

	// Call the service layer to handle file upload and video processing
	insertVideo, err := vc.videoService.UploadedVideo(fileBytes)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{
			"status": contract.ERROR_Download_Video,
			"msg":    contract.Message[contract.ERROR_Upload_Video],
		})

		return
	}

	// MongoDB: Convert the VideoData to base64 encoding
	videoDataBase64 := base64.StdEncoding.EncodeToString(insertVideo.VideoData)

	// uploaded file successfully
	c.HTML(http.StatusOK, "index.html", gin.H{
		"status":          contract.Success,
		"msg":             contract.Message[contract.Success],
		"videoDataBase64": videoDataBase64,
		"videoID":         insertVideo.VID,
	})
}
