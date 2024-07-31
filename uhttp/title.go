package uhttp

import (
	"time"

	"github.com/msw-x/moon/ufmt"
	"github.com/msw-x/moon/utime"
)

func Title(name string, statusCode int, status string, tm time.Duration, requestBodyLen int, responseBodyLen int, err error) string {
	l := []any{name}
	if statusCode != 0 {
		l = append(l, status)
	}
	if requestBodyLen != 0 {
		l = append(l, ufmt.ByteSizeDense(requestBodyLen))
	}
	if tm != 0 {
		l = append(l, utime.PrettyTruncate(tm))
	}
	if responseBodyLen != 0 {
		l = append(l, ufmt.ByteSizeDense(responseBodyLen))
	}
	if err != nil {
		l = append(l, err)
	}
	return ufmt.JoinSlice(l)
}
