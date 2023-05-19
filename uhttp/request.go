package uhttp

import (
	"net/http"
	"net/url"
	"strings"
)

type Request struct {
	Method string
	Url    string
	Params url.Values
	Header http.Header
	Body   []byte
}

func (o *Request) Path(path string) {
	o.Url = urlJoin(o.Url, path)
}

func (o *Request) RefineUrl() {
	if !strings.Contains(o.Url, "://") {
		o.Url = "https://" + o.Url
	}
}

func (o *Request) Uri() string {
	return o.Url + o.ParamsString()
}

func (o *Request) ParamsString() string {
	s := o.Params.Encode()
	if s != "" {
		s = "?" + s
	}
	return s
}

func (o *Request) HeaderString() string {
	return HeaderString(o.Header)
}

func (o *Request) BodyString() string {
	return string(o.Body)
}
