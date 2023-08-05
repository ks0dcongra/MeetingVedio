package controller

import (
	"fmt"
	"context"
	"net/http"
	"MeetingVideoHelper/app/model"
	"MeetingVideoHelper/database"
	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"

)

type UserController struct {
	
}

func NewUserController() *UserController {
	return &UserController{
		
	}
}

// GET
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

// POST
func CreateVideo(c *gin.Context) {
	var video model.Video
	// c.BindJSON(&video)
	CreateVideoRepo(video)
	c.JSON(http.StatusOK, video)
}
func CreateVideoRepo(video model.Video){
	videoInfo := []model.Video{
		{
			ID:1,
			Title: "Test1",
			VideoData:[]byte("123456"),
		},
		{
			ID:2,
			Title: "Test2",
			VideoData:[]byte("123456"),
		},
		{
			ID:3,
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