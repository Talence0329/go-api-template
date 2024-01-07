package tools

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/valyala/fasthttp"
)

// FileByte : 讀取檔案並轉成byte
func FileByte(path string) (data []byte, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	data, err = ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return data, err
}

// Do : 呼叫API
func (request Request) Do(resDatas ...interface{}) (response *Response, err error) {
	res := fasthttp.AcquireResponse()
	req := request.prepareRequest()

	client := &fasthttp.Client{}
	err = client.Do(req, res)
	if err != nil {
		return response, err
	}

	response = &Response{
		Body:       res.Body(),
		Header:     &res.Header,
		StatusCode: res.StatusCode(),
	}
	if len(resDatas) != 0 && resDatas[0] != nil {
		err = json.Unmarshal(response.Body, resDatas[0])
		if err != nil {
			return response, err
		}
	}

	return response, nil
}

// prepareRequest: 準備req結構
func (request *Request) prepareRequest() *fasthttp.Request {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(request.Url)
	req.Header.SetMethod(request.Method)
	req.Header.SetContentType(request.ContentType)
	for key, val := range request.Headers {
		req.Header.Add(key, val)
	}
	if request.Body != nil {
		req.AppendBody(request.Body)
	}

	return req
}

const LogLevelError = "error"
const LogLevelWarning = "warning"
const LogLevelNotice = "notice"
const LogLevelEmergency = "emergency"

func Log(errLevel string, format string, params ...any) (n int, err error) {
	return fmt.Printf(format, params...)
}
