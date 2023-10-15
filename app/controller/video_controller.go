package controller

import (
	"MeetingVideoHelper/app/model"
	"MeetingVideoHelper/database"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"

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
	// dstPath := "./static/videos/" + handler.Filename
	// dst, err := os.Create(dstPath)
	// if err != nil {
	// 	c.HTML(http.StatusInternalServerError, "index.html", gin.H{
	// 		"err": err.Error(),
	// 	})
	// }
	// defer dst.Close()

	// General: Copy the uploaded file to the created file on the filesystem
	// if _, err := io.Copy(dst, file); err != nil {
	// 	c.HTML(http.StatusInternalServerError, "index.html", gin.H{
	// 		"err": err.Error(),
	// 	})
	// }

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
		VID:       primitive.NewObjectID(),
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
	// 查詢到临时文件
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

	tempVideo.Write(insertVideo.VideoData)
	tempVideo.Close()
	// inputFilePath := tempVideo.Name()

	// 复制原始视频
	// copyCmd := exec.Command("ffmpeg", "-i", inputFilePath, "-c:v", "copy", "copyVideo.mp4")
	// copyCmd.Stdout = os.Stdout
	// copyCmd.Stderr = os.Stderr

	// err = copyCmd.Run()
	// if err != nil {
	// 	fmt.Println("Error copying video:", err)
	// 	return
	// }

	// 逆转复制的视频
	reverseCmd := exec.Command("ffmpeg", "-i", "tempVideo.mp4", "-vf", "reverse", "-af", "areverse", "reversedVideo.mp4")
	reverseCmd.Stdout = os.Stdout
	reverseCmd.Stderr = os.Stderr

	err = reverseCmd.Run()
	if err != nil {
		fmt.Println("Error reversing video:", err)
		return
	}

	// 串接视频
	// finalCmd := exec.Command("ffmpeg", "-safe", "0", "-f", "concat", "copyVideo.mp4", "-c:v ", "copy", "-c:a ", "copy", "copyVideo2.mp4")
	// finalCmd := exec.Command("ffmpeg", "-f", "concat", "-i", "../../tempVideo.mp4", "../../copyVideo.mp4")
	// ffmpeg -i 要被串接的影音檔案路徑1 -i 要被串接的影音檔案路徑2 -i 要被串接的影音檔案路徑3 -filter_complex "[0:v][0:a][1:v][1:a][2:v][2:a]concat=n=3:v=1:a=1[outv][outa]" -map "[outv]" -map "[outa]" 輸出的影片檔案路徑
	finalCmd := exec.Command("ffmpeg", "-i", "./tempVideo.mp4", "-i", "./reversedVideo.mp4", "-filter_complex", "[0:v][0:a][1:v][1:a]concat=n=2:v=1:a=1[outv][outa]", "-map", "[outv]", "-map", "[outa]", "-n", "./flippedVideo.mp4")
	finalCmd.Stdout = os.Stdout
	finalCmd.Stderr = os.Stderr

	err = finalCmd.Run()
	if err != nil {
		fmt.Println("Error concat video:", err)
		return
	}

	// 删除临时文件
	if err := os.Remove("tempVideo.mp4"); err != nil {
		fmt.Println("deleting tempVideo fail", err)
	}

	if err := os.Remove("reversedVideo.mp4"); err != nil {
		fmt.Println("deleting reversedVideo fail:", err)
	}

	// MongoDB operation
	file.Seek(0, 0) // Reset the file pointer to the beginning of the file
	// MongoDB: Read the uploaded file data

	data2, err := os.ReadFile("./flippedVideo.mp4")
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Error Reading the File")
		return
	}

	if err := os.Remove("flippedVideo.mp4"); err != nil {
		fmt.Println("deleting flippedVideo fail", err)
	}

	insertFlippedVideo := model.Video{
		VID:       primitive.NewObjectID(),
		Title:     "flippedVideo",
		VideoData: data2,
		UpdatedAt: time.Now(),
	}

	_, err = database.QmgoConnection.InsertOne(context.Background(), insertFlippedVideo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "MongoDB: Error Inserting the insertFlippedVideo fail")
		return
	}

	// log.Println("insertFlippedVideo.Title: ", insertFlippedVideo.Title)
	log.Println("insertFlippedVideo.ID: ", insertFlippedVideo.VID)
	// ============================== Render to index.html ==============================

	var videoResult model.Video
	// MongoDB: find one specific video from mongodb
	err = database.QmgoConnection.Find(context.Background(), bson.M{"vid": insertFlippedVideo.VID}).One(&videoResult)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "MongoDB: Finding the insertFlippedVideo fail")
		return
	}

	// MongoDB: Convert the VideoData to base64 encoding
	videoBase64 := base64.StdEncoding.EncodeToString(insertFlippedVideo.VideoData)

	// uploaded file successfully
	c.HTML(http.StatusOK, "index.html", gin.H{
		"successMsg":    "Uploaded File Successfully ",
		"err":           "",
		"tempVideoFile": videoBase64,
		"videoID":       insertFlippedVideo.VID,
	})
}
