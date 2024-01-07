package systemparam

import (
	"sync"
	"time"
)

type Config struct {
	RefreshTime time.Duration
}

type SystemParamKey string

// SMTP用的帳號
const SMTP_ACCOUNT SystemParamKey = "smtp_account"

// SMTP用的密碼
const SMTP_PASSWORD SystemParamKey = "smtp_password"

// 是否需要兩段式驗證
const NEED_TWOSTEP SystemParamKey = "need_twostep"

// 會員最大憑證數量
const MEMBER_MAX_LICENSE SystemParamKey = "member_max_license"

var SystemParamList = []SystemParamKey{
	SMTP_ACCOUNT,
	SMTP_PASSWORD,
	NEED_TWOSTEP,
	MEMBER_MAX_LICENSE,
}

type SystemParamExport struct {
	SystemParam map[SystemParamKey]SystemParamData
	mux         sync.RWMutex
}

type SystemParamData struct {
	Key      SystemParamKey `json:"key"`
	Value    string         `json:"value"`
	Category string         `json:"category" `
	LastTime time.Time      `json:"lastTime"`
}
