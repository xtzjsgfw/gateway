package router

import (
	v1 "gateway/controller/v1"
	"gateway/middleware"
	"github.com/gin-gonic/gin"
)

func Init() {
	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	Group := engine.Group("api")
	{
		registerController := v1.RegisterController{}
		loginController := v1.LoginController{}
		Group.POST("/auth/login", loginController.AdminLogin)
		Group.POST("/auth/register", registerController.AdminRegister)

		Group.Use(middleware.JWTAuth())
		{
			// 服务相关
			serviceController := v1.ServiceController{}
			Group.GET("/service/list", serviceController.ServiceList)
			Group.GET("/service/delete", serviceController.ServiceDelete)
			Group.GET("/service/detail", serviceController.ServiceDetail)
			// 流量统计
			Group.GET("/service/stat", serviceController.ServiceStat)
			// 添加新增http服务
			Group.POST("/service/add_http", serviceController.ServiceAddHTTP)
			Group.POST("/service/update_http", serviceController.ServiceUpdateHTTP)
			// 添加新增grpc服务
			Group.POST("/service/add_grpc", serviceController.ServiceAddGRPC)
			Group.POST("/service/update_grpc", serviceController.ServiceUpdateGRPC)
			// 添加新增tcp服务
			Group.POST("/service/add_tcp", serviceController.ServiceAddTCP)
			Group.POST("/service/update_tcp", serviceController.ServiceUpdateTCP)

			// 租户相关
			appController := v1.APPController{}
			Group.GET("/app/list", appController.APPList)
			Group.GET("/app/delete", appController.APPDelete)
			Group.GET("/app/detail", appController.APPDetail)
			Group.GET("/app/stat", appController.APPStat)
			Group.POST("/app/add", appController.APPAdd)
			Group.POST("/app/update", appController.APPUpdate)

			// Dashboard
			dashboardController := v1.DashboardController{}
			Group.GET("/dashboard/panel_group_data", dashboardController.PanelGroupData)
			Group.GET("/dashboard/flow_stat", dashboardController.FlowStat)
			Group.GET("/dashboard/service_stat", dashboardController.ServiceStat)
		}
	}

	engine.Run()
}
