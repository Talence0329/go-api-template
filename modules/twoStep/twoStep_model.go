package twostep

import "github.com/golang-jwt/jwt/v4"

type TwoStepData struct {
	Action string
	UID    string
	Code   string
	Mail   string
	Phone  string
}

// Claims : jwt Claims格式
type Claims struct {
	Action string
	UID    string
	Mail   string
	Phone  string
	jwt.RegisteredClaims
}
