package pb

type EGoodsType string

const (
	EGoodsTypeDaily   EGoodsType = "每日"
	EGoodsTypeSeries1 EGoodsType = "系列1"
	EGoodsTypeSeries2 EGoodsType = "系列2"
	EGoodsTypeSeries3 EGoodsType = "系列3"
	EGoodsTypeSeries4 EGoodsType = "系列4"
)

type ECookieStatus int8

const (
	ECookieStatusNone    ECookieStatus = iota
	ECookieStatusValid   ECookieStatus = 1
	ECookieStatusInvalid ECookieStatus = 2
)
