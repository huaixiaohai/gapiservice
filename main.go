package main

import (
	"github.com/gin-gonic/gin"
	"github.com/huaixiaohai/gapiservice/config"
	"github.com/huaixiaohai/lib/log"
)

func init() {
	log.Init(&log.Cfg{
		Level:    config.C.Log.Level,
		Output:   config.C.Log.Output,
		FilePath: config.C.Log.FilePath,
	})
}

func main(){
	engine := gin.New()
	engine.Run(":8000")
}
