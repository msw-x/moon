package urest

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"reflect"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/msw-x/moon/parse"
	"github.com/msw-x/moon/uhttp"
)

type Request[T any] struct {
	Data T

	r    *http.Request
	ca   string
	body []byte
}

func (o Request[T]) EmptyData() bool {
	return reflect.TypeOf(o.Data).Size() == 0
}

func (o Request[T]) HasBody() bool {
	return len(o.body) > 0
}

func (o Request[T]) HeaderString() string {
	return uhttp.HeaderString(o.r.Header)
}

func (o Request[T]) RemoteAddr() string {
	return o.r.RemoteAddr
}

func (o Request[T]) RemoteHost() string {
	return hostFromAddress(o.ClientAddr())
}

func (o Request[T]) ClientAddr() string {
	if o.ca == "" {
		return o.RemoteAddr()
	}
	return o.ca
}

func (o Request[T]) ClientHost() string {
	return hostFromAddress(o.ClientAddr())
}

func (o Request[T]) Var(name string) string {
	return mux.Vars(o.r)[name]
}

func (o Request[T]) VarInt(name string) (int, error) {
	return strconv.Atoi(o.Var(name))
}

func (o Request[T]) VarIntDef(name string) int {
	i, _ := o.VarInt(name)
	return i
}

func (o Request[T]) VarInt64(name string) (int64, error) {
	return strconv.ParseInt(o.Var(name), 10, 64)
}

func (o Request[T]) VarInt64Def(name string) int64 {
	i, _ := o.VarInt64(name)
	return i
}

func (o Request[T]) VarUuid(name string) (uuid.UUID, error) {
	return uuid.Parse(o.Var(name))
}

func (o Request[T]) ParamExists(name string) bool {
	return o.Param(name) != ""
}

func (o Request[T]) Param(name string) string {
	return o.r.URL.Query().Get(name)
}

func (o Request[T]) ParamBool(name string) (bool, error) {
	return parse.Bool(o.Param(name))
}

func (o Request[T]) ParamBoolDef(name string) (r bool, err error) {
	if o.ParamExists(name) {
		return o.ParamBool(name)
	}
	return
}

func (o Request[T]) ParamInt(name string) (int, error) {
	return strconv.Atoi(o.Param(name))
}

func (o Request[T]) ParamInt64(name string) (int64, error) {
	return strconv.ParseInt(o.Param(name), 10, 64)
}

func (o Request[T]) ParamUuid(name string) (uuid.UUID, error) {
	return uuid.Parse(o.Param(name))
}

func (o Request[T]) Params(v any) error {
	o.r.ParseForm()
	return schema.NewDecoder().Decode(v, o.r.Form)
}

func (o *Request[T]) DataFromParams() error {
	return o.Params(&o.Data)
}

func (o *Request[T]) DataFromJson() error {
	return json.Unmarshal(o.body, &o.Data)
}

func (o *Request[T]) readBody() (err error) {
	o.body, err = ioutil.ReadAll(o.r.Body)
	return
}

func hostFromAddress(s string) string {
	host, _, err := net.SplitHostPort(s)
	if err == nil {
		return host
	}
	return s
}
