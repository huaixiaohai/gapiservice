package job

//
import (
	"github.com/huaixiaohai/gapiservice/config"
	"github.com/huaixiaohai/gapiservice/dao"
)

func NewRefreshCookieJob(
	userRepo *dao.InzoneUserRepo,
	userGroupRepo *dao.InzoneUserGroupRepo,
) *RefreshCookieJob {
	ins := &RefreshCookieJob{
		userRepo:      userRepo,
		userGroupRepo: userGroupRepo,
	}
	return ins
}

type RefreshCookieJob struct {
	userRepo      *dao.InzoneUserRepo
	userGroupRepo *dao.InzoneUserGroupRepo
}

func (a *RefreshCookieJob) GetSpec() string {
	return config.C.Cron.RefreshCookieJob
}

func (a *RefreshCookieJob) Run() {
	//a.userRepo.GetByUUID()
}
