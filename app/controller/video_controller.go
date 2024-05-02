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
			"status": contract.ERROR_Download_Video,
			"msg":    contract.Message[contract.ERROR_Download_Video],
		})
		return
	}

	videoDownload, err := primitive.ObjectIDFromHex(videoID[10:34])
	if err != nil {
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{
			"status": contract.ERROR_Download_Video,
			"msg":    contract.Message[contract.ERROR_Download_Video],
		})
		return
	}

	filePath, err := vc.videoService.GetVideoFile(videoDownload)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{
			"status": contract.ERROR_Download_Video,
			"msg":    contract.Message[contract.ERROR_Download_Video],
		})
		return
	}
	defer os.Remove(filePath)
	filePath = filePath[2:]
	fmt.Println("filePathfilePathfilePath", filePath)
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
			"status": contract.ERROR_Upload_Video,
			"msg":    contract.Message[contract.ERROR_Upload_Video],
		})
		return
	}

	file.Seek(0, 0)

	// convert to bytes
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Error Reading the File")
		return
	}

	file.Close()

	// Call the service layer to handle file upload and video processing
	insertVideo, err := vc.videoService.UploadedVideo(fileBytes)
	if err != nil {
		// handle error
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
		"videoTitle":    insertVideo.Title,
		"videoData":     insertVideo.VideoData,
	})
}

// func UploadFile2(c *gin.Context) {
// 	// ============================== Get upload file ==============================
// 	// maximum upload of 10 MB files
// 	c.Request.ParseMultipartForm(10 << 20)

// 	// get handler for filename, size and headers
// 	file, _, err := c.Request.FormFile("myFile")
// 	if err != nil {
// 		fmt.Println("Error Retrieving the File")
// 		fmt.Println(err)
// 		c.HTML(http.StatusInternalServerError, "index.html", gin.H{
// 			"err": "Error Retrieving the File",
// 			"msg": "Uploaded File Successfully ",
// 		})
// 	}

// 	defer file.Close()

// 	// ============================== Create local file to store temp data ==============================
// 	// reset the file pointer to the beginning of the file
// 	file.Seek(0, 0)
// 	// convert to byte
// 	data, err := io.ReadAll(file)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, "Error Reading the File")
// 		return
// 	}

// 	tempVideo, err := os.Create("./tempVideo.mp4")
// 	if err != nil {
// 		fmt.Println("Error creating temporary file:", err)
// 		return
// 	}

// 	tempVideo.Write(data)
// 	tempVideo.Close()

// 	// reverse vedio
// 	reverseCmd := exec.Command("ffmpeg", "-i", "tempVideo.mp4", "-vf", "reverse", "-af", "areverse", "reversedVideo.mp4")
// 	reverseCmd.Stdout = os.Stdout
// 	reverseCmd.Stderr = os.Stderr

// 	err = reverseCmd.Run()
// 	if err != nil {
// 		fmt.Println("Error reversing video:", err)
// 		return
// 	}

// 	// ffmpeg -i 要被串接的影音檔案路徑1 -i 要被串接的影音檔案路徑2 -i 要被串接的影音檔案路徑3 -filter_complex "[0:v][0:a][1:v][1:a][2:v][2:a]concat=n=3:v=1:a=1[outv][outa]" -map "[outv]" -map "[outa]" 輸出的影片檔案路徑
// 	finalCmd := exec.Command("ffmpeg", "-i", "./tempVideo.mp4", "-i", "./reversedVideo.mp4", "-filter_complex", "[0:v][0:a][1:v][1:a]concat=n=2:v=1:a=1[outv][outa]", "-map", "[outv]", "-map", "[outa]", "-n", "./flippedVideo.mp4")
// 	finalCmd.Stdout = os.Stdout
// 	finalCmd.Stderr = os.Stderr

// 	err = finalCmd.Run()
// 	if err != nil {
// 		fmt.Println("Error concat video:", err)
// 		return
// 	}
// 	// 删除臨時文件
// 	if err := os.Remove("tempVideo.mp4"); err != nil {
// 		fmt.Println("deleting tempVideo fail", err)
// 	}

// 	if err := os.Remove("reversedVideo.mp4"); err != nil {
// 		fmt.Println("deleting reversedVideo fail:", err)
// 	}

// 	// MongoDB: Read the uploaded file data
// 	flippedVideoData, err := os.ReadFile("./flippedVideo.mp4")
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, "Error Reading the File")
// 		return
// 	}

// 	if err := os.Remove("flippedVideo.mp4"); err != nil {
// 		fmt.Println("deleting flippedVideo fail", err)
// 	}

// 	insertFlippedVideo := model.Video{
// 		VID:       primitive.NewObjectID(),
// 		Title:     "flippedVideo",
// 		VideoData: flippedVideoData,
// 		UpdatedAt: time.Now(),
// 	}

// 	_, err = database.QmgoConnection.InsertOne(context.Background(), insertFlippedVideo)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, "MongoDB: Error Inserting the insertFlippedVideo fail")
// 		return
// 	}

// 	// log.Println("insertFlippedVideo.Title: ", insertFlippedVideo.Title)
// 	log.Println("insertFlippedVideo.ID: ", insertFlippedVideo.VID)

// 	// TODO:MongoDB的Find，用不到，備份用
// 	// var videoResult model.Video
// 	// MongoDB: find one specific video from mongodb
// 	// err = database.QmgoConnection.Find(context.Background(), bson.M{"vid": insertFlippedVideo.VID}).One(&videoResult)
// 	// if err != nil {
// 	// 	c.JSON(http.StatusInternalServerError, "MongoDB: Finding the insertFlippedVideo fail")
// 	// 	return
// 	// }

// 	// MongoDB: Convert the VideoData to base64 encoding
// 	videoBase64 := base64.StdEncoding.EncodeToString(insertFlippedVideo.VideoData)

// 	// uploaded file successfully
// 	c.HTML(http.StatusOK, "index.html", gin.H{
// 		"Msg":           "Uploaded File Successfully ",
// 		"err":           "",
// 		"tempVideoFile": videoBase64,
// 		"videoID":       insertFlippedVideo.VID,
// 		"videoTitle":    insertFlippedVideo.Title,
// 		"videoData":     insertFlippedVideo.VideoData,
// 	})
// }
