package http_proxy

import (
	"gateway/extend/code"
	"gateway/extend/utils"
	"gateway/middleware/http_proxy"
	"github.com/gin-gonic/gin"
)

func HttpInit(middlewares ...gin.HandlerFunc) *gin.Engine {
	route := gin.New()

	route.Use(middlewares...)
	route.Use(http_proxy.HTTPAccessModeMiddleware())
	route.GET("/ping", func(c *gin.Context) {
		utils.ResponseFormat(c, code.PongCode, nil)
	})

	return route
}

func HttpsInit(middlewares ...gin.HandlerFunc) *gin.Engine {
	route := gin.New()

	route.Use(middlewares...)
	route.GET("/ping", func(c *gin.Context) {
		utils.ResponseFormat(c, code.PongCode, "https")
	})

	return route
}
