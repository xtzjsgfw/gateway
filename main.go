package main

import (
	"gateway/extend/conf"
	"gateway/extend/redis"
	"gateway/models"
	"gateway/router"
	"gateway/validator"
)

func main() {
	conf.Init()

	models.Init()

	redis.Init()

	validator.Init()

	router.Init()
}

