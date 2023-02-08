package app

import (
	"github.com/gin-gonic/gin"
	"github.com/wakimobi/go-wakicore/src/controllers"
)

func mapUrls() *gin.Engine {
	router := gin.Default()

	router.SetTrustedProxies([]string{"127.0.0.1"})

	v1 := router.Group("v1")
	v1.GET("/mo", controllers.MOHandler)
	v1.GET("/dr", controllers.DRHandler)

	v1.POST("/insert", controllers.InsertHandler)

	return router

}
