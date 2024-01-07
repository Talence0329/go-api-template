package cookies

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

// Config : cookies的設定值
type Config struct {
	KeyJWTToken Key       `default:"jwtToken"`
	MaxAge      int       `json:"max_age"`
	Expires     time.Time `json:"expires"`
	Secure      bool      `json:"secure"`
	HTTPOnly    bool      `json:"http_only"`
	Domain      string    `json:"domain"`
	SameSite    string    `json:"same_site"`
}

var defaultConfig = Config{
	KeyJWTToken: "jwtToken",
	MaxAge:      7200,
	Secure:      true,
	HTTPOnly:    true,
	SameSite:    fiber.CookieSameSiteLaxMode,
}

// Key : 使用於cookies的key
type Key = string
