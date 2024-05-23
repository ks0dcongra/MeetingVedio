package routes

import (
	"MeetingVideoHelper/app/controller"
	"MeetingVideoHelper/app/middleware"

	"net/http"

	"github.com/gin-gonic/gin"
)

func ApiRoutes(router *gin.Engine) {
	// Set HTML template
	router.Static("/static", "./static")
	router.LoadHTMLGlob("view/*")

	downloadRoutes := router.Group("/download")
	{
		// 為這個路由組應用 RateLimiter 中間件
		downloadRoutes.Use(middleware.IPRateLimiter)
		downloadRoutes.GET("/:filename", func(c *gin.Context) {
			fileName := c.Param("filename")
			filePath := "./static/videos/" + fileName
			c.Header("Content-Disposition", "attachment; filename="+fileName)
			c.Header("Content-Type", "video/mp4")
			c.File(filePath)
		})
		downloadRoutes.POST("", controller.NewVideoController().DownloadVideo)
	}

	uploadRoutes := router.Group("/upload")
	{
		// 為這個路由組應用 RateLimiter 中間件
		uploadRoutes.Use(middleware.IPRateLimiter)
		uploadRoutes.POST("", controller.NewVideoController().UploadVideo)
	}

	// render index.html
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", map[string]string{"title": "home"})
	})
}
