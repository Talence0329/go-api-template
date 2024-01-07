package fbauth

import (
	"backend/modules/tools"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

var cfg = Config{}

func Init(newCfg Config) {
	cfg = newCfg
}

// GetFBUserInfos : 以Token取得會員資訊
func GetFBUserInfo(token string) (ProfileData, error) {
	fields := "id,name,email"

	req := tools.Request{
		Url:    fmt.Sprintf("%s?fields=%s&access_token=%s", UrlFBMe, fields, token),
		Method: fiber.MethodGet,
	}

	resData := ProfileData{}

	if res, err := req.Do(&resData); err != nil {
		return resData, err
	} else {
		if resData.ID == "" {
			return resData, errors.New("與FB取得玩家資訊失敗")
		}
		fmt.Println(string(res.Body))
		return resData, nil
	}
}
