package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/huaixiaohai/gapiservice/api"
	"github.com/huaixiaohai/gapiservice/auth"
	"github.com/huaixiaohai/gapiservice/config"
	"github.com/huaixiaohai/lib/log"
)

func NewApp(
	inzoneUserGroupApi *api.InzoneUserGroupApi,
	userApi *api.UserApi,
) *App {
	engine := gin.New()
	server := &http.Server{
		Addr:              config.C.Http.Addr,
		Handler:           engine,
		TLSConfig:         nil,
		ReadTimeout:       0,
		ReadHeaderTimeout: 0,
		WriteTimeout:      0,
		IdleTimeout:       0,
		MaxHeaderBytes:    0,
		TLSNextProto:      nil,
		ConnState:         nil,
		ErrorLog:          nil,
		BaseContext:       nil,
		ConnContext:       nil,
	}
	app := &App{
		engine:             engine,
		server:             server,
		inzoneUserGroupApi: inzoneUserGroupApi,
		userApi:            userApi,
	}
	return app
}

type App struct {
	engine *gin.Engine
	server *http.Server

	userApi            *api.UserApi
	inzoneUserGroupApi *api.InzoneUserGroupApi
}

func (a *App) Run() {
	a.registerRouter()
	go func() {
		var err error
		if config.C.Http.CertFile != "" && config.C.Http.KeyFile != "" {
			a.server.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
			err = a.server.ListenAndServeTLS(config.C.Http.CertFile, config.C.Http.KeyFile)
		} else {
			err = a.server.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
}

func (a *App) WaitQuit() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	i := <-c
	log.Println("Received interrupt[%v], shutting down...", i)
}

func (a *App) registerRouter() {
	g := a.engine
	g.NoRoute()
	g.NoMethod()
	g.Use(gin.Logger())
	g.Use(gin.Recovery())

	g.POST("/api/v1/login", wrapper(a.userApi.Login))
	g.GET("/api/v1/user/get", userAuthMiddleware(), wrapper(a.userApi.Get))

	g.POST("/api/v1/inzone/user_group/create", userAuthMiddleware(), wrapper(a.inzoneUserGroupApi.Create))
	g.POST("/api/v1/inzone/user_group/update", userAuthMiddleware(), wrapper(a.inzoneUserGroupApi.Update))
	g.GET("/api/v1/inzone/user_group/get", userAuthMiddleware(), wrapper(a.inzoneUserGroupApi.Get))
	g.GET("/api/v1/inzone/user_group/list", userAuthMiddleware(), wrapper(a.inzoneUserGroupApi.List))
	g.DELETE("/api/v1/inzone/user_group/delete", userAuthMiddleware(), wrapper(a.inzoneUserGroupApi.Delete))
}

func wrapper(f interface{}) func(*gin.Context) {
	fc := reflect.ValueOf(f)
	typ := fc.Type()
	if typ.Kind() != reflect.Func {
		log.Panicf("not function")
	}
	if typ.NumIn() != 2 {
		log.Panicf("number of params not equels to 2")
	}
	if typ.In(0).String() != "*gin.Context" {
		log.Panicf("first parameter should be of type *gin.Context")
	}
	if typ.NumOut() != 2 {
		log.Panicf("number of return values not equels to 2")
	}
	fullName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	strs := strings.Split(fullName, ".")
	viewName := strs[len(strs)-1]

	return func(c *gin.Context) {
		c.Set("view_name", viewName)

		//traceID := c.GetHeader(traceHeader)
		//c.Set(traceHeader, traceID)
		//c.Header(traceHeader, traceID)

		//_ = c.Request.ParseForm()
		//for _, p := range c.Params {
		//	c.Request.Form[p.Key] = []string{p.Value}
		//}

		req := reflect.New(typ.In(1).Elem()).Interface()
		c.Set("request", req)
		err := c.ShouldBind(req)
		if err != nil {
			// 参数不对处理
			return
		}
		in := []reflect.Value{
			reflect.ValueOf(c),
			reflect.ValueOf(req),
		}
		res := fc.Call(in)
		if !res[1].IsNil() {
			//var err error
			//switch e := res[1].Interface().(type) {
			//case *errors.ResponseError:
			//	err = res[1].Interface().(error)
			//	break
			//case error:
			//	err = errors.New500Response(e.Error())
			//	break
			//default:
			//	err = errors.New500Response("未知错误")
			//	break
			//}

			Res(c, nil, res[1].Interface().(error))
			return
		}

		//// 重定向相应
		//if typ.Out(0).String() == "*RedirectResponse" {
		//	c.Redirect(302, res[0].Elem().FieldByName("Url").String())
		//	return
		//}
		Res(c, res[0].Interface(), nil)
	}
}

func Res(c *gin.Context, resp interface{}, err error) {
	type Body struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Result  interface{} `json:"result"`
	}

	body := &Body{
		Result: resp,
	}

	if err != nil {
		body.Message = err.Error()
		if body.Message == "auth" {
			body.Code = 401
		} else {
			body.Code = 500
		}
	} else {
		body.Code = 0
	}

	buf, e := json.Marshal(body)
	if e != nil {
		log.Error(e)
	}

	//c.Set(ResBodyKey, buf)
	c.Data(body.Code, "application/json; charset=utf-8", buf)
	c.Abort()
}

func userAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		var token string
		authorization := c.GetHeader("Authorization")
		prefix := "Bearer "
		if authorization != "" && strings.HasPrefix(authorization, prefix) {
			token = authorization[len(prefix):]
		}

		userID, userName, expireAt, err := auth.ParseToken(token)
		if err != nil || expireAt <= time.Now().Local().Unix() {
			Res(c, nil, errors.New("auth"))
			return
		}
		println(userID, userName)
		c.Set("UserID", userID)
		c.Next()
	}
}
