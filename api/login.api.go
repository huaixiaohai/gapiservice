package api

import "github.com/gin-gonic/gin"

type LoginApi struct {

}

type LoginReq struct {

}

type LoginResp struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresAt   int64  `json:"expires_at"`
}

func (a *LoginApi) Login(ctx *gin.Context, req *LoginReq) (*LoginResp, error) {
	return &LoginResp{
		AccessToken: "sadgisjafdhikjsdolsadisud",
		TokenType: "Bear",
		ExpiresAt: 36000,
	}, nil
}
