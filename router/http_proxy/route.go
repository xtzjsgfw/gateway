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

	route.GET("/ping", func(c *gin.Context) {
		utils.ResponseFormat(c, code.PongCode, nil)
	})

	route.Use(
		http_proxy.HTTPAccessModeMiddleware(),
		http_proxy.HTTPReverseProxyMiddleware(),
	)

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
