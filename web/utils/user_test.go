package utils

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"lottery/web/def"
	"testing"
	"time"
)

func TestUserToken(t *testing.T) {
	claims := def.CustomClaims{
		1,
		"xiaolin1",
		"462441355@qq.com",
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
			Issuer:    "lottery",
			IssuedAt:  time.Now().Unix(),
		},
	}
	token, _ := CreateToken(claims)
	claimsData, err := JwtVerifyToken(token)
	if err != nil {

	}
	fmt.Print()

	if claimsData.Email != "462441355@qq.com" {
		t.Errorf("解析失败")
	}
}
