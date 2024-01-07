package main

import (
	"backend/basic/access"
	"backend/basic/cookies"
	"backend/basic/database"
	"backend/config"
	basicauth "backend/modules/auth/basicAuth"
	fbauth "backend/modules/auth/fbAuth"
	googleauth "backend/modules/auth/googleAuth"
	errorhandler "backend/modules/errorHandler"
	gcsmamnger "backend/modules/gcsManager"
	jwthandler "backend/modules/jwtHandler"
	systemparam "backend/modules/systemParam"
	twostep "backend/modules/twoStep"
	"backend/routers"
	"backend/src/auth"
	"fmt"
	"log"
	"runtime"
	"strconv"
	"time"

	_ "github.com/lib/pq"
	"github.com/lithammer/shortuuid/v4"

	"github.com/gofiber/fiber/v2"
)

var (
	version  = "0.0.0"
	commitID = "dev"
)

// ServiceName : 服務名稱 | 不可修改，會影響資料一致性
const ServiceName = "backend"

// resourcePath : [resource]本地靜態資源檔案路徑
const resourcePath = "resource"

func init() {
	SystemInit()
	FunctionInit()
}

// SystemInit : 系統初始化，流程：取得ENV與連線設定值
func SystemInit() {
	// 從env取得參數
	if err := config.EnvInit(); err != nil {
		log.Fatalf("[pkg->config][func->EnvInit] %s", err)
	}
	env := config.GetEnv()
	// 初始化cookies，帶入設定值
	cookies.Init(cookies.Config{
		KeyJWTToken: env.JwtTokenKey,
		MaxAge:      env.MaxAge,
		Secure:      env.Secure,
		SameSite:    "Lax",
		HTTPOnly:    true,
	})

	database.Init()

	systemparam.Init(systemparam.Config{
		RefreshTime: 5 * time.Minute,
	})

	access.Init(access.Config{
		JWTTokenKey: env.JwtTokenKey,
		JWTSecret:   env.JwtSecretkey,
		JWTExpires:  env.JwtExpires,
	})
	// 初始化jwtHandler，設定金鑰以及其他相關設定
	if err := jwthandler.Init(jwthandler.Config{
		TokenLookupKey: fmt.Sprintf("header:%s", env.JwtTokenKey),
		Secret:         env.JwtSecretkey,
		Expires:        env.JwtExpires,
		LocalsTokenKey: env.JwtTokenKey,
		OnSuccess:      access.JWTOnSuccess,
		OnJWTError:     access.JWTOnError,
	}); err != nil {
		log.Fatalf("[func->jwthandler.Init] %s", err)
	}
	gcsmamnger.Init(gcsmamnger.Config{
		BucketName: "8mb-file",
	})
}

// FunctionInit : 主功能初始化
func FunctionInit() {
	env := config.GetEnv()
	if err := auth.Init(auth.Config{
		IsTwoStep:    env.IsTwoStep,
		IsUseLogAll:  env.IsUseLogAll,
		PasswordHash: env.PasswordHash,
	}); err != nil {
		log.Fatalf("[func->auth.Init] %s", err)
	}
	basicauth.Init(basicauth.Config{
		UseEncodeType: basicauth.HashTypeSHA512,
		PasswordHash:  "at104",
	})
	googleauth.Init(googleauth.Config{
		ClientID:  env.GoogleClientID,
		SecretKey: env.GoogleClientSecret,
	})
	fbauth.Init(fbauth.Config{
		FBlientID:   env.FBClientID,
		FBSecretKey: env.FBClientSecret,
	})
	twostep.Init(twostep.Config{
		SMTPServer:   "smtp.gmail.com:587",
		From:         systemparam.SMTP_ACCOUNT.Get(),
		Password:     systemparam.SMTP_PASSWORD.Get(),
		MaxCount:     5,
		Secret:       env.JwtTwoStepSecretkey,
		Expires:      1800,
		Secure:       env.Secure,
		CookiePrefix: "twostepToken",
		FuncUID:      shortuuid.New,
		FuncCode:     twostep.GenSixRandomNum,
	})
}

func main() {
	env := config.GetEnv()

	// 設定fiber的config
	app := fiber.New(fiber.Config{
		ErrorHandler: errorhandler.New(errorhandler.Config{
			ResError: true,
		}),
		ReadTimeout:  env.ReadTimeout,
		WriteTimeout: env.WriteTimeout,
		IdleTimeout:  time.Second * 1,
	})

	// 修改底層fasthttp
	// app.Server().MaxKeepaliveDuration = time.Second * 1

	app.Get("/version", func(c *fiber.Ctx) error {
		versionData := fiber.Map{
			"environment":        env.Environment,
			"Service":            ServiceName,
			"Version":            version,
			"commitID":           commitID,
			"OpenConnections":    c.App().Server().GetOpenConnectionsCount(),
			"CurrentConcurrency": c.App().Server().GetCurrentConcurrency(),
			"goVersion":          runtime.Version(),
		}
		return c.Status(fiber.StatusOK).JSON(versionData)
	})

	app.Get("/test/:number", func(c *fiber.Ctx) error {
		shortuuidList := []string{}
		if number, err := strconv.Atoi(c.Params("number")); err != nil {
			return c.SendStatus(fiber.StatusOK)
		} else {
			for i := 0; i < number; i++ {
				shortuuidList = append(shortuuidList, shortuuid.New())
			}
		}

		return c.Status(fiber.StatusOK).JSON(shortuuidList)
	})

	// 設定middleware以及各主功能的路由
	routers.Set(app)

	app.Listen(fmt.Sprintf(":%v", env.Port))
}
