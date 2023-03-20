package api

import (
	"github.com/gin-gonic/gin"
	"github.com/huaixiaohai/gapiservice/pb"
)

type InzoneUserApi struct {
}

func (a *InzoneUserApi) Create(ctx *gin.Context, req *pb.InzoneUser) (*pb.ID, error) {
	panic("")
}

func (a *InzoneUserApi) Update(ctx *gin.Context, req *pb.InzoneUser) (*pb.ID, error) {
	panic("")
}

func (a *InzoneUserApi) Delete(ctx *gin.Context, req *pb.ID) (*pb.Empty, error) {
	panic("")
}

func (a *InzoneUserApi) Get(ctx *gin.Context, req *pb.ID) (*pb.InzoneUser, error) {
	panic("")
}

func (a *InzoneUserApi) List(ctx *gin.Context, req *pb.InzoneUserListReq) (*pb.InzoneUserListResp, error) {
	panic("")
}
