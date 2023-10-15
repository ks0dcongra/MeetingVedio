package routes

import (
	"MeetingVideoHelper/app/controller"
	"log"

	"MeetingVideoHelper/app/model"
	"MeetingVideoHelper/database"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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
		c.HTML(200, "index.html", map[string]string{"title": "home"})
	})
	// download file
	router.GET("/download/:filename", func(c *gin.Context) {
		// 取得路由參數中的檔案名稱
		fileName := c.Param("filename")
		// 指定下載的檔案路徑
		filePath := "./static/videos/" + fileName

		// 設定標頭讓瀏覽器知道這是個下載檔案的請求
		c.Header("Content-Disposition", "attachment; filename="+fileName)
		c.Header("Content-Type", "video/mp4")
		// 透過 Gin 的 `File` 方法來提供靜態檔案服務
		c.File(filePath)
	})

	router.POST("/downloads", func(c *gin.Context) {
		// 取得路由參數中的檔案名稱
		// fileName := c.Param("filename")
		// 指定下載的檔案路徑
		// filePath := "./static/videos/" + fileName
		videoID := c.PostForm("videoID")
		log.Println("videoID", videoID)
		var videoResult model.Video
		// MongoDB: find one specific video from mongodb
		err := database.QmgoConnection.Find(context.Background(), bson.M{"vid": videoID}).One(&videoResult)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "MongoDB: Finding the FlippedVideo fail")
			return
		}

		// 設定標頭讓瀏覽器知道這是個下載檔案的請求
		// c.Header("Content-Disposition", "attachment; filename="+fileName)
		c.Header("Content-Type", "video/mp4")
		// 透過 Gin 的 `File` 方法來提供靜態檔案服務
		// c.File(filePath)
	})

	// upload file
	router.POST("/upload", controller.UploadFile)
	// create video
	// router.POST("/createVideos", controller.CreateVideo)
	// find all videos
	// router.GET("/videos", controller.FindAllVideo)
}
