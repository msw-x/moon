package uhttp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/msw-x/moon/parse"
	"github.com/msw-x/moon/ufmt"
)

type Responce struct {
	Request    Request
	Time       time.Duration
	Status     string
	StatusCode int
	Header     http.Header
	Body       []byte
	Error      error
}

func (o *Responce) Ok() bool {
	return o.StatusCode == http.StatusOK && o.Error == nil
}

func (o *Responce) BodyExists() bool {
	return len(o.Body) > 0
}

func (o *Responce) Text() string {
	return string(o.Body)
}

func (o *Responce) Json(v any) error {
	return json.Unmarshal(o.Body, v)
}

func (o *Responce) HeaderTo(v any) error {
	return HeaderTo(o.Header, v)
}

func (o *Responce) HeaderExists(key string) bool {
	return o.HeaderValue(key) != ""
}

func (o *Responce) HeaderValue(key string) string {
	return o.Header.Get(key)
}

func (o *Responce) HeaderInt64(key string) (int64, error) {
	return parse.Int64(o.HeaderValue(key))
}

func (o *Responce) HeaderFloat64(key string) (float64, error) {
	return parse.Float64(o.HeaderValue(key))
}

func (o *Responce) RefineError(text string, err error) {
	o.Error = fmt.Errorf("%s: %v", text, err)
}

func (o *Responce) Title() (s string) {
	s = fmt.Sprintf("%s[%s]", o.Request.Method, o.Request.Url)
	if o.StatusCode != 0 {
		s = ufmt.Join(s, o.Status, o.Time.Truncate(time.Millisecond), ufmt.ByteSizeDense(len(o.Body)))
	}
	if o.Error != nil {
		s = ufmt.Join(s, o.Error)
	}
	return
}

func (o *Responce) Format(f Format) string {
	l := []string{o.Title()}
	push := func(ok bool, name string, value string) {
		if ok && value != "" {
			l = append(l, fmt.Sprintf("%s: %s", name, value))
		}
	}
	push(f.RequestParams, "request-params", o.Request.ParamsString())
	push(f.RequestHeader, "request-header", o.Request.HeaderString())
	push(f.RequestBody, "request-body", o.Request.BodyString())
	push(f.ResponceHeader, "responce-header", o.HeaderString())
	push(f.ResponceBody, "responce-body", string(o.Body))
	return ufmt.NotableJoinSliceWith("\n", l)
}

func (o *Responce) HeaderString() string {
	return HeaderString(o.Header)
}
