package dto

import (
	"fmt"
	"gateway/extend/code"
	"gateway/extend/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

var err error

type ServiceListInput struct {
	Info     string `json:"info"`
	PageNum  int    `json:"page_num" binding:"required"`
	PageSize int    `json:"page_size" binding:"required"`
}

type ServiceListItemOutPut struct {
	ID          uint   `json:"id" form:"id"`
	ServiceName string `json:"service_name" form:"service_name"`
	ServiceDesc string `json:"service_desc" form:"service_desc"`
	LoadType    int    `json:"load_type" form:"load_type"`
	ServiceAddr string `json:"service_addr" form:"service_addr"`
	QPS         int64  `json:"qps" form:"qps"`
	QPD         int64  `json:"qpd" form:"qpd"`
	TotalNode   int    `json:"total_node" form:"total_node"`
}

type ServiceListOutPut struct {
	Total int64                   `json:"total" form:"total"`
	List  []ServiceListItemOutPut `json:"list"`
}

type ServiceDeleteInput struct {
	ServiceID string `json:"service_id"`
}

type ServiceAddHTTPInput struct {
	ServiceName string `json:"service_name" form:"service_name" validate:"required,checkServiceName"` //服务名
	ServiceDesc string `json:"service_desc" form:"service_desc" validate:"required,max=255,min=1"`    //服务描述

	RuleType       int    `json:"rule_type" form:"rule_type" validate:"max=1,min=0"`                       //接入类型
	Rule           string `json:"rule" form:"rule" validate:"required,check_service_rule"`                 //域名或者前缀
	NeedHttps      int    `json:"need_https" form:"need_https" validate:"max=1,min=0"`                     //支持https
	NeedStripUri   int    `json:"need_strip_uri" form:"need_strip_uri" validate:"max=1,min=0"`             //启用strip_uri
	NeedWebsocket  int    `json:"need_websocket" form:"need_websocket" validate:"max=1,min=0"`             //是否支持websocket
	UrlRewrite     string `json:"url_rewrite" form:"url_rewrite" validate:"check_url_rewrite"`             //url重写功能
	HeaderTransfor string `json:"header_transfor" form:"header_transfor" validate:"check_header_transfor"` //header转换

	OpenAuth          int    `json:"open_auth" form:"open_auth"  validate:"max=1,min=0"`              //关键词
	BlackList         string `json:"black_list" form:"black_list" validate:""`                        //黑名单ip
	WhiteList         string `json:"white_list" form:"white_list" validate:""`                        //白名单ip
	ClientipFlowLimit int    `json:"clientip_flow_limit" form:"clientip_flow_limit" validate:"min=0"` //客户端ip限流
	ServiceFlowLimit  int    `json:"service_flow_limit" form:"service_flow_limit"validate:"min=0"`    //服务端限流

	RoundType              int    `json:"round_type" form:"round_type" validate:"max=3,min=0"`                       //轮询方式
	IpList                 string `json:"ip_list" form:"ip_list" validate:"required,check_iplist"`                   //ip列表
	WeightList             string `json:"weight_list" form:"weight_list" validate:"required,check_weightlist"`       //权重列表
	UpstreamConnectTimeout int    `json:"upstream_connect_timeout" form:"upstream_connect_timeout" validate:"min=0"` //建立连接超时, 单位s
	UpstreamHeaderTimeout  int    `json:"upstream_header_timeout" form:"upstream_header_timeout" validate:"min=0"`   //获取header超时, 单位s
	UpstreamIdleTimeout    int    `json:"upstream_idle_timeout" form:"upstream_idle_timeout"validate:"min=0"`        //链接最大空闲时间, 单位s
	UpstreamMaxIdle        int    `json:"upstream_max_idle" form:"upstream_max_idle" validate:"min=0"`               //最大空闲链接数
}

type ServiceAddGRPCInput struct {
	ServiceName string `json:"service_name" form:"service_name" validate:"required,checkServiceName"` //服务名
	ServiceDesc string `json:"service_desc" form:"service_desc" validate:"required,max=255,min=1"`    //服务描述

	Port           int    `json:"port" form:"port" comment:"端口，需要设置8001-8999范围内" validate:"required,min=8001,max=8999"`
	HeaderTransfor string `json:"header_transfor" form:"header_transfor" comment:"metadata转换" validate:""`

	OpenAuth          int    `json:"open_auth" form:"open_auth"  validate:"max=1,min=0"` //关键词
	BlackList         string `json:"black_list" form:"black_list" validate:""`           //黑名单ip
	WhiteList         string `json:"white_list" form:"white_list" validate:""`           //白名单ip
	WhiteHostName     string `json:"white_host_name" form:"white_host_name" comment:"白名单主机，以逗号间隔" validate:""`
	ClientipFlowLimit int    `json:"clientip_flow_limit" form:"clientip_flow_limit" validate:"min=0"` //客户端ip限流
	ServiceFlowLimit  int    `json:"service_flow_limit" form:"service_flow_limit"validate:"min=0"`    //服务端限流

	RoundType  int    `json:"round_type" form:"round_type" validate:"max=3,min=0"`                     //轮询方式
	IpList     string `json:"ip_list" form:"ip_list" validate:"required,check_iplist"`                 //ip列表
	WeightList string `json:"weight_list" form:"weight_list" validate:"required,check_weightlist"`     //权重列表
	ForbidList string `json:"forbid_list" form:"forbid_list" comment:"禁用IP列表" validate:"valid_iplist"` // 禁用列表
}

type ServiceAddTCPInput struct {
	ServiceName       string `json:"service_name" form:"service_name" comment:"服务名称" validate:"required,valid_service_name"`
	ServiceDesc       string `json:"service_desc" form:"service_desc" comment:"服务描述" validate:"required"`

	Port              int    `json:"port" form:"port" comment:"端口，需要设置8001-8999范围内" validate:"required,min=8001,max=8999"`
	HeaderTransfor    string `json:"header_transfor" form:"header_transfor" comment:"header头转换" validate:""`

	OpenAuth          int    `json:"open_auth" form:"open_auth" comment:"是否开启权限验证" validate:""`
	BlackList         string `json:"black_list" form:"black_list" comment:"黑名单IP，以逗号间隔，白名单优先级高于黑名单" validate:"valid_iplist"`
	WhiteList         string `json:"white_list" form:"white_list" comment:"白名单IP，以逗号间隔，白名单优先级高于黑名单" validate:"valid_iplist"`
	WhiteHostName     string `json:"white_host_name" form:"white_host_name" comment:"白名单主机，以逗号间隔" validate:"valid_iplist"`
	ClientIPFlowLimit int    `json:"clientip_flow_limit" form:"clientip_flow_limit" comment:"客户端IP限流" validate:""`
	ServiceFlowLimit  int    `json:"service_flow_limit" form:"service_flow_limit" comment:"服务端限流" validate:""`

	RoundType         int    `json:"round_type" form:"round_type" comment:"轮询策略" validate:""`
	IpList            string `json:"ip_list" form:"ip_list" comment:"IP列表" validate:"required,valid_ipportlist"`
	WeightList        string `json:"weight_list" form:"weight_list" comment:"权重列表" validate:"required,valid_weightlist"`
	ForbidList        string `json:"forbid_list" form:"forbid_list" comment:"禁用IP列表" validate:"valid_iplist"`
}

type ServiceUpdateHTTPInput struct {
	ServiceID   uint   `json:"service_id" form:"id" validate:"required"`                              //服务id
	ServiceName string `json:"service_name" form:"service_name" validate:"required,checkServiceName"` //服务名
	ServiceDesc string `json:"service_desc" form:"service_desc" validate:"required,max=255,min=1"`    //服务描述

	RuleType       int    `json:"rule_type" form:"rule_type" validate:"max=1,min=0"`                       //接入类型
	Rule           string `json:"rule" form:"rule" validate:"required,check_service_rule"`                 //域名或者前缀
	NeedHttps      int    `json:"need_https" form:"need_https" validate:"max=1,min=0"`                     //支持https
	NeedStripUri   int    `json:"need_strip_uri" form:"need_strip_uri" validate:"max=1,min=0"`             //启用strip_uri
	NeedWebsocket  int    `json:"need_websocket" form:"need_websocket" validate:"max=1,min=0"`             //是否支持websocket
	UrlRewrite     string `json:"url_rewrite" form:"url_rewrite" validate:"check_url_rewrite"`             //url重写功能
	HeaderTransfor string `json:"header_transfor" form:"header_transfor" validate:"check_header_transfor"` //header转换

	OpenAuth          int    `json:"open_auth" form:"open_auth"  validate:"max=1,min=0"`              //关键词
	BlackList         string `json:"black_list" form:"black_list" validate:""`                        //黑名单ip
	WhiteList         string `json:"white_list" form:"white_list" validate:""`                        //白名单ip
	ClientipFlowLimit int    `json:"clientip_flow_limit" form:"clientip_flow_limit" validate:"min=0"` //客户端ip限流
	ServiceFlowLimit  int    `json:"service_flow_limit" form:"service_flow_limit"validate:"min=0"`    //服务端限流

	RoundType              int    `json:"round_type" form:"round_type" validate:"max=3,min=0"`                       //轮询方式
	IpList                 string `json:"ip_list" form:"ip_list" validate:"required,check_iplist"`                   //ip列表
	WeightList             string `json:"weight_list" form:"weight_list" validate:"required,check_weightlist"`       //权重列表
	UpstreamConnectTimeout int    `json:"upstream_connect_timeout" form:"upstream_connect_timeout" validate:"min=0"` //建立连接超时, 单位s
	UpstreamHeaderTimeout  int    `json:"upstream_header_timeout" form:"upstream_header_timeout" validate:"min=0"`   //获取header超时, 单位s
	UpstreamIdleTimeout    int    `json:"upstream_idle_timeout" form:"upstream_idle_timeout"validate:"min=0"`        //链接最大空闲时间, 单位s
	UpstreamMaxIdle        int    `json:"upstream_max_idle" form:"upstream_max_idle" validate:"min=0"`               //最大空闲链接数
}

type ServiceUpdateGRPCInput struct {
	ServiceID   uint   `json:"service_id" form:"id" validate:"required"`                              //服务id
	ServiceAddGRPCInput *ServiceAddGRPCInput
}

type ServiceUpdateTCPInput struct {
	ServiceID   uint   `json:"service_id" form:"id" validate:"required"`                              //服务id
	ServiceAddTCPInput *ServiceAddTCPInput
}

type ServiceStatOutput struct {
	Today     []int64 `json:"today"`
	Yesterday []int64 `json:"yesterday"`
}

func (params *ServiceListInput) BindValidParam(c *gin.Context) error {
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

func (params *ServiceDeleteInput) BindValidParam(c *gin.Context) error {
	params.ServiceID = c.Query("service_id")
	return nil
}

func (params *ServiceAddHTTPInput) BindValidParam(c *gin.Context) error {
	if err := c.ShouldBindJSON(params); err != nil {
		utils.ResponseFormat(c, code.ServiceInsideError, nil)
		fmt.Println(err)
		return err
	}
	return nil
}

func (params *ServiceUpdateHTTPInput) BindValidParam(c *gin.Context) error {
	if err := c.ShouldBindJSON(params); err != nil {
		utils.ResponseFormat(c, code.ServiceInsideError, nil)
		fmt.Println(err)
		return err
	}
	return nil
}

func (params *ServiceAddGRPCInput) BindValidParam(c *gin.Context) error {
	if err := c.ShouldBindJSON(params); err != nil {
		utils.ResponseFormat(c, code.ServiceInsideError, nil)
		fmt.Println(err)
		return err
	}
	return nil
}

func (params *ServiceAddTCPInput) BindValidParam(c *gin.Context) error {
	if err := c.ShouldBindJSON(params); err != nil {
		utils.ResponseFormat(c, code.ServiceInsideError, nil)
		fmt.Println(err)
		return err
	}
	return nil
}
