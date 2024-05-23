package main

import (
	"MeetingVideoHelper/database"
	"MeetingVideoHelper/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	database.MongoDBinit()
	database.RedisDBinit()
	
	mainServer := gin.New()

	routes.ApiRoutes(mainServer)

	if err := mainServer.Run(":9528"); err != nil {
		log.Fatal("HTTP service failed: ", err)
	}
}
