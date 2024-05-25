package uhttp

import (
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/msw-x/moon/uerr"
	"github.com/msw-x/moon/ufmt"
	"github.com/msw-x/moon/ulog"
)

type Router struct {
	log            *ulog.Log
	id             string
	path           string
	router         *mux.Router
	xRemoteAddress string
	logRequest     bool
	wsErrorLevel   ulog.Level
}

type OnRequest func(http.ResponseWriter, *http.Request)
type OnWebsocket func(*websocket.Conn)

func NewRouter() *Router {
	return &Router{
		router:         mux.NewRouter(),
		xRemoteAddress: XForwardedFor,
		wsErrorLevel:   ulog.LevelError,
	}
}

func (o Router) Branch(path string) *Router {
	o.path = o.nextPath(path)
	return &o
}

func (o *Router) WithId(id any) *Router {
	o.id = fmt.Sprint(id)
	return o
}

func (o *Router) WithLogRequest(logRequest bool) *Router {
	o.logRequest = logRequest
	return o
}

func (o *Router) WithWebSocketErrorLevel(level ulog.Level) *Router {
	o.wsErrorLevel = level
	return o
}

func (o *Router) WithXRemoteAddress(s string) *Router {
	o.xRemoteAddress = s
	return o
}

func (o *Router) IsRoot() bool {
	return o.path == ""
}

func (o *Router) HandleFunc(path string, onRequest OnRequest) *mux.Route {
	return o.router.HandleFunc(path, onRequest)
}

func (o *Router) Handle(method string, path string, onRequest OnRequest) error {
	if onRequest == nil {
		return errors.New("router on-request func is nil")
	}
	o.init()
	path = o.nextPath(path)
	o.log.Debug(RouteName(method, path))
	o.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		name := o.RequestName(r)
		defer uerr.Recover(func(err string) {
			o.log.Error(name, err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err))
		})
		if o.logRequest {
			if r.ContentLength > 0 {
				o.log.Debug(name, ufmt.ByteSize(r.ContentLength))
			} else {
				o.log.Debug(name)
			}
		}
		onRequest(w, r)
	}).Methods(method)
	return nil
}

func (o *Router) Get(path string, onRequest OnRequest) {
	o.Handle(http.MethodGet, path, onRequest)
}

func (o *Router) Put(path string, onRequest OnRequest) {
	o.Handle(http.MethodPut, path, onRequest)
}

func (o *Router) Post(path string, onRequest OnRequest) {
	o.Handle(http.MethodPost, path, onRequest)
}

func (o *Router) Delete(path string, onRequest OnRequest) {
	o.Handle(http.MethodDelete, path, onRequest)
}

func (o *Router) Options(path string, onRequest OnRequest) {
	o.Handle(http.MethodOptions, path, onRequest)
}

func (o *Router) Files(f fs.FS) {
	o.init()
	path := o.nextPath("")
	o.log.Debugf("%s[fs]", RouteName(http.MethodGet, path))
	fs := http.FileServer(http.FS(f))
	if o.IsRoot() {
		o.router.PathPrefix(path).Handler(fs).Methods(http.MethodGet)
	} else {
		o.router.PathPrefix(path).Handler(http.StripPrefix(strings.TrimSuffix(path, "/"), fs)).Methods(http.MethodGet)
	}
}

func (o *Router) Spa(fs fs.FS) {
	o.init()
	path := o.nextPath("")
	o.log.Debugf("%s[spa]", RouteName(http.MethodGet, path))
	o.router.PathPrefix(path).Handler(NewSpaHandler(fs).WithPath(path)).Methods(http.MethodGet)
}

func (o *Router) WebSocket(path string, onWebsocket OnWebsocket) {
	o.init()
	up := websocket.Upgrader{
		ReadBufferSize:  0,
		WriteBufferSize: 0,
		CheckOrigin: func(*http.Request) bool {
			return true
		},
	}
	method := http.MethodGet
	o.log.Debug(WebSocketName(RouteName(method, o.nextPath(path))))
	o.Handle(method, path, func(w http.ResponseWriter, r *http.Request) {
		defer uerr.Recover(func(err string) {
			o.log.Error(WebSocketName(o.RequestName(r)), err)
		})
		conn, err := up.Upgrade(w, r, nil)
		if err == nil {
			onWebsocket(conn)
		} else {
			o.log.Print(o.wsErrorLevel, WebSocketName(o.RequestName(r)), err)
		}
	})
}

func (o *Router) ReverseProxy(r *ReverseProxy) {
	o.init()
	r.Connect(o)
}

func (o *Router) Log() *ulog.Log {
	o.init()
	return o.log
}

func (o *Router) Router() *mux.Router {
	return o.router
}

func (o *Router) RequestName(r *http.Request) string {
	return ProxyRequestName(r, o.xRemoteAddress)
}

func (o *Router) init() {
	if o.log == nil {
		o.log = ulog.New("router").WithID(o.id)
	}
}

func (o *Router) nextPath(path string) string {
	if path == "" {
		if o.IsRoot() {
			return "/"
		}
		return o.path
	}
	return o.path + "/" + path
}
