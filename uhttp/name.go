package uhttp

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/msw-x/moon/ufmt"
	"github.com/msw-x/moon/utime"
)

func ClientAddress(r *http.Request, xRemoteAddress string) string {
	if r == nil {
		return ""
	}
	var remoteAddress string
	if xRemoteAddress != "" && r.Header != nil {
		remoteAddress = r.Header.Get(xRemoteAddress)
	}
	if remoteAddress == "" {
		remoteAddress = r.RemoteAddr
	}
	return remoteAddress
}

func ClientRequestName(r Request) string {
	return fmt.Sprintf("%s[%s]", r.Method, r.Url)
}

func RouteName(method, uri string) string {
	return ufmt.JoinWith(":", strings.ToUpper(method), uri)
}

func RequestName(r *http.Request) string {
	return ProxyRequestName(r, "")
}

func ProxyRequestName(r *http.Request, xRemoteAddress string) string {
	if r == nil {
		return "?"
	}
	return ufmt.JoinWith("'", ClientAddress(r, xRemoteAddress), RouteName(r.Method, r.URL.Path))
}

func WebSocketName(name string) string {
	return ufmt.JoinWith("@", name, "ws")
}

func Title(name string, statusCode int, status string, tm time.Duration, bodyLen int, err error) string {
	l := []any{name}
	if statusCode != 0 {
		l = append(l, status)
	}
	if tm != 0 {
		l = append(l, utime.PrettyTruncate(tm))
	}
	if bodyLen != 0 {
		l = append(l, ufmt.ByteSizeDense(bodyLen))
	}
	if err != nil {
		l = append(l, err)
	}
	return ufmt.JoinSlice(l)
}
