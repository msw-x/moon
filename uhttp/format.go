package uhttp

import (
	"fmt"

	"github.com/msw-x/moon/ufmt"
)

type Format struct {
	RequestParams     bool
	RequestHeader     bool
	ResponceHeader    bool
	RequestBody       bool
	ResponceBody      bool
	RequestBodyTrim   bool
	ResponceBodyTrim  bool
	RequestBodyLimit  int
	ResponceBodyLimit int
}

type FormatProvider struct {
	Title          func() string
	RequestParams  func() string
	RequestHeader  func() string
	RequestBody    func() string
	ResponceHeader func() string
	ResponceBody   func() string
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
	push(f.ResponceHeader, 0, false, "responce-header", o.ResponceHeader)
	push(f.ResponceBody, f.ResponceBodyLimit, f.ResponceBodyTrim, "responce-body", o.ResponceBody)
	return ufmt.NotableJoinSliceWith("\n", l)
}
