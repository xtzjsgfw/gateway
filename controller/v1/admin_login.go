package v1

import (
	"gateway/dto"
	"gateway/extend/code"
	"gateway/extend/conf"
	myJWT "gateway/extend/jwt"
	"gateway/extend/redis"
	"gateway/extend/utils"
	"gateway/service"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"time"
)

type LoginController struct{}

func (lc *LoginController) AdminLogin(c *gin.Context) {
	params := &dto.AdminLoginInput{}
	err := params.BindValidParam(c)
	if err != nil {
		return
	}

	// 验证是否存在用户
	adminService := &service.AdminService{
		Mobile: params.Mobile,
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

	if utils.MakeSha1(params.Mobile+params.Password) != adminInfo.Password {
		utils.ResponseFormat(c, code.UserOrPassError, nil)
		return
	}

	// 生成Token
	jwtInstance := myJWT.NewJWT()
	nowTime := time.Now()
	expireTime := time.Duration(conf.ServerConf.JwtExpire)
	claims := myJWT.CustomClaims{
		ID:     adminInfo.ID,
		Mobile: adminInfo.Mobile,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: nowTime.Add(expireTime * time.Hour).Unix(),
			Issuer:    "gateway",
		},
	}
	token, err := jwtInstance.CreateToken(claims)
	if err != nil {
		utils.ResponseFormat(c, code.ServiceInsideError, nil)
		return
	}

	// 将Token存入redis
	rdb := redis.GetRedis()
	rdb.Set("TOKEN:"+params.Mobile, token, time.Duration(conf.ServerConf.JwtExpire)*time.Hour)

	utils.ResponseFormat(c, code.Success, gin.H{
		"id":    adminInfo.ID,
		"token": token,
	})
}
