package uhttp

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/msw-x/moon/ufmt"
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
	if method == "" {
		method = "*"
	}
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

func ProxyRequestNameDefault(r *http.Request) string {
	return ProxyRequestName(r, XForwardedFor)
}

func WebSocketName(name string) string {
	return ufmt.JoinWith("@", name, "ws")
}
