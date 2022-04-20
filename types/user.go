package types

import (
	"github.com/golang-jwt/jwt"
)

type PublicUser struct {
	UserId string `json:"user_id"`
	DisplayName  string `json:"name"`
	UserName string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	UserId string `json:"user_id"`
	DisplayName  string `json:"name"`
	UserName string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
	Role string
}

type CustomJWTClaims struct {
	UserId string `json:"user_id"`
	DisplayName string `json:"display_name"`
	Email string `json:"email"`
	Role string `json:"role"`
	ExpiresAt int64  `json:"expires_at"`
	jwt.StandardClaims
}