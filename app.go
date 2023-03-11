package main

import (
	"github.com/gin-gonic/gin"
	"github.com/huaixiaohai/lib/log"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	engine *gin.Engine
}

func (a *App) Run() {

}

func (a *App) WaitQuit() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	i := <-c
	log.Println("Received interrupt[%v], shutting down...", i)
}

func (a *App) registerRouter() {
	//a.engine.pos
}