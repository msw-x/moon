package urest

import (
	"net/http"
)

func AuthHandle[Account any, Session any, RequestData any, ResponseData any](
	ctx *AuthContext[Account, Session],
	method string,
	url string,
	handle func(AuthRequest[Account, Session, RequestData], *Response[ResponseData]),
) {
	var request AuthRequest[Account, Session, RequestData]
	Handle(ctx.Base, method, url, func(h http.Header, w *Response[ResponseData]) {
		request.Account, request.Session, w.Error = ctx.Auth(h)
		if !w.Ok() {
			w.SetUnauthorized()
		}
	}, func(r Request[RequestData], w *Response[ResponseData]) {
		request.Request = r
		handle(request, w)
	})
}

func AuthGet[Account any, Session any, RequestData any, ResponseData any](
	ctx *AuthContext[Account, Session],
	url string,
	handle func(AuthRequest[Account, Session, RequestData], *Response[ResponseData]),
) {
	AuthHandle(ctx, http.MethodGet, url, handle)
}

func AuthPost[Account any, Session any, RequestData any, ResponseData any](
	ctx *AuthContext[Account, Session],
	url string,
	handle func(AuthRequest[Account, Session, RequestData], *Response[ResponseData]),
) {
	AuthHandle(ctx, http.MethodPost, url, handle)
}

func AuthPut[Account any, Session any, RequestData any, ResponseData any](
	ctx *AuthContext[Account, Session],
	url string,
	handle func(AuthRequest[Account, Session, RequestData], *Response[ResponseData]),
) {
	AuthHandle(ctx, http.MethodPut, url, handle)
}

func AuthDelete[Account any, Session any, RequestData any, ResponseData any](
	ctx *AuthContext[Account, Session],
	url string,
	handle func(AuthRequest[Account, Session, RequestData], *Response[ResponseData]),
) {
	AuthHandle(ctx, http.MethodDelete, url, handle)
}
