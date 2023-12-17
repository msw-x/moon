package urest

import (
	"net/http"

	"github.com/msw-x/moon/uhttp"
	"github.com/msw-x/moon/utime"
)

func Handle[RequestData any, ResponceData any](
	ctx Context,
	method string,
	url string,
	header func(http.Header, *Responce[ResponceData]),
	handle func(Request[RequestData], *Responce[ResponceData]),
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
		var responce Responce[ResponceData]
		execHeader(r.Header, &responce, header)
		request.r = r
		responce.w = w
		if responce.Ok() {
			responce.Error = request.readBody()
			if responce.Ok() {
				if request.HasBody() {
					responce.Error = request.DataFromJson()
					if !responce.Ok() {
						responce.RefineBadRequest("unmarshal json")
					}
				} else if !request.EmptyData() {
					if method == http.MethodPut || method == http.MethodPost {
						responce.BadRequest("body is empty")
					} else {
						request.DataFromParams()
					}
				}
				if responce.Ok() {
					execHandle(request, &responce, handle)
				}
			} else {
				responce.w = nil
				responce.RefineError("read body")
			}
		}
		responce.send()
		Trace(ctx, request, responce, tm.Time())
	})
}

func Get[RequestData any, ResponceData any](ctx Context, url string, handle func(Request[RequestData], *Responce[ResponceData])) {
	Handle(ctx, http.MethodGet, url, nil, handle)
}

func Post[RequestData any, ResponceData any](ctx Context, url string, handle func(Request[RequestData], *Responce[ResponceData])) {
	Handle(ctx, http.MethodPost, url, nil, handle)
}

func Put[RequestData any, ResponceData any](ctx Context, url string, handle func(Request[RequestData], *Responce[ResponceData])) {
	Handle(ctx, http.MethodPut, url, nil, handle)
}

func Delete[RequestData any, ResponceData any](ctx Context, url string, handle func(Request[RequestData], *Responce[ResponceData])) {
	Handle(ctx, http.MethodDelete, url, nil, handle)
}
