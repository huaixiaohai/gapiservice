package model

import (
	"github.com/huaixiaohai/gapiservice/pb"
)

type Novel struct {
	Model
	Name string `gorm:"type:varchar(30),uniqueIndex"`
	Desc string
}

func NovelFrom(one *pb.Novel) *Novel {
	if one == nil {
		return nil
	}

	return &Novel{
		Model: Model{
			ID: one.ID,
		},
		Name: one.Name,
		Desc: one.Desc,
	}
}

func NovelTo(one *Novel) *pb.Novel {
	if one == nil {
		return nil
	}

	return &pb.Novel{
		ID:   one.ID,
		Name: one.Name,
		Desc: one.Desc,
	}
}

func NovelListTo(data []*Novel) []*pb.Novel {
	res := make([]*pb.Novel, len(data))
	for k, v := range data {
		res[k] = NovelTo(v)
	}
	return res
}
