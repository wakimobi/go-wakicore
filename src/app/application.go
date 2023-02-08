package app

import (
	"github.com/gin-gonic/gin"
)

func StartApplication() *gin.Engine {
	return mapUrls()
}
