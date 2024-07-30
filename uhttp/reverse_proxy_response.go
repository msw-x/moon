package uhttp

import (
	"net/http"
	"strconv"
	"time"
)

type ReverseProxyResponce struct {
	Request    *http.Request
	Header     http.Header
	StatusCode int
	Status     string
	Time       time.Duration
	Error      error

	router *Router
}

func (o ReverseProxyResponce) Ok() bool {
	return o.Error == nil && o.StatusCode == http.StatusOK
}

func (o ReverseProxyResponce) Format(f Format) string {
	return FormatProvider{
		Title: func() string {
			return Title(o.router.RequestName(o.Request), o.StatusCode, strconv.Itoa(o.StatusCode), o.Time, 0, o.Error)
		},
		RequestParams: func() string {
			return ParamsString(o.Request.URL.Query())
		},
		RequestHeader: func() string {
			return HeaderString(o.Request.Header)
		},
		RequestBody: func() string {
			return ""
		},
		ResponseHeader: func() string {
			return HeaderString(o.Header)
		},
		ResponseBody: func() string {
			return ""
		},
	}.Format(f)
}
