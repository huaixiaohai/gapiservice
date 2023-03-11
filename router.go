package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"reflect"
	"runtime"
	"strings"
)

func route(engine *gin.Engine) {
	engine.GET("/api/v1/log")
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
		if err := ginx.Parse(c, req); err != nil {
			//if err := c.ShouldBind(req); err != nil {
			ginx.Res(c, nil, err)
			return
		}
		in := []reflect.Value{
			reflect.ValueOf(c),
			reflect.ValueOf(req),
		}
		res := fc.Call(in)
		if !res[1].IsNil() {
			var err error
			switch e := res[1].Interface().(type) {
			case *errors.ResponseError:
				err = res[1].Interface().(error)
				break
			case error:
				err = errors.New500Response(e.Error())
				break
			default:
				err = errors.New500Response("未知错误")
				break
			}

			ginx.Res(c, nil, err)
			return
		}

		//// 重定向相应
		//if typ.Out(0).String() == "*RedirectResponse" {
		//	c.Redirect(302, res[0].Elem().FieldByName("Url").String())
		//	return
		//}
		ginx.Res(c, res[0].Interface(), nil)
	}
}