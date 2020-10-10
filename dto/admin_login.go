package dto

import (
	"errors"
	"gateway/extend/code"
	"gateway/extend/utils"
	"github.com/gin-gonic/gin"
)

type AdminLoginInput struct {
	Mobile   string `json:"mobile" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (params *AdminLoginInput) Valid(c *gin.Context) error {
	err := c.ShouldBindJSON(params)
	if err != nil {
		utils.ResponseFormat(c, code.RequestParamError, nil)
		return err
	}
	// 验证用户名
	if len(params.Mobile) < 0 || len(params.Mobile) > 20 {
		utils.ResponseFormat(c, code.UsernameLengthError, nil)
		return errors.New("手机号格式错误")
	}

	// 验证密码
	if len(params.Password) < 6 || len(params.Password) > 20 {
		utils.ResponseFormat(c, code.PasswordLengthError, nil)
		return errors.New("密码格式错误")
	}
	return nil
}
