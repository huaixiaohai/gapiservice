package job

//
import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/huaixiaohai/gapiservice/pb"

	"github.com/huaixiaohai/lib/log"

	"github.com/huaixiaohai/gapiservice/config"
	"github.com/huaixiaohai/gapiservice/dao"
	"github.com/huaixiaohai/gapiservice/inzone"
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

type LuckInzoneUser struct {
	Label string
	Users []*pb.InzoneUser
}

func (a *GetLuckListJob) Run() {
	ctx := context.Background()
	inzoneUser, err := a.userRepo.GetByUUID(ctx, config.C.CookieUserUUID)
	if err != nil {
		log.Error(err)
		return
	}

	f := func(getLuck func(cookie string) ([]*inzone.Luck, error)) error {
		var lucks []*inzone.Luck
		i := 0
		for i < 5 {
			time.Sleep(time.Second)
			i++
			lucks, err = getLuck(inzoneUser.Cookie)
			if err != nil {
				log.Error(err)
				continue
			}
			break
		}
		uuids := make([]string, 0)
		for _, v := range lucks {
			uuids = append(uuids, v.UUIDs...)
		}

		var inzoneUsers []*pb.InzoneUser
		inzoneUsers, err = a.userRepo.GetByUUIDs(ctx, uuids)

		inzoneUserM := make(map[string]*pb.InzoneUser)
		for _, v := range inzoneUsers {
			inzoneUserM[v.UUID] = v
		}

		luckInzoneUsers := make([]*LuckInzoneUser, 0)
		data := make([]*pb.InzoneUser, 0)
		for _, v := range lucks {
			for _, uuid := range v.UUIDs {
				if inzoneUserM[uuid] != nil {
					data = append(data, inzoneUserM[uuid])
				}
			}
			luckInzoneUsers = append(luckInzoneUsers, &LuckInzoneUser{
				Label: v.Label,
				Users: data,
			})
		}

		var groups []*pb.InzoneUserGroup
		groups, err = a.userGroupRepo.List(ctx, &dao.InzoneUserGroupListReq{})
		if err != nil {
			log.Error(err)
			return err
		}

		for _, group := range groups {
			go a.pushLuckMsg(group, luckInzoneUsers)
		}
		return nil
	}

	_ = f(inzone.GetDailyLuckUsers)
	_ = f(inzone.GetSeriesLuckUsers)
}

func (a *GetLuckListJob) pushLuckMsg(group *pb.InzoneUserGroup, luckInzoneUsers []*LuckInzoneUser) {
	text := ""
	for _, v := range luckInzoneUsers {
		count := 0
		for _, user := range v.Users {
			if user.GroupID == group.ID {
				text += "   " + user.Name + "  " + user.Phone + "  " + user.Remark + "\n"
				count++
			}
		}
		text = fmt.Sprintf("%s  (%d人)\n", v.Label, count) + text
		text = fmt.Sprintf("今日消息\n") + text
	}
	log.Println(group.Name)
	log.Println(text)
	msg := fmt.Sprintf(`{
	   "msgtype": "text",
	   "text": {
	       "content": "%s"
	   }
	}`, text)

	r := bytes.NewBuffer([]byte(msg))
	_, err := http.Post(group.DingTalkHook, "application/json", r)
	if err != nil {
		fmt.Println("发送失败", err.Error())
	} else {
		fmt.Println("发送成功")
	}
}
