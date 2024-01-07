package apiprotocol

import "cloud.google.com/go/logging"

// BaseResponse : 基礎回應資料結構
type BaseResponse struct {
	RetStatus RetStatus `json:"retStatus"`
}

// RetStatus : 回傳狀態
type RetStatus struct {
	Code       Code             `json:"code"`
	Msg        string           `json:"msg,omitempty"`
	SystemTime int64            `json:"systemTime"`
	Error      string           `json:"error,omitempty"`
	Level      logging.Severity `json:"level,omitempty"`
}

// RetStatusContent : 回傳狀態內容
type RetStatusContent struct {
	Msg   string           `json:"msg"`
	Level logging.Severity `json:"level,omitempty"`
}

// APIResponse : API統一回傳結構
type APIResponse struct {
	Data      any       `json:"data"`
	RetStatus RetStatus `json:"retStatus"`
}

// StatusResponse : 狀態回傳結構
type StatusResponse struct {
	RetStatus RetStatus `json:"retStatus"`
}

type Code int

// Success10000 : 正常回應
const Success10000 Code = 10000

// Function 200000~599999
// auth 200100~200199
const (
	// auth區段指標
	Auth Code = iota + 200100
	// 登入失敗
	Auth200101
	// 查無此使用者
	Auth200102
	// 此使用者已被停用
	Auth200103
	// 密碼錯誤
	Auth200104
	// 建立會員失敗
	Auth200105
	// 此帳號已被使用
	Auth200106
	// 無法取得使用者資訊
	Auth200107
	// 產生密碼錯誤
	Auth200108
	// 驗證碼錯誤
	Auth200109
	// 需進行兩段式驗證
	Auth200110
	// 兩段式驗證token與要求的內容不符合
	Auth200111
	// 兩段式驗證信件送出失敗
	Auth200112
	// 兩段式驗證token已超過次數
	Auth200113
	// 兩段式驗證Code輸入錯誤
	Auth200114
	// 此使用者已刪除
	Auth200115
	// 兩段式驗證與執行的動作不符
	Auth200116
)

// Request 900100~
const (
	Request Code = iota + 900100
	// Request 資料格式不吻合
	Request900101
)

// Response 900200~
const (
	Response900200 Code = iota + 900200
)

// Response 900300~
const (
	Response900300 Code = iota + 900300
)

// Database 900500~
const (
	Database Code = iota + 900500
	// 資料庫連線失敗
	Database900501
	// 資料庫錯誤
	Database900502
)

// JWT 900600~
const (
	JWT Code = iota + 900600
	// JWT產生失敗
	JWT900601
)

// JSON 900700~900199
const (
	// JSON區段指標
	JSON Code = iota + 900700
	// JSON parse to byte Fail
	JSON900701
	// JSON parse to struct Fail
	JSON900702
)

// API 900800~900899
const (
	// API區段指標
	API Code = iota + 900800
	// API Call Error
	API900801
)

// 未知錯誤
const Uknown999999 Code = 999999

var retStatusList = map[Code]RetStatusContent{
	Success10000:   {Msg: "Success"},
	Auth200101:     {Msg: "登入失敗"},
	Auth200102:     {Msg: "查無此使用者"},
	Auth200103:     {Msg: "此使用者已被停用"},
	Auth200104:     {Msg: "密碼錯誤"},
	Auth200105:     {Msg: "建立會員失敗"},
	Auth200106:     {Msg: "此帳號已被使用"},
	Auth200107:     {Msg: "無法取得使用者資訊"},
	Auth200108:     {Msg: "產生密碼錯誤"},
	Auth200109:     {Msg: "驗證碼錯誤"},
	Auth200110:     {Msg: "需進行兩段式驗證"},
	Auth200111:     {Msg: "兩段式驗證token與要求的內容不符合"},
	Auth200112:     {Msg: "兩段式驗證信件送出失敗"},
	Auth200113:     {Msg: "兩段式驗證token已超過次數"},
	Auth200114:     {Msg: "兩段式驗證Code輸入錯誤"},
	Auth200115:     {Msg: "此使用者已刪除"},
	Auth200116:     {Msg: "兩段式驗證與執行的動作不符"},
	Request900101:  {Msg: "資料格式不吻合", Level: logging.Error},
	Database900501: {Msg: "資料庫連線失敗", Level: logging.Error},
	Database900502: {Msg: "資料庫錯誤", Level: logging.Error},
	JWT900601:      {Msg: "JWT產生失敗", Level: logging.Error},
	JSON900701:     {Msg: "JSON parse to byte Fail"},
	JSON900702:     {Msg: "JSON parse to struct Fail"},
	API900801:      {Msg: "API Call Error"},
	Uknown999999:   {Msg: "未知錯誤", Level: logging.Error},
}

// 狀態碼開頭為1，系統相關回應
// 狀態碼開頭為2~5，功能相關回應
// 狀態碼開頭為6，第三方模組相關回應
// 狀態碼開頭為7，外部服務相關回應
// 狀態碼開頭為9，錯誤狀況相關回應
