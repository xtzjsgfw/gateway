package v1

import (
	"gateway/dto"
	"gateway/extend/code"
	"gateway/extend/utils"
	"gateway/service"
	"github.com/gin-gonic/gin"
)

type RegisterController struct{}

// AdminRegister godoc
// @Summary 管理员注册
// @Description 管理员注册
// @Accept json
// @Produce json
// @Tags 管理员接口
// @ID /admin_register/AdminRegister
// @Param body body auth.AdminRegisterInput{} tr ue "账号注册请求参数"
// @Success 200 {string} json "{"status":200, "code": 2000001, msg:"请求处理成功"}"
// @Failure 400 {string} json "{"status":400, "code": 4000001, msg:"请求参数有误"}"
// @Failure 500 {string} json "{"status":500, "code": 5000001, msg:"服务器内部错误"}"
func (rc *RegisterController) AdminRegister(c *gin.Context) {
	// 解析并验证参数
	adminRegisterInput := &dto.AdminRegisterInput{}
	err := adminRegisterInput.BindValidParam(c)
	if err != nil {
		return
	}

	// 查询是否有对应用户
	registerService := &service.AdminService{
		Mobile: adminRegisterInput.Mobile,
		Password: adminRegisterInput.Password,
		Username: adminRegisterInput.Mobile,
	}
	adminInfo, err := registerService.QueryByMobile(registerService.Mobile)
	if err != nil {
		utils.ResponseFormat(c, code.ServiceInsideError, nil)
		return
	}
	if adminInfo != nil {
		utils.ResponseFormat(c, code.UserIsExistError, nil)
		return
	}

	// 存储，新建用户
	userID, err := registerService.StoreUser(registerService.Mobile, registerService.Password)
	if err != nil {
		utils.ResponseFormat(c, code.ServiceInsideError, nil)
		return
	}
	if userID > 0 {
		utils.ResponseFormat(c, code.Success, gin.H{
			"ID": userID,
		})
	}

}
