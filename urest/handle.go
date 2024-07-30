package urest

import (
	"net/http"

	"github.com/msw-x/moon/uhttp"
	"github.com/msw-x/moon/utime"
)

func Handle[RequestData any, ResponseData any](
	ctx *Context,
	method string,
	url string,
	header func(http.Header, *Response[ResponseData]),
	handle func(Request[RequestData], *Response[ResponseData]),
) {
	if ctx.allowCors {
		ctx.router.Options(url, func(w http.ResponseWriter, r *http.Request) {
			allowCors(w)
		})
	}
	ctx.router.Handle(method, url, func(w http.ResponseWriter, r *http.Request) {
		tm := utime.NewStopwatch()
		var f uhttp.FormatProvider
		f.RequestParams = func() string {
			return uhttp.ParamsString(r.URL.Query())
		}
		f.RequestHeader = func() string {
			return uhttp.HeaderString(r.Header)
		}
		if ctx.allowCors {
			allowCors(w)
		}
		var request Request[RequestData]
		var response Response[ResponseData]
		execHeader(r.Header, &response, header)
		request.r = r
		response.w = w
		if response.Ok() {
			response.Error = request.readBody()
			if response.Ok() {
				if request.HasBody() {
					response.Error = request.DataFromJson()
					if !response.Ok() {
						response.RefineBadRequest("unmarshal json")
					}
				} else if !request.EmptyData() {
					if method == http.MethodPut || method == http.MethodPost {
						response.BadRequest("body is empty")
					} else {
						request.DataFromParams()
					}
				}
				if response.Ok() {
					execHandle(request, &response, handle)
				}
			} else {
				response.w = nil
				response.RefineError("read body")
			}
		}
		response.send()
		Trace(ctx, request, response, tm.Time())
	})
}

func Get[RequestData any, ResponseData any](
	ctx *Context,
	url string,
	handle func(Request[RequestData], *Response[ResponseData]),
) {
	Handle(ctx, http.MethodGet, url, nil, handle)
}

func Post[RequestData any, ResponseData any](
	ctx *Context,
	url string,
	handle func(Request[RequestData], *Response[ResponseData]),
) {
	Handle(ctx, http.MethodPost, url, nil, handle)
}

func Put[RequestData any, ResponseData any](
	ctx *Context,
	url string,
	handle func(Request[RequestData], *Response[ResponseData]),
) {
	Handle(ctx, http.MethodPut, url, nil, handle)
}

func Delete[RequestData any, ResponseData any](
	ctx *Context,
	url string,
	handle func(Request[RequestData], *Response[ResponseData]),
) {
	Handle(ctx, http.MethodDelete, url, nil, handle)
}
