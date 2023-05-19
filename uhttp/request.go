package uhttp

import (
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/msw-x/moon/ufmt"
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
	if len(o.Header) == 0 {
		return ""
	}
	var l []string
	for name, vals := range o.Header {
		l = append(l, fmt.Sprintf("%s%v", name, vals))
	}
	sort.Slice(l, func(i, j int) bool {
		return l[i] < l[j]
	})
	return ufmt.JoinSliceWith("\n", l)
}

func (o *Request) BodyString() string {
	return string(o.Body)
}
