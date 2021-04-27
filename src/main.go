package main

import (
	"lingye-gin/src/config"
	"lingye-gin/src/middleware"
)

func main() {
	// 初始化yaml配置
	new(config.ApplicationProperties).Init()
	// 初始化redis
	new(config.RedisPool).Init()
	// 初始化gin
	new(middleware.GinEngine).Start()
}
