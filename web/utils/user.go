package utils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"lottery/conf"
	"lottery/web/def"
)

var (
	TokenExpired     error = errors.New("Token is expired")
	TokenNotValidYet error = errors.New("Token not active yet")
	TokenMalformed   error = errors.New("That's not even a token")
	TokenInvalid     error = errors.New("Couldn't handle this token:")
)

// 生成token
func CreateToken(claims def.CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(conf.MySecret)
}

// 解析token
func JwtVerifyToken(tokenString string) (*def.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &def.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return conf.MySecret, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if claims, ok := token.Claims.(*def.CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, TokenInvalid
}

// 生成token
//func JwtGetToken() string {
//	// 生成加密串过程
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
//		"nick_name": "xiaolin",
//		"email":     "462441355@qq.com",
//		"id":        "1",
//		"iss":       "lottery",
//		"iat":       time.Now().Unix(),
//		"jti":       "9527",
//		"exp":       time.Now().Add(10 * time.Hour * time.Duration(1)).Unix(),
//	})
//	//  把token已约定的加密方式和加密秘钥加密，当然也可以使用不对称加密
//	tokenString, _ := token.SignedString(conf.MySecret)
//	//  登录时候，把tokenString返回给客户端，然后需要登录的页面就在header上面附此字符串
//	//  eg: header["Authorization"] = "bears "+tokenString
//
//	return tokenString
//}

// 解析token
//func JwtVerifyToken(token string) {
//	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
//		//自己加密的秘钥或者说盐值
//		return conf.MySecret, nil
//	})
//	if err != nil {
//		fmt.Printf("解析错误")
//	}
//	fmt.Printf("解析成功\n")
//	data := parsedToken.Claims
//	fmt.Println(data)
//}
