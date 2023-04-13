package api

import (
	"github.com/google/wire"
	"github.com/huaixiaohai/gapiservice/dao"
)

var Set = wire.NewSet(
	InzoneUserGroupApiSet,
	InzoneUserApiSet,
	InzoneApiSet,
	UserApiSet,

	Ting13ApiSet,
) //

var transaction = dao.Transaction
