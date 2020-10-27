package models

import (
	"github.com/jinzhu/gorm"
	"strings"
)

// 网关负载表
type GatewayServiceLoadBalance struct {
	gorm.Model
	ServiceID              uint   `gorm:"type: bigint(20); not null; default: 0; comment: '服务id'"`
	CheckMethod            int    `gorm:"type: tinyint(20); not null; default: 0; comment: '检查方法 0=tcpchk,检测端口是否握手成功'"`
	CheckTimeout           int    `gorm:"type: int(10); not null; default: 0; comment: 'check超时时间,单位s'"`
	CheckInterval          int    `gorm:"type: int(11); not null; default: 0; comment: '检查间隔, 单位s'"`
	RoundType              int    `gorm:"type: tinyint(4); not null; default: 2; comment: '轮询方式 0=random 1=round-robin 2=weight_round-robin 3=ip_hash'"`
	IpList                 string `gorm:"type: varchar(2000); not null; default: ''; comment: 'ip列表'"`
	WeightList             string `gorm:"type: varchar(2000); not null; default: ''; comment: '权重列表'"`
	ForbidList             string `gorm:"type: varchar(2000); not null; default: ''; comment: '禁用ip列表'"`
	UpstreamConnectTimeout int    `gorm:"type: int(11); not null; default: 0; comment: '建立连接超时, 单位s'"`
	UpstreamHeaderTimeout  int    `gorm:"type: int(11); not null; default: 0; comment: '获取header超时, 单位s'"`
	UpstreamIdleTimeout    int    `gorm:"type: int(10); not null; default: 0; comment: '链接最大空闲时间, 单位s'"`
	UpstreamMaxIdle        int    `gorm:"type: int(11); not null; default: 0; comment: '最大空闲链接数'"`
}

func (g *GatewayServiceLoadBalance) TableName() string {
	return "gateway_service_load_balance"
}

func (g *GatewayServiceLoadBalance) Find(search *GatewayServiceLoadBalance) (*GatewayServiceLoadBalance, error) {
	model := &GatewayServiceLoadBalance{}
	err := DB.Where(search).Find(model).Error
	return model, err
}

func (g *GatewayServiceLoadBalance) GetIPListByModel() []string {
	return strings.Split(g.IpList, ",")
}

func (g *GatewayServiceLoadBalance) GetWeightListByModel() []string {
	return strings.Split(g.WeightList, ",")
}

func (g *GatewayServiceLoadBalance) Save(DB *gorm.DB) error {
	return DB.Create(g).Error
}
