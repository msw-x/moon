package uhttp

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/msw-x/moon/refl"
	"github.com/msw-x/moon/ujson"
	"github.com/msw-x/moon/ulog"
	"github.com/msw-x/moon/ustring"
)

type Performer struct {
	Request Request
	c       *http.Client
	trace   func(Responce)
}

func (o *Performer) Do() (r Responce) {
	r.Request = o.Request
	ts := time.Now()
	r.Request.RefineUrl()
	request, err := http.NewRequest(r.Request.Method, r.Request.Uri(), bytes.NewReader(r.Request.Body))
	request.Header = r.Request.Header
	if err == nil {
		responce, err := o.c.Do(request)
		if err == nil {
			defer responce.Body.Close()
			r.Header = responce.Header
			r.Status = responce.Status
			r.StatusCode = responce.StatusCode
			r.Body, err = io.ReadAll(responce.Body)
			if err != nil {
				ok := err.Error() == "context deadline exceeded (Client.Timeout or context cancellation while reading body)" && len(r.Body) > 0
				if ok {
					ulog.Trace("read body: context deadline exceeded")
				} else {
					r.RefineError("read body", err)
				}
			}
		} else {
			r.RefineError("do request", err)
		}
	} else {
		r.RefineError("init request", err)
	}
	r.Time = time.Since(ts)
	if o.trace != nil {
		o.trace(r)
	}
	return
}

func (o *Performer) Param(name string, value any) *Performer {
	return o.param(name, value, false)
}

func (o *Performer) Params(s any) *Performer {
	refl.WalkOnTagsAny(s, UrlTag, func(v any, name string, flags []string) {
		o.param(name, v, OmitEmpty(flags))
	})
	return o
}

func (o *Performer) Header(name string, value any) *Performer {
	return o.header(name, value, false)
}

func (o *Performer) Headers(s any) *Performer {
	refl.WalkOnTagsAny(s, HeaderTag, func(v any, name string, flags []string) {
		o.header(name, v, OmitEmpty(flags))
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
	o.Request.Body = body
	return o
}

func (o *Performer) Json(v any) *Performer {
	o.ContentTypeJson()
	body, _ := ujson.MarshalLowerCase(v)
	return o.Body(body)
}

func (o *Performer) Trace(trace func(Responce)) *Performer {
	o.trace = trace
	return o
}

func (o *Performer) param(name string, value any, omitempty bool) *Performer {
	if o.Request.Params == nil {
		o.Request.Params = make(url.Values)
	}
	if v, omit := Marshal(value, omitempty); !omit {
		name = ustring.TitleLowerCase(name)
		o.Request.Params.Set(name, v)
	}
	return o
}

func (o *Performer) header(name string, value any, omitempty bool) *Performer {
	if o.Request.Header == nil {
		o.Request.Header = make(http.Header)
	}
	if v, omit := Marshal(value, omitempty); !omit {
		o.Request.Header.Set(name, v)
	}
	return o
}
