package routes

import (
	"MeetingVideoHelper/app/controller"
	// "fmt"
	// "log"
	// "os"

	// "MeetingVideoHelper/app/model"
	// "MeetingVideoHelper/database"
	// "context"
	// "net/http"

	"github.com/gin-gonic/gin"
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/bson/primitive"
)

func ApiRoutes(router *gin.Engine) {
	// Set HTML template
	router.Static("/static", "./static")
	router.LoadHTMLGlob("view/*")

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

	// router.POST("/downloads", func(c *gin.Context) {
	// 	videoID := c.PostForm("videoID")
	// 	log.Println("videoID", videoID[10:34])
	// 	videoDownload, err := primitive.ObjectIDFromHex(videoID[10:34])
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	// videoDownload, _ := primitive.ObjectIDFromHex(videoID)
	// 	log.Println("videoDownload", videoDownload)

	// 	var videoResult model.Video

	// 	// MongoDB: find one specific video from mongodb
	// 	err = database.QmgoConnection.Find(context.Background(), bson.M{"vid": videoDownload}).One(&videoResult)
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, "MongoDB: Finding the FlippedVideo fail")
	// 		return
	// 	}

	// 	tempVideo, err := os.CreateTemp("./", "tempVideo*.mp4")
	// 	if err != nil {
	// 		fmt.Println("Error creating temporary file:", err)
	// 		return
	// 	}

	// 	tempVideo.Write(videoResult.VideoData)
	// 	tempVideo.Close()
	// 	filePath := tempVideo.Name()
	// 	log.Println("filePathfilePathfilePathfilePathfilePathfilePath",filePath)
	// 	// 設定標頭讓瀏覽器知道這是個下載檔案的請求
	// 	c.Header("Content-Disposition", "attachment; filename="+filePath)
	// 	c.Header("Content-Type", "video/mp4")
	// 	// 透過 Gin 的 `File` 方法來提供靜態檔案服務
	// 	c.File(filePath)

	// 	e := os.Remove(filePath)
	// 	if e != nil {
	// 		log.Fatal(e)
	// 	}

	// 	c.HTML(200, "index.html", map[string]string{"title": "home"})
	// })

	router.POST("/downloads", controller.NewVideoController().DownloadVideo)

	// upload file
	router.POST("/upload", controller.NewVideoController().UploadVideo)
}
