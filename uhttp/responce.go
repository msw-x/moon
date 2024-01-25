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

func (o *Responce) NotOkError() error {
	if o.Ok() {
		return nil
	}
	if o.Error == nil {
		return errors.New(o.GetStatus())
	}
	return o.Error
}

func (o *Responce) GetStatus() string {
	if o.Status == "" {
		return strconv.Itoa(o.StatusCode)
	}
	return o.Status
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

func (o *Responce) Title() string {
	return Title(ClientRequestName(o.Request), o.StatusCode, o.Status, o.Time, len(o.Body), o.Error)
}

func (o *Responce) Format(f Format) string {
	return FormatProvider{
		Title:          o.Title,
		RequestParams:  o.Request.ParamsString,
		RequestHeader:  o.Request.HeaderString,
		RequestBody:    o.Request.BodyString,
		ResponceHeader: o.HeaderString,
		ResponceBody:   o.BodyString,
	}.Format(f)
}

func (o *Responce) HeaderString() string {
	return HeaderString(o.Header)
}

func (o *Responce) BodyString() string {
	return string(o.Body)
}
