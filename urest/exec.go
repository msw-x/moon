package urest

import (
	"fmt"
	"net/http"

	"github.com/msw-x/moon/uerr"
)

func execHeader[ResponceData any](h http.Header, w *Responce[ResponceData], f func(http.Header, *Responce[ResponceData])) {
	defer uerr.Recover(func(s string) {
		w.Error = fmt.Errorf("header: %s", s)
	})
	if f != nil {
		f(h, w)
	}
}

func execHandle[RequestData any, ResponceData any](r Request[RequestData], w *Responce[ResponceData], f func(Request[RequestData], *Responce[ResponceData])) {
	defer uerr.Recover(func(s string) {
		w.Error = fmt.Errorf("handle: %s", s)
	})
	if f != nil {
		f(r, w)
	}
}
