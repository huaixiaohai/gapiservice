package dao

import (
	"context"
	"time"

	"github.com/huaixiaohai/lib/log"

	"github.com/google/wire"
	"github.com/huaixiaohai/gapiservice/dao/model"
	"github.com/huaixiaohai/gapiservice/pb"
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

func (a *InzoneUserRepo) Delete(ctx context.Context, id string) error {
	err := getSession(ctx).Model(&model.InzoneUser{}).Where("id=?", id).Delete(&model.InzoneUser{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (a *InzoneUserRepo) DeleteByCID(ctx context.Context, cid string) error {
	err := getSession(ctx).Model(&model.InzoneUser{}).Where("cid=?", cid).Delete(&model.InzoneUser{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (a *InzoneUserRepo) Get(ctx context.Context, id string) (*pb.InzoneUser, error) {
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

func (a *InzoneUserRepo) GetValidFirst(ctx context.Context) (*pb.InzoneUser, error) {
	one := &model.InzoneUser{}
	err := getSession(ctx).Model(&model.InzoneUser{}).Where("cookie_status=?", pb.ECookieStatusValid).First(one).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return model.InzoneUserTo(one), nil
}

func (a *InzoneUserRepo) GetByUUID(ctx context.Context, uuid string) (*pb.InzoneUser, error) {
	one := &model.InzoneUser{}
	err := getSession(ctx).Model(&model.InzoneUser{}).Where("uuid=?", uuid).First(one).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return model.InzoneUserTo(one), nil
}

func (a *InzoneUserRepo) GetByUUIDs(ctx context.Context, uuids []string) ([]*pb.InzoneUser, error) {
	records := make([]*model.InzoneUser, 0)
	err := getSession(ctx).Model(&model.InzoneUser{}).Where("uuid in ?", uuids).Find(&records).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return model.InzoneUserListTo(records), nil
}

type InzoneUserListReq struct {
	PageSize     int64
	PageIndex    int64
	Name         string
	Phone        string
	GroupID      string
	CookieStatus pb.ECookieStatus
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
	if req.Name != "" {
		s.Where("phone like ?", "%"+req.Phone+"%")
	}
	if req.GroupID != "" {
		s.Where("group_id = ?", req.GroupID)
	}
	if req.CookieStatus != pb.ECookieStatusNone {
		s.Where("cookie_status = ?", req.CookieStatus)
	}
	s.Order("id desc")
	if req.PageSize != 0 {
		s.Limit(int(req.PageSize)).Offset(int((req.PageIndex - 1) * req.PageSize))
	}
	return s
}

//func (a *InzoneUserRepo) GetUsers(ctx context.Context) ([]*pb.InzoneUser, error) {
//	records := make([]*model.InzoneUser, 0)
//	err := getSession(ctx).Model(&model.InzoneUser{}).Where("? - UNIX_TIMESTAMP(refresh_cookie_at) > ?", time.Now().Local().Unix(), 2400).Find(&records).Error
//	if err == gorm.ErrRecordNotFound {
//		return nil, nil
//	}
//	if err != nil {
//		return nil, err
//	}
//	return model.InzoneUserListTo(records), nil
//}

func (a *InzoneUserRepo) GetIDsByCookieStatus(ctx context.Context, cookieStatus pb.ECookieStatus) ([]string, error) {
	ids := make([]string, 0)
	err := getSession(ctx).Model(&model.InzoneUser{}).Select("id").Where("cookie_status=?", cookieStatus).Order("cookie_refresh_at").Find(&ids).Error
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return ids, err
}

func (a *InzoneUserRepo) GetLRUCookieUser(ctx context.Context) (*pb.InzoneUser, error) {
	one := &model.InzoneUser{}
	err := getSession(ctx).Model(&model.InzoneUser{}).Where("cookie_status=?", pb.ECookieStatusValid).Order("cookie_refresh_at").First(one).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return model.InzoneUserTo(one), nil
}

func (a *InzoneUserRepo) UpdateCookie(ctx context.Context, cid string, cookie string, cookieStatus pb.ECookieStatus) error {
	if cookieStatus == pb.ECookieStatusValid {
		err := getSession(ctx).Model(&model.InzoneUser{}).Where("cid=?", cid).Updates(map[string]interface{}{"cookie_refresh_at": time.Now().Local(), "cookie_status": cookieStatus, "cookie": cookie}).Error
		if err != nil {
			return err
		}
	} else {
		err := getSession(ctx).Model(&model.InzoneUser{}).Where("cid=?", cid).Updates(map[string]interface{}{"cookie_status": pb.ECookieStatusInvalid}).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *InzoneUserRepo) UpdatePhone(ctx context.Context, phone, name, cid string) error {
	return getSession(ctx).Model(&model.InzoneUser{}).Where("cid=?", cid).Updates(map[string]interface{}{"phone": phone, "name": name}).Error

}
