package urest

import (
	"strconv"
	"time"

	"github.com/msw-x/moon/uhttp"
)

func Trace[RequestData any, ResponceData any](ctx Context, r Request[RequestData], w Responce[ResponceData], tm time.Duration) {
	ctx.Trace(uhttp.FormatProvider{
		Title: func() string {
			name := ctx.router.RequestName(r.r)
			statusCode := w.Status
			if w.Ok() && statusCode == 0 {
				statusCode = 200
			}
			return uhttp.Title(name, statusCode, strconv.Itoa(statusCode), tm, len(w.body), w.Error)
		},
		RequestParams: func() string {
			return uhttp.ParamsString(r.r.URL.Query())
		},
		RequestHeader: func() string {
			return uhttp.HeaderString(r.r.Header)
		},
		RequestBody: func() string {
			return string(r.body)
		},
		ResponceHeader: func() string {
			return uhttp.HeaderString(w.w.Header())
		},
		ResponceBody: func() string {
			return string(w.body)
		},
	}, w.Ok() || w.muteError)
}
