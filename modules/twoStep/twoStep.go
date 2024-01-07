package twostep

import (
	"errors"
	"net/smtp"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/lithammer/shortuuid/v4"
)

type Config struct {
	SMTPServer   string
	From         string
	Password     string
	MaxCount     int
	Secret       string
	Secure       bool
	Expires      int
	CookiePrefix string
	FuncUID      func() string
	FuncCode     func() string
	Keyfunc      jwt.Keyfunc
}

var cfg = Config{
	CookiePrefix: "TwoStep_",
	FuncUID:      shortuuid.New,
	FuncCode:     GenSixRandomNum,
}

var plainAuth smtp.Auth

func Init(initConfig Config) {
	cfg = initConfig
	cfg.Keyfunc = func(t *jwt.Token) (any, error) { return []byte(cfg.Secret), nil }

	plainAuth = smtp.PlainAuth("", cfg.From, cfg.Password, "smtp.gmail.com")
}

// New : 產生一組信件驗證碼
func NewMail(action string, mail string) TwoStepData {
	return TwoStepData{
		Action: action,
		UID:    cfg.FuncUID(),
		Code:   cfg.FuncCode(),
		Mail:   mail,
	}
}

// VerifySuccess : 驗證成功後的行為，會以uid與信件或電話產生token後存入cookies
func (d TwoStepData) VerifySuccess(c *fiber.Ctx) error {
	claims := Claims{
		Action: d.Action,
		UID:    d.UID,
		Mail:   d.Mail,
		Phone:  d.Phone,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(cfg.Expires))),
		},
	}
	// 產生token
	if jwtToken, err := generateToken(claims); err != nil {
		return err
	} else {
		// 設定token至header
		c.Response().Header.Add(cfg.CookiePrefix, jwtToken)
	}

	return nil
}

// GetTokenData : 從token取得資訊
func GetTokenData(c *fiber.Ctx, action string) (*Claims, error) {
	auth := string(c.Request().Header.Peek(cfg.CookiePrefix))
	if auth == "" {
		return nil, errors.New("missing or malformed JWT")
	}

	tokenData, err := jwt.Parse(auth, cfg.Keyfunc)
	if err != nil {
		return nil, err
	}

	if !tokenData.Valid {
		return nil, errors.New("不合法的token")
	}

	claims, err := parseTokenToClaims(tokenData.Raw)
	if err != nil {
		return nil, err
	}

	return claims, nil
}
