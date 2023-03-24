// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package main

import (
	"github.com/huaixiaohai/gapiservice/api"
	"github.com/huaixiaohai/gapiservice/dao"
)

// Injectors from wire.go:

// GetAppInstance 生成注入器
func GetAppInstance() (*App, error) {
	inzoneUserGroupRepo := dao.NewInzoneUserGroupRepo()
	inzoneUserGroupApi := api.NewInzoneUserGroupApi(inzoneUserGroupRepo)
	inzoneUserRepo := dao.NewInzoneUserRepo()
	inzoneUserApi := api.NewInzoneUserApi(inzoneUserRepo)
	userApi := api.NewUserApi()
	app := NewApp(inzoneUserGroupApi, inzoneUserApi, userApi)
	return app, nil
}
