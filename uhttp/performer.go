package uhttp

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/msw-x/moon/rt"
	"github.com/msw-x/moon/ujson"
	"github.com/msw-x/moon/ustring"
)

type Performer struct {
	c *http.Client
	r Request
	t func(Responce)
}

func (o *Performer) Do() (r Responce) {
	r.Request = o.r
	ts := time.Now()
	r.Request.RefineUrl()
	request, err := http.NewRequest(r.Request.Method, r.Request.Uri(), bytes.NewReader(r.Request.Body))
	if err == nil {
		responce, err := o.c.Do(request)
		if err == nil {
			defer responce.Body.Close()
			r.Status = responce.Status
			r.StatusCode = responce.StatusCode
			r.Body, err = io.ReadAll(responce.Body)
			if err != nil {
				r.RefineError("read body", err)
			}
		} else {
			r.RefineError("do request", err)
		}
	} else {
		r.RefineError("init request", err)
	}
	r.Time = time.Since(ts)
	if o.t != nil {
		o.t(r)
	}
	return
}

func (o *Performer) Param(name string, value any) *Performer {
	if o.r.Params == nil {
		o.r.Params = make(url.Values)
	}
	name = ustring.TitleLowerCase(name)
	o.r.Params.Set(name, fmt.Sprint(value))
	return o
}

func (o *Performer) Params(s any) *Performer {
	rt.PlainValues(s, "url", func(v any, name string, flags []string) {
		o.Param(name, v)
	})
	return o
}

func (o *Performer) Header(name string, value any) *Performer {
	if o.r.Header == nil {
		o.r.Header = make(http.Header)
	}
	o.r.Header.Set(name, fmt.Sprint(value))
	return o
}

func (o *Performer) Headers(s any) *Performer {
	rt.PlainValues(s, "http", func(v any, name string, flags []string) {
		o.Header(name, v)
	})
	return o
}

func (o *Performer) ContentType(v string) *Performer {
	return o.Header("Content-Type", v)
}

func (o *Performer) ContentTypeJson() *Performer {
	return o.ContentType("application/json")
}

func (o *Performer) Auth(v string) *Performer {
	return o.Header("Authorization", v)
}

func (o *Performer) AuthBearer(token string) *Performer {
	return o.Auth("Bearer " + token)
}

func (o *Performer) Body(body []byte) *Performer {
	o.r.Body = body
	return o
}

func (o *Performer) Json(v any) *Performer {
	o.ContentTypeJson()
	body, _ := ujson.MarshalLowerCase(v)
	return o.Body(body)
}

func (o *Performer) Trace(t func(Responce)) *Performer {
	o.t = t
	return o
}
