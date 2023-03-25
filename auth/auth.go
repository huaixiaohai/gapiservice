package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type TokenClaims struct {
	UserID   string
	UserName string
	*jwt.StandardClaims
}

func GenToken(userID, userName string) (string, int64, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(3000 * time.Second).Local().Unix()
	issuer := "frank"
	claims := &TokenClaims{
		UserID:   userID,
		UserName: userName,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expireTime,
			Issuer:    issuer,
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("golang"))
	return token, expireTime, err

}

func ParseToken(token string) (userID, userName string, expiresTime int64, err error) {
	var tokenClaims *jwt.Token
	tokenClaims, err = jwt.ParseWithClaims(token, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("golang"), nil
	})
	if err != nil {
		return
	}

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*TokenClaims); ok && tokenClaims.Valid {
			userID = claims.UserID
			userName = claims.UserName
			expiresTime = claims.ExpiresAt
			return
		}
	}

	return
}
