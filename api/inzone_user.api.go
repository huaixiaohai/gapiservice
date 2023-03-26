package api

import (
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
	return &InzoneUserApi{
		userRepo:      userRepo,
		userGroupRepo: userGroupRepo,
	}
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
	return &pb.ID{ID: req.ID}, a.userRepo.Create(ctx, req)
}

func (a *InzoneUserApi) Update(ctx *gin.Context, req *pb.InzoneUser) (*pb.Empty, error) {
	var err error
	req.UUID, err = pb.GetUUID(req.Name, req.Phone)
	if err != nil {
		return nil, err
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
