package uhttp

import (
	"bufio"
	"bytes"
	"net/http"
	"net/url"
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
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	o.Header.Write(o.Header)
	return string(b)
}

func (o *Request) BodyString() string {
	return string(o.Body)
}
