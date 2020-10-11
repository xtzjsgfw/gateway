package models

import (
	"gateway/dto"
	"github.com/jinzhu/gorm"
)

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

func (g *GatewayApp) TableName() string {
	return "gateway_app"
}

func (g *GatewayApp) APPList(db *gorm.DB, params *dto.APPListInput) ([]GatewayApp, int, error) {
	var list []GatewayApp
	var total int

	pageNum := params.PageNum
	pageSize := params.PageSize
	offset := (pageNum - 1) * pageSize
	query := db.Table(g.TableName()).Select("*").Where("is_delete = ?", 0)
	if params.Info != "" {
		query = query.Where("Name LIKE ?", "%"+params.Info+"%")
	}
	if err = query.Limit(pageSize).Offset(offset).Find(&list).Error; err != gorm.ErrRecordNotFound && err != nil {
		return nil, 0, err
	}

	query.Limit(pageSize).Offset(offset).Count(&total)
	return list, total, nil
}

func (g *GatewayApp) Find(db *gorm.DB, search *GatewayApp) (*GatewayApp, error) {
	appModel := &GatewayApp{}
	err := db.Where(search).Find(appModel).Error
	if err != gorm.ErrRecordNotFound && err != nil {
		return nil, err
	}
	if appModel.ID > 0 {
		return appModel, nil
	}
	return nil, nil
}

func (g *GatewayApp) Save(db *gorm.DB) error {
	return db.Save(g).Error
}