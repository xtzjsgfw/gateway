package http_proxy

import (
	"gateway/extend/code"
	"gateway/extend/utils"
	"gateway/models"
	"github.com/gin-gonic/gin"
)

//匹配接入方式 基于请求信息
func HTTPAccessModeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		service, err := models.ServiceManagerHandler.HTTPAccessMode(c)
		if err != nil {
			utils.ResponseFormat(c, code.RulePixOrDomIsNotExistError, nil)
			c.Abort()
			return
		}
		c.Set("service", service)
		c.Next()
	}
}
