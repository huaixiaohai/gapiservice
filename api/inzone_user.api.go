package api

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/huaixiaohai/gapiservice/config"

	"github.com/robfig/cron/v3"

	"github.com/huaixiaohai/gapiservice/inzone"

	"github.com/huaixiaohai/lib/log"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/huaixiaohai/gapiservice/dao"
	"github.com/huaixiaohai/gapiservice/pb"
	"github.com/huaixiaohai/lib/snowflake"
)

var InzoneUserApiSet = wire.NewSet(NewInzoneUserApi)

func NewInzoneUserApi(
	userRepo *dao.InzoneUserRepo,
	userGroupRepo *dao.InzoneUserGroupRepo,
) *InzoneUserApi {
	ins := &InzoneUserApi{
		userRepo:      userRepo,
		userGroupRepo: userGroupRepo,
	}

	go ins.refreshCookie(context.Background())

	c := cron.New(cron.WithSeconds())
	fmt.Println(config.C.Cron, config.C.Cron.GetLuckUserJob)
	_, err := c.AddFunc(config.C.Cron.GetLuckUserJob, ins.LuckJob)
	if err != nil {
		panic(err)
	}
	c.Start()
	return ins
}

type InzoneUserApi struct {
	userRepo      *dao.InzoneUserRepo
	userGroupRepo *dao.InzoneUserGroupRepo
}

func (a *InzoneUserApi) Create(ctx *gin.Context, req *pb.InzoneUser) (*pb.ID, error) {
	req.ID = snowflake.MustID()
	var err error
	req.UUID, err = pb.GetUUID(req.Name, req.Phone)
	if err != nil {
		return nil, err
	}

	req.CookieStatus = pb.ECookieStatusInvalid

	return &pb.ID{ID: req.ID}, a.userRepo.Create(ctx, req)
}

func (a *InzoneUserApi) Update(ctx *gin.Context, req *pb.InzoneUser) (*pb.Empty, error) {
	user, err := a.userRepo.Get(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return &pb.Empty{}, nil
	}

	user.UUID, err = pb.GetUUID(req.Name, req.Phone)
	if err != nil {
		return nil, err
	}

	user.Name = req.Name
	user.GroupID = req.GroupID
	user.Phone = req.Phone
	user.CID = req.CID

	return &pb.Empty{}, a.userRepo.Update(ctx, req)
}

func (a *InzoneUserApi) Delete(ctx *gin.Context, req *pb.ID) (*pb.Empty, error) {
	return &pb.Empty{}, a.userRepo.Delete(ctx, req.ID)
}

func (a *InzoneUserApi) Get(ctx *gin.Context, req *pb.ID) (*pb.InzoneUser, error) {
	return a.userRepo.Get(ctx, req.ID)
}

func (a *InzoneUserApi) List(ctx *gin.Context, req *pb.InzoneUserListReq) (*pb.InzoneUserListResp, error) {

	daoReq := &dao.InzoneUserListReq{
		Name:         req.Name,
		Phone:        req.Phone,
		GroupID:      req.GroupID,
		CookieStatus: req.CookieStatus,
	}
	total, err := a.userRepo.Count(ctx, daoReq)
	if err != nil {
		return nil, err
	}

	daoReq.PageSize = req.PageSize
	daoReq.PageIndex = req.PageIndex
	data, err := a.userRepo.List(ctx, daoReq)
	if err != nil {
		return nil, err
	}
	for _, v := range data {
		group, err := a.userGroupRepo.Get(ctx, v.GroupID)
		if err != nil {
			return nil, err
		}
		if group != nil {
			v.GroupName = group.Name
		}
	}
	return &pb.InzoneUserListResp{
		List:  data,
		Total: total,
	}, nil
}

func (a *InzoneUserApi) UpdateCookie(ctx *gin.Context, req *pb.Empty) (*pb.Empty, error) {
	cookie, err := ctx.Request.Cookie("PHPSESSID")
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	var cid string
	cid, err = inzone.GetCID(cookie.String())
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	time.Sleep(time.Millisecond * 100)
	phone, err := inzone.GetPhone(cookie.String())
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	time.Sleep(time.Millisecond * 100)
	name, err := inzone.GetUserName(cookie.String())
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	inzoneUser, err := a.userRepo.Get(ctx, cid)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	if inzoneUser == nil {
		err = a.userRepo.Create(ctx, &pb.InzoneUser{
			ID:              cid,
			Name:            name,
			Phone:           phone,
			Remark:          "",
			GroupID:         "",
			GroupName:       "",
			Cookie:          cookie.String(),
			CookieRefreshAt: time.Now().Local().Unix(),
			CookieStatus:    pb.ECookieStatusValid,
			UUID:            "",
			CID:             cid,
		})
	} else {
		err = a.userRepo.Update(ctx, &pb.InzoneUser{
			ID:              cid,
			Name:            name,
			Phone:           phone,
			Remark:          "",
			GroupID:         "",
			GroupName:       "",
			Cookie:          cookie.String(),
			CookieRefreshAt: time.Now().Local().Unix(),
			CookieStatus:    pb.ECookieStatusValid,
			UUID:            "",
			CID:             cid,
		})
	}
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return &pb.Empty{}, nil
}

// 刷新cookie
func (a *InzoneUserApi) refreshCookie(ctx context.Context) {
	sleepTime := time.Second
	for {
		time.Sleep(sleepTime)
		startTime := time.Now().Local().Unix()
		ids, err := a.userRepo.GetIDsByCookieStatus(ctx, pb.ECookieStatusInvalid)
		if err != nil {
			log.Error(err)
			sleepTime = time.Minute
			continue
		}
		if len(ids) <= 0 {
			sleepTime = time.Minute
			continue
		}

		for _, id := range ids {
			inzoneUser, err := a.userRepo.Get(ctx, id)
			if err != nil {
				log.Error(err)
				continue
			}
			if inzoneUser == nil {
				continue
			}
			cookieStatus := pb.ECookieStatusInvalid
			if inzone.IsValid(inzoneUser.Cookie) {
				cookieStatus = pb.ECookieStatusValid
			}
			err = a.userRepo.UpdateCookie(ctx, inzoneUser.CID, inzoneUser.Cookie, cookieStatus)
			if err != nil {
				log.Error(err)
				continue
			}

			time.Sleep(time.Millisecond * 300)

		}
		println(startTime + 2400 - time.Now().Local().Unix())
		sleepTime = time.Duration(startTime+3000-time.Now().Local().Unix()) * time.Second
	}
}

func (a *InzoneUserApi) LuckJob() {
	users, err := a.userRepo.List(context.Background(), &dao.InzoneUserListReq{
		CookieStatus: pb.ECookieStatusValid,
	})
	if err != nil {
		log.Error(err)
		return
	}
	for _, user := range users {
		b, err := inzone.IsLuck(user.Cookie)
		if err != nil {
			log.Error(err)
			continue
		}
		if b {
			pushLuckMsg("", user.Name, user.Phone)
		}
	}
}

func pushLuckMsg(hook, name, phone string) {
	hook = "https://oapi.dingtalk.com/robot/send?access_token=014a2ccdb00864a4db8fdc3f63b507b4cb3e8bde3b6d94cbf7711c4e25dacf69"

	text := fmt.Sprintf("中奖用户:%s, %s \n", phone, name)

	fmt.Println(text)
	//for _, v := range records {
	//	text += fmt.Sprintf("%s  (%d人)\n", v.Goods, len(v.Users))
	//	for _, user := range v.Users {
	//		text += "   " + user.Name + "  " + user.Phone + "  " + user.Remark + "\n"
	//	}
	//}

	msg := fmt.Sprintf(`{
	   "msgtype": "text",
	   "text": {
	       "content": "%s"
	   }
	}`, text)

	r := bytes.NewBuffer([]byte(msg))
	_, err := http.Post(hook, "application/json", r)
	if err != nil {
		fmt.Println("发送失败", err.Error())
	} else {
		fmt.Println("发送成功")
	}
}
