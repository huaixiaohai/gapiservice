package api

import (
	"fmt"

	"github.com/huaixiaohai/gapiservice/api/job"

	"github.com/google/wire"
	"github.com/huaixiaohai/gapiservice/dao"
	"github.com/robfig/cron/v3"
)

var InzoneApiSet = wire.NewSet(NewInzoneApi)

func NewInzoneApi(
	userRepo *dao.InzoneUserRepo,
	userGroupRepo *dao.InzoneUserGroupRepo,
) *InzoneApi {
	ins := &InzoneApi{
		userRepo:      userRepo,
		userGroupRepo: userGroupRepo,
	}
	ins.init()
	return ins
}

type InzoneApi struct {
	c *cron.Cron

	userRepo      *dao.InzoneUserRepo
	userGroupRepo *dao.InzoneUserGroupRepo
}

func (a *InzoneApi) init() {
	a.c = cron.New(cron.WithSeconds())
	// 任务列表
	jobs := []job.IJob{
		job.NewGetLuckListJob(a.userRepo, a.userGroupRepo), // 每天获取获奖名单
	}

	for _, v := range jobs {
		_, err := a.c.AddJob(v.GetSpec(), v)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
	}
	a.c.Start()
}
