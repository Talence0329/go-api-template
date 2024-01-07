package errorhandler

import (
	"backend/basic/apiprotocol"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func New(config ...Config) fiber.ErrorHandler {
	cfg := setConfig(config...)
	return func(c *fiber.Ctx, e error) error {
		if e != nil {
			switch err := e.(type) {
			case *apiprotocol.BaseResponse:
				code := err.RetStatus.Code
				if cfg.Log {
					fmt.Printf("[HTTP_%d/%s] %s%s[msg] %s [error] %s", code, e, c.Hostname(), c.OriginalURL(), err.Msg(), err.Error())
				}
				if !cfg.ResError {
					err.RetStatus.Error = ""
				}
				return c.JSON(err)
			case *fiber.Error:
				if cfg.Log {
					fmt.Printf("[HTTP_%d/%s] %s%s %s", err.Code, e, c.Hostname(), c.OriginalURL(), err.Message)
				}
				if cfg.ResError {
					return c.Status(err.Code).JSON(e)
				} else {
					return c.SendStatus(err.Code)
				}
			default:
				if cfg.Log {
					fmt.Printf("[HTTP_%d] %s", fiber.StatusInternalServerError, err.Error())
				}
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.ErrInternalServerError)
			}
		}
		return c.Next()
	}
}

// Config defines the config for middleware.
type Config struct {
	// ResponseTraceHeader : 回應 cloud trace header
	ResCloudTraceHeader bool
	ResError            bool
	Log                 bool
}

// configDefault is the default config
var configDefault = Config{
	ResCloudTraceHeader: true,
	ResError:            true,
	Log:                 true,
}

// Helper function to set default values
func setConfig(config ...Config) Config {
	// Return default config if nothing provided
	if len(config) < 1 {
		return configDefault
	}
	// Override default config
	cfg := config[0]
	return cfg
}
