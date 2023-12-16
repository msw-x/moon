package urest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/msw-x/moon/ujson"
	"github.com/msw-x/moon/webs"
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
		if ctx.allowCors {
			allowCors(w)
		}
		var request Request[RequestData]
		var responce Responce[ResponceData]
		processHeader(r.Header, &responce, header)
		if responce.Ok() {
			request.r = r
			if reflect.TypeOf(request.Data).Size() != 0 {
				if method == http.MethodGet || method == http.MethodDelete {
					request.DataFromParams()
				} else {
					body, _ := ioutil.ReadAll(r.Body)
					if len(body) > 0 {
						//if ctx.logResponceBody {
						//	ctx.log.Debugf("%s:%s %s", method, r.URL, string(body))
						//}
						responce.Error = json.Unmarshal(body, &request.Data)
						if !responce.Ok() {
							responce.RefineBadRequest("unmarshal json")
						}
					} else {
						responce.BadRequest("request is empty")
					}
				}
			}
			if responce.Ok() {
				processHandle(request, &responce, handle)
			}
		}
		if !responce.Ok() {
			name := webs.RequestNameX(r, xRemoteAddress)
			if responce.muteError {
				ctx.log.Debug("[mute]", name, responce.Error)
			} else {
				ctx.log.Error(name, responce.Error)
			}
			if responce.Status == 0 {
				responce.Status = http.StatusInternalServerError
			}
		}
		if responce.Status > 0 {
			w.WriteHeader(responce.Status)
		}
		var body []byte
		contentType := "application/json"
		if responce.Ok() {
			if reflect.TypeOf(responce.Data).Size() != 0 {
				switch v := any(responce.Data).(type) {
				case void:
				case text:
					contentType = "text/plain"
					body = []byte(v)
				case image:
					contentType = "image/" + v.Type
					body = v.Data
				default:
					body, _ = ujson.MarshalLowerCase(v)
				}
			}
		} else {
			v := struct {
				Error     string `json:",omitempty"`
				ErrorCode int    `json:",omitempty"`
			}{
				Error:     fmt.Sprint(responce.Error),
				ErrorCode: responce.ErrorCode,
			}
			if v.Error != "" || v.ErrorCode != 0 {
				body, _ = ujson.MarshalLowerCase(v)
			}
		}
		w.Header().Set("Content-Type", contentType)
		w.Write(body)
	})
}
