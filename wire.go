//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/google/wire"
	"github.com/huaixiaohai/gapiservice/api"
	"github.com/huaixiaohai/gapiservice/dao"
)

// GetAppInstance 生成注入器
func GetAppInstance() (*App, error) {
	wire.Build(
		NewApp,
		api.Set,
		dao.Set,
	)
	return &App{}, nil
}
