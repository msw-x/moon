package uhttp

import "github.com/msw-x/moon/ufmt"

func urlJoin(s ...any) string {
	return ufmt.NotableJoinWith("/", s...)
}
