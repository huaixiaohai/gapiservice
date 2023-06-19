package model

import (
	"time"

	"github.com/huaixiaohai/gapiservice/pb"
)

type InzoneUser struct {
	Model
	Name            string `gorm:"type:varchar(30)"`
	Phone           string `gorm:"type:varchar(30)"`
	Remark          string
	GroupID         string
	CID             string `gorm:"column:cid;type:varchar(50);uniqueIndex"`
	Cookie          string
	CookieRefreshAt time.Time
	CookieStatus    pb.ECookieStatus
}

func InzoneUserFrom(one *pb.InzoneUser) *InzoneUser {
	if one == nil {
		return nil
	}

	return &InzoneUser{
		Model: Model{
			ID: one.ID,
		},
		Name:            one.Name,
		Phone:           one.Phone,
		Remark:          one.Remark,
		GroupID:         one.GroupID,
		Cookie:          one.Cookie,
		CookieStatus:    one.CookieStatus,
		CookieRefreshAt: time.Unix(one.CookieRefreshAt, 0),
		CID:             one.CID,
	}
}

func InzoneUserTo(one *InzoneUser) *pb.InzoneUser {
	if one == nil {
		return nil
	}

	return &pb.InzoneUser{
		ID:              one.ID,
		Name:            one.Name,
		Phone:           one.Phone,
		Remark:          one.Remark,
		GroupID:         one.GroupID,
		Cookie:          one.Cookie,
		CookieStatus:    one.CookieStatus,
		CookieRefreshAt: one.CookieRefreshAt.Local().Unix(),
		CID:             one.CID,
	}
}

func InzoneUserListTo(data []*InzoneUser) []*pb.InzoneUser {
	res := make([]*pb.InzoneUser, len(data))
	for k, v := range data {
		res[k] = InzoneUserTo(v)
	}
	return res
}
