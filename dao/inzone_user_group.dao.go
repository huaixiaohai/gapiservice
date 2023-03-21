package dao

import (
	"context"

	"github.com/google/wire"
	"github.com/huaixiaohai/gapiservice/dao/model"
	"github.com/huaixiaohai/gapiservice/pb"
	"gorm.io/gorm"
)

var InzoneUserGroupRepoSet = wire.NewSet(NewInzoneUserGroupRepo)

func NewInzoneUserGroupRepo() *InzoneUserGroupRepo {
	return &InzoneUserGroupRepo{}
}

type InzoneUserGroupRepo struct {
}

func (a *InzoneUserGroupRepo) Create(ctx context.Context, one *pb.InzoneUserGroup) error {
	err := getSession(ctx).Model(&model.InzoneUserGroup{}).Create(model.InzoneUserGroupFrom(one)).Error
	if err != nil {
		return err
	}
	return nil
}

func (a *InzoneUserGroupRepo) Update(ctx context.Context, one *pb.InzoneUserGroup) error {
	err := getSession(ctx).Model(&model.InzoneUserGroup{}).Where("id=?", one.ID).Updates(model.InzoneUserGroupFrom(one)).Error
	if err != nil {
		return err
	}
	return nil
}

func (a *InzoneUserGroupRepo) Delete(ctx context.Context, id string) error {
	err := getSession(ctx).Model(&model.InzoneUserGroup{}).Where("id=?", id).Delete(&model.InzoneUserGroup{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (a *InzoneUserGroupRepo) Get(ctx context.Context, id string) (*pb.InzoneUserGroup, error) {
	one := &model.InzoneUserGroup{}
	err := getSession(ctx).Model(&model.InzoneUserGroup{}).Where("id=?", id).First(one).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return model.InzoneUserGroupTo(one), nil
}

type InzoneUserGroupListReq struct {
	PageSize  int64
	PageIndex int64
	Name      string
}

// List 返回任务列表，按照优先级排序
func (a *InzoneUserGroupRepo) List(ctx context.Context, req *InzoneUserGroupListReq) ([]*pb.InzoneUserGroup, error) {
	records := make([]*model.InzoneUserGroup, 0)
	err := a.listReq(ctx, req).Find(&records).Error
	if err != nil {
		return nil, err
	}
	return model.InzoneUserGroupListTo(records), nil
}

func (a *InzoneUserGroupRepo) Count(ctx context.Context, req *InzoneUserGroupListReq) (int64, error) {
	var count int64
	return count, a.listReq(ctx, req).Count(&count).Error
}

func (a *InzoneUserGroupRepo) listReq(ctx context.Context, req *InzoneUserGroupListReq) *gorm.DB {
	s := getSession(ctx).Model(&model.InzoneUserGroup{})
	if req.Name != "" {
		s.Where("name like ?", "%"+req.Name+"%")
	}
	return s.Order("id desc").Limit(int(req.PageSize)).Offset(int((req.PageIndex - 1) * req.PageSize))
}
