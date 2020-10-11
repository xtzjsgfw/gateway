package v1

import (
	"fmt"
	"gateway/dto"
	"gateway/extend/code"
	"gateway/extend/utils"
	"gateway/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"time"
)

type APPController struct{}

func (ac *APPController) APPList(c *gin.Context) {
	params := &dto.APPListInput{}
	if err := params.BindValidParam(c); err != nil {
		return
	}

	db := models.GetDB()

	app := &models.GatewayApp{}
	list, total, err := app.APPList(db, params)
	if err != nil {
		utils.ResponseFormat(c, code.ServiceInsideError, nil)
		return
	}

	outputList := []dto.APPListItemOutput{}
	for _, item := range list {
		realQps := 0
		realQpd := 0
		outputList = append(outputList, dto.APPListItemOutput{
			ID:       item.ID,
			AppID:    item.AppID,
			Name:     item.Name,
			Secret:   item.Secret,
			WhiteIPS: item.WhiteIPs,
			Qpd:      item.QPD,
			Qps:      item.QPS,
			RealQpd:  realQpd,
			RealQps:  realQps,
		})
	}
	output := dto.APPListOutput{
		Total: total,
		List:  outputList,
	}
	utils.ResponseFormat(c, code.Success, output)
}

func (ac *APPController) APPDetail(c *gin.Context) {
	params := &dto.APPDetailInput{}
	if err := params.BindValidParam(c); err != nil {
		return
	}

	db := models.GetDB()

	// 没有判断没有找到的情况
	app := &models.GatewayApp{Model: gorm.Model{ID: uint(params.ID)}}
	app, err := app.Find(db, app)
	if err != nil {
		fmt.Println(err)
		utils.ResponseFormat(c, code.ServiceInsideError, nil)
		return
	}
	if app == nil {
		utils.ResponseFormat(c, code.APPIsNotExistError, nil)
		return
	}

	utils.ResponseFormat(c, code.Success, app)
}

func (ac *APPController) APPDelete(c *gin.Context) {
	params := &dto.APPDeleteInput{}
	if err := params.BindValidParam(c); err != nil {
		return
	}

	db := models.GetDB()

	app := &models.GatewayApp{Model: gorm.Model{ID: uint(params.ID)}}
	app, err := app.Find(db, app)
	if err != nil {
		fmt.Println(err)
		utils.ResponseFormat(c, code.ServiceInsideError, nil)
		return
	}
	if app == nil {
		utils.ResponseFormat(c, code.APPIsNotExistError, nil)
		return
	}
	app.IsDelete = 1

	// 保存
	err = app.Save(db)
	if err != nil {
		utils.ResponseFormat(c, code.DBSaveError, nil)
		return
	}
	utils.ResponseFormat(c, code.Success, gin.H{"id": app.ID})
}

func (ac *APPController) APPStat(c *gin.Context) {
	params := &dto.APPStatInput{}
	if err := params.BindValidParam(c); err != nil {
		return
	}
	todayList := []int64{}
	for i := 0; i < time.Now().Hour(); i++ {
		todayList = append(todayList, 20)
	}
	yesterdayList := []int64{}
	for i := 0; i < 23; i++ {
		yesterdayList = append(yesterdayList, 20)
	}

	output := dto.APPStatOutput{
		Today:     todayList,
		Yesterday: yesterdayList,
	}

	utils.ResponseFormat(c, code.Success, output)
}

func (ac *APPController) APPAdd(c *gin.Context) {
	params := &dto.APPAddInput{}
	if err := params.BindValidParam(c); err != nil {
		return
	}

	db := models.GetDB()

	// 查找是否存在租户
	search := &models.GatewayApp{AppID: params.AppID}
	search, err := search.Find(db, search)
	if err != nil {
		fmt.Println(err)
		utils.ResponseFormat(c, code.ServiceInsideError, nil)
		return
	}
	if search != nil {
		utils.ResponseFormat(c, code.APPIsExistError, nil)
		return
	}
	if params.Secret == "" {
		params.Secret = utils.MakeSha1(params.AppID)
	}

	// 新建租户
	app := &models.GatewayApp{
		AppID:    params.AppID,
		Name:     params.Name,
		Secret:   params.Secret,
		WhiteIPs: params.WhiteIPS,
		QPS:      params.Qps,
		QPD:      params.Qpd,
		IsDelete: 0,
	}
	if err := app.Save(db); err != nil {
		utils.ResponseFormat(c, code.DBSaveError, nil)
		return
	}

	utils.ResponseFormat(c, code.Success, gin.H{"id": app.ID})
}

func (ac *APPController) APPUpdate(c *gin.Context) {
	params := &dto.APPUpdateInput{}
	if err := params.BindValidParam(c); err != nil {
		return
	}

	db := models.GetDB()

	// 查找是否存在租户
	search := &models.GatewayApp{AppID: params.AppID}
	info, err := search.Find(db, search)
	if err != nil {
		fmt.Println(err)
		utils.ResponseFormat(c, code.ServiceInsideError, nil)
		return
	}
	if info == nil {
		utils.ResponseFormat(c, code.APPIsNotExistError, nil)
		return
	}

	// 更新操作
	if params.Secret == "" {
		params.Secret = utils.MakeSha1(params.AppID)
	}

	info.Name = params.Name
	info.Secret = params.Secret
	info.WhiteIPs = params.WhiteIPS
	info.QPS = params.Qps
	info.QPD = params.Qpd


	if err := info.Save(db); err != nil {
		utils.ResponseFormat(c, code.DBSaveError, nil)
		return
	}
	utils.ResponseFormat(c, code.Success, gin.H{"id": info.ID})
}
