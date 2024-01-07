package validatehandler

import (
	"backend/basic/apiprotocol"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// New :  檢查格式是否如設定
func New(st any) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		if err := c.BodyParser(st); err != nil {
			return apiprotocol.Request900101.ToRes().Err(err.Error()).ToErr()
		}
		validate := validator.New()
		errStr := "[validatehandler] "
		if err := validate.Struct(st); err != nil {
			switch e := err.(type) {
			case validator.ValidationErrors:
				for _, ee := range e {
					errStr += fmt.Sprintf("%v %v %v, ", ee.StructNamespace(), ee.Tag(), ee.Param())
				}
				return apiprotocol.Request900101.ToRes().Err(errStr).ToErr()
			default:
				return apiprotocol.Request900101.ToRes().Err(e.Error()).ToErr()
			}
		}

		return c.Next()
	}
}
