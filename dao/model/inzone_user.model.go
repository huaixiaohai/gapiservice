package model

import "github.com/huaixiaohai/gapiservice/pb"

type InzoneUser struct {
	Model
	Name    string `gorm:"type:varchar(20);uniqueIndex:_name_phone"`
	Phone   string `gorm:"type:varchar(20);uniqueIndex:_name_phone"`
	Remark  string
	GroupID string
	Cookie  string
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
	}
}

func InzoneUserListTo(data []*InzoneUser) []*pb.InzoneUser {
	res := make([]*pb.InzoneUser, len(data))
	for k, v := range data {
		res[k] = InzoneUserTo(v)
	}
	return res
}
