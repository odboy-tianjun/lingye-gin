package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"lingye-gin/src/config"
	"strings"
)

type GinEngine struct {
}

func (v *GinEngine) Start() {
	if strings.Compare(config.AppProps.App.Mode, "debug") == 0 {
		gin.SetMode(gin.DebugMode)
	} else if strings.Compare(config.AppProps.App.Mode, "release") == 0 {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.Default()
	engine.Use(config.LoggerToFile())
	// 设置路由
	new(GinRouter).Init(engine)
	config.Logger.Info("GinWebServer Init...")
	_ = engine.Run(fmt.Sprintf("0.0.0.0:%d", config.AppProps.Server.Port))
	config.Logger.Info("GinWebServer Starting...")
}
