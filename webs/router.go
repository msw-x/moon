package webs

import (
	"fmt"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/msw-x/moon"
	"github.com/msw-x/moon/ufmt"
	"github.com/msw-x/moon/ulog"
)

type Router struct {
	log        *ulog.Log
	id         string
	path       string
	router     *mux.Router
	logRequest bool
}

type OnRequest func(http.ResponseWriter, *http.Request)
type OnWebsocket func(*websocket.Conn)

func NewRouter() (ret *Router) {
	return Router{
		router: mux.NewRouter(),
	}.Branch("")
}

func (o Router) Branch(path string) *Router {
	o.path += path
	if !strings.HasSuffix(o.path, "/") {
		o.path += "/"
	}
	return &o
}

func (o *Router) WithID(id any) *Router {
	o.id = fmt.Sprint(id)
	return o
}

func (o *Router) WithLogRequest(logRequest bool) *Router {
	o.logRequest = logRequest
	return o
}

func (o *Router) IsRoot() bool {
	return o.path == "/"
}

func (o *Router) Handle(method string, path string, onRequest OnRequest) {
	if onRequest == nil {
		panic("router on-request func is nil")
	}
	o.init()
	uri := o.uri(path)
	if o.logRequest {
		o.log.Debug(RouteName(method, uri))
	}
	o.router.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		name := RequestName(r)
		defer moon.Recover(func(err string) {
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

func (o *Router) Files(files fs.FS) {
	if o.logRequest {
		o.log.Debugf("%s[files]", RouteName(http.MethodGet, o.path))
	}
	fs := http.FileServer(http.FS(files))
	if o.IsRoot() {
		o.router.PathPrefix(o.path).Handler(fs)
	} else {
		o.router.PathPrefix(o.path).Handler(http.StripPrefix(strings.TrimSuffix(o.path, "/"), fs))
	}
}

func (o *Router) WebSocket(path string, onWebsocket OnWebsocket) {
	up := websocket.Upgrader{
		ReadBufferSize:  0,
		WriteBufferSize: 0,
	}
	method := http.MethodGet
	o.log.Debug(WebSocketName(RouteName(method, o.uri(path))))
	o.Handle(method, path, func(w http.ResponseWriter, r *http.Request) {
		defer moon.Recover(func(err string) {
			o.log.Error(WebSocketName(RequestName(r)), err)
		})
		conn, err := up.Upgrade(w, r, nil)
		moon.Strict(err, "upgrade")
		onWebsocket(conn)
	})
}

func (o *Router) Log() *ulog.Log {
	o.init()
	return o.log
}

func (o *Router) Router() *mux.Router {
	return o.router
}

func (o *Router) init() {
	if o.log == nil {
		o.log = ulog.New("router").WithID(o.id)
	}
}

func (o *Router) uri(path string) string {
	return o.path + path
}

func RouteName(method, uri string) string {
	return ufmt.JoinWith(":", strings.ToUpper(method), uri)
}

func RequestName(r *http.Request) string {
	if r == nil {
		return "?"
	}
	return ufmt.JoinWith("?", r.RemoteAddr, RouteName(r.Method, r.RequestURI))
}

func WebSocketName(name string) string {
	return ufmt.JoinWith("@", name, "ws")
}
