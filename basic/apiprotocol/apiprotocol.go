package apiprotocol

import "time"

func (code Code) ToRes() *BaseResponse {
	return &BaseResponse{
		RetStatus: RetStatus{
			Code:       code,
			SystemTime: time.Now().UnixMilli(),
		},
	}
}

func (code Code) ToRet() *RetStatus {
	msg := "未知的錯誤"
	retStatus, exist := retStatusList[code]
	if exist {
		msg = retStatus.Msg
	}
	return &RetStatus{
		Code:       code,
		Msg:        msg,
		SystemTime: time.Now().UnixMilli(),
	}

}

func (br *BaseResponse) Err(err string) *BaseResponse {
	br.RetStatus.Error = err
	return br
}

func (br *BaseResponse) ToErr() error {
	if br.RetStatus.Code == Success10000 {
		return nil
	}
	return br
}

func (br *BaseResponse) Msg() string {
	return br.RetStatus.Msg
}
func (br *BaseResponse) Error() string {
	return br.RetStatus.Error
}
