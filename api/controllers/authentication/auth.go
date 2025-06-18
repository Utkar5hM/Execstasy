package authentication

import (
	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
	Name     string `json:"name"`
	Id       int    `json:"id"`
}

type JwtCustomClaims struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Id       int    `json:"id"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}
