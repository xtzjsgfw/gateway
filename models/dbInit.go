package models

import (
	"fmt"
	"gateway/extend/conf"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

var DB *gorm.DB
var err error

func Init() {
	connectStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.DBConf.User,
		conf.DBConf.Password,
		conf.DBConf.Host,
		conf.DBConf.Port,
		conf.DBConf.DBName)
	DB, err = gorm.Open(conf.DBConf.DBType, connectStr)
	if err != nil {
		fmt.Println("连接数据库失败, err : ", err)
		time.Sleep(10 * time.Second)
		DB, err = gorm.Open(conf.DBConf.DBType, connectStr)
		if err != nil {
			panic(err.Error())
		}
	}
	if DB.Error != nil {
		fmt.Println("连接数据库失败--DB:", DB.Error)
	}
	DB.LogMode(conf.DBConf.Debug)
	DB.SingularTable(true)
	DB.DB().SetMaxIdleConns(conf.DBConf.MaxIdleConns)
	DB.DB().SetMaxOpenConns(conf.DBConf.MaxOpenConns)
	DB.AutoMigrate(&GatewayAdmin{}, &GatewayApp{}, &GatewayServiceAccessControl{},
		&GatewayApp{}, &GatewayServiceGrpcRule{}, &GatewayServiceHttpRule{},
		&GatewayServiceTcpRule{}, &GatewayServiceInfo{},
		&GatewayServiceLoadBalance{}, &GatewayServiceAccessControl{})
}

func GetDB() *gorm.DB {
	return DB
}
