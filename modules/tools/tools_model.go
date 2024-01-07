package tools

import "github.com/valyala/fasthttp"

type Request struct {
	Url         string
	Method      string
	ContentType string
	Headers     map[string]string
	Body        []byte
}

type Response struct {
	Body       []byte
	Header     *fasthttp.ResponseHeader
	StatusCode int
}
