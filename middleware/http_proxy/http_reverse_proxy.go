package http_proxy

import (
	"fmt"
	"gateway/reverse_proxy"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"time"
)

func HTTPReverseProxyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//service, ok := c.Get("service")
		//if !ok {
		//	utils.ResponseFormat(c, code.GetServiceFailed, nil)
		//	c.Abort()
		//	return
		//}
		//fmt.Printf("%T\n", service) // *models.ServiceDetail
		//fmt.Println(utils.Obj2Json(service))

		fmt.Println("step1")
		//创建 reverseproxy
		//使用 reverseproxy.ServerHTTP(c.Request,c.Response)
		proxy := reverse_proxy.NewLoadBalanceReverseProxy(c, transport)
		proxy.ServeHTTP(c.Writer, c.Request)
		fmt.Println("step2")
		c.Abort()
		return
	}
}

var transport = &http.Transport{
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second, //连接超时
		KeepAlive: 30 * time.Second, //长连接超时时间
	}).DialContext,
	MaxIdleConns:          100,              //最大空闲连接
	IdleConnTimeout:       90 * time.Second, //空闲超时时间
	TLSHandshakeTimeout:   10 * time.Second, //tls握手超时时间
	ExpectContinueTimeout: 1 * time.Second,  //100-continue 超时时间
}
