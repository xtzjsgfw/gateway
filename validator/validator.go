package validator

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"regexp"
	"strings"
)

func Init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("check_service_name", checkServiceName)
		v.RegisterValidation("check_service_rule", checkServiceRule)
		v.RegisterValidation("check_url_rewrite", checkUrlRewrite)
		v.RegisterValidation("check_header_transfor", checkHeaderTransfor)
		v.RegisterValidation("check_iplist", checkIpList)
		v.RegisterValidation("check_weightlist", checkWeightList)
	}

}

func checkServiceName(f validator.FieldLevel) bool {
	serviceName := f.Field().String()
	re := `[a-zA-Z0-9]{6,128}`
	r := regexp.MustCompile(re)
	return r.MatchString(serviceName)
}

func checkServiceRule(f validator.FieldLevel) bool {
	serviceName := f.Field().String()
	re := `\S{6,128}`
	r := regexp.MustCompile(re)
	return r.MatchString(serviceName)
}

func checkUrlRewrite(f validator.FieldLevel) bool {
	urlRewrite := f.Field().String()
	splitUrlRewrite := strings.Split(urlRewrite, "\n")
	for _, ms := range splitUrlRewrite {
		if len(strings.Split(ms, " ")) != 2 {
			return false
		}
	}
	return true
}

func checkHeaderTransfor(f validator.FieldLevel) bool {
	urlRewrite := f.Field().String()
	splitUrlRewrite := strings.Split(urlRewrite, "\n")
	for _, ms := range splitUrlRewrite {
		if len(strings.Split(ms, " ")) != 3 {
			return false
		}
	}
	return true
}

func checkIpList(f validator.FieldLevel) bool {
	ipList := f.Field().String()
	splitIpList := strings.Split(ipList, "\n")
	for _, ms := range splitIpList {
		if matched, _ := regexp.Match(`^\S+\:\d+$`, []byte(ms)); !matched {
			return false
		}
	}
	return true
}

func checkWeightList(f validator.FieldLevel) bool {
	weightList := f.Field().String()
	splitWeightList := strings.Split(weightList, "\n")
	for _, ms := range splitWeightList {
		if matched, _ := regexp.Match(`^\d+`, []byte(ms)); !matched {
			return false
		}
	}
	return true
}
