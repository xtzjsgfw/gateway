package reverse_proxy

import (
	"gateway/extend/code"
	"gateway/extend/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httputil"
	"strings"
)

func NewLoadBalanceReverseProxy(c *gin.Context, transport *http.Transport) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = "127.0.0.1:2003"
		req.URL.Path = ""
		req.Host = "127.0.0.1:2003"

		//fmt.Println("reqreqreq::::::::::::::::::", req)
		if _, ok := req.Header["User-Agent"]; !ok {
			req.Header.Set("User-Agent", "user-agent")
		}
	}

	modifyFunc := func(resp *http.Response) error {
		if strings.Contains(resp.Header.Get("Connection"), "Upgrade") {
			return nil
		}
		return nil
	}

	errFunc := func(w http.ResponseWriter, r *http.Request, err error) {
		utils.ResponseFormat(c, code.ReverseProxyError, nil)
	}
	return &httputil.ReverseProxy{Director: director, ModifyResponse: modifyFunc, ErrorHandler: errFunc}
}
