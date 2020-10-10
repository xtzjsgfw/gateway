package v1

import (
	"gateway/dto"
	"gateway/extend/code"
	"gateway/extend/utils"
	"gateway/service"
	"github.com/gin-gonic/gin"
)

type LoginController struct{}

func (lc *LoginController) AdminLogin(c *gin.Context) {
	adminLoginInput := &dto.AdminLoginInput{}
	err := adminLoginInput.Valid(c)
	if err != nil {
		return
	}

	// 验证是否存在用户
	adminService := &service.AdminService{
		Mobile:   adminLoginInput.Mobile,
		Password: adminLoginInput.Password,
	}
	adminInfo, err := adminService.QueryByMobile(adminService.Mobile)
	if err != nil {
		utils.ResponseFormat(c, code.ServiceInsideError, nil)
		return
	}
	if adminInfo == nil {
		utils.ResponseFormat(c, code.UserIsNotExistError, nil)
		return
	}

	token := "Token"
	// 验证密码是否正确
	if utils.MakeSha1(adminService.Mobile+adminService.Password) == adminInfo.Password {
		utils.ResponseFormat(c, code.Success, gin.H{
			"token": token,
		})
		return
	}else {
		utils.ResponseFormat(c, code.UserOrPassError, nil)
		return
	}
}
