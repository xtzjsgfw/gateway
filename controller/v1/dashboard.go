package v1

import (
	"gateway/dto"
	"gateway/extend/code"
	"gateway/extend/utils"
	"gateway/models"
	"gateway/public"
	"github.com/gin-gonic/gin"
	"time"
)

type DashboardController struct{}

// 指标统计
func (d *DashboardController) PanelGroupData(c *gin.Context) {
	db := models.GetDB()

	// 服务数
	serviceInfo := &models.GatewayServiceInfo{}
	_, serviceNum, err := serviceInfo.PageList(db, &dto.ServiceListInput{PageNum: 1, PageSize: 1})
	if err != nil {
		utils.ResponseFormat(c, code.ServiceInsideError, nil)
		return
	}

	// 租户数
	app := &models.GatewayApp{}
	_, appNum, err := app.APPList(db, &dto.APPListInput{PageNum: 1, PageSize: 1})
	if err != nil {
		utils.ResponseFormat(c, code.ServiceInsideError, nil)
		return
	}

	out := &dto.PanelGroupDataOutput{
		ServiceNum:      serviceNum,
		APPNum:          appNum,
		CurrentQPS:      0,
		TodayRequestNum: 0,
	}
	utils.ResponseFormat(c, code.Success, out)
}

// 流量统计
func (d *DashboardController) FlowStat(c *gin.Context) {
	todayList := []int64{}
	for i := 0; i < time.Now().Hour(); i++ {
		todayList = append(todayList, 0)
	}

	yesterdayList := []int64{}
	for i := 0; i < 23; i++ {
		yesterdayList = append(yesterdayList, 1)
	}

	out := &dto.FlowStatOutput{
		Today:     todayList,
		Yesterday: yesterdayList,
	}
	utils.ResponseFormat(c, code.Success, out)
}

// 服务统计
func (d *DashboardController) ServiceStat(c *gin.Context) {
	db := models.GetDB()
	// 服务类型占比
	serviceInfo := &models.GatewayServiceInfo{}
	list, err := serviceInfo.GroupByLoadType(db)
	if err != nil {
		utils.ResponseFormat(c, code.ServiceInsideError, nil)
		return
	}

	// 从list中取出Legend
	legend := []string{}
	for index, item := range list {
		name, ok := public.LoadTypeMap[item.LoadType]
		if !ok {
			utils.ResponseFormat(c, code.ServiceInsideError, nil)
			return
		}
		list[index].Name = name
		legend = append(legend, name)
	}

	out := &dto.DashServiceStatOutput{
		Legend: legend,
		Data:   list,
	}
	utils.ResponseFormat(c, code.Success, out)
}
