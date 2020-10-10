package dto

import (
	"errors"
	"gateway/extend/code"
	"gateway/extend/utils"
	"github.com/gin-gonic/gin"
)

type AdminRegisterInput struct {
	Mobile   string `json:"mobile" binding:"required"`
	Password string `json:"password" binding:"required"`
	Salt     string `json:"salt" binding:"required"`
}

func (params *AdminRegisterInput) BindValidParam(c *gin.Context) error {
	err := c.ShouldBindJSON(params)
	if err != nil {
		utils.ResponseFormat(c, code.ServiceInsideError, "解析参数出错")
		return err
	}
	// 校验手机号
	if len(params.Mobile) != 11 {
		utils.ResponseFormat(c, code.MobileLengthError, nil)
		return errors.New("手机号错误")
	}
	// 校验密码
	if len(params.Password) < 6 || len(params.Password) > 20 {
		utils.ResponseFormat(c, code.PasswordLengthError, nil)
		return errors.New("密码错误")
	}
	return nil
}
