package main

import (
	"github.com/google/wire"
	"github.com/huaixiaohai/gapiservice/api"
	"github.com/huaixiaohai/gapiservice/dao"
)

// BuildInjector 生成注入器
func GetAppInstance() (*App, error) {
	wire.Build(
		NewApp,
		api.Set,
		dao.Set,
	)
	return &App{}, nil
}
