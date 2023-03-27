package pb

import "errors"

type InzoneUser struct {
	ID              string        `json:"id"`
	Name            string        `json:"name"`
	Phone           string        `json:"phone"`
	Remark          string        `json:"remark"`
	GroupID         string        `json:"group_id"`
	GroupName       string        `json:"group_name"`
	Cookie          string        `json:"cookie"`
	CookieRefreshAt int64         `json:"cookie_refresh_at"`
	CookieStatus    ECookieStatus `json:"cookie_status"`
	UUID            string        `json:"uuid"`
}

type InzoneUserListReq struct {
}

type InzoneUserListResp struct {
	List []*InzoneUser `json:"list"`
}

func GetUUID(name, phone string) (string, error) {
	if len(name) != 6 && len(name) != 9 {
		return "", errors.New("名称长度不正确")
	}
	if len(phone) != 11 {
		return "", errors.New("手机号长度不正确")
	}

	var uuid string
	var na string
	if len(name) == 6 {
		na = name[0:3] + "*"
	} else {
		na = name[0:3] + "**"
	}
	uuid = na + phone[0:3] + "*****" + phone[8:]
	return uuid, nil
}
