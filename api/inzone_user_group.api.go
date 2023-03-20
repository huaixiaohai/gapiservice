package api

import (
	"github.com/gin-gonic/gin"
	"github.com/huaixiaohai/gapiservice/pb"
)

type InzoneUserGroupApi struct {
}

func (a *InzoneUserGroupApi) Create(ctx *gin.Context, req *pb.InzoneUserGroup) (*pb.ID, error) {
	panic("")
}

func (a *InzoneUserGroupApi) Update(ctx *gin.Context, req *pb.InzoneUserGroup) (*pb.ID, error) {
	panic("")
}

func (a *InzoneUserGroupApi) Delete(ctx *gin.Context, req *pb.ID) (*pb.Empty, error) {
	panic("")
}

func (a *InzoneUserGroupApi) Get(ctx *gin.Context, req *pb.ID) (*pb.InzoneUserGroup, error) {
	panic("")
}

func (a *InzoneUserGroupApi) List(ctx *gin.Context, req *pb.InzoneUserGroupListReq) (*pb.InzoneUserGroupListResp, error) {
	panic("")
}
