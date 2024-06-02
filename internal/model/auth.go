package model

import "github.com/golang-jwt/jwt/v5"

type MyClaim struct {
	UserData
	jwt.RegisteredClaims
}

type UserData struct {
	ID    uint64 `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}
