package pb

type InzoneUser struct {
	ID              string
	Name            string
	Phone           string
	Remark          string
	GroupID         string
	Cookie          string
	CookieUpdatedAt int64
}

type InzoneUserListReq struct {
}

type InzoneUserListResp struct {
}
