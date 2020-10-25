package main

import (
	"flag"
	"gateway/extend/conf"
	"gateway/extend/redis"
	"gateway/models"
	"gateway/router"
	"gateway/router/http_proxy"
	"gateway/validator"
	"os"
	"os/signal"
	"syscall"
)

//endpoint dashboard后台管理  server代理服务器
//config ./conf/prod/ 对应配置文件夹
var (
	endpoint = flag.String("endpoint", "", "input endpoint dashboard or server")
	//config = flag.String("config", "", "input config file like ./config/dev")
)

func main() {
	flag.Parse()

	if *endpoint == "" {
		flag.Usage()
		os.Exit(1)
	}
	if *endpoint == "dashboard" {
		conf.Init()

		models.Init()

		redis.Init()

		validator.Init()

		router.Init()
	} else if *endpoint == "server" {
		conf.Init()

		models.Init()

		redis.Init()

		validator.Init()

		models.ServiceManagerHandler.LoadOnce()

		go func() {
			http_proxy.HttpServerRun()
		}()
		go func() {
			http_proxy.HttpsServerRun()
		}()

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		http_proxy.HttpServerStop()
		http_proxy.HttpsServerStop()
	}
}
