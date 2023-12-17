package uhttp

import "net/url"

func ParamsString(params url.Values) string {
	s := params.Encode()
	if s != "" {
		s = "?" + s
	}
	return s
}
