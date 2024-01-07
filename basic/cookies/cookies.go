package cookies

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

var cfg Config

// Init : 設定cookies的config
func Init(c ...Config) {
	if len(c) != 0 {
		cfg = c[0]
	} else {
		cfg = defaultConfig
	}
}

// SetJWT : 設定jwt，紀錄JWT token
func SetJWT(c *fiber.Ctx, value string) {
	c.Cookie(&fiber.Cookie{
		Name:     cfg.KeyJWTToken,
		Value:    value,
		MaxAge:   cfg.MaxAge,
		HTTPOnly: cfg.HTTPOnly,
		Secure:   cfg.Secure,
		Domain:   cfg.Domain,
		SameSite: cfg.SameSite,
	})
}

// ClearJWT : 清除jwt，紀錄JWT token
func ClearJWT(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     cfg.KeyJWTToken,
		Expires:  time.Now().Add(-(time.Hour * 2)),
		HTTPOnly: cfg.HTTPOnly,
		Secure:   cfg.Secure,
		Domain:   cfg.Domain,
		SameSite: cfg.SameSite,
	})
}
