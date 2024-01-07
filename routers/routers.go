package routers

import (
	"backend/config"
	"backend/src/auth"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/helmet/v2"
)

var excludePath = map[string]bool{
	"/api/auth/accountLogin":    true,
	"/api/auth/accountRegister": true,
	"/api/auth/googleLogin":     true,
	"/api/auth/fbLogin":         true,
}

// Set : 設定全部的路由 middleware、功能
func Set(r fiber.Router) {
	setMiddlewareRouter(r)
	setFuncRouter(r)
}

// setMiddlewareRouter : 設定功能的路由
func setMiddlewareRouter(r fiber.Router) {
	r.Use(
		recover.New(), // panic的抓取動作，須放置於其他路由上方
		func(c *fiber.Ctx) error { // 設定每次request建立一個log物件，並在最後處理或印出log
			if len(string(c.Request().Body())) > 200 {
				fmt.Printf("[req][path->%v][body]%v \n", c.Path(), string(c.Request().Body())[:200])
			}
			fmt.Printf("[req][path->%v][body]%v \n", c.Path(), string(c.Request().Body()))
			defer func() {
				fmt.Printf("[res][path->%v][code->%v] \n", c.Path(), c.Response().StatusCode())
			}()
			return c.Next()
		},
		cors.New(cors.Config{ // 設定cors
			AllowOrigins:     config.GetEnv().AllowOrigins,
			AllowCredentials: config.GetEnv().AllowCredentials,
		}),
		helmet.New(), // 設定header安全機制(參考fiber官方設定)
	)
}

// setFuncRouter : 設定功能的路由
func setFuncRouter(r fiber.Router) {
	api := r.Group("/api")
	{
		auth.SetRouter(api)
	}
}
