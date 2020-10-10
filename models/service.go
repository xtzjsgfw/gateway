package models

type ServiceDetail struct {
	Info          *GatewayServiceInfo          `json:"info"`
	HTTPRule      *GatewayServiceHttpRule      `json:"http_rule"`
	TCPRule       *GatewayServiceTcpRule       `json:"tcp_rule"`
	GRPCRule      *GatewayServiceGrpcRule      `json:"grpc_rule"`
	LoadBalance   *GatewayServiceLoadBalance   `json:"load_balance"`
	AccessControl *GatewayServiceAccessControl `json:"access_control"`
}

