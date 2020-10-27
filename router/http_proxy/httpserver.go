package http_proxy

import (
	"context"
	"fmt"
	"gateway/cert_file"
	"gateway/extend/conf"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

var (
	HttpSrvHandler  *http.Server
	HttpsSrvHandler *http.Server
)

func HttpServerRun() {
	r := HttpInit(gin.Recovery(), gin.Logger())

	HttpSrvHandler = &http.Server{
		Addr:           conf.HttpConf.Addr,
		Handler:        r,
		ReadTimeout:    time.Duration(conf.HttpConf.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(conf.HttpConf.WriteTimeout) * time.Second,
		MaxHeaderBytes: conf.HttpConf.MaxHeaderBytes,
	}
	log.Printf(" [INFO] http_proxy_run %v\n", conf.HttpConf.Addr)
	if err := HttpSrvHandler.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Println("HTTP服务出错")
	}
}

func HttpsServerRun() {
	r := HttpsInit(gin.Recovery(), gin.Logger())

	HttpsSrvHandler = &http.Server{
		Addr:           conf.HttpsConf.Addr,
		Handler:        r,
		ReadTimeout:    time.Duration(conf.HttpsConf.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(conf.HttpsConf.WriteTimeout) * time.Second,
		MaxHeaderBytes: conf.HttpsConf.MaxHeaderBytes,
	}
	log.Printf(" [INFO] http_proxy_run %s\n", conf.HttpsConf.Addr)
	if err := HttpSrvHandler.ListenAndServeTLS(cert_file.Path("server.crt"), cert_file.Path("server.key")); err != nil && err != http.ErrServerClosed {
		fmt.Println("HTTPS服务出错")
	}
}

func HttpServerStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := HttpSrvHandler.Shutdown(ctx); err != nil {
		log.Printf(" [ERROR] http_proxy_stop err:%v\n", err)
	}
	log.Printf(" [INFO] http_proxy_stop %v stopped\n", conf.HttpConf.Addr)
}

func HttpsServerStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := HttpsSrvHandler.Shutdown(ctx); err != nil {
		log.Fatalf(" [ERROR] https_proxy_stop err:%v\n", err)
	}
	log.Printf(" [INFO] https_proxy_stop %v stopped\n", conf.HttpsConf.Addr)
}
