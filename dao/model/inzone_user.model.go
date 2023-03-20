package model

import "github.com/huaixiaohai/gapiservice/pb"

type InzoneUser struct {
	Model

}

func InzoneUserFrom(one *pb.InzoneUser) *InzoneUser {
	if one == nil {
		return nil
	}

	return &InzoneUser{
		Model: Model{
			ID: one.ID,
		},

	}
}

func InzoneUserTo(one *InzoneUser) *pb.InzoneUser {
	if one == nil {
		return nil
	}


	return &pb.InzoneUser{
		ID:          one.ID,
	}
}

func InzoneUserListTo(data []*InzoneUser) []*pb.InzoneUser {
	res := make([]*pb.InzoneUser, len(data))
	for k, v := range data {
		res[k] = InzoneUserTo(v)
	}
	return res
}
