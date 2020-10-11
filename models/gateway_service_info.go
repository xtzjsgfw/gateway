package models

import (
	"gateway/dto"
	"github.com/jinzhu/gorm"
)

//	网关基本信息表
type GatewayServiceInfo struct {
	gorm.Model
	LoadType    int    `gorm:"type:tinyint(4); not null; default: 0; comment: '负载类型 0=http 1=tcp 2=grpc'"`
	ServiceName string `gorm:"type: varchar(255); not null; default: ''; comment: '服务名称 6-128 数字字母下划线'"`
	ServiceDesc string `gorm:"type: varchar(255); not null; default: ''; comment: '服务描述'"`
	IsDelete    int    `gorm:"type:tinyint(4); not null; default: 0; comment: '是否删除'"`
}

func (g *GatewayServiceInfo) TableName() string {
	return "gateway_service_info"
}

func (g *GatewayServiceInfo) PageList(db *gorm.DB, params *dto.ServiceListInput) ([]GatewayServiceInfo, int64, error) {
	total := int64(0)
	list := []GatewayServiceInfo{}
	offset := (params.PageNum - 1) * params.PageSize
	query := db.Table(g.TableName()).Where("is_delete = ?", 0)

	if params.Info != "" {
		query = query.Where("(service_name like ? or service_desc like ?)",
			"%"+params.Info+"%", "%"+params.Info+"%")
	}
	if err = query.Limit(params.PageSize).Offset(offset).Find(&list).Error; err != gorm.ErrRecordNotFound && err != nil {
		return nil, 0, err
	}

	query.Limit(params.PageSize).Offset(offset).Count(&total)
	return list, total, nil
}

func (g *GatewayServiceInfo) ServiceDetail(db *gorm.DB, search *GatewayServiceInfo) (*ServiceDetail, error) {
	id := search.ID
	httpRule := &GatewayServiceHttpRule{ServiceID: id}
	httpRule, err = httpRule.Find(db, httpRule)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	tcpRule := &GatewayServiceTcpRule{ServiceID: id}
	tcpRule, err = tcpRule.Find(db, tcpRule)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	grpcRule := &GatewayServiceGrpcRule{ServiceID: id}
	grpcRule, err = grpcRule.Find(db, grpcRule)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	accessControl := &GatewayServiceAccessControl{ServiceID: id}
	accessControl, err = accessControl.Find(accessControl)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	loadBalance := &GatewayServiceLoadBalance{ServiceID: id}
	loadBalance, err = loadBalance.Find(loadBalance)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	detail := &ServiceDetail{
		Info:          search,
		HTTPRule:      httpRule,
		TCPRule:       tcpRule,
		GRPCRule:      grpcRule,
		LoadBalance:   loadBalance,
		AccessControl: accessControl,
	}
	return detail, nil
}

func (g *GatewayServiceInfo) Find(db *gorm.DB, search *GatewayServiceInfo) (*GatewayServiceInfo, error) {
	model := &GatewayServiceInfo{}
	err := db.Where(search).Find(&model).Error
	if err != gorm.ErrRecordNotFound && err != nil {
		return nil, err
	}
	if model.ID > 0 {
		return model, nil
	}
	// 什么都没找到
	return nil, nil
}

func (g *GatewayServiceInfo) Save(DB *gorm.DB) error {
	return DB.Save(g).Error
}