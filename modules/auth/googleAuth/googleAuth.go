package googleauth

import (
	"backend/modules/tools"
	"fmt"
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/tidwall/gjson"
)

var cfg = Config{}

func Init(newCfg Config) {
	cfg = newCfg
}

func LoginInit(redirectUri string) string {
	return fmt.Sprintf("%s?client_id=%s&response_type=%s&scope=%s/%s&redirect_uri=%s",
		UrlLogin,
		cfg.ClientID,
		ResponseTypeCode,
		UrlScope,
		ScopeEmail,
		redirectUri,
	)
}

func GetAccessToken(code string, redirectUri string) (token string, err error) {
	data := url.Values{
		"code":          {code},
		"client_id":     {cfg.ClientID},
		"client_secret": {cfg.SecretKey},
		"grant_type":    {GrantTypeAuthCode},
		"redirect_uri":  {redirectUri},
	}

	req := tools.Request{
		Url:         UrlToken,
		Method:      fiber.MethodPost,
		ContentType: fiber.MIMEApplicationForm,
		Body:        []byte(data.Encode()),
	}

	if res, err := req.Do(); err != nil {
		return token, err
	} else {
		token = gjson.GetBytes(res.Body, "access_token").String()
	}

	return token, nil
}

// GetGoogleUserInfo : 以Token取得會員資訊
func GetGoogleUserInfo(token string) (ScopeEmailData, error) {
	req := tools.Request{
		Url:    fmt.Sprintf("%s?alt=json&access_token=%s", UrlUserInfo, token),
		Method: fiber.MethodGet,
	}

	resData := ScopeEmailData{}

	if _, err := req.Do(&resData); err != nil {
		return resData, err
	} else {
		return resData, nil
	}
}
