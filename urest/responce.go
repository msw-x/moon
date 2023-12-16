package urest

import (
	"errors"
	"fmt"
	"net/http"
)

type Responce[T any] struct {
	Status    int
	Data      T
	Error     error
	ErrorCode int

	muteError bool
}

func (o *Responce[T]) Ok() bool {
	return o.Error == nil && o.Status == 0
}

func (o *Responce[T]) MuteError() {
	o.muteError = true
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
