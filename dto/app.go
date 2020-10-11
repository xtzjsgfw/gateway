package dto

import (
	"gateway/extend/code"
	"gateway/extend/utils"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

type APPListInput struct {
	Info     string `json:"info"`
	PageNum  int    `json:"page_num" validate:"required"`
	PageSize int    `json:"page_size" validate:"required"`
}

type APPListOutput struct {
	Total int                 `json:"total" form:"total"`
	List  []APPListItemOutput `json:"list" form:"list"`
}

type APPListItemOutput struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	AppID     string    `json:"app_id" gorm:"column:app_id" description:"租户id	"`
	Name      string    `json:"name" gorm:"column:name" description:"租户名称	"`
	Secret    string    `json:"secret" gorm:"column:secret" description:"密钥"`
	WhiteIPS  string    `json:"white_ips" gorm:"column:white_ips" description:"ip白名单，支持前缀匹配		"`
	Qpd       int       `json:"qpd" gorm:"column:qpd" description:"日请求量限制"`
	Qps       int       `json:"qps" gorm:"column:qps" description:"每秒请求量限制"`
	RealQpd   int       `json:"real_qpd" description:"日请求量限制"`
	RealQps   int       `json:"real_qps" description:"每秒请求量限制"`
	UpdatedAt time.Time `json:"create_at" gorm:"column:create_at" description:"添加时间	"`
	CreatedAt time.Time `json:"update_at" gorm:"column:update_at" description:"更新时间"`
	IsDelete  int8      `json:"is_delete" gorm:"column:is_delete" description:"是否已删除；0：否；1：是"`
}

type APPDetailInput struct {
	ID int `json:"id" validate:"required"`
}

type APPDeleteInput struct {
	ID int `json:"id" validate:"required"`
}

type APPStatInput struct {
	ID int `json:"id" validate:"required"`
}

type APPStatOutput struct {
	Today     []int64 `json:"today" form:"today" comment:"今日统计" validate:"required"`
	Yesterday []int64 `json:"yesterday" form:"yesterday" comment:"昨日统计" validate:"required"`
}

type APPAddInput struct {
	AppID    string `json:"app_id" form:"app_id" comment:"租户id" validate:"required"`
	Name     string `json:"name" form:"name" comment:"租户名称" validate:"required"`
	Secret   string `json:"secret" form:"secret" comment:"密钥" validate:""`
	WhiteIPS string `json:"white_ips" form:"white_ips" comment:"ip白名单，支持前缀匹配"`
	Qpd      int    `json:"qpd" form:"qpd" comment:"日请求量限制" validate:""`
	Qps      int    `json:"qps" form:"qps" comment:"每秒请求量限制" validate:""`
}

type APPUpdateInput struct {
	ID int `json:"id" validate:"required"`
	APPAddInput
}

func (params *APPListInput) BindValidParam(c *gin.Context) error {
	params.Info = c.Query("info")
	pageNum, err := strconv.Atoi(c.Query("page_num"))
	if err != nil {
		utils.ResponseFormat(c, code.ServiceInsideError, nil)
		return err
	}
	params.PageNum = pageNum
	pageSize, err := strconv.Atoi(c.Query("page_size"))
	if err != nil {
		utils.ResponseFormat(c, code.ServiceInsideError, nil)
		return err
	}
	params.PageSize = pageSize
	return nil
}

func (params *APPDetailInput) BindValidParam(c *gin.Context) error {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		utils.ResponseFormat(c, code.ServiceInsideError, nil)
		return err
	}
	params.ID = id
	return nil
}

func (params *APPDeleteInput) BindValidParam(c *gin.Context) error {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		utils.ResponseFormat(c, code.ServiceInsideError, nil)
		return err
	}
	params.ID = id
	return nil
}

func (params *APPStatInput) BindValidParam(c *gin.Context) error {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		utils.ResponseFormat(c, code.ServiceInsideError, nil)
		return err
	}
	params.ID = id
	return nil
}

func (params *APPAddInput) BindValidParam(c *gin.Context) error {
	if err := c.ShouldBindJSON(params); err != nil {
		utils.ResponseFormat(c, code.ServiceInsideError, nil)
		return err
	}
	return nil
}

func (params *APPUpdateInput) BindValidParam(c *gin.Context) error {
	if err := c.ShouldBindJSON(params); err != nil {
		utils.ResponseFormat(c, code.ServiceInsideError, nil)
		return err
	}
	return nil
}
