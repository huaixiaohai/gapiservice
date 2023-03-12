package main

import (
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
	app := NewApp()
	app.Run()
	app.WaitQuit()
}
