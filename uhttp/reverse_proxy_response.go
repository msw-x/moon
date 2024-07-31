package uhttp

import (
	"bytes"
	"net/http"
	"strconv"
	"time"
)

type ReverseProxyResponse struct {
	r           *http.Request
	w           http.ResponseWriter
	requestBody []byte
	ts          time.Time
	buf         bytes.Buffer
	statusCode  int
	err         error
	tracer      *Tracer[ReverseProxyResponse]
	router      *Router
}

func NewReverseProxyResponse(r *http.Request, w http.ResponseWriter, tracer *Tracer[ReverseProxyResponse], router *Router) *ReverseProxyResponse {
	o := new(ReverseProxyResponse)
	o.r = r
	o.w = w
	o.ts = time.Now()
	o.tracer = tracer
	o.router = router
	if tracer != nil && tracer.RequireRequestBody() {
		o.requestBody, o.err = DumpBody(r)
	}
	return o
}

func (o *ReverseProxyResponse) Header() http.Header {
	return o.w.Header()
}

func (o *ReverseProxyResponse) Write(v []byte) (int, error) {
	if o.tracer != nil && o.tracer.RequireResponseBody() {
		o.buf.Write(v)
	}
	return o.w.Write(v)
}

func (o *ReverseProxyResponse) WriteHeader(statusCode int) {
	o.statusCode = statusCode
	o.w.WriteHeader(statusCode)
}

func (o *ReverseProxyResponse) SetError(err error) {
	o.err = err
}

func (o ReverseProxyResponse) ErrorFree() bool {
	return o.err == nil
}

func (o ReverseProxyResponse) Ok() bool {
	return o.err == nil && o.statusCode == http.StatusOK
}

func (o ReverseProxyResponse) Format(f Format) string {
	responseBody := o.buf.Bytes()
	return FormatProvider{
		Title: func() string {
			var err error
			responseBodyLen := len(responseBody)
			if responseBodyLen == 0 {
				err = o.err
			}
			return Title(o.router.RequestName(o.r), o.statusCode, strconv.Itoa(o.statusCode), time.Since(o.ts), responseBodyLen, err)
		},
		RequestParams: func() string {
			return ParamsString(o.r.URL.Query())
		},
		RequestHeader: func() string {
			return HeaderString(o.r.Header)
		},
		RequestBody: func() string {
			return string(o.requestBody)
		},
		ResponseHeader: func() string {
			return HeaderString(o.Header())
		},
		ResponseBody: func() string {
			return string(responseBody)
		},
	}.Format(f)
}

func (o ReverseProxyResponse) Close() {
	if o.err != nil {
		if o.statusCode == 0 || o.statusCode == http.StatusOK {
			o.WriteHeader(http.StatusBadGateway)
		}
		o.Write([]byte(o.err.Error()))
	}
	if o.tracer != nil {
		o.tracer.Trace(o)
	}
}
