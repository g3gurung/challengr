package middleware

import jwt "github.com/dgrijalva/jwt-go"

//JWTUser struct is a schema for
type JWTUser struct {
	ID             int64   `json:"id"`
	FacebookUserID string  `json:"facebook_user_id"`
	Weight         float32 `json:"weight"`
	Role           string  `json:"role"`

	jwt.StandardClaims
}
