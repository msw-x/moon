package uhttp

import (
	"errors"
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
	path           string
	router         *mux.Router
	xRemoteAddress string
	wsErrorLevel   ulog.Level
}

type OnRequest func(http.ResponseWriter, *http.Request)
type OnWebsocket func(*websocket.Conn)

func NewRouter() *Router {
	return &Router{
		log:            ulog.Empty(),
		router:         mux.NewRouter(),
		xRemoteAddress: XForwardedFor,
		wsErrorLevel:   ulog.LevelError,
	}
}

func (o Router) Branch(path string) *Router {
	o.path = o.nextPath(path)
	return &o
}

func (o *Router) WithLog(log *ulog.Log) *Router {
	o.log = log
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

func (o *Router) Handle(method string, path string, onRequest OnRequest) error {
	path = o.nextPath(path)
	o.log.Debug(RouteName(method, path))
	return o.handle(method, path, onRequest)
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
	path := o.nextPath("")
	o.log.Debugf("%s[spa]", RouteName(http.MethodGet, path))
	o.router.PathPrefix(path).Handler(NewSpaHandler(fs).WithPath(path)).Methods(http.MethodGet)
}

func (o *Router) WebSocket(path string, onWebsocket OnWebsocket) {
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

func (o *Router) ReverseProxy(path string, proxy *ReverseProxy) {
	proxy.Init()
	path = o.nextPath(path)
	o.log.Debug(RouteName("PROXY", ufmt.NotableJoinWith("->", path, proxy.Target())))
	o.handle("", path+"{path:.*}", func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = mux.Vars(r)["path"]
		proxy.ServeHTTP(w, r)
	})
}

func (o *Router) Log() *ulog.Log {
	return o.log
}

func (o *Router) Router() *mux.Router {
	return o.router
}

func (o *Router) RequestName(r *http.Request) string {
	return ProxyRequestName(r, o.xRemoteAddress)
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

func (o *Router) handle(method string, path string, onRequest OnRequest) error {
	if onRequest == nil {
		return errors.New("router on-request func is nil")
	}
	route := o.router.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		name := o.RequestName(r)
		defer uerr.Recover(func(err string) {
			o.log.Error(name, err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err))
		})
		onRequest(w, r)
	})
	if method != "" && method != "*" {
		route.Methods(method)
	}
	return nil
}
