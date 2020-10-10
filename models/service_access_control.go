package models

import "github.com/jinzhu/gorm"

// 网关权限控制表
type GatewayServiceAccessControl struct {
	gorm.Model
	ServiceID         uint   `gorm:"type:bigint(20); not null; default: 0; comment: '服务id'"`
	OpenAuth          int    `gorm:"type:bigint(4); not null; default: 0; comment: '是否开启权限 0=不开启 1=开启'"`
	BlackList         string `gorm:"type:varchar(1000); not null; default: ''; comment: '黑名单ip'"`
	WhiteList         string `gorm:"type:varchar(1000); not null; default: ''; comment: '白名单ip'"`
	WhiteHostName     string `gorm:"type:varchar(1000); not null; default: ''; comment: '白名单主机'"`
	ClientIpFlowLimit int    `gorm:"type:int(11); not null; default: 0; comment: '客户端ip限流'"`
	ServiceFlowLimit  int    `gorm:"type:int(20); not null; default: 0; comment: '服务端限流'"`
}

func (g *GatewayServiceAccessControl) TableName() string {
	return "gateway_service_access_control"
}

func (g *GatewayServiceAccessControl) Find(search *GatewayServiceAccessControl) (*GatewayServiceAccessControl, error) {
	model := &GatewayServiceAccessControl{}
	err := DB.Where(search).Find(model).Error
	return model, err
}

func (g *GatewayServiceAccessControl) Save(DB *gorm.DB) error {
	return DB.Create(g).Error
}
