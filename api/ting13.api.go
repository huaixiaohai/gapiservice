package api

import (
	"context"

	"github.com/huaixiaohai/lib/log"

	"github.com/huaixiaohai/lib/snowflake"

	"github.com/huaixiaohai/gapiservice/dao"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/huaixiaohai/gapiservice/pb"
)

var Ting13ApiSet = wire.NewSet(NewTing13Api)

func NewTing13Api() *Ting13Api {
	ins := &Ting13Api{}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return ins
}

type Ting13Api struct {
	novelRepo         *dao.NovelRepo
	novelResourceRepo *dao.NovelResourceRepo
}

func (a *Ting13Api) Download(ctx *gin.Context, req *pb.ID) (*pb.Empty, error) {

}

func (a *Ting13Api) Create(ctx *gin.Context, req *pb.ID) (*pb.Empty, error) {
	return &pb.Empty{}, transaction(ctx, func(ctx context.Context) error {
		novel := &pb.Novel{
			ID:   snowflake.MustID(),
			Name: "万族之劫",
			Desc: "一种侃侃演播",
		}
		err := a.novelRepo.Create(ctx, novel)
		if err != nil {
			log.Error(err)
			return err
		}

		return nil
	})
}

func (a *Ting13Api) List(ctx *gin.Context, req *pb.ID) (*pb.Empty, error) {
	//"https: //www.ting13.com/youshengxiaoshuo/15687/45.html"

}
