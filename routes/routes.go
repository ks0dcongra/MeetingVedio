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

	router.GET("/index", func(c *gin.Context) {
		c.HTML(200, "index.html", map[string]string{"title":"home"})
	})

	router.POST("/createVideos", controller.CreateVideo)
	router.GET("/videos", controller.FindAllVideo)
}
