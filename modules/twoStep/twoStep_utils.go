package twostep

import (
	"fmt"
	"log"
	"math/rand"
	"net/smtp"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// GenSixRandomNum : 取得六碼隨機數字
func GenSixRandomNum() string {
	t := time.Now().UnixNano() + 104
	r := rand.New(rand.NewSource(t))
	ri := r.Int()
	ret := ""
	for i := 0; i < 6; i++ {
		ret = ret + strconv.Itoa(ri%10)
		ri = ri / 10
	}

	return ret
}

// SendMail : 寄出信件
func SendMail(to string, subject string, body string) {
	from := cfg.From

	// 建立 MIME header
	header := make(map[string]string)
	header["From"] = from
	header["To"] = to
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/plain; charset=\"utf-8\""

	// 組裝郵件內容
	content := ""
	for k, v := range header {
		content += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	content += "\r\n" + body

	err := smtp.SendMail(cfg.SMTPServer, plainAuth, from, []string{to}, []byte(content))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}
}

func (d TwoStepData) SendMail(subject string, body string) error {
	from := cfg.From

	// 建立 MIME header
	header := make(map[string]string)
	header["From"] = fmt.Sprintf("%s <%s>", "8mb System", from)
	header["To"] = d.Mail
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=\"utf-8\""

	// 組裝郵件內容
	content := ""
	for k, v := range header {
		content += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	content += "\r\n" + body

	return smtp.SendMail(cfg.SMTPServer, plainAuth, from, []string{d.Mail}, []byte(content))
}

// generateToken : 產生金鑰
func generateToken(claims jwt.Claims) (jwtToken string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtToken, err = token.SignedString([]byte(cfg.Secret))
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}

// parseTokenToClaims : 從JWT中取得資訊
func parseTokenToClaims(token string) (*Claims, error) {
	claims := &Claims{}

	if _, err := jwt.ParseWithClaims(token, claims, cfg.Keyfunc); err != nil {
		return nil, err
	}
	return claims, nil
}
