package models

import "github.com/jinzhu/gorm"

type GatewayApp struct {
	gorm.Model
	AppID    string `gorm:"type: varchar(255); not null; default: ''; comment: '租户id'"`
	Name     string `gorm:"type: varchar(255); not null; default: ''; comment: '租户名称'"`
	Secret   string `gorm:"type: varchar(255); not null; default: ''; comment: '密钥'"`
	WhiteIPs string `gorm:"type: varchar(1000); not null; default: ''; comment: 'ip白名单，支持前缀匹配'"`
	QPD      int    `gorm:"type: bigint(20); not null; default: 0; comment: '日请求量限制'"`
	QPS      int    `gorm:"type: bigint(20); not null; default: 0; comment: '秒请求量限制'"`
	IsDelete int    `gorm:"type: tinyint(4); not null; default: 0; comment: '是否删除'"`
}
