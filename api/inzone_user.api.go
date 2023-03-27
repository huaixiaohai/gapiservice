package api

import (
	"context"
	"time"

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

	if inzone.IsValid(req.Cookie) {
		req.CookieRefreshAt = time.Now().Local().Unix()
		req.CookieStatus = pb.ECookieStatusValid
	} else {
		req.CookieStatus = pb.ECookieStatusInvalid
	}

	return &pb.ID{ID: req.ID}, a.userRepo.Create(ctx, req)
}

func (a *InzoneUserApi) Update(ctx *gin.Context, req *pb.InzoneUser) (*pb.Empty, error) {
	var err error
	req.UUID, err = pb.GetUUID(req.Name, req.Phone)
	if err != nil {
		return nil, err
	}

	if inzone.IsValid(req.Cookie) {
		req.CookieRefreshAt = time.Now().Local().Unix()
		req.CookieStatus = pb.ECookieStatusValid
	} else {
		req.CookieStatus = pb.ECookieStatusInvalid
	}

	return &pb.Empty{}, a.userRepo.Update(ctx, req)
}

func (a *InzoneUserApi) Delete(ctx *gin.Context, req *pb.ID) (*pb.Empty, error) {
	return &pb.Empty{}, a.userRepo.Delete(ctx, req.ID)
}

func (a *InzoneUserApi) Get(ctx *gin.Context, req *pb.ID) (*pb.InzoneUser, error) {
	return a.userRepo.Get(ctx, req.ID)
}

func (a *InzoneUserApi) List(ctx *gin.Context, req *pb.InzoneUserListReq) (*pb.InzoneUserListResp, error) {
	data, err := a.userRepo.List(ctx, &dao.InzoneUserListReq{})
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
		List: data,
	}, nil
}

// 刷新cookie
func (a *InzoneUserApi) refreshCookie(ctx context.Context) {
	sleepTime := time.Second
	for {
		time.Sleep(sleepTime)
		count, err := a.userRepo.Count(ctx, &dao.InzoneUserListReq{})
		if err != nil {
			log.Error(err)
			sleepTime = time.Second
			continue
		}
		if count == 0 {
			sleepTime = time.Minute
			continue
		}
		n := count/3600 + 1
		sleepTime = time.Millisecond * time.Duration(1000/n)

		inzoneUser, err := a.userRepo.GetLRUCookieUser(ctx)
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
		err = a.userRepo.UpdateCookie(ctx, inzoneUser.ID, cookieStatus)
		if err != nil {
			log.Error(err)
			continue
		}
	}
}

//func (a *InzoneUserApi) RefreshCookie
