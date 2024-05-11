package main

import (
	"MeetingVideoHelper/database"
	"MeetingVideoHelper/routes"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// connect database
	database.MongoDBinit()

	mainServer := gin.New()
	// 定義router呼叫格式與跨域限制
	mainServer.Use(cors.New(cors.Config{
		// 只允许来自 "http://localhost:8000" 的请求访问该服务器。
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		// AddAllowHeaders 允許添加自定義標頭
		AllowHeaders: []string{"Origin"},
		// 允許添加自定義公開標頭
		ExposeHeaders:    []string{"Content-Type", "application/javascript"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 連接Router
	routes.ApiRoutes(mainServer)

	if err := mainServer.Run(":9528"); err != nil {
		log.Fatal("HTTP service failed: ", err)
	}
}
