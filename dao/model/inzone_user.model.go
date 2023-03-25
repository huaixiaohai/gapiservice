package model

import "github.com/huaixiaohai/gapiservice/pb"

type InzoneUser struct {
	Model
	Name    string `gorm:"type:varchar(30)"`
	Phone   string `gorm:"type:varchar(30)"`
	Remark  string
	GroupID string
	Cookie  string
	UniID   string `gorm:"type:varchar(30);uniqueIndex"`
}

func InzoneUserFrom(one *pb.InzoneUser) *InzoneUser {
	if one == nil {
		return nil
	}

	return &InzoneUser{
		Model: Model{
			ID: one.ID,
		},
		Name:    one.Name,
		Phone:   one.Phone,
		Remark:  one.Remark,
		GroupID: one.GroupID,
		Cookie:  one.Cookie,
		UniID:   one.UniID,
	}
}

func InzoneUserTo(one *InzoneUser) *pb.InzoneUser {
	if one == nil {
		return nil
	}

	return &pb.InzoneUser{
		ID:      one.ID,
		Name:    one.Name,
		Phone:   one.Phone,
		Remark:  one.Remark,
		GroupID: one.GroupID,
		Cookie:  one.Cookie,
		UniID:   one.UniID,
	}
}

func InzoneUserListTo(data []*InzoneUser) []*pb.InzoneUser {
	res := make([]*pb.InzoneUser, len(data))
	for k, v := range data {
		res[k] = InzoneUserTo(v)
	}
	return res
}
