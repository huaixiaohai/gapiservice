package model

import (
	"github.com/huaixiaohai/gapiservice/pb"
)

type NovelResource struct {
	Model
	Number  string
	Name    string `gorm:"type:varchar(30)"`
	NovelID string `gorm:"index"`
}

func NovelResourceFrom(one *pb.NovelResource) *NovelResource {
	if one == nil {
		return nil
	}

	return &NovelResource{
		Model: Model{
			ID: one.ID,
		},
		Name:    one.Name,
		NovelID: one.NovelID,
	}
}

func NovelResourceTo(one *NovelResource) *pb.NovelResource {
	if one == nil {
		return nil
	}

	return &pb.NovelResource{
		ID:      one.ID,
		Name:    one.Name,
		NovelID: one.NovelID,
	}
}

func NovelResourceListTo(data []*NovelResource) []*pb.NovelResource {
	res := make([]*pb.NovelResource, len(data))
	for k, v := range data {
		res[k] = NovelResourceTo(v)
	}
	return res
}
