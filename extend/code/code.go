package code

import "net/http"

type Code struct {
	Status  int    `json:"status"`  // HTTP状态
	Code    int    `json:"code"`    // 业务码
	Message string `json:"message"` // 业务响应信息
}

var (
	// 6xxx 登陆相关
	UsernameLengthError = &Code{http.StatusBadRequest, 6001, "用户名长度错误"}
	PasswordLengthError = &Code{http.StatusBadRequest, 6002, "密码长度错误"}
	MobileLengthError   = &Code{http.StatusBadRequest, 6003, "手机号长度错误"}
	UserIsExistError    = &Code{http.StatusBadRequest, 6004, "该用户已存在"}
	UserIsNotExistError = &Code{http.StatusBadRequest, 6005, "该用户不存在"}
	UserOrPassError     = &Code{http.StatusBadRequest, 6006, "用户名或密码错误"}

	TokenIsNotExistError = &Code{http.StatusUnauthorized, 6007, "Token不存在"}
	TokenParseError      = &Code{http.StatusUnauthorized, 6008, "Token解析出错"}
	TokenInvalid         = &Code{http.StatusUnauthorized, 6009, "Token无效"}

	Success = &Code{http.StatusOK, 2000, "请求成功"}

	// 7xxx Service相关
	ServiceIsExistError    = &Code{http.StatusBadRequest, 7001, "已经存在同名服务"}
	ServiceIsNotExistError = &Code{http.StatusBadRequest, 7002, "服务不存在"}
	// 接入前缀/域名已存在
	RulePixDomIsExistError      = &Code{http.StatusBadRequest, 7003, "接入前缀/域名已存在"}
	RulePixOrDomIsNotExistError = &Code{http.StatusBadRequest, 7003, "接入前缀/域名不存在"}
	WeightAndIpNumNotEqualError = &Code{http.StatusBadRequest, 7004, "Ip列表数量与权重数量不等"}
	DBSaveError                 = &Code{http.StatusBadRequest, 7005, "数据库保存错误"}

	PortOccupiedError = &Code{http.StatusBadRequest, 7006, "端口被占用"}

	// 8xxx app相关
	APPIsNotExistError = &Code{http.StatusBadRequest, 8001, "该租户不存在"}
	APPIsExistError    = &Code{http.StatusBadRequest, 8001, "已存在该租户"}

	RequestParamError  = &Code{http.StatusBadRequest, 400001, "请求参数有误"}
	ServiceInsideError = &Code{http.StatusInternalServerError, 5000, "服务器内部错误"}

	PongCode = &Code{http.StatusOK, 2001, "Pong"}
)
