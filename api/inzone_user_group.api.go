package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/huaixiaohai/gapiservice/dao"
	"github.com/huaixiaohai/gapiservice/pb"
	"github.com/huaixiaohai/lib/snowflake"
)

var InzoneUserGroupApiSet = wire.NewSet(NewInzoneUserGroupApi)

func NewInzoneUserGroupApi(groupRepo *dao.InzoneUserGroupRepo) *InzoneUserGroupApi {
	return &InzoneUserGroupApi{
		groupRepo: groupRepo,
	}
}

type InzoneUserGroupApi struct {
	groupRepo *dao.InzoneUserGroupRepo
}

func (a *InzoneUserGroupApi) Create(ctx *gin.Context, req *pb.InzoneUserGroup) (*pb.ID, error) {
	req.ID = snowflake.MustID()
	return &pb.ID{ID: req.ID}, a.groupRepo.Create(ctx, req)
}

func (a *InzoneUserGroupApi) Update(ctx *gin.Context, req *pb.InzoneUserGroup) (*pb.Empty, error) {
	return &pb.Empty{}, a.groupRepo.Update(ctx, req)
}

func (a *InzoneUserGroupApi) Delete(ctx *gin.Context, req *pb.ID) (*pb.Empty, error) {
	return &pb.Empty{}, a.groupRepo.Delete(ctx, req.ID)
}

func (a *InzoneUserGroupApi) Get(ctx *gin.Context, req *pb.ID) (*pb.InzoneUserGroup, error) {
	return a.groupRepo.Get(ctx, req.ID)
}

func (a *InzoneUserGroupApi) List(ctx *gin.Context, req *pb.InzoneUserGroupListReq) (*pb.InzoneUserGroupListResp, error) {
	data, err := a.groupRepo.List(ctx, &dao.InzoneUserGroupListReq{})
	if err != nil {
		return nil, err
	}
	return &pb.InzoneUserGroupListResp{
		List: data,
	}, nil
}
