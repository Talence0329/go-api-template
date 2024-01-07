// 以JWT為本體實作
package access

import (
	"backend/basic/locals"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

var cfg Config

func Init(initConfig Config) {
	cfg = initConfig
}

func JWTOnSuccess(c *fiber.Ctx) error {
	jwtData := getJwtData(c)
	if jwtData == nil {
		return c.SendStatus(http.StatusUnauthorized)
	}

	if !jwtData.Valid {
		return c.SendStatus(http.StatusUnauthorized)
	}

	claims, err := parseTokenToClaims(jwtData.Raw)
	if err != nil {
		return c.SendStatus(http.StatusUnauthorized)
	}

	// 沒有帳號的話直接回傳失敗
	if claims.UID == "" {
		return c.SendStatus(http.StatusUnauthorized)
	} else {
		locals.SetMemberUID(c, claims.UID)
	}

	// 如果剩下不到一半的存活時間就給新的token
	if claims.ExpiresAt.Before(time.Now().Add(-1 * time.Second * time.Duration(cfg.JWTExpires/2))) {
		// 更新token exp
		jwtToken, err := generateToken(claims.UID, time.Now().Add(time.Second*time.Duration(cfg.JWTExpires)))
		if err != nil {
			return c.SendStatus(http.StatusUnauthorized)
		}
		c.Response().Header.Add(cfg.JWTTokenKey, jwtToken)
	}

	return c.Next()
}

func JWTOnError(c *fiber.Ctx, err error) error {
	fmt.Println(err)
	return c.SendStatus(http.StatusUnauthorized)
}

// LoginSuccess : 登入成功後的行為，會以uid產生token後存入cookies
func LoginSuccess(c *fiber.Ctx, uid string) error {
	// 產生token
	jwtToken, err := generateToken(uid, time.Now().Add(time.Second*time.Duration(cfg.JWTExpires)))
	if err != nil {
		return err
	}

	// 設定token至header
	c.Response().Header.Add(cfg.JWTTokenKey, jwtToken)

	return nil
}

// generateToken : 產生金鑰
func generateToken(uid string, exp time.Time) (jwtToken string, err error) {
	claims := Claims{
		UID: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtToken, err = token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}

// parseTokenToClaims : 從JWT中取得使用者資訊
func parseTokenToClaims(token string) (*Claims, error) {
	claims := &Claims{}

	if _, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) { return []byte(cfg.JWTSecret), nil }); err != nil {
		return nil, err
	}
	return claims, nil
}

// getJwtData : 透過Locals取得jwt
func getJwtData(ctx *fiber.Ctx) *jwt.Token {
	if ctx.Locals(locals.KeyJWTToken) != nil {
		switch ctx.Locals(locals.KeyJWTToken).(type) {
		case *jwt.Token:
			return ctx.Locals(locals.KeyJWTToken).(*jwt.Token)
		}
	}
	return nil
}
