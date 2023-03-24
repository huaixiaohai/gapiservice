package pb

type InzoneUserGroup struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	DingTalkHook string `json:"ding_talk_hook"`
}

type InzoneUserGroupListReq struct {
}

type InzoneUserGroupListResp struct {
	List []*InzoneUserGroup `json:"list"`
}
