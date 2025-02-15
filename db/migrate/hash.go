package migrate

import "github.com/msw-x/moon/ustring"

func Hash(s string) string {
	return ustring.Sha1(s)
}
