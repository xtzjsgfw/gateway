package models

import (
	"errors"
	"gateway/dto"
	"gateway/public"
	"github.com/gin-gonic/gin"
	"strings"
	"sync"
)

type ServiceDetail struct {
	Info          *GatewayServiceInfo          `json:"info"`
	HTTPRule      *GatewayServiceHttpRule      `json:"http_rule"`
	TCPRule       *GatewayServiceTcpRule       `json:"tcp_rule"`
	GRPCRule      *GatewayServiceGrpcRule      `json:"grpc_rule"`
	LoadBalance   *GatewayServiceLoadBalance   `json:"load_balance"`
	AccessControl *GatewayServiceAccessControl `json:"access_control"`
}

var ServiceManagerHandler *ServiceManager

func init() {
	ServiceManagerHandler = NewServiceManager()
}

type ServiceManager struct {
	ServiceMap   map[string]*ServiceDetail
	ServiceSlice []*ServiceDetail
	Locker       sync.RWMutex
	init         sync.Once
	err          error
}

func NewServiceManager() *ServiceManager {
	return &ServiceManager{
		ServiceMap:   map[string]*ServiceDetail{},
		ServiceSlice: []*ServiceDetail{},
		Locker:       sync.RWMutex{},
		init:         sync.Once{},
	}
}

func (s *ServiceManager) HTTPAccessMode(c *gin.Context) (*ServiceDetail, error) {
	//1、前缀匹配 /abc ==> serviceSlice.rule
	// c.Request.URL.Path

	//2、域名匹配 www.test.com ==> serviceSlice.rule
	// c.Request.Host
	host := c.Request.Host
	host = host[0:strings.Index(host, ":")]
	path := c.Request.URL.Path

	for _, serviceItem := range s.ServiceSlice {
		if serviceItem.Info.LoadType != public.LoadTypeHTTP {
			continue
		}
		if serviceItem.HTTPRule.RuleType == public.HTTPRuleTypeDomain {
			if serviceItem.HTTPRule.Rule == host {
				return serviceItem, nil
			}
		}
		if serviceItem.HTTPRule.RuleType == public.HTTPRuleTypePrefixURL {
			if strings.HasPrefix(path, serviceItem.HTTPRule.Rule) {
				return serviceItem, nil
			}
		}
	}
	return nil, errors.New("not matched service")
}

func (s *ServiceManager) LoadOnce() error {
	s.init.Do(func() {
		db := GetDB()
		gatewayServiceInfo := &GatewayServiceInfo{}
		params := &dto.ServiceListInput{PageNum: 1, PageSize: 9999}
		list, _, err := gatewayServiceInfo.PageList(db, params)
		if err != nil {
			s.err = err
			return
		}
		s.Locker.Lock()
		defer s.Locker.Unlock()

		for _, listItem := range list {
			tmpItem := listItem
			serviceDetail, err := tmpItem.ServiceDetail(db, &tmpItem)
			if err != nil {
				s.err = err
				return
			}
			s.ServiceMap[listItem.ServiceName] = serviceDetail
			s.ServiceSlice = append(s.ServiceSlice, serviceDetail)
		}
	})
	return s.err
}
