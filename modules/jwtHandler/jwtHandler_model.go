package jwthandler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

var cfg Config

type Config struct {
	TokenLookupKey string
	// Secret : jwt使用的密鑰
	Secret string `default:"jwtjwt"`
	// Expires : jwt存活時間長度
	Expires int `default:"7200"`
	// LocalsTokenKey : 存於Locals的jwtToken keyname
	LocalsTokenKey string `default:"jwtToken"`
	// OnSuccess : jwt流程成功後執行流程
	OnSuccess func(*fiber.Ctx) error
	// OnJWTError : jwt流程中發生錯誤流程
	OnJWTError func(*fiber.Ctx, error) error
}

// Claims : jwt Claims格式
type Claims struct {
	UID string
	jwt.RegisteredClaims
}
