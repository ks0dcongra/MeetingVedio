package controller

import (
	"MeetingVideoHelper/app/model"
	"MeetingVideoHelper/database"
	"context"
	"encoding/base64"
	"fmt"
	"io"

	"os/exec"

	"net/http"
	"os"

	"time"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserController struct {
}

func NewUserController() *UserController {
	return &UserController{}
}

// upload file
func UploadFile(c *gin.Context) {
	// ============================== Get upload file ==============================

	var videoResult model.Video

	// General: maximum upload of 10 MB files
	c.Request.ParseMultipartForm(10 << 20)

	// General: get handler for filename, size and headers
	file, handler, err := c.Request.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{
			"err": "Error Retrieving the File",
		})
	}

	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// ============================== Store in local directory ==============================
	// General: Create file in the temp folder
	dstPath := "./static/videos/" + handler.Filename
	dst, err := os.Create(dstPath)
	defer dst.Close()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{
			"err": err.Error(),
		})
	}

	// General: Copy the uploaded file to the created file on the filesystem
	if _, err := io.Copy(dst, file); err != nil {
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{
			"err": err.Error(),
		})
	}

	// ============================== Insert to MongoDB ==============================
	// MongoDB operation
	file.Seek(0, 0) // Reset the file pointer to the beginning of the file
	// MongoDB: Read the uploaded file data
	data, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Error Reading the File")
		return
	}

	// MongoDB: insert data to mongodb
	insertVideo := model.Video{
		ID:        primitive.NewObjectID(),
		Title:     handler.Filename,
		VideoData: data,
		UpdatedAt: time.Now(),
	}

	// _, err = database.QmgoConnection.InsertOne(context.Background(), video)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, "MongoDB: Error Inserting the File")
	// 	return
	// } else {
	// ============================== Reverse video ==============================
	// 写入视频数据到临时文件
	// err = database.QmgoConnection.Find(context.Background(), bson.M{"title": insertVideo.Title, "vid": insertVideo.ID}).One(&videoResult)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, "MongoDB: Error Finding the File")
	// 	return
	// }

	tempVideo, err := os.Create("./tempVideo.mp4")
	if err != nil {
		fmt.Println("Error creating temporary file:", err)
		return
	}

	defer tempVideo.Close()
	tempVideo.Write(insertVideo.VideoData)

	inputFilePath := tempVideo.Name()

	// 复制原始视频
	copyCmd := exec.Command("ffmpeg", "-i", inputFilePath, "-c:v", "copy", "copyVideo.mp4")
	copyCmd.Stdout = os.Stdout
	copyCmd.Stderr = os.Stderr

	err = copyCmd.Run()
	if err != nil {
		fmt.Println("Error copying video:", err)
		return
	}

	// 逆转复制的视频
	// reverseCmd := exec.Command("ffmpeg", "-i", "copyVideo.mp4", "-vf", "reverse", "-af", "areverse", "reversedCopyVideo.mp4")
	// reverseCmd.Stdout = os.Stdout
	// reverseCmd.Stderr = os.Stderr

	// err = reverseCmd.Run()
	// if err != nil {
	// 	fmt.Println("Error reversing video:", err)
	// 	return
	// }

	// 串接视频
	// finalCmd := exec.Command("ffmpeg", "-safe", "0", "-f", "concat", "copyVideo.mp4", "-c:v ", "copy", "-c:a ", "copy", "copyVideo2.mp4")
	// finalCmd := exec.Command("ffmpeg", "-f", "concat", "-i", "../../tempVideo.mp4", "../../copyVideo.mp4")
	finalCmd := exec.Command("ffmpeg", "-i", "./tempVideo.mp4", "-i", "./copyVideo.mp4", "-filter_complex", "[0:v][0:a][1:v][1:a]concat=n=2:v=1:a=1[outv][outa]", "-map", "[outv]", "-map", "[outa]", "./flipped.mp4")
	finalCmd.Stdout = os.Stdout
	finalCmd.Stderr = os.Stderr
	// ffmpeg -i 要被串接的影音檔案路徑1 -i 要被串接的影音檔案路徑2 -i 要被串接的影音檔案路徑3 -filter_complex "[0:v][0:a][1:v][1:a][2:v][2:a]concat=n=3:v=1:a=1[outv][outa]" -map "[outv]" -map "[outa]" 輸出的影片檔案路徑

	err = finalCmd.Run()
	if err != nil {
		fmt.Println("Error concat video:", err)
		return
	}

	// 删除临时文件
	// if err := os.Remove("copy.mp4"); err != nil {
	//     fmt.Println("Error deleting temporary file:", err)
	// }

	// if err := os.Remove("flipped.mp4"); err != nil {
	//     fmt.Println("Error deleting temporary file:", err)
	// }\
	// _, err = database.QmgoConnection.InsertOne(context.Background(), bson.M{
	// 	"title":      "Flipped Video",
	// 	"video_data": outputBuffer.Bytes(),
	// })
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, "MongoDB: Error Inserting the File")
	// 	return
	// }
	// }

	// ============================== Render to index.html ==============================
	// MongoDB: find one specific video from mongodb
	err = database.QmgoConnection.Find(context.Background(), bson.M{"title": handler.Filename, "vid": insertVideo.ID}).One(&videoResult)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "MongoDB: Error Finding the File")
		return
	}

	// MongoDB: Convert the VideoData to base64 encoding
	videoBase64 := base64.StdEncoding.EncodeToString(videoResult.VideoData)

	// uploaded file successfully
	c.HTML(http.StatusOK, "index.html", gin.H{
		"successMsg":    "Uploaded File Successfully ",
		"err":           "",
		"tempVideoFile": videoBase64,
	})
}
