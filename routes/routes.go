package routes

import (
	"MeetingVideoHelper/app/controller"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ApiRoutes(router *gin.Engine) {
	// Set HTML template
	router.Static("/static", "./static")
	router.LoadHTMLGlob("view/*")

	// rander index.html
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", map[string]string{"title": "home"})
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

	router.POST("/downloads", controller.NewVideoController().DownloadVideo)

	// upload file
	router.POST("/upload",controller.NewVideoController().UploadVideo)
}
