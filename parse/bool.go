package parse

import (
	"github.com/msw-x/moon"
	"strings"
)

func Bool(s string) bool {
	s = strings.ToLower(s)
	switch s {
	case "0", "false", "no", "off", "disable":
		return false
	case "1", "true", "yes", "on", "enable":
		return true
	}
	moon.Panic("parse bool:", s)
	return false
}
