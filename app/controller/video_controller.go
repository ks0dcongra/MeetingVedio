package controller

import (
	"MeetingVideoHelper/app/contract"
	"MeetingVideoHelper/app/service"
	"encoding/base64"
	"fmt"
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

	if videoID == "" {
		c.HTML(http.StatusBadRequest, "index.html", gin.H{
			"status": contract.ERROR_Download_Video_NotExist,
			"msg":    fmt.Sprintf("%s: %s", contract.Message[contract.ERROR_Download_Video], contract.Message[contract.ERROR_Download_Video_NotExist]),
		})
		return
	}

	videoDownload, err := primitive.ObjectIDFromHex(videoID[10:34])
	if err != nil {
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{
			"status": contract.ERROR_Video_To_Hex,
			"msg":    fmt.Sprintf("%s: %s", contract.Message[contract.ERROR_Download_Video], contract.Message[contract.ERROR_Video_To_Hex]),
		})
		return
	}

	filePath, err := vc.videoService.GetVideoFile(videoDownload)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{
			"status": contract.ERROR_Download_Video,
			"msg":    fmt.Sprintf("%s: %s", contract.Message[contract.ERROR_Download_Video], err.Error()),
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
	// ============================== Get upload file ==============================
	// maximum upload of 10 MB files
	c.Request.ParseMultipartForm(10 << 20)

	// get the upload file
	file, _, err := c.Request.FormFile("myFile")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{
			"status": contract.ERROR_Upload_Video_NotExist,
			"msg":    fmt.Sprintf("%s: %s", contract.Message[contract.ERROR_Upload_Video], contract.Message[contract.ERROR_Upload_Video_NotExist]),
		})
		return
	}

	file.Seek(0, 0)

	// convert to bytes
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{
			"status": contract.ERROR_Video_Format,
			"msg":    fmt.Sprintf("%s: %s", contract.Message[contract.ERROR_Upload_Video], contract.Message[contract.ERROR_Video_Format]),
		})
		return
	}

	file.Close()

	// Call the service layer to handle file upload and video processing
	insertVideo, err := vc.videoService.UploadedVideo(fileBytes)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{
			"status": contract.ERROR_Download_Video,
			"msg":    fmt.Sprintf("%s: %s", contract.Message[contract.ERROR_Upload_Video], err.Error()),
		})

		return
	}

	// MongoDB: Convert the VideoData to base64 encoding
	videoDataBase64 := base64.StdEncoding.EncodeToString(insertVideo.VideoData)

	// uploaded file successfully
	c.HTML(http.StatusOK, "index.html", gin.H{
		"status":        contract.Success,
		"msg":           contract.Message[contract.Success],
		"tempVideoFile": videoDataBase64,
		"videoID":       insertVideo.VID,
		"videoData":     insertVideo.VideoData,
	})
}
