package uhttp

import (
	"fmt"

	"github.com/msw-x/moon/ufmt"
)

type Format struct {
	RequestParams     bool
	RequestHeader     bool
	ResponseHeader    bool
	RequestBody       bool
	ResponseBody      bool
	RequestBodyTrim   bool
	ResponseBodyTrim  bool
	RequestBodyLimit  int
	ResponseBodyLimit int
}

type FormatProvider struct {
	Title          func() string
	RequestParams  func() string
	RequestHeader  func() string
	RequestBody    func() string
	ResponseHeader func() string
	ResponseBody   func() string
}

func (o FormatProvider) Format(f Format) string {
	l := []string{o.Title()}
	push := func(ok bool, limit int, trim bool, name string, f func() string) {
		if ok && f != nil {
			v := f()
			if v != "" {
				if limit != 0 && len(v) > limit {
					m := fmt.Sprintf("trace limit exceeded: %s / %s", ufmt.ByteSize(len(v)), ufmt.ByteSize(limit))
					if trim {
						v = v[0:limit] + "...\n" + m
					} else {
						v = m
					}
				}
				l = append(l, fmt.Sprintf("%s: %s", name, v))
			}
		}
	}
	push(f.RequestParams, 0, false, "request-params", o.RequestParams)
	push(f.RequestHeader, 0, false, "request-header", o.RequestHeader)
	push(f.RequestBody, f.RequestBodyLimit, f.RequestBodyTrim, "request-body", o.RequestBody)
	push(f.ResponseHeader, 0, false, "response-header", o.ResponseHeader)
	push(f.ResponseBody, f.ResponseBodyLimit, f.ResponseBodyTrim, "response-body", o.ResponseBody)
	return ufmt.NotableJoinSliceWith("\n", l)
}
