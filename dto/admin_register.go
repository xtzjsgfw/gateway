package dto

import (
	"gateway/extend/code"
	"gateway/extend/utils"
	"github.com/gin-gonic/gin"
)

type AdminRegisterInput struct {
	Mobile   string `json:"mobile" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (params *AdminRegisterInput) BindValidParam(c *gin.Context) error {
	err := c.ShouldBindJSON(params)
	if err != nil {
		utils.ResponseFormat(c, code.ServiceInsideError, "解析参数出错")
		return err
	}
	return nil
}
