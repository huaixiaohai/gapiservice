package api

import (
	"errors"

	"github.com/google/wire"

	"github.com/gin-gonic/gin"
	"github.com/huaixiaohai/gapiservice/auth"
)

var UserApiSet = wire.NewSet(NewUserApi)

func NewUserApi() *UserApi {
	return &UserApi{}
}

type UserApi struct {
}

type User struct {
	ID       string
	Name     string
	Password string
}

var admin = &User{
	ID:       "10",
	Name:     "admin",
	Password: "admin123",
}

type LoginReq struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type LoginResp struct {
	Token     string `json:"token"`
	TokenType string `json:"token_type"`
	ExpiresAt int64  `json:"expires_at"`
}

func (a *UserApi) Login(ctx *gin.Context, req *LoginReq) (*LoginResp, error) {

	if req.UserName != admin.Name || req.Password != admin.Password {
		return nil, errors.New("用戶名或密码不正确")
	}

	token, expiresAt, err := auth.GenToken(admin.ID, admin.Password)
	if err != nil {
		return nil, err
	}
	return &LoginResp{
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}

type Empty struct {
}

type GetUserResp struct {
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
}

func (a *UserApi) Get(ctx *gin.Context, req *Empty) (*GetUserResp, error) {

	return &GetUserResp{
		UserID:   admin.ID,
		UserName: admin.Name,
	}, nil
}
