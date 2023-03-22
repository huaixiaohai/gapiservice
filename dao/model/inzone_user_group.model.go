package model

import "github.com/huaixiaohai/gapiservice/pb"

type InzoneUserGroup struct {
	Model
	Name         string
	DingTalkHook string
}

func InzoneUserGroupFrom(one *pb.InzoneUserGroup) *InzoneUserGroup {
	if one == nil {
		return nil
	}

	return &InzoneUserGroup{
		Model: Model{
			ID: one.ID,
		},
		Name:         one.Name,
		DingTalkHook: one.DingTalkHook,
	}
}

func InzoneUserGroupTo(one *InzoneUserGroup) *pb.InzoneUserGroup {
	if one == nil {
		return nil
	}

	return &pb.InzoneUserGroup{
		ID:           one.ID,
		Name:         one.Name,
		DingTalkHook: one.DingTalkHook,
	}
}

func InzoneUserGroupListTo(data []*InzoneUserGroup) []*pb.InzoneUserGroup {
	res := make([]*pb.InzoneUserGroup, len(data))
	for k, v := range data {
		res[k] = InzoneUserGroupTo(v)
	}
	return res
}
