package urest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/msw-x/moon/ujson"
)

type Response[T any] struct {
	Status    int
	Data      T
	Error     error
	ErrorCode int

	w           http.ResponseWriter
	body        []byte
	contentType string
	muteError   bool
	notTrace    bool
}

func (o *Response[T]) Ok() bool {
	return o.Error == nil && o.Status == 0
}

func (o *Response[T]) MuteError() {
	o.muteError = true
}

func (o *Response[T]) NotTrace() {
	o.notTrace = true
}

func (o *Response[T]) EmptyData() bool {
	return reflect.TypeOf(o.Data).Size() == 0
}

func (o *Response[T]) RefineError(prefix string) bool {
	if o.Error == nil {
		return false
	}
	o.Error = fmt.Errorf("%s: %v", prefix, o.Error)
	return true
}

func (o *Response[T]) RefineErrorf(s string, args ...any) bool {
	return o.RefineError(fmt.Sprintf(s, args...))
}

func (o *Response[T]) RefineBadRequest(prefix string) bool {
	if o.RefineError(prefix) {
		o.SetBadRequest()
		return true
	}
	return false
}

func (o *Response[T]) SetBadRequest() {
	o.Status = http.StatusBadRequest
}

func (o *Response[T]) SetForbidden() {
	o.Status = http.StatusForbidden
}

func (o *Response[T]) SetUnauthorized() {
	o.Status = http.StatusUnauthorized
}

func (o *Response[T]) SetNotAcceptable() {
	o.Status = http.StatusNotAcceptable
}

func (o *Response[T]) SetNotImplemented() {
	o.Status = http.StatusNotImplemented
}

func (o *Response[T]) SetNotFound() {
	o.Status = http.StatusNotFound
}

func (o *Response[T]) SetServiceUnavailable() {
	o.Status = http.StatusServiceUnavailable
}

func (o *Response[T]) BadRequest(s string) {
	o.SetBadRequest()
	o.Error = errors.New(s)
}

func (o *Response[T]) Forbidden(s string) {
	o.SetForbidden()
	o.Error = errors.New(s)
}

func (o *Response[T]) Unauthorized(s string) {
	o.SetUnauthorized()
	o.Error = errors.New(s)
}

func (o *Response[T]) NotAcceptable(s string) {
	o.SetNotAcceptable()
	o.Error = errors.New(s)
}

func (o *Response[T]) NotImplemented(s string) {
	o.SetNotImplemented()
	o.Error = errors.New(s)
}

func (o *Response[T]) NotFound(s string) {
	o.SetNotFound()
	o.Error = errors.New(s)
}

func (o *Response[T]) ServiceUnavailable(s string) {
	o.SetServiceUnavailable()
	o.Error = errors.New(s)
}

func (o *Response[T]) BadRequestf(s string, args ...any) {
	o.BadRequest(fmt.Sprintf(s, args...))
}

func (o *Response[T]) Forbiddenf(s string, args ...any) {
	o.Forbidden(fmt.Sprintf(s, args...))
}

func (o *Response[T]) Unauthorizedf(s string, args ...any) {
	o.Unauthorized(fmt.Sprintf(s, args...))
}

func (o *Response[T]) NotAcceptablef(s string, args ...any) {
	o.NotAcceptable(fmt.Sprintf(s, args...))
}

func (o *Response[T]) NotImplementedf(s string, args ...any) {
	o.NotImplemented(fmt.Sprintf(s, args...))
}

func (o *Response[T]) ServiceUnavailablef(s string, args ...any) {
	o.ServiceUnavailable(fmt.Sprintf(s, args...))
}

func (o *Response[T]) makeContent() {
	if o.Ok() {
		if !o.EmptyData() {
			switch v := any(o.Data).(type) {
			case Void:
			case Text:
				o.contentType = "text/plain"
				o.body = []byte(v)
			case Image:
				o.contentType = "image/" + v.Type
				o.body = v.Data
			case json.RawMessage:
				o.contentType = "application/json"
				o.body = v
			default:
				o.contentType = "application/json"
				ujson.InitNilSlice(&o.Data)
				o.body, _ = ujson.MarshalLowerCase(o.Data)
			}
		}
	} else {
		o.contentType = "application/json"
		v := struct {
			Error     string `json:",omitempty"`
			ErrorCode int    `json:",omitempty"`
		}{
			Error:     fmt.Sprint(o.Error),
			ErrorCode: o.ErrorCode,
		}
		if !(o.Error == nil || v.Error == "") || v.ErrorCode != 0 {
			o.body, _ = ujson.MarshalLowerCase(v)
		}
	}
	return
}

func (o *Response[T]) send() {
	if o.w != nil {
		o.makeContent()
		if !o.Ok() && o.Status == 0 {
			o.Status = http.StatusInternalServerError
		}
		if o.Status > 0 {
			o.w.WriteHeader(o.Status)
		}
		if o.contentType != "" {
			o.w.Header().Set("Content-Type", o.contentType)
		}
		if len(o.body) > 0 {
			_, err := o.w.Write(o.body)
			if err != nil {
				o.Error = err
			}
		}
	}
}
