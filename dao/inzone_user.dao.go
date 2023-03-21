package dao

import (
	"context"

	"github.com/google/wire"
	"github.com/huaixiaohai/gapiservice/dao/model"
	"github.com/huaixiaohai/gapiservice/pb"
	"github.com/huaixiaohai/lib/log"

	"gorm.io/gorm"
)

var InzoneUserRepoSet = wire.NewSet(NewInzoneUserRepo)

func NewInzoneUserRepo() *InzoneUserRepo {
	return &InzoneUserRepo{}
}

type InzoneUserRepo struct {
}

func (a *InzoneUserRepo) Create(ctx context.Context, one *pb.InzoneUser) error {
	err := getSession(ctx).Model(&model.InzoneUser{}).Create(model.InzoneUserFrom(one)).Error
	if err != nil {
		return err
	}
	return nil
}

func (a *InzoneUserRepo) Update(ctx context.Context, one *pb.InzoneUser) error {
	err := getSession(ctx).Model(&model.InzoneUser{}).Where("id=?", one.ID).Updates(model.InzoneUserFrom(one)).Error
	if err != nil {
		return err
	}
	return nil
}

func (a *InzoneUserRepo) Delete(ctx context.Context, id uint64) error {
	log.Error("待实现")
	return nil
}

func (a *InzoneUserRepo) Get(ctx context.Context, id uint64) (*pb.InzoneUser, error) {
	one := &model.InzoneUser{}
	err := getSession(ctx).Model(&model.InzoneUser{}).Where("id=?", id).First(one).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return model.InzoneUserTo(one), nil
}

type InzoneUserListReq struct {
	PageSize  int64
	PageIndex int64
	Name      string
}

// List 返回任务列表，按照优先级排序
func (a *InzoneUserRepo) List(ctx context.Context, req *InzoneUserListReq) ([]*pb.InzoneUser, error) {
	records := make([]*model.InzoneUser, 0)
	err := a.listReq(ctx, req).Find(&records).Error
	if err != nil {
		return nil, err
	}
	return model.InzoneUserListTo(records), nil
}

func (a *InzoneUserRepo) Count(ctx context.Context, req *InzoneUserListReq) (int64, error) {
	var count int64
	return count, a.listReq(ctx, req).Count(&count).Error
}

func (a *InzoneUserRepo) listReq(ctx context.Context, req *InzoneUserListReq) *gorm.DB {
	s := getSession(ctx).Model(&model.InzoneUser{})
	if req.Name != "" {
		s.Where("name like ?", "%"+req.Name+"%")
	}
	return s.Order("id desc").Limit(int(req.PageSize)).Offset(int((req.PageIndex - 1) * req.PageSize))
}
