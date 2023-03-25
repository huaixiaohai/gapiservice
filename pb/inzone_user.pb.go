package pb

type InzoneUser struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Remark    string `json:"remark"`
	GroupID   string `json:"group_id"`
	GroupName string `json:"group_name"`
	Cookie    string `json:"cookie"`
	UniID     string `json:"uni_id"`
}

type InzoneUserListReq struct {
}

type InzoneUserListResp struct {
	List []*InzoneUser `json:"list"`
}
