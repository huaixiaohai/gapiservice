package dao

import (
	"context"

	"github.com/google/wire"
	"github.com/huaixiaohai/gapiservice/dao/model"
	"github.com/huaixiaohai/gapiservice/pb"
	"gorm.io/gorm"
)

var NovelResourceRepoSet = wire.NewSet(NewNovelResourceRepo)

func NewNovelResourceRepo() *NovelResourceRepo {
	return &NovelResourceRepo{}
}

type NovelResourceRepo struct {
}

func (a *NovelResourceRepo) Create(ctx context.Context, one *pb.NovelResource) error {
	err := getSession(ctx).Model(&model.NovelResource{}).Create(model.NovelResourceFrom(one)).Error
	if err != nil {
		return err
	}
	return nil
}

func (a *NovelResourceRepo) Update(ctx context.Context, one *pb.NovelResource) error {
	err := getSession(ctx).Model(&model.NovelResource{}).Where("id=?", one.ID).Updates(model.NovelResourceFrom(one)).Error
	if err != nil {
		return err
	}
	return nil
}

func (a *NovelResourceRepo) Delete(ctx context.Context, id string) error {
	err := getSession(ctx).Model(&model.NovelResource{}).Where("id=?", id).Delete(&model.NovelResource{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (a *NovelResourceRepo) Get(ctx context.Context, id string) (*pb.NovelResource, error) {
	one := &model.NovelResource{}
	err := getSession(ctx).Model(&model.NovelResource{}).Where("id=?", id).First(one).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return model.NovelResourceTo(one), nil
}

type NovelResourceListReq struct {
	PageSize  int64
	PageIndex int64
	Name      string
}

// List 返回任务列表，按照优先级排序
func (a *NovelResourceRepo) List(ctx context.Context, req *NovelResourceListReq) ([]*pb.NovelResource, error) {
	records := make([]*model.NovelResource, 0)
	err := a.listReq(ctx, req).Find(&records).Error
	if err != nil {
		return nil, err
	}
	return model.NovelResourceListTo(records), nil
}

func (a *NovelResourceRepo) Count(ctx context.Context, req *NovelResourceListReq) (int64, error) {
	var count int64
	return count, a.listReq(ctx, req).Count(&count).Error
}

func (a *NovelResourceRepo) listReq(ctx context.Context, req *NovelResourceListReq) *gorm.DB {
	s := getSession(ctx).Model(&model.NovelResource{})
	return s.Order("id desc").Limit(int(req.PageSize)).Offset(int((req.PageIndex - 1) * req.PageSize))
}
