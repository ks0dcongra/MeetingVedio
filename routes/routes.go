package routes

import "github.com/gin-gonic/gin"

func ApiRoutes(router *gin.Engine) {
	// 设置静态文件目录（可选）
	router.Static("/static", "./static")
	// 设置HTML模板文件目录（可选）
	router.LoadHTMLGlob("view/*")
	// userApi := router.Group("user/api")
	// create
	// userApi.POST("create", controller.NewUserController().CreateUser())

	// 獲得CSRF Token 與 攻擊CSRF之網頁
	router.GET("/index", func(c *gin.Context) {
		c.HTML(200, "index.html", map[string]string{"title":"home"})
	})
}
