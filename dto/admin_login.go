package dto

import (
	"gateway/extend/code"
	"gateway/extend/utils"
	"github.com/gin-gonic/gin"
)

type AdminLoginInput struct {
	Mobile   string `json:"mobile" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (params *AdminLoginInput) BindValidParam(c *gin.Context) error {
	err := c.ShouldBindJSON(params)
	if err != nil {
		utils.ResponseFormat(c, code.RequestParamError, nil)
		return err
	}
	return nil
}
