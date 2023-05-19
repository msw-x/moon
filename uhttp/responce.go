package uhttp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/msw-x/moon/ufmt"
)

type Responce struct {
	Request    Request
	Time       time.Duration
	Status     string
	StatusCode int
	Body       []byte
	Error      error
}

func (o *Responce) Ok() bool {
	return o.StatusCode == http.StatusOK && o.Error == nil
}

func (o *Responce) Json(v any) error {
	return json.Unmarshal(o.Body, v)
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

func (o *Responce) Format(f Format) (s string) {
	delim := "\n"
	push := func(ok bool, value string) {
		if ok && value != "" {
			s += delim + value
		}
	}
	s = o.Title()
	push(f.Params, o.Request.ParamsString())
	push(f.Header, o.Request.HeaderString())
	push(f.RequestBody, o.Request.BodyString())
	push(f.ResponceBody, string(o.Body))
	return
}
