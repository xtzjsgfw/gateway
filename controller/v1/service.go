package v1

import (
	"fmt"
	"gateway/dto"
	"gateway/extend/code"
	"gateway/extend/conf"
	"gateway/extend/utils"
	"gateway/models"
	"gateway/public"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"strconv"
	"strings"
	"time"
)

type ServiceController struct{}

// 返回服务列表
func (sc *ServiceController) ServiceList(c *gin.Context) {
	// 解析参数
	params := &dto.ServiceListInput{}
	if err := params.BindValidParam(c); err != nil {
		return
	}
	db := models.GetDB()

	//从db中分页读取基本信息
	gatewayServiceInfo := &models.GatewayServiceInfo{}
	list, total, err := gatewayServiceInfo.PageList(db, params)
	if err != nil {
		utils.ResponseFormat(c, code.ServiceInsideError, nil)
		return
	}
	
	//格式化输出信息
	outList := []dto.ServiceListItemOutPut{}
	for _, listItem := range list {
		serviceDetail, err := listItem.ServiceDetail(db, &listItem)
		if err != nil {
			utils.ResponseFormat(c, code.ServiceInsideError, nil)
			return
		}
		//1、http后缀接入 clusterIP+clusterPort+path
		//2、http域名接入 domain
		//3、tcp、grpc接入 clusterIP+servicePort
		serviceAddr := "unknow"
		clusterIP := conf.ClusterConf.IP
		clusterPort := conf.ClusterConf.Port
		clusterSSLPort := conf.ClusterConf.SSLPort
		if serviceDetail.Info.LoadType == public.LoadTypeHTTP &&
			serviceDetail.HTTPRule.RuleType == public.HTTPRuleTypePrefixURL &&
			serviceDetail.HTTPRule.NeedHttps == 0 {
			serviceAddr = fmt.Sprintf("%s:%s%s", clusterIP, clusterPort, serviceDetail.HTTPRule.Rule)
		}
		if serviceDetail.Info.LoadType == public.LoadTypeHTTP &&
			serviceDetail.HTTPRule.RuleType == public.HTTPRuleTypePrefixURL &&
			serviceDetail.HTTPRule.NeedHttps == 1 {
			serviceAddr = fmt.Sprintf("%s:%s%s", clusterIP, clusterSSLPort, serviceDetail.HTTPRule.Rule)
		}
		if serviceDetail.Info.LoadType == public.LoadTypeHTTP &&
			serviceDetail.HTTPRule.RuleType == public.HTTPRuleTypeDomain {
			serviceAddr = serviceDetail.HTTPRule.Rule
		}
		if serviceDetail.Info.LoadType == public.LoadTypeTCP {
			serviceAddr = fmt.Sprintf("%s:%d", clusterIP, serviceDetail.TCPRule.Port)
		}
		if serviceDetail.Info.LoadType == public.LoadTypeGRPC {
			serviceAddr = fmt.Sprintf("%s:%d", clusterIP, serviceDetail.GRPCRule.Port)
		}

		ipList := serviceDetail.LoadBalance.GetIPListByModel()

		outItem := dto.ServiceListItemOutPut{
			ID:          listItem.ID,
			LoadType:    listItem.LoadType,
			ServiceName: listItem.ServiceName,
			ServiceDesc: listItem.ServiceDesc,
			ServiceAddr: serviceAddr,
			QPS:         0,
			QPD:         0,
			TotalNode:   len(ipList),
		}
		outList = append(outList, outItem)
	}

	out := &dto.ServiceListOutPut{
		Total: total,
		List:  outList,
	}
	utils.ResponseFormat(c, code.Success, out)
}

// 删除服务，标记is_delete为1
func (sc *ServiceController) ServiceDelete(c *gin.Context) {
	params := &dto.ServiceDeleteInput{}
	if err := params.BindValidParam(c); err != nil {
		utils.ResponseFormat(c, code.ServiceInsideError, nil)
		return
	}

	db := models.GetDB()

	id, err := strconv.Atoi(params.ServiceID)
	if err != nil {
		utils.ResponseFormat(c, code.ServiceInsideError, nil)
		return
	}
	gatewayServiceInfo := &models.GatewayServiceInfo{Model: gorm.Model{ID: uint(id)}}
	serviceInfo, err := gatewayServiceInfo.Find(db, gatewayServiceInfo)
	if err != nil {
		utils.ResponseFormat(c, code.ServiceInsideError, nil)
		return
	}
	serviceInfo.IsDelete = 1

	if err := serviceInfo.Save(db); err != nil {
		utils.ResponseFormat(c, code.DBSaveError, nil)
		return
	}
	utils.ResponseFormat(c, code.Success, gin.H{
		"service_id": params.ServiceID,
	})
}

// 服务详情
func (sc *ServiceController) ServiceDetail(c *gin.Context) {
	params := &dto.ServiceDeleteInput{}
	if err := params.BindValidParam(c); err != nil {
		utils.ResponseFormat(c, code.ServiceInsideError, nil)
		return
	}
	fmt.Println(params)

	db := models.GetDB()

	id, err := strconv.Atoi(params.ServiceID)
	if err != nil {
		utils.ResponseFormat(c, code.ServiceInsideError, nil)
		return
	}

	gatewayServiceInfo := &models.GatewayServiceInfo{Model: gorm.Model{ID: uint(id)}}
	serviceInfo, err := gatewayServiceInfo.Find(db, gatewayServiceInfo)
	if err != nil {
		utils.ResponseFormat(c, code.ServiceInsideError, nil)
		return
	}
	serviceDetail, err := serviceInfo.ServiceDetail(db, serviceInfo)
	if err != nil {
		utils.ResponseFormat(c, code.ServiceInsideError, nil)
		return
	}
	utils.ResponseFormat(c, code.Success, serviceDetail)
}

// 流量统计
func (sc *ServiceController) ServiceStat(c *gin.Context) {
	params := &dto.ServiceDeleteInput{}
	if err := params.BindValidParam(c); err != nil {
		utils.ResponseFormat(c, code.ServiceInsideError, nil)
		return
	}

	//db := models.GetDB()

	id, err := strconv.Atoi(params.ServiceID)
	if err != nil {
		utils.ResponseFormat(c, code.ServiceInsideError, nil)
		return
	}

	fmt.Println(id)

	//gatewayServiceInfo := &models.GatewayServiceInfo{Model: gorm.Model{ID: uint(id)}}
	//serviceInfo, err := gatewayServiceInfo.Find(db, gatewayServiceInfo)
	//if err != nil {
	//	utils.ResponseFormat(c, code.ServiceInsideError, nil)
	//	return
	//}
	//serviceDetail, err := serviceInfo.ServiceDetail(serviceInfo)
	//if err != nil {
	//	utils.ResponseFormat(c, code.ServiceInsideError, nil)
	//	return
	//}
	todayList := []int64{}
	for i := 0; i < time.Now().Hour(); i++ {
		todayList = append(todayList, 0)
	}
	yesterdayList := []int64{}
	for i := 0; i < 23; i++ {
		yesterdayList = append(yesterdayList, 0)
	}
	serviceStatOutput := dto.ServiceStatOutput{
		Today:     todayList,
		Yesterday: yesterdayList,
	}

	utils.ResponseFormat(c, code.Success, serviceStatOutput)
}

// 添加HTTP服务
func (sc *ServiceController) ServiceAddHTTP(c *gin.Context) {
	params := &dto.ServiceAddHTTPInput{}
	err := params.BindValidParam(c)
	if err != nil {
		return
	}

	fmt.Println("step 1")
	// 判断ip列表和权重列表数量是否相对应
	if len(strings.Split(params.IpList, "\n")) != len(strings.Split(params.WeightList, "\n")) {
		utils.ResponseFormat(c, code.WeightAndIpNumNotEqualError, nil)
		return
	}

	db := models.GetDB()
	// 开启事务
	db.Begin()

	fmt.Println("step 2")
	// 查询HTTP服务是否存在
	serviceInfo := &models.GatewayServiceInfo{
		ServiceName: params.ServiceName,
		IsDelete:    0,
	}
	serviceInfo, err = serviceInfo.Find(db, serviceInfo)
	if serviceInfo != nil && err == nil {
		db.Rollback()
		utils.ResponseFormat(c, code.ServiceIsExistError, nil)
		return
	}
	if err != nil {
		db.Rollback()
		utils.ResponseFormat(c, code.ServiceInsideError, nil)
		return
	}

	fmt.Println("step 3")
	// 查询HTTPRule是否存在
	httpUrl := &models.GatewayServiceHttpRule{RuleType: params.RuleType, Rule: params.Rule}
	if _, err := httpUrl.Find(db, httpUrl); err == nil {
		db.Rollback()
		utils.ResponseFormat(c, code.RulePixDomIsExistError, nil)
		return
	}

	fmt.Println("step 4")
	serviceModel := &models.GatewayServiceInfo{
		ServiceName: params.ServiceName,
		ServiceDesc: params.ServiceDesc,
		LoadType:    public.LoadTypeHTTP,
	}
	if err = serviceModel.Save(db); err != nil {
		fmt.Println(err)
		db.Rollback()
		utils.ResponseFormat(c, code.DBSaveError, nil)
		return
	}
	// service.ID

	fmt.Println("step 5")
	httpRule := &models.GatewayServiceHttpRule{
		ServiceID:      serviceModel.ID,
		RuleType:       params.RuleType,
		Rule:           params.Rule,
		NeedHttps:      params.NeedHttps,
		NeedStripUri:   params.NeedStripUri,
		NeedWebsocket:  params.NeedWebsocket,
		UrlRewrite:     params.UrlRewrite,
		HeaderTransfor: params.HeaderTransfor,
	}
	if err = httpRule.Save(db); err != nil {
		fmt.Println(err)
		db.Rollback()
		utils.ResponseFormat(c, code.DBSaveError, nil)
		return
	}

	fmt.Println("step 6")
	accessControl := &models.GatewayServiceAccessControl{
		ServiceID:         serviceModel.ID,
		OpenAuth:          params.OpenAuth,
		BlackList:         params.BlackList,
		WhiteList:         params.WhiteList,
		ClientIpFlowLimit: params.ClientipFlowLimit,
		ServiceFlowLimit:  params.ServiceFlowLimit,
	}
	if err = accessControl.Save(db); err != nil {
		fmt.Println(err)
		db.Rollback()
		utils.ResponseFormat(c, code.DBSaveError, nil)
		return
	}

	fmt.Println("step 7")
	loadBalance := &models.GatewayServiceLoadBalance{
		ServiceID:              serviceModel.ID,
		RoundType:              params.RoundType,
		IpList:                 params.IpList,
		WeightList:             params.WeightList,
		UpstreamConnectTimeout: params.UpstreamConnectTimeout,
		UpstreamHeaderTimeout:  params.UpstreamHeaderTimeout,
		UpstreamIdleTimeout:    params.UpstreamIdleTimeout,
		UpstreamMaxIdle:        params.UpstreamMaxIdle,
	}
	if err = loadBalance.Save(db); err != nil {
		fmt.Println(err)
		db.Rollback()
		utils.ResponseFormat(c, code.DBSaveError, nil)
		return
	}

	// 提交
	fmt.Println("step 8")
	db.Commit()

	fmt.Println("step 9")
	utils.ResponseFormat(c, code.Success, gin.H{
		"service_id": serviceModel.ID,
	})
}

// 更新HTTP服务
func (sc *ServiceController) ServiceUpdateHTTP(c *gin.Context) {
	params := &dto.ServiceUpdateHTTPInput{}
	err := params.BindValidParam(c)
	if err != nil {
		return
	}

	// 判断ip列表和权重列表数量是否相对应
	if len(strings.Split(params.IpList, "\n")) != len(strings.Split(params.WeightList, "\n")) {
		utils.ResponseFormat(c, code.WeightAndIpNumNotEqualError, nil)
		return
	}

	db := models.GetDB()
	// 开启事务
	db.Begin()

	// 查询HTTP服务是否存在
	serviceInfo := &models.GatewayServiceInfo{ServiceName: params.ServiceName}
	serviceInfo, err = serviceInfo.Find(db, serviceInfo)
	if serviceInfo == nil && err == nil {
		fmt.Println(err)
		utils.ResponseFormat(c, code.ServiceIsNotExistError, nil)
		return
	}
	if serviceInfo == nil && err != nil {
		utils.ResponseFormat(c, code.ServiceInsideError, "淦")
		db.Rollback()
		return
	}

	serviceDetail, err := serviceInfo.ServiceDetail(db, serviceInfo)
	if err != nil {
		fmt.Println(err)
		utils.ResponseFormat(c, code.ServiceIsNotExistError, nil)
		return
	}

	fmt.Printf("%v", serviceDetail.Info)

	info := serviceDetail.Info
	info.ServiceDesc = params.ServiceDesc
	if err = db.Save(info).Error; err != nil {
		fmt.Println(err)
		db.Rollback()
		utils.ResponseFormat(c, code.DBSaveError, nil)
		return
	}

	httpRule := serviceDetail.HTTPRule
	httpRule.NeedHttps = params.NeedHttps
	httpRule.NeedStripUri = params.NeedStripUri
	httpRule.NeedWebsocket = params.NeedWebsocket
	httpRule.UrlRewrite = params.UrlRewrite
	httpRule.HeaderTransfor = params.HeaderTransfor
	if err = db.Save(httpRule).Error; err != nil {
		fmt.Println(err)
		db.Rollback()
		utils.ResponseFormat(c, code.DBSaveError, nil)
		return
	}

	accessControl := serviceDetail.AccessControl
	accessControl.OpenAuth = params.OpenAuth
	accessControl.BlackList = params.BlackList
	accessControl.WhiteList = params.WhiteList
	accessControl.ClientIpFlowLimit = params.ClientipFlowLimit
	accessControl.ServiceFlowLimit = params.ServiceFlowLimit
	if err = db.Save(accessControl).Error; err != nil {
		fmt.Println(err)
		db.Rollback()
		utils.ResponseFormat(c, code.DBSaveError, nil)
		return
	}

	loadBalance := serviceDetail.LoadBalance
	loadBalance.RoundType = params.RoundType
	loadBalance.IpList = params.IpList
	loadBalance.WeightList = params.WeightList
	loadBalance.UpstreamConnectTimeout = params.UpstreamConnectTimeout
	loadBalance.UpstreamHeaderTimeout = params.UpstreamHeaderTimeout
	loadBalance.UpstreamIdleTimeout = params.UpstreamIdleTimeout
	loadBalance.UpstreamMaxIdle = params.UpstreamMaxIdle
	if err = db.Save(loadBalance).Error; err != nil {
		fmt.Println(err)
		db.Rollback()
		utils.ResponseFormat(c, code.DBSaveError, nil)
		return
	}

	// 提交
	db.Commit()

	utils.ResponseFormat(c, code.Success, gin.H{
		"msg": "更新成功",
	})
}

// 添加GRPC服务
func (sc *ServiceController) ServiceAddGRPC(c *gin.Context) {
	params := &dto.ServiceAddGRPCInput{}
	err := params.BindValidParam(c)
	if err != nil {
		return
	}

	db := models.GetDB()

	fmt.Println("step 1")
	// 验证服务是否存在
	infoSearch := &models.GatewayServiceInfo{
		ServiceName: params.ServiceName,
		IsDelete:    0,
	}
	if _, err = infoSearch.Find(db, infoSearch); err == nil {
		utils.ResponseFormat(c, code.ServiceIsExistError, nil)
		return
	}

	fmt.Println("step 2")
	// 验证端口是否被占用
	tcpRuleSearch := &models.GatewayServiceTcpRule{
		Port: params.Port,
	}
	if _, err = tcpRuleSearch.Find(db, tcpRuleSearch); err == nil {
		utils.ResponseFormat(c, code.PortOccupiedError, nil)
		return
	}
	grpcRuleSearch := &models.GatewayServiceGrpcRule{
		Port: params.Port,
	}
	if _, err = grpcRuleSearch.Find(db, grpcRuleSearch); err == nil {
		utils.ResponseFormat(c, code.PortOccupiedError, nil)
		return
	}

	fmt.Println("step 3")
	// 判断ip列表和权重列表数量是否相对应
	if len(strings.Split(params.IpList, "\n")) != len(strings.Split(params.WeightList, "\n")) {
		utils.ResponseFormat(c, code.WeightAndIpNumNotEqualError, nil)
		return
	}

	// 开启事务
	db.Begin()

	fmt.Println("step 4")
	// info表存储
	info := &models.GatewayServiceInfo{
		LoadType:    public.LoadTypeGRPC,
		ServiceName: params.ServiceName,
		ServiceDesc: params.ServiceDesc,
	}
	if err = info.Save(db); err != nil {
		fmt.Println(err)
		db.Rollback()
		utils.ResponseFormat(c, code.DBSaveError, nil)
		return
	}

	fmt.Println("step 5")
	// grpcRule表存储
	grpcRule := &models.GatewayServiceGrpcRule{
		ServiceID:      info.ID,
		Port:           params.Port,
		HeaderTransfor: params.HeaderTransfor,
	}
	if err = grpcRule.Save(db); err != nil {
		fmt.Println(err)
		db.Rollback()
		utils.ResponseFormat(c, code.DBSaveError, nil)
		return
	}

	fmt.Println("step 6")
	// 权限控制表存储
	accessControl := &models.GatewayServiceAccessControl{
		ServiceID:         info.ID,
		OpenAuth:          params.OpenAuth,
		BlackList:         params.BlackList,
		WhiteList:         params.WhiteList,
		WhiteHostName:     params.WhiteHostName,
		ClientIpFlowLimit: params.ClientipFlowLimit,
		ServiceFlowLimit:  params.ServiceFlowLimit,
	}
	if err = accessControl.Save(db); err != nil {
		fmt.Println(err)
		db.Rollback()
		utils.ResponseFormat(c, code.DBSaveError, nil)
		return
	}

	fmt.Println("step 7")
	// loadBalance表存储
	loadBalance := &models.GatewayServiceLoadBalance{
		ServiceID:  info.ID,
		RoundType:  params.RoundType,
		IpList:     params.IpList,
		WeightList: params.WeightList,
		ForbidList: params.ForbidList,
	}
	if err = loadBalance.Save(db); err != nil {
		fmt.Println(err)
		db.Rollback()
		utils.ResponseFormat(c, code.DBSaveError, nil)
		return
	}

	fmt.Println("step 8")
	// 提交
	db.Commit()
	fmt.Println("step 9")
	utils.ResponseFormat(c, code.Success, gin.H{
		"service_id": info.ID,
	})
}

// 更新GRPC服务
func (sc *ServiceController) ServiceUpdateGRPC(c *gin.Context) {
	params := &dto.ServiceUpdateHTTPInput{}
	err := params.BindValidParam(c)
	if err != nil {
		return
	}

	// 判断ip列表和权重列表数量是否相对应
	if len(strings.Split(params.IpList, "\n")) != len(strings.Split(params.WeightList, "\n")) {
		utils.ResponseFormat(c, code.WeightAndIpNumNotEqualError, nil)
		return
	}

	db := models.GetDB()
	// 开启事务
	db.Begin()

	// 查询HTTP服务是否存在
	serviceInfo := &models.GatewayServiceInfo{ServiceName: params.ServiceName}
	serviceInfo, err = serviceInfo.Find(db, serviceInfo)
	if err != nil {
		fmt.Println(err)
		utils.ResponseFormat(c, code.ServiceIsNotExistError, nil)
		return
	}

	serviceDetail, err := serviceInfo.ServiceDetail(db, serviceInfo)
	if err != nil {
		fmt.Println(err)
		utils.ResponseFormat(c, code.ServiceIsNotExistError, nil)
		return
	}

	fmt.Printf("%v", serviceDetail.Info)

	info := serviceDetail.Info
	info.ServiceDesc = params.ServiceDesc
	if err = db.Save(info).Error; err != nil {
		fmt.Println(err)
		db.Rollback()
		utils.ResponseFormat(c, code.DBSaveError, nil)
		return
	}

	httpRule := serviceDetail.HTTPRule
	httpRule.NeedHttps = params.NeedHttps
	httpRule.NeedStripUri = params.NeedStripUri
	httpRule.NeedWebsocket = params.NeedWebsocket
	httpRule.UrlRewrite = params.UrlRewrite
	httpRule.HeaderTransfor = params.HeaderTransfor
	if err = db.Save(httpRule).Error; err != nil {
		fmt.Println(err)
		db.Rollback()
		utils.ResponseFormat(c, code.DBSaveError, nil)
		return
	}

	accessControl := serviceDetail.AccessControl
	accessControl.OpenAuth = params.OpenAuth
	accessControl.BlackList = params.BlackList
	accessControl.WhiteList = params.WhiteList
	accessControl.ClientIpFlowLimit = params.ClientipFlowLimit
	accessControl.ServiceFlowLimit = params.ServiceFlowLimit
	if err = db.Save(accessControl).Error; err != nil {
		fmt.Println(err)
		db.Rollback()
		utils.ResponseFormat(c, code.DBSaveError, nil)
		return
	}

	loadBalance := serviceDetail.LoadBalance
	loadBalance.RoundType = params.RoundType
	loadBalance.IpList = params.IpList
	loadBalance.WeightList = params.WeightList
	loadBalance.UpstreamConnectTimeout = params.UpstreamConnectTimeout
	loadBalance.UpstreamHeaderTimeout = params.UpstreamHeaderTimeout
	loadBalance.UpstreamIdleTimeout = params.UpstreamIdleTimeout
	loadBalance.UpstreamMaxIdle = params.UpstreamMaxIdle
	if err = db.Save(loadBalance).Error; err != nil {
		fmt.Println(err)
		db.Rollback()
		utils.ResponseFormat(c, code.DBSaveError, nil)
		return
	}

	// 提交
	db.Commit()

	utils.ResponseFormat(c, code.Success, gin.H{
		"msg": "更新成功",
	})
}

// 添加TCP服务
func (sc *ServiceController) ServiceAddTCP(c *gin.Context) {
	params := &dto.ServiceAddTCPInput{}
	err := params.BindValidParam(c)
	if err != nil {
		return
	}

	db := models.GetDB()

	// 验证服务是否存在
	infoSearch := &models.GatewayServiceInfo{
		ServiceName: params.ServiceName,
		IsDelete:    0,
	}
	if _, err = infoSearch.Find(db, infoSearch); err == nil {
		utils.ResponseFormat(c, code.ServiceIsExistError, nil)
		return
	}

	// 验证端口是否被占用
	tcpRuleSearch := &models.GatewayServiceTcpRule{
		Port: params.Port,
	}
	if _, err = tcpRuleSearch.Find(db, tcpRuleSearch); err == nil {
		utils.ResponseFormat(c, code.PortOccupiedError, nil)
		return
	}
	grpcRuleSearch := &models.GatewayServiceGrpcRule{
		Port: params.Port,
	}
	if _, err = grpcRuleSearch.Find(db, grpcRuleSearch); err == nil {
		utils.ResponseFormat(c, code.PortOccupiedError, nil)
		return
	}

	// 判断ip列表和权重列表数量是否相对应
	if len(strings.Split(params.IpList, "\n")) != len(strings.Split(params.WeightList, "\n")) {
		utils.ResponseFormat(c, code.WeightAndIpNumNotEqualError, nil)
		return
	}

	// 开启事务
	db.Begin()

	// info表存储
	info := &models.GatewayServiceInfo{
		LoadType:    public.LoadTypeGRPC,
		ServiceName: params.ServiceName,
		ServiceDesc: params.ServiceDesc,
	}
	if err = info.Save(db); err != nil {
		fmt.Println(err)
		db.Rollback()
		utils.ResponseFormat(c, code.DBSaveError, nil)
		return
	}

	// grpcRule表存储
	grpcRule := &models.GatewayServiceGrpcRule{
		ServiceID:      info.ID,
		Port:           params.Port,
		HeaderTransfor: params.HeaderTransfor,
	}
	if err = grpcRule.Save(db); err != nil {
		fmt.Println(err)
		db.Rollback()
		utils.ResponseFormat(c, code.DBSaveError, nil)
		return
	}

	// 权限控制表存储
	accessControl := &models.GatewayServiceAccessControl{
		ServiceID:         info.ID,
		OpenAuth:          params.OpenAuth,
		BlackList:         params.BlackList,
		WhiteList:         params.WhiteList,
		WhiteHostName:     params.WhiteHostName,
		ClientIpFlowLimit: params.ClientIPFlowLimit,
		ServiceFlowLimit:  params.ServiceFlowLimit,
	}
	if err = accessControl.Save(db); err != nil {
		fmt.Println(err)
		db.Rollback()
		utils.ResponseFormat(c, code.DBSaveError, nil)
		return
	}

	// loadBalance表存储
	loadBalance := &models.GatewayServiceLoadBalance{
		ServiceID:  info.ID,
		RoundType:  params.RoundType,
		IpList:     params.IpList,
		WeightList: params.WeightList,
		ForbidList: params.ForbidList,
	}
	if err = loadBalance.Save(db); err != nil {
		fmt.Println(err)
		db.Rollback()
		utils.ResponseFormat(c, code.DBSaveError, nil)
		return
	}

	// 提交
	db.Commit()
	utils.ResponseFormat(c, code.Success, gin.H{
		"service_id": info.ID,
	})
}

// 更新TCP服务
func (sc *ServiceController) ServiceUpdateTCP(c *gin.Context) {
	params := &dto.ServiceUpdateHTTPInput{}
	err := params.BindValidParam(c)
	if err != nil {
		return
	}

	// 判断ip列表和权重列表数量是否相对应
	if len(strings.Split(params.IpList, "\n")) != len(strings.Split(params.WeightList, "\n")) {
		utils.ResponseFormat(c, code.WeightAndIpNumNotEqualError, nil)
		return
	}

	db := models.GetDB()
	// 开启事务
	db.Begin()

	// 查询HTTP服务是否存在
	serviceInfo := &models.GatewayServiceInfo{ServiceName: params.ServiceName}
	serviceInfo, err = serviceInfo.Find(db, serviceInfo)
	if err != nil {
		fmt.Println(err)
		utils.ResponseFormat(c, code.ServiceIsNotExistError, nil)
		return
	}

	serviceDetail, err := serviceInfo.ServiceDetail(db, serviceInfo)
	if err != nil {
		fmt.Println(err)
		utils.ResponseFormat(c, code.ServiceIsNotExistError, nil)
		return
	}

	fmt.Printf("%v", serviceDetail.Info)

	info := serviceDetail.Info
	info.ServiceDesc = params.ServiceDesc
	if err = db.Save(info).Error; err != nil {
		fmt.Println(err)
		db.Rollback()
		utils.ResponseFormat(c, code.DBSaveError, nil)
		return
	}

	httpRule := serviceDetail.HTTPRule
	httpRule.NeedHttps = params.NeedHttps
	httpRule.NeedStripUri = params.NeedStripUri
	httpRule.NeedWebsocket = params.NeedWebsocket
	httpRule.UrlRewrite = params.UrlRewrite
	httpRule.HeaderTransfor = params.HeaderTransfor
	if err = db.Save(httpRule).Error; err != nil {
		fmt.Println(err)
		db.Rollback()
		utils.ResponseFormat(c, code.DBSaveError, nil)
		return
	}

	accessControl := serviceDetail.AccessControl
	accessControl.OpenAuth = params.OpenAuth
	accessControl.BlackList = params.BlackList
	accessControl.WhiteList = params.WhiteList
	accessControl.ClientIpFlowLimit = params.ClientipFlowLimit
	accessControl.ServiceFlowLimit = params.ServiceFlowLimit
	if err = db.Save(accessControl).Error; err != nil {
		fmt.Println(err)
		db.Rollback()
		utils.ResponseFormat(c, code.DBSaveError, nil)
		return
	}

	loadBalance := serviceDetail.LoadBalance
	loadBalance.RoundType = params.RoundType
	loadBalance.IpList = params.IpList
	loadBalance.WeightList = params.WeightList
	loadBalance.UpstreamConnectTimeout = params.UpstreamConnectTimeout
	loadBalance.UpstreamHeaderTimeout = params.UpstreamHeaderTimeout
	loadBalance.UpstreamIdleTimeout = params.UpstreamIdleTimeout
	loadBalance.UpstreamMaxIdle = params.UpstreamMaxIdle
	if err = db.Save(loadBalance).Error; err != nil {
		fmt.Println(err)
		db.Rollback()
		utils.ResponseFormat(c, code.DBSaveError, nil)
		return
	}

	// 提交
	db.Commit()

	utils.ResponseFormat(c, code.Success, gin.H{
		"msg": "更新成功",
	})
}
