package models

import (
	"github.com/jinzhu/gorm"
)

// 网关路由匹配表 GRPC
type GatewayServiceGrpcRule struct {
	gorm.Model
	ServiceID uint `gorm:"type: bigint(20); not null; default: 0; comment: '服务id'"`
	Port      int `gorm:"type: int(5); not null; default: 0; comment: '端口'"`
	HeaderTransfor string `gorm:"type: varchar(5000); not null; default: '';
	comment: 'header转换支持增加(add)、删除(del)、修改(edit) 格式: add headname headvalue 多个逗号间隔'"`
}

// 网关路由匹配表 HTTP
type GatewayServiceHttpRule struct {
	gorm.Model
	ServiceID     uint    `gorm:"type: bigint(20); not null; default: 0; comment: '服务id'"`
	RuleType      int    `gorm:"type: tinyint(4); not null; default: 0; comment: '匹配类型 0=url前缀url_prefix 1=域名domain'"`
	Rule          string `gorm:"type: varchar(255); not null; default: ''; comment: 'type=domain表示域名，type=url_prefix时表示url前缀'"`
	NeedHttps     int    `gorm:"type: tinyint(4); not null; default: 0; comment: '支持https 1=支持'"`
	NeedStripUri  int    `gorm:"type: tinyint(4); not null; default: 0; comment: '启用strip_uri 1=启用'"`
	NeedWebsocket int    `gorm:"type: tinyint(4); not null; default: 0; comment: '是否支持websocket 1=支持'"`
	UrlRewrite string `gorm:"type: varchar(5000); not null; default: '';
	comment: 'url重写功能 格式：^/gatekeeper/test_service(.*) $1 多个逗号间隔'"`
	HeaderTransfor string `gorm:"type: varchar(5000); not null; default: '';
	comment: 'header转换支持增加(add)、删除(del)、修改(edit) 格式: add headname headvalue 多个逗号间隔'"`
}

//	网关路由匹配表
type GatewayServiceTcpRule struct {
	gorm.Model
	ServiceID uint `gorm:"type: bigint(20); not null; comment: '服务id'"`
	Port      int `gorm:"type: int(5); not null; default: 0; comment: '端口号'"`
}


func (g *GatewayServiceHttpRule) TableName() string {
	return "gateway_service_http_rule"
}

func (g *GatewayServiceGrpcRule) TableName() string {
	return "gateway_service_grpc_rule"
}

func (g *GatewayServiceTcpRule) TableName() string {
	return "gateway_service_tcp_rule"
}


func (g *GatewayServiceHttpRule) Find(db *gorm.DB, search *GatewayServiceHttpRule) (*GatewayServiceHttpRule, error) {
	model := &GatewayServiceHttpRule{}
	err := db.Where(search).Find(model).Error
	return model, err
}

func (g *GatewayServiceHttpRule) Save(db *gorm.DB) error {
	return db.Create(g).Error
}

func (g *GatewayServiceTcpRule) Find(db *gorm.DB,search *GatewayServiceTcpRule) (*GatewayServiceTcpRule, error) {
	model := &GatewayServiceTcpRule{}
	err := db.Where(search).Find(model).Error
	return model, err
}

func (g *GatewayServiceTcpRule) Save(db *gorm.DB) error {
	return db.Create(g).Error
}

func (g *GatewayServiceGrpcRule) Find(db *gorm.DB, search *GatewayServiceGrpcRule) (*GatewayServiceGrpcRule, error) {
	model := &GatewayServiceGrpcRule{}
	err := db.Where(search).Find(model).Error
	return model, err
}

func (g *GatewayServiceGrpcRule) Save(db *gorm.DB) error {
	return db.Create(g).Error
}