package routes

import (
	"MeetingVideoHelper/app/controller"

	"github.com/gin-gonic/gin"
)

func ApiRoutes(router *gin.Engine) {
	// Set HTML template
	router.Static("/static", "./static")
	router.LoadHTMLGlob("view/*")
	
	// Set Group
	// userApi := router.Group("user/api")
	// userApi.POST("create", controller.NewUserController().CreateUser())

	// rander index.html
	router.GET("/index", func(c *gin.Context) {
		c.HTML(200, "index.html", map[string]string{"title":"home"})
	})
	// upload file
	router.POST("/upload", controller.UploadFile)
	// create video
	router.POST("/createVideos", controller.CreateVideo)
	// find all videos
	router.GET("/videos", controller.FindAllVideo)
}
