package basicauth

import (
	"crypto/sha512"
	"encoding/hex"
	"errors"
)

// 最基本的帳號密碼驗證
var cfg Config

type Config struct {
	UseEncodeType HashType
	PasswordHash  string
}

func Init(newCfg Config) {
	cfg = newCfg
}

func Check(password string, truePassword string) error {
	// 加密輸入的密碼
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return err
	}

	// 比對加密後密碼是否相同
	if truePassword != hashedPassword {
		return errors.New("compare password fail")
	}

	return nil
}

// HashPassword : 加密密碼
func HashPassword(password string) (string, error) {
	switch cfg.UseEncodeType {
	case HashTypeSHA512:
		// 加密輸入的密碼
		hasher := sha512.New()
		_, err := hasher.Write([]byte(password + cfg.PasswordHash))
		if err != nil {
			return "", err
		}

		return hex.EncodeToString(hasher.Sum(nil)), nil
	default:
		// 未指定加密，不處理後送出
		return password, nil
	}

}
