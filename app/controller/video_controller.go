package controller

import (
	"MeetingVideoHelper/app/model"
	"MeetingVideoHelper/database"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserController struct {
	
}

func NewUserController() *UserController {
	return &UserController{
		
	}
}

// upload file
func UploadFile(c *gin.Context) {
	// General: maximum upload of 10 MB files
	c.Request.ParseMultipartForm(10 << 20)

	// General: get handler for filename, size and headers
	file, handler, err := c.Request.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{
			"err":  "Error Retrieving the File",
		})
	}

	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// General: Create file in the temp folder
	dstPath := "./static/videos/" + handler.Filename
	dst, err := os.Create(dstPath)
	defer dst.Close()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{
			"err":  err.Error(),
		})
	}

	// General: Copy the uploaded file to the created file on the filesystem
	if _, err := io.Copy(dst, file); err != nil {
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{
			"err":  err.Error(),
		})
	}

	// MongoDB operation
	file.Seek(0, 0) // Reset the file pointer to the beginning of the file
	// MongoDB: Read the uploaded file data
	data, err := io.ReadAll(file)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error Reading the File")
		return
	}
	
	// MongoDB: insert data to mongodb
	video := model.Video{
		ID:primitive.NewObjectID(),
		Title: handler.Filename,
		VideoData: data,
	}
	_, err = database.QmgoConnection.InsertOne(context.Background(), video)
	if err != nil {
		c.String(http.StatusInternalServerError, "MongoDB: Error Inserting the File")
		return
	}

	// MongoDB: find one specific video from mongodb
	var videoResult model.Video
	err = database.QmgoConnection.Find(context.Background(), bson.M{"title": handler.Filename, "vid":video.ID}).One(&videoResult)
	if err != nil {
		c.String(http.StatusInternalServerError, "MongoDB: Error Finding the File")
		return
	}

	// MongoDB: Convert the VideoData to base64 encoding
	videoBase64 := base64.StdEncoding.EncodeToString(videoResult.VideoData)
	dstPath = videoBase64
	
	// General: successfully uploaded file
	c.HTML(http.StatusOK, "index.html", gin.H{
		"err":  "Successfully Uploaded File",
		"tempVideoFile": dstPath,
	})
}

// FindAllVideo GET Mongo
func FindAllVideo(c *gin.Context) {
	videos := FindAllVideoRepo()
	c.JSON(http.StatusOK, videos)
}

func FindAllVideoRepo() []model.Video {
	var videos []model.Video
	database.QmgoConnection.Find(context.TODO(), bson.M{}).All(&videos)
	for _, video := range videos {
		fmt.Printf("%+v",video)
	}
	return videos
}

// CreateVideo POST Mongo
func CreateVideo(c *gin.Context) {
	var video model.Video
	CreateVideoRepo(video)
	c.JSON(http.StatusOK, video)
}
func CreateVideoRepo(video model.Video){
	videoInfo := []model.Video{
		{
			ID:primitive.NewObjectID(),
			Title: "Test1",
			VideoData:[]byte("123456"),
		},
		{
			ID:primitive.NewObjectID(),
			Title: "Test2",
			VideoData:[]byte("123456"),
		},
		{
			ID:primitive.NewObjectID(),
			Title: "Test3",
			VideoData:[]byte("123456"),
		},
	}
	result, err := database.QmgoConnection.InsertMany(context.TODO(), videoInfo)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n",result)
	return
}

// Create User
// func (h *UserController) CreateUser() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		requestData := new(model.Student)
// 		if err := c.ShouldBindJSON(&requestData); err != nil {
// 			fmt.Println("Error:" + err.Error())
// 			c.JSON(http.StatusNotAcceptable, responses.Status(responses.ParameterErr, nil))
// 			return
// 		}
// 		student_id, status := service.NewUserService().CreateUser(requestData)
// 		if status != responses.Success {
// 			c.JSON(http.StatusNotFound, responses.Status(responses.Error, nil))
// 			return
// 		}
// 		c.JSON(http.StatusOK, responses.Status(responses.Success, student_id))
// 	}
// }