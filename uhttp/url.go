package uhttp

func urlJoin(s ...string) string {
	return ustring.NotableJoinWith("/", s...)
}
