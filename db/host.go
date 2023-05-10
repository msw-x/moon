package db

import (
	"fmt"
	"strings"
)

func host(host string) string {
	if strings.Contains(host, ":") {
		return host
	}
	return fmt.Sprintf("%s:5432", host)
}
