package models

import (
	"github.com/jinzhu/gorm"
)

// 网关管理员表
type GatewayAdmin struct {
	gorm.Model
	Username string `gorm:"type:varchar(255); default: ''; not null; comment: '用户名'"`
	Mobile   string `gorm:"type:char(11); default: ''; not null; comment: '手机号'"`
	Salt     string `gorm:"type:varchar(50); not null; default: ''; comment: '盐'"`
	Password string `gorm:"type: varchar(255); not null; default: ''; comment: '密码'"`
	IsDelete int    `gorm:"type:tinyint(4); not null; default: 0; comment: '是否删除'"`
}

func (admin *GatewayAdmin) FindOne(condition map[string]interface{}) (*GatewayAdmin, error) {
	var adminInfo GatewayAdmin
	result := DB.Where(condition).First(&adminInfo)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return nil, result.Error
	}
	if adminInfo.ID > 0 {
		return &adminInfo, nil
	}
	return nil, nil
}

func (admin *GatewayAdmin) Insert() (UserID uint, err error) {
	result := DB.Create(&admin)
	UserID = admin.ID
	if result.Error != nil {
		err = result.Error
	}
	return
}
