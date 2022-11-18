package parse

import (
	"fmt"
	"strings"

	"github.com/msw-x/moon"
)

func Bool(s string) (bool, error) {
	s = strings.ToLower(s)
	switch s {
	case "0", "false", "no", "off", "disable":
		return false, nil
	case "1", "true", "yes", "on", "enable":
		return true, nil
	}
	return false, fmt.Errorf("parse bool: %S", s)
}

func BoolStrict(s string) bool {
	b, err := Bool(s)
	moon.Strict(err)
	return b
}
