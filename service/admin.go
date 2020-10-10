package service

import (
	"gateway/extend/utils"
	"gateway/models"
)

type AdminService struct{
	AdminID uint
	Mobile string
	Salt string
	Password string
	Username string
}


func (rs *AdminService) QueryByMobile(mobile string) (*models.GatewayAdmin, error) {
	gatewayAdmin := &models.GatewayAdmin{}
	condition := map[string]interface{}{
		"mobile": mobile,
	}
	adminInfo, err := gatewayAdmin.FindOne(condition)
	return adminInfo, err
}

func (rs *AdminService) StoreUser(mobile, password string) (userID uint, err error) {
	gatewayAdmin := &models.GatewayAdmin{
		Mobile:   mobile,
		Password: password,
		Username: mobile,
	}
	gatewayAdmin.Password = utils.MakeSha1(gatewayAdmin.Mobile + gatewayAdmin.Password)
	userID, err = gatewayAdmin.Insert()
	return userID, err
}
