package middleware

import (
	"MeetingVideoHelper/app/contract"
	"MeetingVideoHelper/database"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func IPRateLimiter(c *gin.Context) {
	err := database.RedisConn.SetNX(context.TODO(), c.ClientIP(), 0, time.Hour).Err()

	database.RedisConn.Incr(context.TODO(), c.ClientIP())
	if err != nil {
		log.Fatal(err)
	}

	val, err := database.RedisConn.Get(context.TODO(), c.ClientIP()).Int()
	if err != nil {
		log.Fatal(err)
	}

	if val > 10 {
		c.Abort()
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{
			"status": contract.ERROR_IP_Request_Limit,
			"msg":    contract.Message[contract.ERROR_IP_Request_Limit],
		})
		return
	} else {
		c.Next()
	}
}
