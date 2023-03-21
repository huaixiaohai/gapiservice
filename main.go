package main

import (
	"github.com/huaixiaohai/gapiservice/config"
	"github.com/huaixiaohai/gapiservice/dao"
	"github.com/huaixiaohai/lib/log"
)

func init() {
	log.Init(&log.Cfg{
		Level:    config.C.Log.Level,
		Output:   config.C.Log.Output,
		FilePath: config.C.Log.FilePath,
	})
	err := dao.Init()
	if err != nil {
		panic(err)
	}
}

func main() {
	app, err := GetAppInstance()
	if err != nil {
		panic(err)
	}
	app.Run()
	app.WaitQuit()
}
