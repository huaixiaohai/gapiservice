package api

type InZoneGroup struct {
	ID     string
	Name   string
	DTHook string
}

type InZoneUser struct {
	ID        string
	GroupID   string
	Name      string
	Phone     string
	Remark    string
	Cookie    string
	IsExpired string
}

type InZone struct {
}

//func (a *InZone)
