package def

import "github.com/dgrijalva/jwt-go"

type CustomClaims struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	jwt.StandardClaims
}