package dao

import (
	"context"

	"github.com/google/wire"
	"github.com/huaixiaohai/gapiservice/dao/model"
	"github.com/huaixiaohai/gapiservice/pb"
	"gorm.io/gorm"
)

var NovelRepoSet = wire.NewSet(NewNovelRepo)

func NewNovelRepo() *NovelRepo {
	return &NovelRepo{}
}

type NovelRepo struct {
}

func (a *NovelRepo) Create(ctx context.Context, one *pb.Novel) error {
	err := getSession(ctx).Model(&model.Novel{}).Create(model.NovelFrom(one)).Error
	if err != nil {
		return err
	}
	return nil
}

func (a *NovelRepo) Update(ctx context.Context, one *pb.Novel) error {
	err := getSession(ctx).Model(&model.Novel{}).Where("id=?", one.ID).Updates(model.NovelFrom(one)).Error
	if err != nil {
		return err
	}
	return nil
}

func (a *NovelRepo) Delete(ctx context.Context, id string) error {
	err := getSession(ctx).Model(&model.Novel{}).Where("id=?", id).Delete(&model.Novel{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (a *NovelRepo) Get(ctx context.Context, id string) (*pb.Novel, error) {
	one := &model.Novel{}
	err := getSession(ctx).Model(&model.Novel{}).Where("id=?", id).First(one).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return model.NovelTo(one), nil
}

type NovelListReq struct {
	PageSize  int64
	PageIndex int64
	Name      string
}

// List 返回任务列表，按照优先级排序
func (a *NovelRepo) List(ctx context.Context, req *NovelListReq) ([]*pb.Novel, error) {
	records := make([]*model.Novel, 0)
	err := a.listReq(ctx, req).Find(&records).Error
	if err != nil {
		return nil, err
	}
	return model.NovelListTo(records), nil
}

func (a *NovelRepo) Count(ctx context.Context, req *NovelListReq) (int64, error) {
	var count int64
	return count, a.listReq(ctx, req).Count(&count).Error
}

func (a *NovelRepo) listReq(ctx context.Context, req *NovelListReq) *gorm.DB {
	s := getSession(ctx).Model(&model.Novel{})
	return s.Order("id desc").Limit(int(req.PageSize)).Offset(int((req.PageIndex - 1) * req.PageSize))
}
