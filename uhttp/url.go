package uhttp

import "github.com/msw-x/moon/ufmt"

func UrlJoin(s ...any) string {
	return ufmt.NotableJoinWith("/", s...)
}
