package job

//
import (
	"bytes"
	"context"
	"fmt"
	"net/http"

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

func (a *GetLuckListJob) Run() {
	ctx := context.Background()

	i := 0
	for i < 5 {
		i++
		err := func() error {
			inzoneUser, err := a.userRepo.GetByUUID(ctx, config.C.CookieUserUUID)
			if err != nil {
				log.Error(err)
				return err
			}
			var res [][]string
			res, err = inzone.GetLuckUsers(inzoneUser.Cookie)
			if err != nil {
				log.Error(err)
				return err
			}

			uuids := make([]string, 0)
			for _, v := range res {
				uuids = append(uuids, v...)
			}

			var inzoneUsers []*pb.InzoneUser
			inzoneUsers, err = a.userRepo.GetByUUIDs(ctx, uuids)

			inzoneUserM := make(map[string]*pb.InzoneUser)
			for _, v := range inzoneUsers {
				inzoneUserM[v.UUID] = v
			}

			luckInzoneUser := make([][]*pb.InzoneUser, 0)
			data := make([]*pb.InzoneUser, 0)
			for _, v := range res {
				for _, uuid := range v {
					if inzoneUserM[uuid] != nil {
						data = append(data, inzoneUserM[uuid])
					}
				}
				luckInzoneUser = append(luckInzoneUser, data)
			}

			var groups []*pb.InzoneUserGroup
			groups, err = a.userGroupRepo.List(ctx, &dao.InzoneUserGroupListReq{})
			if err != nil {
				log.Error(err)
				return err
			}

			for _, group := range groups {
				go a.pushLuckMsg(group, luckInzoneUser)
			}

			return nil
		}()
		if err != nil {
			log.Error(fmt.Sprintf("第%d次错误err: ", i), err)
			break
		}
	}

}

func (a *GetLuckListJob) pushLuckMsg(group *pb.InzoneUserGroup, luckInzoneUser [][]*pb.InzoneUser) {
	text := fmt.Sprintf("今日消息\n")
	for k, v := range luckInzoneUser {
		label := "每日"
		if k > 0 {
			label = fmt.Sprintf("系列%d", k)
		}
		text += fmt.Sprintf("%s  (%d人)\n", label, len(v))
		for _, user := range v {
			text += "   " + user.Name + "  " + user.Phone + "  " + user.Remark + "\n"
		}
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
