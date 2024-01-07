package access

import "github.com/golang-jwt/jwt/v4"

type Config struct {
	JWTTokenKey string
	JWTSecret   string
	JWTExpires  int
}

// Claims : jwt Claims格式
type Claims struct {
	UID string
	jwt.RegisteredClaims
}
