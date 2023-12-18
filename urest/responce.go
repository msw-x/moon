package urest

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/msw-x/moon/ujson"
)

type Responce[T any] struct {
	Status    int
	Data      T
	Error     error
	ErrorCode int

	w           http.ResponseWriter
	body        []byte
	contentType string
	muteError   bool
}

func (o *Responce[T]) Ok() bool {
	return o.Error == nil && o.Status == 0
}

func (o *Responce[T]) MuteError() {
	o.muteError = true
}

func (o *Responce[T]) EmptyData() bool {
	return reflect.TypeOf(o.Data).Size() == 0
}

func (o *Responce[T]) RefineError(prefix string) bool {
	if o.Error == nil {
		return false
	}
	o.Error = fmt.Errorf("%s: %v", prefix, o.Error)
	return true
}

func (o *Responce[T]) RefineErrorf(s string, args ...any) bool {
	return o.RefineError(fmt.Sprintf(s, args...))
}

func (o *Responce[T]) RefineBadRequest(prefix string) bool {
	if o.RefineError(prefix) {
		o.SetBadRequest()
		return true
	}
	return false
}

func (o *Responce[T]) SetBadRequest() {
	o.Status = http.StatusBadRequest
}

func (o *Responce[T]) SetForbidden() {
	o.Status = http.StatusForbidden
}

func (o *Responce[T]) SetUnauthorized() {
	o.Status = http.StatusUnauthorized
}

func (o *Responce[T]) SetNotAcceptable() {
	o.Status = http.StatusNotAcceptable
}

func (o *Responce[T]) BadRequest(s string) {
	o.SetBadRequest()
	o.Error = errors.New(s)
}

func (o *Responce[T]) Forbidden(s string) {
	o.SetForbidden()
	o.Error = errors.New(s)
}

func (o *Responce[T]) Unauthorized(s string) {
	o.SetUnauthorized()
	o.Error = errors.New(s)
}

func (o *Responce[T]) NotAcceptable(s string) {
	o.SetNotAcceptable()
	o.Error = errors.New(s)
}

func (o *Responce[T]) BadRequestf(s string, args ...any) {
	o.BadRequest(fmt.Sprintf(s, args...))
}

func (o *Responce[T]) Forbiddenf(s string, args ...any) {
	o.Forbidden(fmt.Sprintf(s, args...))
}

func (o *Responce[T]) Unauthorizedf(s string, args ...any) {
	o.Unauthorized(fmt.Sprintf(s, args...))
}

func (o *Responce[T]) NotAcceptablef(s string, args ...any) {
	o.NotAcceptable(fmt.Sprintf(s, args...))
}

func (o *Responce[T]) makeContent() {
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
			default:
				o.contentType = "application/json"
				ujson.InitNilArray(&o.Data)
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
		if v.Error != "" || v.ErrorCode != 0 {
			o.body, _ = ujson.MarshalLowerCase(v)
		}
	}
	return
}

func (o *Responce[T]) send() {
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
