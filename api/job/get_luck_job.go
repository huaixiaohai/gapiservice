package job

//
import (
	"github.com/huaixiaohai/gapiservice/config"
	"github.com/huaixiaohai/gapiservice/dao"
)

func NewGetLuckListJob(
	userRepo *dao.InzoneUserRepo,
	userGroupRepo *dao.InzoneUserGroupRepo,
) *GetLuckListJob {
	ins := &GetLuckListJob{
		userRepo:      userRepo,
		userGroupRepo: userGroupRepo,
	}
	return ins
}

type GetLuckListJob struct {
	userRepo      *dao.InzoneUserRepo
	userGroupRepo *dao.InzoneUserGroupRepo
}

func (a *GetLuckListJob) GetSpec() string {
	return config.C.Cron.GetLuckListJob
}

func (a *GetLuckListJob) Run() {

}
