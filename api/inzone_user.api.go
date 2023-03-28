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
	data, err := a.userRepo.List(ctx, &dao.InzoneUserListReq{
		PageSize:     req.PageSize,
		PageIndex:    req.PageIndex,
		Name:         req.Name,
		Phone:        req.Phone,
		GroupID:      req.GroupID,
		CookieStatus: req.CookieStatus,
	})
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
	err = a.userRepo.UpdateCookie(ctx, cid, cookie.String(), pb.ECookieStatusValid)
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
		ids, err := a.userRepo.GetIDsByCookieStatus(ctx, pb.ECookieStatusValid)
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
