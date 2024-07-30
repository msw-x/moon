package urest

import (
	"fmt"
	"net/http"

	"github.com/msw-x/moon/uerr"
)

func execHeader[ResponseData any](h http.Header, w *Response[ResponseData], f func(http.Header, *Response[ResponseData])) {
	defer uerr.Recover(func(s string) {
		w.Error = fmt.Errorf("header: %s", s)
	})
	if f != nil {
		f(h, w)
	}
}

func execHandle[RequestData any, ResponseData any](r Request[RequestData], w *Response[ResponseData], f func(Request[RequestData], *Response[ResponseData])) {
	defer uerr.Recover(func(s string) {
		w.Error = fmt.Errorf("handle: %s", s)
	})
	if f != nil {
		f(r, w)
	}
}
