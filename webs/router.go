package webs

import (
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
	log               *ulog.Log
	id                string
	path              string
	router            *mux.Router
	xRemoteAddress    string
	logRequest        bool
	upgradeErrorLevel ulog.Level
}

type OnRequest func(http.ResponseWriter, *http.Request)
type OnWebsocket func(*websocket.Conn)

func NewRouter() *Router {
	return &Router{
		router:            mux.NewRouter(),
		upgradeErrorLevel: ulog.LevelError,
	}
}

func (o Router) Branch(path string) *Router {
	o.path = o.uri(path)
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

func (o *Router) WithUpgradeErrorLevel(level ulog.Level) *Router {
	o.upgradeErrorLevel = level
	return o
}

func (o *Router) WithXremoteAddress(s string) *Router {
	o.xRemoteAddress = s
	return o
}

func (o *Router) IsRoot() bool {
	return o.path == ""
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
		name := o.requestName(r)
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
	uri := o.uri("")
	if o.logRequest {
		o.log.Debugf("%s[files]", RouteName(http.MethodGet, uri))
	}
	fs := http.FileServer(http.FS(files))
	if o.IsRoot() {
		o.router.PathPrefix(uri).Handler(fs)
	} else {
		o.router.PathPrefix(uri).Handler(http.StripPrefix(strings.TrimSuffix(uri, "/"), fs))
	}
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
	o.log.Debug(WebSocketName(RouteName(method, o.uri(path))))
	o.Handle(method, path, func(w http.ResponseWriter, r *http.Request) {
		defer uerr.Recover(func(err string) {
			o.log.Error(WebSocketName(o.requestName(r)), err)
		})
		conn, err := up.Upgrade(w, r, nil)
		if err == nil {
			onWebsocket(conn)
		} else {
			o.log.Print(o.upgradeErrorLevel, WebSocketName(o.requestName(r)), err)
		}
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
	if path == "" {
		if o.IsRoot() {
			return "/"
		}
		return o.path
	}
	return o.path + "/" + path
}

func (o *Router) requestName(r *http.Request) string {
	return RequestNameX(r, o.xRemoteAddress)
}

func RouteName(method, uri string) string {
	return ufmt.JoinWith(":", strings.ToUpper(method), uri)
}

func RequestName(r *http.Request) string {
	return RequestNameX(r, "")
}

func RequestNameX(r *http.Request, xRemoteAddress string) string {
	if r == nil {
		return "?"
	}
	var remoteAddress string
	if xRemoteAddress != "" && r.Header != nil {
		remoteAddress = r.Header.Get(xRemoteAddress)
	}
	if remoteAddress == "" {
		remoteAddress = r.RemoteAddr
	}
	return ufmt.JoinWith("'", remoteAddress, RouteName(r.Method, r.RequestURI))
}

func WebSocketName(name string) string {
	return ufmt.JoinWith("@", name, "ws")
}
