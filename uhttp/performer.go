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
	trace   func(Response)
	errors  OnErrors
}

/// timeout на каждый запрос

/// отмена выполнения запроса?

/// Skip io.EOF if Content-Lenght equal to len(Body)
/// SkipEof()

func (o *Performer) Do() (r Response) {
	o.errors.init(&r)
	r.Request = o.Request
	ts := time.Now()
	r.Request.RefineUrl()
	request, err := http.NewRequest(r.Request.Method, r.Request.Uri(), bytes.NewReader(r.Request.Body))
	request.Header = r.Request.Header
	if err == nil {
		var response *http.Response
		// https://pkg.go.dev/net/http#Client.Do
		// Any returned error will be of type *url.Error.
		// The url.Error value's Timeout method will report true if the request timed out.
		response, err = o.c.Do(request)
		if err == nil {
			// If the returned error is nil, the Response will contain a non-nil Body which the user is expected to close.
			// If the Body is not both read to EOF and closed, the Client's underlying RoundTripper (typically Transport)
			// may not be able to re-use a persistent TCP connection to the server for a subsequent "keep-alive" request.
			defer response.Body.Close()
			r.Header = response.Header
			r.Status = response.Status
			r.StatusCode = response.StatusCode
			r.Body, err = io.ReadAll(response.Body)
			if err != nil {
				o.errors.readBody(err)
				/// but if len(r.Body) == http.Content-Length ?
			}
		} else {
			o.errors.doRequest(err)
		}
	} else {
		o.errors.initRequest(err)
	}
	r.Time = time.Since(ts)
	if o.trace != nil {
		o.trace(r) /// print Content-Length
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
	ujson.InitNilSlice(&v)
	body, _ := ujson.MarshalLowerCase(v)
	return o.Body(body)
}

func (o *Performer) Trace(trace func(Response)) *Performer {
	o.trace = trace
	return o
}

func (o *Performer) OnInitRequestError(f OnError) *Performer {
	o.errors.InitRequest = f
	return o
}

func (o *Performer) OnDoRequestError(f OnError) *Performer {
	o.errors.DoRequest = f
	return o
}

func (o *Performer) OnReadBodyError(f OnError) *Performer {
	o.errors.ReadBody = f
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

func traceReadAll(r io.Reader) ([]byte, error) {
	l := make([]time.Duration, 100)
	b := make([]byte, 0, 512)
	for {
		t := time.Now()
		n, err := r.Read(b[len(b):cap(b)])
		l = append(l, time.Since(t))
		b = b[:len(b)+n]
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			if err.Error() == "context deadline exceeded (Client.Timeout or context cancellation while reading body)" {
				ulog.Trace("read body timeout:", l)
			}
			return b, err
		}
		if len(b) == cap(b) {
			b = append(b, 0)[:len(b)]
		}
	}
}
