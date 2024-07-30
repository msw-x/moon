package uhttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/msw-x/moon/parse"
)

type Response struct {
	Request    Request
	Time       time.Duration
	Status     string
	StatusCode int
	Header     http.Header
	Body       []byte
	Error      error
}

func (o Response) Ok() bool {
	return o.StatusCode == http.StatusOK && o.Error == nil
}

func (o Response) NotOkError() error {
	if o.Ok() {
		return nil
	}
	if o.Error == nil {
		return errors.New(o.GetStatus())
	}
	return o.Error
}

func (o Response) GetStatus() string {
	if o.Status == "" {
		return strconv.Itoa(o.StatusCode)
	}
	return o.Status
}

func (o Response) BodyExists() bool {
	return len(o.Body) > 0
}

func (o Response) Text() string {
	return string(o.Body)
}

func (o Response) Json(v any) error {
	return json.Unmarshal(o.Body, v)
}

func (o Response) HeaderTo(v any) error {
	return HeaderTo(o.Header, v)
}

func (o Response) HeaderExists(key string) bool {
	return o.HeaderValue(key) != ""
}

func (o Response) HeaderValue(key string) string {
	return o.Header.Get(key)
}

func (o Response) HeaderInt64(key string) (int64, error) {
	return parse.Int64(o.HeaderValue(key))
}

func (o Response) HeaderFloat64(key string) (float64, error) {
	return parse.Float64(o.HeaderValue(key))
}

func (o Response) RefineError(text string, err error) {
	o.Error = fmt.Errorf("%s: %v", text, err)
}

func (o Response) Title() string {
	return Title(ClientRequestName(o.Request), o.StatusCode, o.Status, o.Time, len(o.Body), o.Error)
}

func (o Response) Format(f Format) string {
	return FormatProvider{
		Title:          o.Title,
		RequestParams:  o.Request.ParamsString,
		RequestHeader:  o.Request.HeaderString,
		RequestBody:    o.Request.BodyString,
		ResponseHeader: o.HeaderString,
		ResponseBody:   o.BodyString,
	}.Format(f)
}

func (o Response) HeaderString() string {
	return HeaderString(o.Header)
}

func (o Response) BodyString() string {
	return string(o.Body)
}
